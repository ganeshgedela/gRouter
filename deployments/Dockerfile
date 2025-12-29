# Build stage
FROM golang:1.24-alpine AS builder

# Install Bazel
RUN apk add --no-cache bash curl git

# Install Bazelisk (Bazel version manager)
RUN curl -Lo /usr/local/bin/bazel https://github.com/bazelbuild/bazelisk/releases/download/v1.19.0/bazelisk-linux-amd64 && \
    chmod +x /usr/local/bin/bazel

WORKDIR /workspace

# Copy workspace files
COPY WORKSPACE BUILD.bazel deps.bzl ./
COPY go.mod go.sum ./

# Copy source code
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY configs/ configs/

# Build with Bazel
RUN bazel build //cmd/service:service

# Extract binary
RUN cp bazel-bin/cmd/service/service_/service /service

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /service /app/service

# Copy config files
COPY configs/ /app/configs/

# Create non-root user
RUN addgroup -g 1000 grouter && \
    adduser -D -u 1000 -G grouter grouter && \
    chown -R grouter:grouter /app

USER grouter

# Expose ports if needed (adjust as necessary)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD pgrep service || exit 1

# Run the service
ENTRYPOINT ["/app/service"]
CMD ["--config", "/app/configs/config.yaml"]
