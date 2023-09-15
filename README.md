# Podman Launcher

This project is a simple golang wrapper that uses `embed` to ship the latest release
of https://github.com/mgoltzsche/podman-static/

That repo builds and releases all podman components as statically linked binaries
this will let us to easily ship the container manager without needing all the
dependency resolution of a package manager.

This project will take care of shipping the release (together with `crun`) and
setting it up properly in order to work completely from $HOME, and **without overlapping**
with a native `podman` installation.

Rootful `podman` works (if needed), and will unpack a copy of the binaries in /root for it to
work.

## Installation

Download the binary, make it executable and put it in your $PATH

Optionally, you can name it `podman` in order to make it easier to type/use

## Usage

This launcher is transparent, so you will use it with all `podman`'s flags and so on

## Use in your project

You can use the `podman-launcher` as a library in your project, if you depend on `podman` and
want to embed it as a dependency.

You'll need to embed the `assets.tar.gz` (that you'll find in the release page) in your
application, and pass it to the `launcher.Config` struct for it to work

Example code:

```go
package main

import (
    _ "embed"

	"github.com/89luca89/podman-launcher/pkg/launcher"
)

var assets []byte

func main() {
    conf := launcher.NewLauncher("/home/luca-linux/.podman-launcher", "/var/tmp", assets)

    command := []string{
        "podman",
        "run", "--rm", "-ti",
        "alpine:latest",
        "/bin/sh"
    }

	err := conf.Run(command)
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			os.Exit(exiterr.ExitCode())
		}
	}
}
```

## Upgrade

To update, download the new release, and with the new binary run `podman-launcher upgrade`
to upgrade the embedded `podman` package.

## Dependencies

On the system, the only dependencies needed are the one that `podman` needs.
Specifically `iptables` and `ip6tables` for the bridge to work (not needed if using host's network namespace). 

For rootless setup to work you need `newuidmap` and `newgidmap` binaries (usually
part of the `shadow` package) and correctly set the `/etc/subuid` and `/etc/subgid`

Refer to the official documentation for further info: https://github.com/containers/podman/blob/main/docs/tutorials/rootless_tutorial.md

# Compile

```console
make clean
make
```

`make download` will download the latest bundles of `crun` and `podman-static` and
prepare them for the launcher.

`make podman-launcher` will actually compile the `main.go` and embed the targz in it.

# Use Cases

It's a nice-to-have for systems like the Steamdeck or where you're not allowed
to modify the system in any way.

Thought to be a nice fallback container engine option for Distrobox (https://github.com/89luca89/distrobox)
