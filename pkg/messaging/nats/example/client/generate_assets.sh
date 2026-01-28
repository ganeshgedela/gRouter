#!/bin/bash
set -e

# Directory for assets
ASSETS_DIR="$(dirname "$0")/assets"
mkdir -p "$ASSETS_DIR"
cd "$ASSETS_DIR"

echo "=== Generating TLS Certificates ==="

# 1. Generate CA key and cert
if [ ! -f ca.key ]; then
    openssl genrsa -out ca.key 2048
    openssl req -new -x509 -days 3650 -key ca.key -out ca.pem -subj "/CN=NATS-Example-CA"
fi

# 2. Generate Server key and cert
if [ ! -f server.key ]; then
    openssl genrsa -out server.key 2048
    openssl req -new -key server.key -out server.csr -subj "/CN=localhost"
    echo "subjectAltName=DNS:localhost,IP:127.0.0.1" > extfile.cnf
    openssl x509 -req -days 365 -in server.csr -CA ca.pem -CAkey ca.key -CAcreateserial -out server.pem -extfile extfile.cnf
fi

# 3. Generate Client key and cert
if [ ! -f client-key.pem ]; then
    openssl genrsa -out client-key.pem 2048
    openssl req -new -key client-key.pem -out client.csr -subj "/CN=client"
    openssl x509 -req -days 365 -in client.csr -CA ca.pem -CAkey ca.key -CAcreateserial -out client-cert.pem
fi

echo "✅ TLS Certificates generated in $ASSETS_DIR"


echo "=== Generating NATS Credentials (Requires 'nsc') ==="

if command -v nsc &> /dev/null; then
    # Setup locally
    export NKEYS_PATH=$(pwd)/nsc/nkeys
    export NSC_HOME=$(pwd)/nsc/config
    mkdir -p "$NKEYS_PATH" "$NSC_HOME"

    # Initialize operator
    nsc init -n my-operator --dir "$NSC_HOME"
    nsc edit operator --service-url nats://localhost:4222

    # Create Account
    nsc add account -n my-account
    
    # Create User
    nsc add user -n my-user -a my-account

    # Generate Creds file
    nsc generate creds -a my-account -n my-user -o user.creds
    
    echo "✅ NATS Credentials generated: $(pwd)/user.creds"
else
    echo "⚠️ 'nsc' tool not found. Skipping Creds generation."
fi
