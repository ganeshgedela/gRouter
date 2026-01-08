package main

import (
	"log"
	"time"

	"grouter/pkg/messaging/grpc"
	"grouter/pkg/messaging/grpc/interceptors"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// Example showing how to create a gRPC server with all middleware.
func ExampleServer() {
	logger, _ := zap.NewProduction()
	tracer := otel.Tracer("example")

	// Create server with all interceptors
	server := grpc.NewServer(logger,
		grpc.WithPort(9090),
		grpc.WithReflection(),

		// Add interceptors in recommended order
		grpc.WithUnaryInterceptor(interceptors.Recovery(logger)),        // 1. Panic recovery
		grpc.WithUnaryInterceptor(interceptors.Logging(logger)),         // 2. Logging
		grpc.WithUnaryInterceptor(interceptors.Metrics()),               // 3. Metrics
		grpc.WithUnaryInterceptor(interceptors.Tracing(tracer)),         // 4. Tracing
		grpc.WithUnaryInterceptor(interceptors.RateLimiter(100, 10)),    // 5. Rate limiting
		grpc.WithUnaryInterceptor(interceptors.Timeout(30*time.Second)), // 6. Timeout

		// Stream interceptors
		grpc.WithStreamInterceptor(interceptors.RecoveryStream(logger)),
		grpc.WithStreamInterceptor(interceptors.LoggingStream(logger)),
		grpc.WithStreamInterceptor(interceptors.MetricsStream()),
		grpc.WithStreamInterceptor(interceptors.TracingStream(tracer)),
	)

	// Register your gRPC services here
	// pb.RegisterYourServiceServer(server.server, yourImplementation)

	// Start server
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// For graceful shutdown
	// server.Stop(context.Background())
}

// Example showing how to create a gRPC client with all middleware.
func ExampleClient() {
	logger, _ := zap.NewProduction()
	tracer := otel.Tracer("example")

	config := grpc.ClientConfig{
		Target:  "localhost:9090",
		Timeout: 10 * time.Second,
	}

	// Create circuit breaker
	cb := interceptors.NewCircuitBreaker(5, 60*time.Second, 1)

	// Create client with all interceptors
	client, err := grpc.NewClient(config, logger,
		// Add client interceptors
		grpc.WithClientUnaryInterceptor(interceptors.LoggingClient(logger)),
		grpc.WithClientUnaryInterceptor(interceptors.MetricsClient()),
		grpc.WithClientUnaryInterceptor(interceptors.TracingClient(tracer)),
		grpc.WithClientUnaryInterceptor(interceptors.CircuitBreakerInterceptor(cb)),
		grpc.WithClientUnaryInterceptor(interceptors.Retry(interceptors.DefaultRetryConfig())),
		grpc.WithClientUnaryInterceptor(interceptors.TimeoutClient(10*time.Second)),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Use client.GetConn() to create service clients
	// pb.NewYourServiceClient(client.GetConn())
}

func main() {
	// Examples are shown above
	log.Println("See example functions for usage patterns")
}
