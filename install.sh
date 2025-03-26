#!/bin/sh
set -e

# Define variables
REPO="bakito/toolbox"
INSTALL_DIR="$HOME/bin"
BIN_NAME="toolbox"

# Ensure the installation directory exists
mkdir -p "$INSTALL_DIR"

# Get the latest release tag
LATEST_TAG=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | cut -d '"' -f 4)

# Construct download URL

ARCHIVE_NAME="${BIN_NAME}_${LATEST_TAG#v}_linux_amd64.tar.gz"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_TAG/$ARCHIVE_NAME"

# Temporary directory for extraction
TMP_DIR=$(mktemp -d)

# Download the archive
echo "Downloading $BIN_NAME version $LATEST_TAG..."
curl -L "$DOWNLOAD_URL" -o "$TMP_DIR/$ARCHIVE_NAME"

# Extract the archive
tar -xzf "$TMP_DIR/$ARCHIVE_NAME" -C "$TMP_DIR"

# Move the binary to the install directory
mv "$TMP_DIR/$BIN_NAME" "$INSTALL_DIR/"

# Make it executable
chmod +x "$INSTALL_DIR/$BIN_NAME"

# Clean up
rm -rf "$TMP_DIR"

# Print success message
echo "$BIN_NAME installed successfully to $INSTALL_DIR"
echo "Add $INSTALL_DIR to your PATH if not already included"
