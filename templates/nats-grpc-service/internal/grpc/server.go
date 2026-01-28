package grpc

import (
	"context"
	"fmt"

	pb "grouter/templates/nats-grpc-service/api/proto"

	"go.uber.org/zap"
)

// HelloServer implements the HelloService defined in proto
type HelloServer struct {
	pb.UnimplementedHelloServiceServer
	logger *zap.Logger
}

// NewHelloServer creates a new server instance
func NewHelloServer(logger *zap.Logger) *HelloServer {
	return &HelloServer{
		logger: logger,
	}
}

// SayHello implements the RPC method
func (s *HelloServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	s.logger.Info("received hello request", zap.String("name", req.Name))
	return &pb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s from gRouter Production gRPC Service!", req.Name),
	}, nil
}
