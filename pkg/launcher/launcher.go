// SPDX-License-Identifier: GPL-3.0-only
//
// This file is part of the podman-launcher project:
//
//	https://github.com/89luca89/podman-launcher
//
// # Copyright (C) 2023 podman-launcher contributors
//
// podman-launcher is free software; you can redistribute it and/or modify it
// under the terms of the GNU General Public License version 3
// as published by the Free Software Foundation.
//
// podman-launcher is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with podman-launcher; if not, see <http://www.gnu.org/licenses/>.

// Package launcher is a library to help ship the latest release of
//
//	https://github.com/mgoltzsche/podman-static/
//
// This repo builds and releases all podman components as statically linked binaries
// this will let us to easily ship the container manager without needing all the
// dependency resolution of a package manager.
// To make it work properly we need also to setup some variables, configs and paths.
// This program will take care of that, and will make sure that the podman configuration
// does not overlap with the one eventually installed by a package manager, making
// this iteration of podman isolated from the rest.
package launcher

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"syscall"
	"time"
)

var policyCommads = []string{
	"build",
	"create",
	"import",
	"load",
	"pull",
	"push",
	"run",
	"save",
	"play",
}

// NewLauncher will return an initialized launcher config with input dirs and payload.
// Refer to https://github.com/89luca89/podman-launcher/releases for the assets.tar.gz
// to embed in your application, to pass here as pack.
func NewLauncher(targetDir, tmpDir string, pack []byte) *Config {
	return &Config{
		pack:                   pack,
		targetDir:              targetDir,
		tmpDir:                 tmpDir,
		runtimeDir:             filepath.Join(tmpDir, "podman-static", strconv.Itoa(os.Getuid())),
		containersConf:         filepath.Join(targetDir, "/conf/containers/containers.conf"),
		containersRegistryConf: filepath.Join(targetDir, "/conf/containers/registries.conf"),
		containersStorageConf:  filepath.Join(targetDir, "/conf/containers/storage.conf"),
		containersPolicyJSON:   filepath.Join(targetDir, "/conf/containers/policy.json"),
	}
}

// Run will execute the payload with attached config.
func (conf *Config) Run(argv []string) error {
	if len(argv) > 1 && argv[1] == "upgrade" {
		err := os.RemoveAll(filepath.Join(conf.targetDir, "bin/podman"))
		if err != nil {
			return err
		}

		argv = []string{argv[0], "info"}
	}

	// Prepare dirs, and ensure we've unpacked everything
	err := conf.prepareFiles()
	if err != nil {
		return err
	}

	// set the --root and --runroot flags accordingly
	args := []string{
		"--root", filepath.Join(conf.targetDir, "share/containers/storage"),
		"--runroot", filepath.Join(conf.runtimeDir, "containers"),
	}

	// There isn't a config to inject the default signature policy in a place
	// other than /etc/containers/policy.jon
	//
	// So we will need to add the "--signature-policy" flag in the commands that
	// support it.
	for _, command := range policyCommads {
		if slices.Contains(argv, command) {
			index := slices.Index(argv, command)
			argv = slices.Insert(argv, index+1, []string{"--signature-policy", conf.containersPolicyJSON}...)

			break
		}
	}

	// then we just forward all the flags to the child podman command
	args = append(args, argv[1:]...)

	command := filepath.Join(conf.targetDir, "bin/podman")

	cmd := exec.Command(command, args...)

	// Forward current environment
	env := os.Environ()
	// Setup our ENV to point to our custom files:
	//		https://docs.podman.io/en/latest/markdown/podman.1.html#environment-variables
	env = append(env, "CONTAINERS_CONF="+conf.containersConf)
	env = append(env, "CONTAINERS_REGISTRIES_CONF="+conf.containersRegistryConf)
	env = append(env, "CONTAINERS_STORAGE_CONF="+conf.containersStorageConf)
	env = append(env, "PATH="+conf.targetDir+"/bin:"+os.Getenv("PATH"))
	env = append(env, "XDG_RUNTIME_DIR="+conf.runtimeDir)

	cmd.Env = env

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Sometimes podman tty get's stuck until a resize is detected, so we will
	// trigger one if the process did not exit.
	go func() {
		for {
			if cmd.Process != nil {
				syscall.Kill(cmd.Process.Pid, syscall.SIGWINCH)

				break
			}

			time.Sleep(time.Second)
		}
	}()

	defer cleanupContainerPids()

	return cmd.Run()
}
