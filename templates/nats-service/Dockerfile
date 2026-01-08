# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy Go modules manifests
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY pkg/ pkg/
COPY templates/nats-service/ templates/nats-service/
COPY configs/ configs/

# Build the worker binary
# Assuming running from root context
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/worker ./templates/nats-service/cmd/worker

# Final Stage (Distroless)
FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /app/worker .
COPY --from=builder /app/configs ./configs

# Expose ports (if any, e.g. metrics)
EXPOSE 9090

ENTRYPOINT ["/app/worker"]
