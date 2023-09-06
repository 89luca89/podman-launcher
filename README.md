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

Download the binary and put it in your $PATH

Optionally, you can name it `podman` in order to make it easier to type/use

## Usage

This launcher is transparent, so you will use it with all `podman`'s flags and so on

## Upgrade

To update, download the new release, and with the new binary run `podman-launcher upgrade`
to upgrade the embedded `podman` package.

# Compile

```console
make clean
make download
make podman-launcher
```

`make download` will download the latest bundles of `crun` and `podman-static` and
prepare them for the launcher.

`make podman-launcher` will actually compile the `main.go` and embed the targz in it.

# Use Cases

It's a nice-to-have for systems like the Steamdeck or where you're not allowed
to modify the system in any way.

Thought to be a nice fallback container engine option for Distrobox (https://github.com/89luca89/distrobox)
