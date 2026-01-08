# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY pkg/ pkg/
COPY templates/rest-service/ templates/rest-service/
COPY configs/ configs/

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api ./templates/rest-service/cmd/api

# Final Stage
FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /app/api .
COPY --from=builder /app/configs ./configs

EXPOSE 8080

ENTRYPOINT ["/app/api"]
