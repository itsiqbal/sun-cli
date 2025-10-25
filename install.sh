#!/bin/bash

set -e

# Configuration
REPO="YOUR_GITHUB_USERNAME/YOUR_REPO_NAME"  # Change this to your repo
BINARY_NAME="sun"
INSTALL_DIR="/usr/local/bin"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper functions
error() {
    echo -e "${RED}Error: $1${NC}" >&2
    exit 1
}

info() {
    echo -e "${GREEN}$1${NC}"
}

warn() {
    echo -e "${YELLOW}$1${NC}"
}

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case "$OS" in
        darwin)
            OS="darwin"
            ;;
        *)
            error "This installer only supports macOS. OS detected: $OS"
            ;;
    esac
    
    case "$ARCH" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            error "Unsupported architecture: $ARCH"
            ;;
    esac
    
    info "Detected platform: $OS-$ARCH"
}

# Check if running with sudo/root for installation
check_permissions() {
    if [ ! -w "$INSTALL_DIR" ]; then
        warn "Installation directory $INSTALL_DIR is not writable."
        warn "You may need to run this script with sudo or choose a different directory."
        
        # Offer alternative installation location
        read -p "Install to ~/.local/bin instead? (y/n) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            INSTALL_DIR="$HOME/.local/bin"
            mkdir -p "$INSTALL_DIR"
            info "Installing to $INSTALL_DIR"
            
            # Check if in PATH
            if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
                warn "Note: $INSTALL_DIR is not in your PATH."
                warn "Add this to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
                echo "export PATH=\"\$PATH:$INSTALL_DIR\""
            fi
        else
            error "Installation cancelled. Please run with sudo or choose a different directory."
        fi
    fi
}

# Download the binary
download_binary() {
    BINARY_URL="https://github.com/${REPO}/releases/latest/download/${BINARY_NAME}-${OS}-${ARCH}"
    TEMP_FILE=$(mktemp)
    
    info "Downloading $BINARY_NAME from $BINARY_URL..."
    
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$BINARY_URL" -o "$TEMP_FILE" || error "Failed to download binary"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "$BINARY_URL" -O "$TEMP_FILE" || error "Failed to download binary"
    else
        error "Neither curl nor wget found. Please install one of them."
    fi
    
    # Verify download
    if [ ! -s "$TEMP_FILE" ]; then
        error "Downloaded file is empty"
    fi
    
    echo "$TEMP_FILE"
}

# Install the binary
install_binary() {
    local temp_file=$1
    local install_path="$INSTALL_DIR/$BINARY_NAME"
    
    info "Installing to $install_path..."
    
    # Make executable
    chmod +x "$temp_file" || error "Failed to make binary executable"
    
    # Move to installation directory
    mv "$temp_file" "$install_path" || error "Failed to install binary"
    
    info "Installation complete!"
}

# Verify installation
verify_installation() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        info "✓ $BINARY_NAME installed successfully!"
        info "Version: $($BINARY_NAME --version 2>/dev/null || echo 'unknown')"
        info ""
        info "Run '$BINARY_NAME --help' to get started."
    else
        warn "Installation complete, but $BINARY_NAME is not in your PATH."
        warn "You may need to:"
        warn "  1. Restart your shell"
        warn "  2. Add $INSTALL_DIR to your PATH"
        warn "  3. Run: export PATH=\"\$PATH:$INSTALL_DIR\""
    fi
}

# Main installation process
main() {
    info "Installing $BINARY_NAME..."
    echo
    
    detect_platform
    check_permissions
    
    temp_file=$(download_binary)
    install_binary "$temp_file"
    
    echo
    verify_installation
}

# Run main function
main