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
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func untar(reader io.Reader, dst string) error {
	_, err := exec.LookPath("tar")
	if err != nil {
		return fmt.Errorf("missing dependency tar")
	}

	err = os.MkdirAll(dst, 0o755)
	if err != nil {
		return err
	}

	extract := exec.Command("tar", "-xzf", "-", "-C", dst)
	extract.Stdin = reader

	return extract.Run()
}

func cleanupContainerPids() error {
	// in case of stop or rm, let's check if any zombie "container cleanup" processes
	// are left running, and remove them
	pids, err := getCleanupPid("container\000cleanup")
	if err != nil {
		return err
	}

	if len(pids) > 0 {
		for _, pid := range pids {
			err = syscall.Kill(pid, syscall.SIGTERM)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// cleanup the "cleanup container" pids that gets stuck for some reasons.
func getCleanupPid(cmdline string) ([]int, error) {
	idb := []byte(cmdline)
	result := []int{}

	processes, err := os.ReadDir("/proc")
	if err != nil {
		return result, err
	}

	// manually find in /proc a process that has the input cmdline
	for _, proc := range processes {
		mapfile := filepath.Join("/proc", proc.Name(), "/cmdline")

		filedata, err := os.ReadFile(mapfile)
		if err != nil {
			continue
		}

		if bytes.Contains(filedata, idb) {
			pid, err := strconv.ParseInt(proc.Name(), 10, 32)
			if err != nil {
				continue
			}

			result = append(result, int(pid))
		}
	}

	return result, nil
}
