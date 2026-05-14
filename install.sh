#!/usr/bin/env sh
# install.sh — install gen CLI
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/aminshahid573/gen/main/install.sh | sh
#   wget -qO- https://raw.githubusercontent.com/aminshahid573/gen/main/install.sh | sh

set -e

REPO="aminshahid573/gen"
BINARY="gen"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# ── detect OS ────────────────────────────────────────────────────────────────
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
  linux)  OS="linux"   ;;
  darwin) OS="darwin"  ;;
  freebsd) OS="freebsd" ;;
  *) echo "Unsupported OS: $OS" && exit 1 ;;
esac

# ── detect arch ──────────────────────────────────────────────────────────────
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64 | amd64) ARCH="amd64" ;;
  i386 | i686)    ARCH="386"   ;;
  aarch64 | arm64) ARCH="arm64" ;;
  armv7*)          ARCH="armv7" ;;
  armv6*)          ARCH="armv6" ;;
  *) echo "Unsupported arch: $ARCH" && exit 1 ;;
esac

# ── fetch latest tag ─────────────────────────────────────────────────────────
LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
  | grep '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')

VERSION="${LATEST#v}"
TARBALL="${BINARY}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${LATEST}/${TARBALL}"

echo "Installing ${BINARY} ${LATEST} (${OS}/${ARCH})..."

TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

curl -fsSL "$URL" -o "${TMP}/${TARBALL}"
tar -xzf "${TMP}/${TARBALL}" -C "$TMP"

# install (try sudo if needed)
if [ -w "$INSTALL_DIR" ]; then
  mv "${TMP}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
else
  sudo mv "${TMP}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
fi

chmod +x "${INSTALL_DIR}/${BINARY}"
echo "✓ ${BINARY} installed to ${INSTALL_DIR}/${BINARY}"
echo "  Run: ${BINARY} --version"
