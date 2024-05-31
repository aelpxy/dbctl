#!/bin/bash

GITHUB_API_URL="https://api.github.com/repos/aelpxy/dbctl/releases/latest"

download_binary() {
    local os=$1
    local arch=$2
    local latest_release_url=$(curl -s $GITHUB_API_URL | grep "browser_download_url" | grep "${os}-${arch}" | cut -d '"' -f 4)
    local output_file="dbctl-latest.tar.gz"

    if [ -z "$latest_release_url" ]; then
        if [ "$arch" == "x86_64" ]; then
            arch="amd64"
            latest_release_url=$(curl -s $GITHUB_API_URL | grep "browser_download_url" | grep "${os}-${arch}" | cut -d '"' -f 4)
        elif [ "$arch" == "aarch64" ]; then
            arch="arm64"
            latest_release_url=$(curl -s $GITHUB_API_URL | grep "browser_download_url" | grep "${os}-${arch}" | cut -d '"' -f 4)
        fi

        if [ -z "$latest_release_url" ]; then
            echo "Error: Could not find the latest release for ${os}-${arch}"
            exit 1
        fi
    fi

    echo "Downloading latest dbctl binary for ${os}-${arch}..."
    curl -L -o "$output_file" "$latest_release_url"
    tar -xzf "$output_file"
    sudo mv dbctl /usr/local/bin/
    rm "$output_file"
    echo "dbctl installed successfully!"
}

case "$(uname -s)" in
    Linux)
        download_binary "linux" "$(uname -m)"
        ;;
    Darwin)
        download_binary "darwin" "$(uname -m)"
        ;;
    *)
        echo "Unsupported operating system: $(uname -s)"
        exit 1
        ;;
esac