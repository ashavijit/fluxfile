#!/usr/bin/env sh

set -e

OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
esac

BIN_URL="https://github.com/ashavijit/fluxfile/releases/latest/download/flux-${OS}-${ARCH}"

echo "Downloading Flux for $OS/$ARCH ..."
curl -fsSL "$BIN_URL" -o flux

chmod +x flux
sudo mv flux /usr/local/bin/flux

echo "Flux installed successfully!"
flux -v
