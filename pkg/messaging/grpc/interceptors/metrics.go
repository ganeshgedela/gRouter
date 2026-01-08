package interceptors

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// Server metrics
	grpcRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_server_requests_total",
			Help: "Total number of gRPC requests received",
		},
		[]string{"method", "code"},
	)

	grpcRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_server_request_duration_seconds",
			Help:    "Duration of gRPC requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "code"},
	)

	// Client metrics
	grpcClientRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_client_requests_total",
			Help: "Total number of gRPC client requests made",
		},
		[]string{"method", "code"},
	)

	grpcClientRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_client_request_duration_seconds",
			Help:    "Duration of gRPC client requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "code"},
	)
)

// init registers metrics with Prometheus default registry.
func init() {
	prometheus.MustRegister(grpcRequestsTotal)
	prometheus.MustRegister(grpcRequestDuration)
	prometheus.MustRegister(grpcClientRequestsTotal)
	prometheus.MustRegister(grpcClientRequestDuration)
}

// Metrics creates a server interceptor that records Prometheus metrics.
func Metrics() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// Call the handler
		resp, err := handler(ctx, req)

		// Calculate duration
		duration := time.Since(start)

		// Extract status code
		code := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				code = st.Code()
			}
		}

		// Record metrics
		grpcRequestsTotal.WithLabelValues(info.FullMethod, code.String()).Inc()
		grpcRequestDuration.WithLabelValues(info.FullMethod, code.String()).Observe(duration.Seconds())

		return resp, err
	}
}

// MetricsStream creates a stream interceptor that records Prometheus metrics.
func MetricsStream() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		// Call the handler
		err := handler(srv, ss)

		// Calculate duration
		duration := time.Since(start)

		// Extract status code
		code := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				code = st.Code()
			}
		}

		// Record metrics
		grpcRequestsTotal.WithLabelValues(info.FullMethod, code.String()).Inc()
		grpcRequestDuration.WithLabelValues(info.FullMethod, code.String()).Observe(duration.Seconds())

		return err
	}
}

// MetricsClient creates a client interceptor that records Prometheus metrics.
func MetricsClient() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()

		// Call the handler
		err := invoker(ctx, method, req, reply, cc, opts...)

		// Calculate duration
		duration := time.Since(start)

		// Extract status code
		code := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				code = st.Code()
			}
		}

		// Record metrics
		grpcClientRequestsTotal.WithLabelValues(method, code.String()).Inc()
		grpcClientRequestDuration.WithLabelValues(method, code.String()).Observe(duration.Seconds())

		return err
	}
}
