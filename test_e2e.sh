#!/bin/bash
# E2E Test Suite for Factory Pattern Implementation
# Tests NATS, REST, and gRPC service templates

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$SCRIPT_DIR"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🚀 Factory Pattern E2E Test Suite"
echo "===================================="
echo ""

# Function to print test results
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
        exit 1
    fi
}

# Test 1: Build all services
echo "📦 Test 1: Building all services..."
cd "$PROJECT_ROOT/templates/nats-service"
go build -o bin/worker cmd/worker/main.go
print_result $? "NATS worker built"

cd "$PROJECT_ROOT/templates/rest-service"
go build -o bin/api cmd/api/main.go
print_result $? "REST API built"

cd "$PROJECT_ROOT/templates/grpc-service"
go build -o bin/server cmd/server/main.go
print_result $? "gRPC server built"

echo ""

# Test 2: Check NATS server is running
echo "📡 Test 2: Checking NATS server..."
if pgrep -x "nats-server" > /dev/null; then
    print_result 0 "NATS server is running"
else
    echo -e "${YELLOW}⚠️  NATS server not running, skipping NATS worker test${NC}"
fi

echo ""

# Test 3: Test REST service
echo "🌐 Test 3: Testing REST service..."
cd "$PROJECT_ROOT/templates/rest-service"
./bin/api --config internal/config/config.yaml &
REST_PID=$!
sleep 2

# Test health endpoint
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/v1/health)
if [ "$HTTP_CODE" -eq 200 ]; then
    print_result 0 "REST health endpoint returned 200"
else
    print_result 1 "REST health endpoint failed (got $HTTP_CODE)"
fi

# Test liveness
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health/live)
if [ "$HTTP_CODE" -eq 200 ]; then
    print_result 0 "REST liveness probe returned 200"
else
    print_result 1 "REST liveness probe failed (got $HTTP_CODE)"
fi

# Test readiness
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health/ready)
if [ "$HTTP_CODE" -eq 200 ]; then
    print_result 0 "REST readiness probe returned 200"
else
    print_result 1 "REST readiness probe failed (got $HTTP_CODE)"
fi

# Cleanup REST
kill $REST_PID
wait $REST_PID 2>/dev/null || true

echo ""

# Test 4: Test gRPC service
echo "🔌 Test 4: Testing gRPC service..."
cd "$PROJECT_ROOT/templates/grpc-service"
./bin/server --config internal/config/config.yaml &
GRPC_PID=$!
sleep 2

# Build and run gRPC client
go build -o bin/client cmd/client/main.go
./bin/client > /tmp/grpc_test_output.txt 2>&1
GRPC_RESULT=$?

if [ $GRPC_RESULT -eq 0 ] && grep -q "All gRPC tests completed successfully" /tmp/grpc_test_output.txt; then
    print_result 0 "gRPC client successfully called all RPCs"
else
    print_result 1 "gRPC client failed"
fi

# Cleanup gRPC
kill $GRPC_PID
wait $GRPC_PID 2>/dev/null || true

echo ""

# Test 5: Verify factory pattern in logs
echo "🏭 Test 5: Verifying factory pattern..."

cd "$PROJECT_ROOT/templates/rest-service"
./bin/api --config internal/config/config.yaml > /tmp/rest_startup.log 2>&1 &
REST_PID=$!
sleep 2

if grep -q "Built services" /tmp/rest_startup.log && \
   grep -q "users" /tmp/rest_startup.log && \
   grep -q "orders" /tmp/rest_startup.log; then
    print_result 0 "REST service discovered handlers via factory pattern"
else
    print_result 1 "REST service factory pattern failed"
fi

kill $REST_PID
wait $REST_PID 2>/dev/null || true

echo ""
echo "🎉 All E2E tests passed!"
echo ""
echo "Summary:"
echo "  ✅ All services build successfully"
echo "  ✅ REST health endpoints working"
echo "  ✅ gRPC server and client working"
echo "  ✅ Factory pattern discovering services"
echo ""
echo "Factory pattern implementation is production-ready! 🚀"
