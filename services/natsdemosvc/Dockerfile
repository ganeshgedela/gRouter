FROM golang:1.24-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build static binary
RUN cd services/natsdemosvc/cmd/natsdemosvc && CGO_ENABLED=0 go build -o /app/natsdemosvc .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/natsdemosvc .
COPY services/natsdemosvc/internal/config/config.yaml ./config.yaml

EXPOSE 8080

CMD ["./natsdemosvc", "--config", "config.yaml"]
