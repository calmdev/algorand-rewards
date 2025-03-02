#!/bin/bash
# 
# This script builds the application for Linux, Windows, and macOS.
# It uses the fyne-cross utility to cross-compile the application.
#
# Usage: ./scripts/build.sh

# Build the application for Linux.
fyne-cross linux --pull -arch=* ./cmd/algorewards

# Build the application for Windows.
fyne-cross windows --pull -arch=* ./cmd/algorewards

# Build the application for macOS.
if [[ "$(uname)" == "Darwin" ]]; then
    fyne-cross darwin --pull -arch=* ./cmd/algorewards
fi