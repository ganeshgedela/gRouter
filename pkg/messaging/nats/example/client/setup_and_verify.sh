# Configuration
SCRIPT_DIR="$(dirname "$0")"
ASSETS_DIR="$SCRIPT_DIR/assets"

# Always change to the script directory so relative paths in Go code work
cd "$SCRIPT_DIR"
EXAMPLE_DIR="."
ASSETS_DIR="./assets"

# Export NSC environment variables so nsc sees our local generated accounts
export NSC_HOME="$(pwd)/assets/nsc/config"
export NKEYS_PATH="$(pwd)/assets/nsc/nkeys"

echo "============================================"
echo "   NATS Client Auth Verification Helper"
echo "============================================"
echo "1. Verify Token Auth"
echo "2. Verify User/Pass Auth"
echo "3. Verify Creds (NKEY/JWT) Auth"
echo "4. Verify Mutual TLS Auth"
echo "q. Quit"
echo "============================================"
read -p "Select validation scenario: " choice

case $choice in
    1)
        echo "--> Starting NATS Server with Token..."
        nats-server -auth "my-secret-token" -p 4222 &
        SERVER_PID=$!
        sleep 2
        
        echo "--> Running Client..."
        go run $EXAMPLE_DIR/token/main.go
        
        kill $SERVER_PID
        ;;
    2)
        echo "--> Starting NATS Server with User/Pass..."
        nats-server --user myuser --pass mypassword -p 4222 &
        SERVER_PID=$!
        sleep 2
        
        echo "--> Running Client..."
        go run $EXAMPLE_DIR/userpass/main.go
        
        kill $SERVER_PID
        ;;
    3)
        echo "--> Starting NATS Server with Creds (Operator Mode)..."
        # We need to run nats-server pointing to the resolver
        # For simplicity in this demo, since we just generated a creds file but not a full resolver config for the server,
        # we might need to export the account config.
        # However, nsc can generate a server config!
        
        echo "Generating server config from nsc..."
        nsc generate config --mem-resolver --force --config-file $ASSETS_DIR/server.conf
        
        nats-server -c $ASSETS_DIR/server.conf -p 4222 &
        SERVER_PID=$!
        sleep 2
        
        echo "--> Running Client..."
        # Update connection string in creds_auth.go if needed, but it points to localhost:4222
        # Ensure creds path is correct
        go run $EXAMPLE_DIR/creds/main.go
        
        kill $SERVER_PID
        ;;
    4)
        echo "--> Starting NATS Server with TLS..."
        nats-server --tls --tlscert $ASSETS_DIR/server.pem --tlskey $ASSETS_DIR/server.key --tlscacert $ASSETS_DIR/ca.pem -p 4222 &
        SERVER_PID=$!
        sleep 2
        
        echo "--> Running Client..."
        go run $EXAMPLE_DIR/tls/main.go
        
        kill $SERVER_PID
        ;;
    q)
        exit 0
        ;;
    *)
        echo "Invalid option"
        ;;
esac

echo
echo "Done."
