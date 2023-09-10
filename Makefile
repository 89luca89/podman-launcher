.PHONY: all clean download-amd64 podman-launcher-amd64 download-arm64 podman-launcher-arm64

all: clean download-amd64 podman-launcher-amd64 download-arm64 podman-launcher-arm64

CRUN_VERSION="1.8.7"
PODMAN_VERSION="4.6.1"
PODMAN_LAUNCHER_VERSION="0.0.2"

clean:
	@rm -f podman-launcher-*
	@rm -rf assets-*
	@rm -f assets*.tar.gz

launcher: podman-launcher-arm64 podman-launcher-amd64

podman-launcher-amd64:
	@rm -f podman-launcher-amd64
	@cp assets-amd64.tar.gz assets.tar.gz
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod vendor -ldflags="-s -w -X 'main.version=$${RELEASE_VERSION:-$(PODMAN_LAUNCHER_VERSION)}'" -o podman-launcher-amd64 main.go

podman-launcher-arm64:
	@rm -f podman-launcher-arm64
	@cp assets-arm64.tar.gz assets.tar.gz
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -mod vendor -ldflags="-s -w -X 'main.version=$${RELEASE_VERSION:-$(PODMAN_LAUNCHER_VERSION)}'" -o podman-launcher-arm64 main.go

download-arm64:
	rm -rf assets-arm64
	rm -f assets-arm64.tar.gz
	mkdir -p assets-arm64
	curl -L \
		"https://github.com/containers/crun/releases/download/$(CRUN_VERSION)/crun-$(CRUN_VERSION)-linux-arm64" \
		-o "./assets-arm64/crun"
	chmod +x assets-arm64/crun
	curl -L \
		"https://github.com/mgoltzsche/podman-static/releases/download/v$(PODMAN_VERSION)/podman-linux-arm64.tar.gz" \
		-o ./assets-arm64/podman-linux-arm64.tar.gz
	tar xvf ./assets-arm64/podman-linux-arm64.tar.gz -C assets-arm64
	mv assets-arm64/podman-linux-arm64/usr/local/bin/ assets-arm64/bin
	mv assets-arm64/podman-linux-arm64/usr/local/lib/ assets-arm64/lib
	mv assets-arm64/podman-linux-arm64/etc/ assets-arm64/etc
	mv assets-arm64/crun assets-arm64/bin/
	rm -rf assets-arm64/podman-linux-arm64 assets-arm64/podman-linux-arm64.tar.gz
	tar czfv assets-arm64.tar.gz -C assets-arm64 .

download-amd64:
	rm -rf assets-amd64
	rm -f assets-amd64.tar.gz
	mkdir -p assets-amd64
	curl -L \
		"https://github.com/containers/crun/releases/download/$(CRUN_VERSION)/crun-$(CRUN_VERSION)-linux-amd64" \
		-o "./assets-amd64/crun"
	chmod +x assets-amd64/crun
	curl -L \
		"https://github.com/mgoltzsche/podman-static/releases/download/v$(PODMAN_VERSION)/podman-linux-amd64.tar.gz" \
		-o ./assets-amd64/podman-linux-amd64.tar.gz
	tar xvf ./assets-amd64/podman-linux-amd64.tar.gz -C assets-amd64
	mv assets-amd64/podman-linux-amd64/usr/local/bin/ assets-amd64/bin
	mv assets-amd64/podman-linux-amd64/usr/local/lib/ assets-amd64/lib
	mv assets-amd64/podman-linux-amd64/etc/ assets-amd64/etc
	mv assets-amd64/crun assets-amd64/bin/
	rm -rf assets-amd64/podman-linux-amd64 assets-amd64/podman-linux-amd64.tar.gz
	tar czfv assets-amd64.tar.gz -C assets-amd64 .
