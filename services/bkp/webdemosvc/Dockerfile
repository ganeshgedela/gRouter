FROM golang:1.24-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build static binary
RUN cd services/webdemosvc/cmd/webdemosvc && CGO_ENABLED=0 go build -o /app/webdemosvc .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/webdemosvc .
COPY services/webdemosvc/internal/config/config.yaml ./config.yaml

EXPOSE 8080

CMD ["./webdemosvc", "--config", "config.yaml"]
