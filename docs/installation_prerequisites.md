# Installation Prerequisites

This guide details the prerequisites and installation steps required to build and run gRouter.

## 1. Go Programming Language

gRouter requires **Go 1.24.0** or higher.

### Installation via Version Manager (Recommended)

We recommend using a version manager like `gvm` or simply downloading the binary.

```bash
# Download Go 1.24.0
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz

# Remove previous installation
sudo rm -rf /usr/local/go

# Extract to /usr/local
sudo tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz

# Add to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH=$PATH:/usr/local/go/bin
```

Verify installation:
```bash
go version
# Output should be: go version go1.24.0 linux/amd64
```

## 2. Docker & Docker Compose

Docker is used for building container images and running dependencies like NATS locally.

### Installation (Ubuntu)

```bash
# Add Docker's official GPG key:
sudo apt-get update
sudo apt-get install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources:
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update

# Install Docker packages
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```

Verify installation:
```bash
sudo docker run hello-world
```

## 3. Bazel (Bazelisk)

gRouter uses Bazel for hermetic builds. We recommend installing `bazelisk`, which automatically manages Bazel versions.

### Installation

```bash
# Download Bazelisk binary
curl -Lo bazel https://github.com/bazelbuild/bazelisk/releases/download/v1.19.0/bazelisk-linux-amd64

# Make executable
chmod +x bazel

# Move to PATH
sudo mv bazel /usr/local/bin/
```

Verify installation:
```bash
bazel version
```

## 4. NATS Server & CLI

A NATS server is required for the event-driven architecture.

### NATS CLI Installation

```bash
# Download NATS CLI
curl -L https://github.com/nats-io/natscli/releases/download/v0.1.4/nats-0.1.4-linux-amd64.zip -o nats.zip
unzip nats.zip
sudo mv nats-0.1.4-linux-amd64/nats /usr/local/bin/
rm -rf nats.zip nats-0.1.4-linux-amd64
```

Verify installation:
```bash
nats --version
```

### Running NATS Server (Docker)

```bash
docker run -d --name nats -p 4222:4222 -p 8222:8222 nats:latest -js
```

## 5. Build Tools

Ensure basic build tools are installed.

```bash
sudo apt-get update
sudo apt-get install -y build-essential git unzip
```

## Summary Checklist

- [ ] `go version` >= 1.24.0
- [ ] `docker --version`
- [ ] `bazel --version`
- [ ] `nats --version`
