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
package main

import (
	_ "embed"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/89luca89/podman-launcher/pkg/launcher"
)

//go:embed assets.tar.gz
var pack []byte

var (
	targetDir = filepath.Join(os.Getenv("HOME"), ".local/share/podman-static")
	tmpDir    = "/var/tmp"
)

func main() {
	// if specific PODMAN_STATIC_TARGET_DIR is set, then use that instead
	if os.Getenv("PODMAN_STATIC_TARGET_DIR") != "" {
		targetDir = os.Getenv("PODMAN_STATIC_TARGET_DIR")
	}

	// if specific PODMAN_STATIC_TMP_DIR is set, then use that instead
	if os.Getenv("PODMAN_STATIC_TMP_DIR") != "" {
		tmpDir = os.Getenv("PODMAN_STATIC_TMP_DIR")
	}

	conf := launcher.NewLauncher(targetDir, tmpDir, pack)

	err := conf.Run(os.Args)
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			os.Exit(exiterr.ExitCode())
		}
	}
}
