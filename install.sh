#!/bin/bash
set -e

REPO="Aro-M/go-micro-gen"
BINARY_NAME="go-micro-gen"

echo "=============================================="
echo "  🚀 Installing go-micro-gen"
echo "=============================================="

# Detect OS & Architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
	x86_64) ARCH="amd64" ;;
	aarch64|arm64) ARCH="arm64" ;;
	i*86) ARCH="386" ;;
	*) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
esac

if [[ "$OS" == "mingw"* || "$OS" == "cygwin"* || "$OS" == "msys"* ]]; then
	OS="windows"
fi

if [ "$OS" == "windows" ] && [ "$ARCH" == "arm64" ]; then
    echo "❌ Windows ARM64 is not currently supported."
    exit 1
fi

echo "🔍 Detecting latest release for $OS/$ARCH..."
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
	echo "❌ Failed to fetch the latest release from GitHub."
	exit 1
fi

echo "✨ Latest version found: $LATEST_RELEASE"

FILENAME="${BINARY_NAME}-${OS}-${ARCH}"
if [ "$OS" == "windows" ]; then
    FILENAME="${FILENAME}.exe"
fi

DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_RELEASE}/${FILENAME}"

# Check for write permissions to /usr/local/bin
INSTALL_DIR="/usr/local/bin"
SUDO=""
if [ "$OS" != "windows" ]; then
    if [ ! -w "$INSTALL_DIR" ]; then
        echo "💡 Sudo privileges required to install to $INSTALL_DIR (you may be asked for password)"
        SUDO="sudo"
    fi
fi

echo ""
echo "⬇️  Downloading from GitHub..."
# Using curl -# for the beautiful visual progress bar!
curl -# -L "$DOWNLOAD_URL" -o "$FILENAME"
echo ""

if [ "$OS" == "windows" ]; then
    echo "📦 Moving to current directory (Windows)..."
    echo "💡 Please move ${FILENAME} to a folder in your PATH."
else
    echo "📦 Setting executable permissions and moving to $INSTALL_DIR..."
    chmod +x "$FILENAME"
    $SUDO mv "$FILENAME" "$INSTALL_DIR/$BINARY_NAME"
fi

echo "✅ Successfully installed! Run 'go-micro-gen --help' to get started."
