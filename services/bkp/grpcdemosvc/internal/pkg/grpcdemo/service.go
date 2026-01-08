package grpcdemo

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"grouter/pkg/manager"
	gserver "grouter/pkg/messaging/grpc"
	pb "grouter/services/grpcdemosvc/api/proto"

	"github.com/go-viper/mapstructure/v2"
)

func init() {
	manager.RegisterFactory("grpcdemo", func(deps manager.Deps) (manager.Service, error) {
		return NewGRPCDemo(deps), nil
	})
}

// GRPCDemoConfig holds configuration for the GRPCDemo service
type GRPCDemoConfig struct {
	Enabled bool `mapstructure:"enabled"`
	Port    int  `mapstructure:"port"`
}

// GRPCDemo implements manager.Service
type GRPCDemo struct {
	logger *zap.Logger
	server *gserver.Server
	config GRPCDemoConfig

	status    manager.HealthStatus
	lastError error
}

// Ensure UnimplementedGreeterServer is embedded
type greeterServer struct {
	pb.UnimplementedGreeterServer
	logger *zap.Logger
}

// SayHello implements helloworld.GreeterServer
func (s *greeterServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	s.logger.Info("Received request", zap.String("name", in.GetName()))
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

// NewGRPCDemo creates a new GRPCDemo service
func NewGRPCDemo(deps manager.Deps) *GRPCDemo {
	// Parse config from deps.Config.Services["grpcdemo"]
	var cfg GRPCDemoConfig
	if svcCfg, ok := deps.Config.Services["grpcdemo"]; ok {
		if err := mapstructure.Decode(svcCfg, &cfg); err != nil {
			deps.Logger.Error("Failed to decode grpcdemo config", zap.Error(err))
		}
	}

	return &GRPCDemo{
		logger: deps.Logger,
		config: cfg,
		status: manager.StatusCreated,
	}
}

// Name returns the name of the service
func (s *GRPCDemo) Name() string {
	return "grpcdemo"
}

func (s *GRPCDemo) ID() string {
	return "grpcdemo"
}

// Dependencies returns the list of dependencies.
func (s *GRPCDemo) Dependencies() []string {
	return nil
}

func (s *GRPCDemo) Status() manager.HealthStatus {
	return s.status
}

func (s *GRPCDemo) LastError() error {
	return s.lastError
}

// Init initializes the service
func (s *GRPCDemo) Init(ctx context.Context) error {
	s.logger.Info("Initializing GRPCDemo service")
	s.status = manager.StatusInitialized
	// Create gRPC server
	// Port 0 in config might mean default or dynamic?
	// If config.Port is 0, we might want to default it or error.
	port := s.config.Port
	if port == 0 {
		port = 50051 // Default
	}

	s.server = gserver.NewServer(s.logger, gserver.WithPort(port), gserver.WithReflection())

	// Register Greeter service
	s.server.RegisterService(func(registrar grpc.ServiceRegistrar) {
		pb.RegisterGreeterServer(registrar, &greeterServer{logger: s.logger})
	})

	return nil
}

// Start starts the service
func (s *GRPCDemo) Start(ctx context.Context) error {
	s.logger.Info("Starting GRPCDemo service")
	s.status = manager.StatusStarting
	if err := s.server.Start(); err != nil {
		s.status = manager.StatusFailed
		s.lastError = err
		return err
	}
	s.status = manager.StatusRunning
	return nil
}

// Stop stops the service
func (s *GRPCDemo) Stop(ctx context.Context) error {
	s.logger.Info("Stopping GRPCDemo service")
	s.status = manager.StatusStopping
	if s.server != nil {
		s.server.Stop(ctx)
	}
	s.status = manager.StatusStopped
	return nil
}
