#!/bin/bash

# NATS Integration Tests Runner
# This script starts NATS with Docker Compose and runs integration tests

set -e

echo "Starting NATS server with Docker Compose..."
docker-compose -f pkg/messaging/nats/docker-compose.test.yml up -d

echo "Waiting for NATS to be healthy..."
timeout 30s bash -c 'until docker-compose -f pkg/messaging/nats/docker-compose.test.yml ps | grep -q "healthy"; do sleep 1; done'

echo "Running NATS integration tests..."
INTEGRATION_TESTS=true NATS_URL=nats://localhost:4222 go test -v -coverprofile=coverage_nats_integration.out ./pkg/messaging/nats/...

echo "Stopping NATS server..."
docker-compose -f pkg/messaging/nats/docker-compose.test.yml down -v

echo "Coverage report:"
go tool cover -func=coverage_nats_integration.out | tail -1
