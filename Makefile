PKG := ./cmd/dbctl
GOBUILD := CGO_ENABLED=0 go build -trimpath -ldflags="-s -w"

build:
	$(GOBUILD) -o dbctl.out $(PKG)

build_linux:
	env GOOS=linux GOARCH=amd64 $(GOBUILD) -o dbctl-linux-amd64.out $(PKG)
	env GOOS=linux GOARCH=arm64 $(GOBUILD) -o dbctl-linux-arm64.out $(PKG)

build_darwin:
	env GOOS=darwin GOARCH=amd64 $(GOBUILD) -o dbctl-darwin-amd64.out $(PKG)
	env GOOS=darwin GOARCH=arm64 $(GOBUILD) -o dbctl-darwin-arm64.out $(PKG)

build_windows:
	env GOOS=windows GOARCH=amd64 $(GOBUILD) -o dbctl-windows-amd64.exe $(PKG)
	env GOOS=windows GOARCH=arm64 $(GOBUILD) -o dbctl-windows-arm64.exe $(PKG)

build-all: build_linux build_darwin build_windows

test:
	go test -race ./...

lint:
	golangci-lint run ./...

fmt:
	golangci-lint fmt

.PHONY: build build_linux build_darwin build_windows build-all test lint fmt
