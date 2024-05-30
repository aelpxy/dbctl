build:
	go build -ldflags="-s -w" -o dbctl.out .

build_linux:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dbctl-linux-amd64.out .
	env GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dbctl-linux-arm64.out .

build_darwin:
	env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dbctl-darwin-amd64.out .
	env GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dbctl-darwin-arm64.out .

build_windows:
	env GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dbctl-windows-amd64.exe .
	env GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o dbctl-windows-arm64.exe .

build-all: build_linux build_darwin build_windows

.PHONY: build build_linux build_darwin build_windows