#!/bin/bash
set -e

# Cycles Development Environment Setup Script
# This script installs all necessary dependencies for building Cycles

echo "========================================="
echo "Cycles Development Environment Setup"
echo "========================================="
echo ""

# Detect OS
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$ID
    VER=$VERSION_ID
else
    echo "Cannot detect OS. Please install dependencies manually."
    exit 1
fi

echo "Detected OS: $OS $VER"
echo ""

# Function to install dependencies on Ubuntu/Debian
install_ubuntu_deps() {
    echo "Installing system dependencies for Ubuntu/Debian..."
    sudo apt-get update
    sudo apt-get install -y \
        libgl1-mesa-dev \
        libxcursor-dev \
        libxrandr-dev \
        libxinerama-dev \
        libxi-dev \
        libglfw3-dev \
        libxxf86vm-dev \
        pkg-config \
        gcc \
        g++ \
        make
    echo "✓ System dependencies installed"
}

# Function to install dependencies on Fedora/RHEL
install_fedora_deps() {
    echo "Installing system dependencies for Fedora/RHEL..."
    sudo dnf install -y \
        mesa-libGL-devel \
        libXcursor-devel \
        libXrandr-devel \
        libXinerama-devel \
        libXi-devel \
        glfw-devel \
        libXxf86vm-devel \
        pkg-config \
        gcc \
        gcc-c++ \
        make
    echo "✓ System dependencies installed"
}

# Function to install dependencies on Arch Linux
install_arch_deps() {
    echo "Installing system dependencies for Arch Linux..."
    sudo pacman -S --needed --noconfirm \
        mesa \
        libxcursor \
        libxrandr \
        libxinerama \
        libxi \
        glfw-x11 \
        pkg-config \
        gcc \
        make
    echo "✓ System dependencies installed"
}

# Install system dependencies based on OS
case "$OS" in
    ubuntu|debian|pop|linuxmint|zorin|elementary|neon)
        install_ubuntu_deps
        ;;
    fedora|rhel|centos|rocky|almalinux)
        install_fedora_deps
        ;;
    arch|manjaro|endeavouros)
        install_arch_deps
        ;;
    *)
        # Check if it's a derivative of Ubuntu/Debian by checking for apt-get
        if command -v apt-get &> /dev/null; then
            echo "Detected Debian/Ubuntu-based system (using apt-get)"
            install_ubuntu_deps
        elif command -v dnf &> /dev/null; then
            echo "Detected Red Hat-based system (using dnf)"
            install_fedora_deps
        elif command -v pacman &> /dev/null; then
            echo "Detected Arch-based system (using pacman)"
            install_arch_deps
        else
            echo "Unsupported OS: $OS"
            echo "Please install dependencies manually:"
            echo "  - OpenGL development libraries"
            echo "  - X11 development libraries (Xcursor, Xrandr, Xinerama, Xi)"
            echo "  - GLFW development libraries"
            echo "  - pkg-config, gcc, make"
            exit 1
        fi
        ;;
esac

echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed!"
    echo "Please install Go from https://golang.org/dl/"
    echo ""
    echo "Recommended installation:"
    echo "  wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz"
    echo "  sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz"
    echo "  export PATH=\$PATH:/usr/local/go/bin"
    echo ""
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo "✓ Go installed: $GO_VERSION"

# Install Go dependencies
echo ""
echo "Installing Go dependencies..."
go mod download
go mod tidy
echo "✓ Go dependencies installed"

# Verify build
echo ""
echo "Verifying build setup..."
if go build -o /tmp/cycles-test .; then
    echo "✓ Build verification successful"
    rm -f /tmp/cycles-test
else
    echo "✗ Build verification failed"
    exit 1
fi

echo ""
echo "========================================="
echo "Setup Complete!"
echo "========================================="
echo ""
echo "You can now:"
echo "  make build    - Build the application"
echo "  make run      - Build and run"
echo "  make test     - Run tests"
echo "  make clean    - Clean build artifacts"
echo ""
echo "For more information, see DEVELOPER_GUIDE.md"
