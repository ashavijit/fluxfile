#!/usr/bin/env sh

set -e

# Colors
CYAN='\033[36m'
GREEN='\033[32m'
YELLOW='\033[33m'
RED='\033[31m'
BOLD='\033[1m'
RESET='\033[0m'

print_step() {
    printf "  %s %s\n" "$1" "$2"
}

print_header() {
    printf "\n"
    printf "  ${CYAN}${BOLD}=======================================${RESET}\n"
    printf "  ${CYAN}${BOLD}          FLUX INSTALLER               ${RESET}\n"
    printf "  ${CYAN}${BOLD}=======================================${RESET}\n"
    printf "\n"
}

print_header

# Detect OS and architecture
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64)  ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    arm64)   ARCH="arm64" ;;
esac

BIN_URL="https://github.com/ashavijit/fluxfile/releases/latest/download/flux-${OS}-${ARCH}"
INSTALL_DIR="/usr/local/bin"
BIN_PATH="${INSTALL_DIR}/flux"

print_step "${CYAN}[*]${RESET}" "Detected: ${OS}/${ARCH}"

# Remove old binary if exists
if [ -f "$BIN_PATH" ]; then
    print_step "${YELLOW}[!]${RESET}" "Removing old version..."
    sudo rm -f "$BIN_PATH" 2>/dev/null || true
    print_step "${GREEN}[OK]${RESET}" "Old version removed"
fi

# Download new binary
print_step "${CYAN}[*]${RESET}" "Downloading Flux..."
if curl -fsSL "$BIN_URL" -o /tmp/flux; then
    print_step "${GREEN}[OK]${RESET}" "Download complete"
else
    print_step "${RED}[X]${RESET}" "Download failed"
    exit 1
fi

# Install binary
print_step "${CYAN}[*]${RESET}" "Installing to ${INSTALL_DIR}..."
chmod +x /tmp/flux
if sudo mv /tmp/flux "$BIN_PATH"; then
    print_step "${GREEN}[OK]${RESET}" "Installation complete"
else
    print_step "${RED}[X]${RESET}" "Installation failed (try with sudo)"
    exit 1
fi

# Verify installation
printf "\n"
printf "  ${GREEN}${BOLD}=======================================${RESET}\n"
printf "  ${GREEN}${BOLD}       INSTALLATION COMPLETE           ${RESET}\n"
printf "  ${GREEN}${BOLD}=======================================${RESET}\n"
printf "\n"

print_step "${CYAN}[>]${RESET}" "Installed to: ${BIN_PATH}"

if VERSION=$(flux -v 2>&1); then
    print_step "${CYAN}[>]${RESET}" "Version: ${VERSION}"
else
    print_step "${YELLOW}[!]${RESET}" "Could not verify version"
fi

printf "\n"
printf "  ${CYAN}Usage:${RESET}\n"
printf "    flux init            Create new FluxFile\n"
printf "    flux build           Run build task\n"
printf "    flux -l              List all tasks\n"
printf "    flux logs            View execution logs\n"
printf "\n"
