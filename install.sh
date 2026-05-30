#!/bin/sh
# Indiwtf CLI installer
#
# Usage:
#   curl -fsSL https://github.com/indiwtf/indiwtf-cli/raw/main/install.sh | sh
#   wget -qO- https://github.com/indiwtf/indiwtf-cli/raw/main/install.sh | sh
#
set -e

REPO="indiwtf/indiwtf-cli"
BINARY="indiwtf"
BINARY_URL="https://github.com/${REPO}/raw/main/${BINARY}"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
INSTALL_PATH="${INSTALL_DIR}/${BINARY}"

info() { printf '\033[1;34m==>\033[0m %s\n' "$1"; }
err()  { printf '\033[1;31mError:\033[0m %s\n' "$1" >&2; exit 1; }

# Use sudo only when we cannot write to the install directory ourselves.
SUDO=""
if [ ! -w "$INSTALL_DIR" ]; then
  if command -v sudo >/dev/null 2>&1; then
    SUDO="sudo"
  else
    err "No write permission for ${INSTALL_DIR} and sudo is not available. Re-run as root or set INSTALL_DIR to a writable path."
  fi
fi

# Pick an available downloader.
if command -v curl >/dev/null 2>&1; then
  DOWNLOAD="curl -fSL -o"
elif command -v wget >/dev/null 2>&1; then
  DOWNLOAD="wget -O"
else
  err "Neither curl nor wget is installed. Please install one and try again."
fi

info "Downloading ${BINARY} to ${INSTALL_PATH}"
$SUDO $DOWNLOAD "$INSTALL_PATH" "$BINARY_URL"

info "Making ${BINARY} executable"
$SUDO chmod +x "$INSTALL_PATH"

if command -v "$BINARY" >/dev/null 2>&1; then
  info "Installed successfully. Run '${BINARY}' to get started."
else
  info "Installed to ${INSTALL_PATH}, but ${INSTALL_DIR} is not in your PATH."
  info "Add it to your PATH or run the binary directly: ${INSTALL_PATH}"
fi
