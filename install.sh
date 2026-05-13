#!/bin/bash

set -e

REPO="aminshahid573/gen"
LATEST=$(curl -sL "https://api.github.com/repos/${REPO}/releases/latest" | grep -o '"tag_name": "v[^"]*"' | cut -d'"' -f4)

if [ -z "$LATEST" ]; then
    LATEST="v1.0.0"
fi

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
esac

case "$OS" in
    mingw*|msys*|cygwin*) EXT=".exe" ;;
    *) EXT="" ;;
esac

URL="https://github.com/${REPO}/releases/download/${LATEST}/gen_${LATEST#v}_${OS}_${ARCH}.tar.gz"

echo "Downloading gen ${LATEST} for ${OS}/${ARCH}..."

mkdir -p gen_tmp
cd gen_tmp
curl -sL "$URL" | tar xz

if [ -f "gen${EXT}" ]; then
    chmod +x "gen${EXT}"
    sudo mv "gen${EXT}" /usr/local/bin/gen
    echo "gen installed successfully to /usr/local/bin/gen"
else
    echo "Error: Could not find gen binary in the archive"
    exit 1
fi

cd ..
rm -rf gen_tmp