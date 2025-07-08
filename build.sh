#!/bin/bash

# Build script for GNode releases
# Generates binaries for multiple platforms

set -e

VERSION=${1:-"1.0.0"}
OUTPUT_DIR="dist"
BINARY_NAME="gnode"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Platform configurations
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

build_for_platform() {
    local platform=$1
    local os=$(echo $platform | cut -d'/' -f1)
    local arch=$(echo $platform | cut -d'/' -f2)
    
    log_info "Building for $os/$arch..."
    
    local output_name="${BINARY_NAME}"
    if [ "$os" = "windows" ]; then
        output_name="${BINARY_NAME}.exe"
    fi
    
    local output_path="$OUTPUT_DIR/${BINARY_NAME}-${os}-${arch}"
    mkdir -p "$output_path"
    
    # Build binary
    GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build \
        -ldflags "-s -w -X main.version=$VERSION" \
        -o "$output_path/$output_name" \
        ./cmd/gnode
    
    # Create archive
    local archive_name="${BINARY_NAME}-${os}-${arch}.tar.gz"
    tar -czf "$OUTPUT_DIR/$archive_name" -C "$output_path" "$output_name"
    
    # Calculate checksum
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum "$OUTPUT_DIR/$archive_name" >> "$OUTPUT_DIR/checksums.txt"
    elif command -v shasum >/dev/null 2>&1; then
        shasum -a 256 "$OUTPUT_DIR/$archive_name" >> "$OUTPUT_DIR/checksums.txt"
    fi
    
    log_success "Built $archive_name"
}

main() {
    log_info "Building GNode v$VERSION for multiple platforms..."
    
    # Clean and create output directory
    rm -rf "$OUTPUT_DIR"
    mkdir -p "$OUTPUT_DIR"
    
    # Initialize checksums file
    echo "# GNode v$VERSION Checksums" > "$OUTPUT_DIR/checksums.txt"
    
    # Build for each platform
    for platform in "${PLATFORMS[@]}"; do
        build_for_platform "$platform"
    done
    
    # Create installation script
    cp install.sh "$OUTPUT_DIR/install.sh"
    
    log_success "All builds completed!"
    log_info "Output directory: $OUTPUT_DIR"
    log_info "Files created:"
    ls -la "$OUTPUT_DIR"
}

main "$@"