# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY pkg/ pkg/
COPY templates/grpc-service/ templates/grpc-service/
COPY configs/ configs/

# Proto generation usually happens before build or via a generate step.
# For this Dockerfile, we assume source files (including generated pb.go) are present or generated.
# If they are not checked in, we would need protoc here.
# For template, we assume user generates locally and commits, OR we install protoc.
# Let's simplify and assume pre-generated for now, or failing that, this is a reference.

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./templates/grpc-service/cmd/server

# Final Stage
FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/configs ./configs

EXPOSE 9090

ENTRYPOINT ["/app/server"]
