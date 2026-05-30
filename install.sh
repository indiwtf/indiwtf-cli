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
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
INSTALL_PATH="${INSTALL_DIR}/${BINARY}"

info() { printf '\033[1;34m==>\033[0m %s\n' "$1"; }
err()  { printf '\033[1;31mError:\033[0m %s\n' "$1" >&2; exit 1; }

# Pick an available downloader.
if command -v curl >/dev/null 2>&1; then
  HAS_CURL=1
elif command -v wget >/dev/null 2>&1; then
  HAS_CURL=0
else
  err "Neither curl nor wget is installed. Please install one and try again."
fi

# fetch <url> downloads to stdout.
fetch() {
  if [ "$HAS_CURL" = "1" ]; then
    curl -fsSL "$1"
  else
    wget -qO- "$1"
  fi
}

# download <url> <dest> saves to a file (via sudo when needed).
download() {
  if [ "$HAS_CURL" = "1" ]; then
    $SUDO curl -fSL -o "$2" "$1"
  else
    $SUDO wget -O "$2" "$1"
  fi
}

# Use sudo only when we cannot write to the install directory ourselves.
SUDO=""
if [ ! -w "$INSTALL_DIR" ]; then
  if command -v sudo >/dev/null 2>&1; then
    SUDO="sudo"
  else
    err "No write permission for ${INSTALL_DIR} and sudo is not available. Re-run as root or set INSTALL_DIR to a writable path."
  fi
fi

# Resolve the version to install. Honor an explicit VERSION override,
# otherwise query GitHub for the latest release tag.
VERSION="${VERSION:-}"
if [ -z "$VERSION" ]; then
  info "Resolving the latest release"
  VERSION=$(fetch "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name"' \
    | head -n 1 \
    | sed -E 's/.*"tag_name"[[:space:]]*:[[:space:]]*"([^"]+)".*/\1/')
fi

if [ -z "$VERSION" ]; then
  err "Could not determine the latest release. Set VERSION=<tag> and try again."
fi

BINARY_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY}"

info "Downloading ${BINARY} ${VERSION} to ${INSTALL_PATH}"
download "$BINARY_URL" "$INSTALL_PATH"

info "Making ${BINARY} executable"
$SUDO chmod +x "$INSTALL_PATH"

if command -v "$BINARY" >/dev/null 2>&1; then
  info "Installed successfully. Run '${BINARY}' to get started."
else
  info "Installed to ${INSTALL_PATH}, but ${INSTALL_DIR} is not in your PATH."
  info "Add it to your PATH or run the binary directly: ${INSTALL_PATH}"
fi
