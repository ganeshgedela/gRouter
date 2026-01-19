#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
ORANGE='\033[0;33m' # Using orange/brown for warnings as a close alternative since standard yellow can be hard to read on white
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

# Function to check if a command exists
command_exists() {
    command -v "$1" &> /dev/null
}

# Install Go
install_go() {
    local GO_VERSION="1.24.0"
    if command_exists go; then
        CURRENT_VERSION=$(go version | awk '{print $3}')
        log_info "Go is already installed: $CURRENT_VERSION"
        if [[ "$CURRENT_VERSION" != "go$GO_VERSION" ]]; then
             log_warn "Installed version ($CURRENT_VERSION) differs from recommended (go$GO_VERSION)."
        fi
    else
        log_info "Installing Go $GO_VERSION..."
        if [ "$EUID" -ne 0 ]; then
             log_warn "Installation requires sudo access. Please run script with sudo or install manually."
             exit 1
        fi
        
        curl -LO "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
        rm -rf /usr/local/go
        tar -C /usr/local -xzf "go${GO_VERSION}.linux-amd64.tar.gz"
        rm "go${GO_VERSION}.linux-amd64.tar.gz"
        
        # Add to PATH temporarily for this script
        export PATH=$PATH:/usr/local/go/bin
        
        log_warn "Please add the following to your shell profile (~/.bashrc or ~/.zshrc):"
        echo "export PATH=\$PATH:/usr/local/go/bin"
    fi
}

# Install Docker
install_docker() {
    if command_exists docker; then
        log_info "Docker is already installed."
    else
        log_info "Installing Docker..."
        if [ "$EUID" -ne 0 ]; then
             log_warn "Installation requires sudo access. Please run script with sudo or install manually."
             exit 1
        fi

        # Add Docker's official GPG key:
        apt-get update
        apt-get install -y ca-certificates curl
        install -m 0755 -d /etc/apt/keyrings
        curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
        chmod a+r /etc/apt/keyrings/docker.asc

        # Add the repository to Apt sources:
        echo \
          "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
          $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
          tee /etc/apt/sources.list.d/docker.list > /dev/null
        apt-get update
        apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
        
        log_info "Docker installed. You may need to add your user to the 'docker' group."
    fi
}

# Install Bazel (Bazelisk)
install_bazel() {
    if command_exists bazel; then
        log_info "Bazel (or Bazelisk) is already installed."
    else
        log_info "Installing Bazelisk..."
        if [ "$EUID" -ne 0 ]; then
             log_warn "Installation requires sudo access. Please run script with sudo or install manually."
             exit 1
        fi
        
        curl -Lo bazel https://github.com/bazelbuild/bazelisk/releases/download/v1.19.0/bazelisk-linux-amd64
        chmod +x bazel
        mv bazel /usr/local/bin/
    fi
}

# Install NATS CLI
install_nats() {
    if command_exists nats; then
        log_info "NATS CLI is already installed."
    else
        log_info "Installing NATS CLI..."
        if [ "$EUID" -ne 0 ]; then
             log_warn "Installation requires sudo access. Please run script with sudo or install manually."
             exit 1
        fi
        
        curl -L https://github.com/nats-io/natscli/releases/download/v0.1.4/nats-0.1.4-linux-amd64.zip -o nats.zip
        unzip -o nats.zip
        mv nats-0.1.4-linux-amd64/nats /usr/local/bin/
        rm -rf nats.zip nats-0.1.4-linux-amd64
    fi
}

install_dependencies() {
     log_info "Installing dependencies (build-essential, git, unzip)..."
     if [ "$EUID" -eq 0 ]; then
        apt-get update && apt-get install -y build-essential git unzip
     else
         log_warn "Skipping dependency update (not root). Ensure build-essential, git, and unzip are installed."
     fi
}

# Main execution
echo "Starting gRouter prerequisites installation..."

# Check OS
if [ -f /etc/os-release ]; then
    . /etc/os-release
    if [[ "$ID" != "ubuntu" && "$ID" != "debian" && "$ID_LIKE" != *"ubuntu"* && "$ID_LIKE" != *"debian"* ]]; then
        log_warn "This script is designed for Ubuntu/Debian. Your system ($ID) might require different commands."
    fi
fi

install_dependencies
install_go
install_docker
install_bazel
install_nats

echo -e "${GREEN}Prerequisites check/installation complete!${NC}"
echo "Please restart your shell or source your profile to ensure all paths are updated."
