#!/bin/sh
# Claude PP installation script for macOS and Linux

set -e

# Configuration
REPO="DandaAkhilReddy/claude_pp"
INSTALL_DIR="${CLAUDE_PP_INSTALL_DIR:-$HOME/.local/bin}"

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "darwin" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *)       echo "unknown" ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64) echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *)            echo "unknown" ;;
    esac
}

# Get latest release version
get_latest_version() {
    curl -sL "https://api.github.com/repos/${REPO}/releases/latest" | \
        grep '"tag_name":' | \
        sed -E 's/.*"([^"]+)".*/\1/'
}

# Main installation
main() {
    OS=$(detect_os)
    ARCH=$(detect_arch)

    if [ "$OS" = "unknown" ] || [ "$ARCH" = "unknown" ]; then
        echo "Error: Unsupported platform: $(uname -s) $(uname -m)"
        exit 1
    fi

    echo "Detected platform: ${OS}/${ARCH}"

    # Get latest version
    VERSION=$(get_latest_version)
    if [ -z "$VERSION" ]; then
        VERSION="v0.1.0"
        echo "Warning: Could not detect latest version, using ${VERSION}"
    else
        echo "Latest version: ${VERSION}"
    fi

    # Construct download URL
    BINARY="claude_pp-${OS}-${ARCH}"
    if [ "$OS" = "windows" ]; then
        BINARY="${BINARY}.exe"
    fi

    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY}"

    # Create install directory
    mkdir -p "$INSTALL_DIR"

    # Download binary
    echo "Downloading ${BINARY}..."
    if command -v curl >/dev/null 2>&1; then
        curl -sL "$DOWNLOAD_URL" -o "${INSTALL_DIR}/claude_pp"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "$DOWNLOAD_URL" -O "${INSTALL_DIR}/claude_pp"
    else
        echo "Error: Neither curl nor wget found. Please install one of them."
        exit 1
    fi

    # Make executable
    chmod +x "${INSTALL_DIR}/claude_pp"

    echo ""
    echo "Claude PP installed to ${INSTALL_DIR}/claude_pp"

    # Check if install dir is in PATH
    case ":$PATH:" in
        *":${INSTALL_DIR}:"*)
            echo ""
            echo "Run 'claude_pp setup' to configure for your AI coding tool."
            ;;
        *)
            echo ""
            echo "Add ${INSTALL_DIR} to your PATH:"
            echo ""
            echo "  export PATH=\"\$PATH:${INSTALL_DIR}\""
            echo ""
            echo "Then run 'claude_pp setup' to configure for your AI coding tool."
            ;;
    esac
}

main "$@"
