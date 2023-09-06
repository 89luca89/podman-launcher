.PHONY: all podman-launcher download

all: clean download podman-launcher

CRUN_VERSION="1.8.7"
PODMAN_VERSION="4.6.1"
PODMAN_LAUNCHER_VERSION="0.0.1"

clean:
	@rm -f podman-launcher
	@rm -rf assets
	@rm -f assets.tar.gz

podman-launcher:
	@rm -f podman-launcher
	CGO_ENABLED=0 go build -mod vendor -ldflags="-s -w -X 'main.version=$${RELEASE_VERSION:-$(PODMAN_LAUNCHER_VERSION)}'" -o podman-launcher main.go

download:
	rm -rf assets
	rm -f assets.tar.gz
	mkdir -p assets
	curl -L \
		"https://github.com/containers/crun/releases/download/$(CRUN_VERSION)/crun-$(CRUN_VERSION)-linux-amd64" \
		-o "./assets/crun"
	chmod +x assets/crun
	curl -L \
		"https://github.com/mgoltzsche/podman-static/releases/download/v$(PODMAN_VERSION)/podman-linux-amd64.tar.gz" \
		-o ./assets/podman-linux-amd64.tar.gz
	tar xvf ./assets/podman-linux-amd64.tar.gz -C assets
	mv assets/podman-linux-amd64/usr/local/bin/ assets/bin
	mv assets/podman-linux-amd64/usr/local/lib/ assets/lib
	mv assets/podman-linux-amd64/etc/ assets/etc
	mv assets/crun assets/bin/
	rm -rf assets/podman-linux-amd64 assets/podman-linux-amd64.tar.gz
	tar czfv assets.tar.gz -C assets .
