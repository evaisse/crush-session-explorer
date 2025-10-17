#!/bin/bash

# Crush Session Explorer Installer
# This script downloads and installs the latest release of crush-md

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# GitHub repository details
REPO="evaisse/crush-session-explorer"
BINARY_NAME="crush-md"

# Default installation directory
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Utility functions
info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# Detect OS
detect_os() {
    local os
    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        CYGWIN*|MINGW*|MSYS*) os="windows" ;;
        *)          error "Unsupported operating system: $(uname -s)" ;;
    esac
    echo "$os"
}

# Detect architecture
detect_arch() {
    local arch
    case "$(uname -m)" in
        x86_64|amd64)   arch="amd64" ;;
        aarch64|arm64)  arch="arm64" ;;
        *)              error "Unsupported architecture: $(uname -m)" ;;
    esac
    echo "$arch"
}

# Get the latest release version from GitHub API
get_latest_version() {
    local version
    if command -v curl &> /dev/null; then
        version=$(curl -sSfL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    elif command -v wget &> /dev/null; then
        version=$(wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        error "Neither curl nor wget is available. Please install one of them."
    fi
    
    if [ -z "$version" ]; then
        error "Failed to get latest version from GitHub API"
    fi
    
    echo "$version"
}

# Download file using curl or wget
download_file() {
    local url="$1"
    local output="$2"
    
    info "Downloading from: $url"
    
    if command -v curl &> /dev/null; then
        curl -sSfL "$url" -o "$output"
    elif command -v wget &> /dev/null; then
        wget -qO "$output" "$url"
    else
        error "Neither curl nor wget is available. Please install one of them."
    fi
}

# Main installation function
main() {
    info "Starting Crush Session Explorer installation..."
    
    # Detect system information
    local os
    os=$(detect_os)
    local arch
    arch=$(detect_arch)
    info "Detected OS: $os"
    info "Detected Architecture: $arch"
    
    # Get latest version
    local version
    version=$(get_latest_version)
    info "Latest version: $version"
    
    # Build binary name based on OS and arch
    local binary_suffix="${os}-${arch}"
    if [ "$os" = "windows" ]; then
        binary_suffix="${binary_suffix}.exe"
    fi
    local binary_file="${BINARY_NAME}-${binary_suffix}"
    
    # Build download URL
    local download_url="https://github.com/${REPO}/releases/download/${version}/${binary_file}"
    
    # Create temporary directory
    local tmp_dir
    tmp_dir=$(mktemp -d)
    trap 'rm -rf "$tmp_dir"' EXIT
    
    local tmp_file="${tmp_dir}/${binary_file}"
    
    # Download binary
    info "Downloading ${binary_file}..."
    if ! download_file "$download_url" "$tmp_file"; then
        error "Failed to download binary from $download_url"
    fi
    
    # Make binary executable
    chmod +x "$tmp_file"
    
    # Install binary
    info "Installing to ${INSTALL_DIR}/${BINARY_NAME}..."
    
    # Check if we need sudo
    if [ -w "$INSTALL_DIR" ]; then
        mv "$tmp_file" "${INSTALL_DIR}/${BINARY_NAME}"
    else
        warn "Requires elevated privileges to install to ${INSTALL_DIR}"
        sudo mv "$tmp_file" "${INSTALL_DIR}/${BINARY_NAME}"
    fi
    
    # Verify installation
    if ! command -v "$BINARY_NAME" &> /dev/null; then
        warn "Binary installed but not found in PATH"
        warn "You may need to add ${INSTALL_DIR} to your PATH"
        info "Installation completed: ${INSTALL_DIR}/${BINARY_NAME}"
    else
        info "Installation successful!"
        info "Run '${BINARY_NAME} --help' to get started"
        
        # Show version
        local installed_version
        installed_version=$(${BINARY_NAME} --version 2>&1 || echo "version check failed")
        info "Installed version: $installed_version"
    fi
}

# Run main function only if script is executed directly
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    main
fi
