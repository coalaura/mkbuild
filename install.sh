#!/bin/bash
set -e

OS=$(uname -s | tr 'A-Z' 'a-z')

ARCH=$(uname -m)
case "$ARCH" in
	x86_64)
		ARCH=amd64
		;;
	aarch64|arm64)
		ARCH=arm64
		;;
	*)
		echo "Unsupported architecture: $ARCH" >&2
		exit 1
		;;
esac

echo "Resolving latest version..."

VERSION=$(curl -sL https://api.github.com/repos/coalaura/mkbuild/releases/latest | grep -Po '"tag_name": *"\K.*?(?=")')

if ! printf '%s\n' "$VERSION" | grep -Eq '^v[0-9]+\.[0-9]+\.[0-9]+$'; then
	echo "Error: '$VERSION' is not in vMAJOR.MINOR.PATCH format" >&2
	exit 1
fi

rm -f /tmp/mkbuild

BIN="mkbuild_${OS}_${ARCH}"
URL="https://github.com/coalaura/mkbuild/releases/download/${VERSION}/${BIN}"

echo "Downloading ${BIN} (${VERSION})..."

if ! curl -sL "$URL" -o /tmp/mkbuild; then
	echo "Error: failed to download $URL" >&2
	exit 1
fi

trap 'rm -f /tmp/mkbuild' EXIT

chmod +x /tmp/mkbuild

echo "Installing to /usr/local/bin/mkbuild requires sudo"

if ! sudo install -m755 /tmp/mkbuild /usr/local/bin/mkbuild; then
	echo "Error: install failed" >&2
	exit 1
fi

echo "mkbuild $VERSION installed to /usr/local/bin/mkbuild"