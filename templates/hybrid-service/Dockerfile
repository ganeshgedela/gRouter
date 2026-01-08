FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY pkg/ pkg/
COPY templates/ templates/
COPY configs/ configs/

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/hybrid ./templates/hybrid-service/cmd/server

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=builder /app/hybrid .
COPY --from=builder /app/configs ./configs
EXPOSE 8080
ENTRYPOINT ["/app/hybrid"]
