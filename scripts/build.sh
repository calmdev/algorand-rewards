#!/bin/bash
# 
# This script builds the application for Linux, Windows, and macOS.
# It uses the fyne-cross utility to cross-compile the application.
#
# Usage: ./scripts/build.sh [platform] [arch]
# platform: linux, windows, darwin, all (default: current platform)
# arch: amd64, arm64, 386, etc. (default: current architecture)

TARGET_OS=${1:-$(uname | tr '[:upper:]' '[:lower:]')}
TARGET_ARCH=${2:-$(uname -m)}

# Map architecture names to fyne-cross compatible values
case $TARGET_ARCH in
    x86_64)
        TARGET_ARCH="amd64"
        ;;
    arm64|aarch64)
        TARGET_ARCH="arm64"
        ;;
    i386)
        TARGET_ARCH="386"
        ;;
    all)
        if [[ "$(uname)" == "Darwin" ]]; then
            TARGET_ARCH="amd64,arm64"
        else
            TARGET_ARCH="amd64,arm64,386"
        fi
        ;;
    *)
        echo "Unknown architecture: $TARGET_ARCH"
        exit 1
        ;;
esac

build() {
    local platform=$1
    local arch=$2
    fyne-cross $platform --pull -arch=$arch ./cmd/algorewards
}

case $TARGET_OS in
    linux)
        build linux $TARGET_ARCH
        ;;
    windows)
        build windows $TARGET_ARCH
        ;;
    darwin)
        build darwin $TARGET_ARCH
        ;;
    all)
        build linux $TARGET_ARCH
        build windows $TARGET_ARCH
        if [[ "$(uname)" == "Darwin" ]]; then
            build darwin $TARGET_ARCH
        fi
        ;;
    *)
        echo "Unknown platform: $TARGET_OS"
        exit 1
        ;;
esac