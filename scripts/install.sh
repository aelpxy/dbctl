#!/usr/bin/env bash

set -euo pipefail

GITHUB_API_URL="https://api.github.com/repos/aelpxy/dbctl/releases/latest"
INSTALL_DIR="/usr/local/bin"

case "$(uname -s)" in
    Linux) os="linux" ;;
    Darwin) os="darwin" ;;
    *)
        echo "Unsupported operating system: $(uname -s)"
        exit 1
        ;;
esac

case "$(uname -m)" in
    x86_64 | amd64) arch="amd64" ;;
    aarch64 | arm64) arch="arm64" ;;
    *)
        echo "Unsupported architecture: $(uname -m)"
        exit 1
        ;;
esac

release_url=$(curl -fsSL "$GITHUB_API_URL" | grep "browser_download_url" | grep "${os}-${arch}.tar.gz" | cut -d '"' -f 4 | head -n 1)

if [ -z "$release_url" ]; then
    echo "Error: Could not find the latest release for ${os}-${arch}"
    exit 1
fi

tmp_dir=$(mktemp -d)
trap 'rm -rf "$tmp_dir"' EXIT

echo "Downloading latest dbctl binary for ${os}-${arch}..."
curl -fsSL -o "$tmp_dir/dbctl.tar.gz" "$release_url"
tar -xzf "$tmp_dir/dbctl.tar.gz" -C "$tmp_dir"

if [ -w "$INSTALL_DIR" ]; then
    mv "$tmp_dir/dbctl" "$INSTALL_DIR/dbctl"
else
    sudo mv "$tmp_dir/dbctl" "$INSTALL_DIR/dbctl"
fi

echo "dbctl installed successfully!"
