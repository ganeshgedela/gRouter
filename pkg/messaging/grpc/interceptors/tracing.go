package interceptors

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	// TracerName is the name of the tracer used for gRPC instrumentation.
	TracerName = "grouter/grpc"
)

// Tracing creates a server interceptor that adds OpenTelemetry tracing.
func Tracing(tracer trace.Tracer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Start a new span
		ctx, span := tracer.Start(ctx, info.FullMethod,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("rpc.system", "grpc"),
				attribute.String("rpc.service", extractServiceName(info.FullMethod)),
				attribute.String("rpc.method", extractMethodName(info.FullMethod)),
			),
		)
		defer span.End()

		// Call the handler
		resp, err := handler(ctx, req)

		// Set span status based on error
		if err != nil {
			span.RecordError(err)
			if st, ok := status.FromError(err); ok {
				span.SetStatus(codes.Error, st.Message())
				span.SetAttributes(attribute.Int("rpc.grpc.status_code", int(st.Code())))
			} else {
				span.SetStatus(codes.Error, err.Error())
			}
		} else {
			span.SetStatus(codes.Ok, "")
			span.SetAttributes(attribute.Int("rpc.grpc.status_code", int(grpcCodes.OK)))
		}

		return resp, err
	}
}

// TracingStream creates a stream interceptor that adds OpenTelemetry tracing.
func TracingStream(tracer trace.Tracer) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Start a new span
		ctx, span := tracer.Start(ss.Context(), info.FullMethod,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("rpc.system", "grpc"),
				attribute.String("rpc.service", extractServiceName(info.FullMethod)),
				attribute.String("rpc.method", extractMethodName(info.FullMethod)),
				attribute.Bool("rpc.grpc.is_client_stream", info.IsClientStream),
				attribute.Bool("rpc.grpc.is_server_stream", info.IsServerStream),
			),
		)
		defer span.End()

		// Wrap the stream with the new context
		wrappedStream := &tracedServerStream{
			ServerStream: ss,
			ctx:          ctx,
		}

		// Call the handler
		err := handler(srv, wrappedStream)

		// Set span status based on error
		if err != nil {
			span.RecordError(err)
			if st, ok := status.FromError(err); ok {
				span.SetStatus(codes.Error, st.Message())
				span.SetAttributes(attribute.Int("rpc.grpc.status_code", int(st.Code())))
			} else {
				span.SetStatus(codes.Error, err.Error())
			}
		} else {
			span.SetStatus(codes.Ok, "")
			span.SetAttributes(attribute.Int("rpc.grpc.status_code", int(grpcCodes.OK)))
		}

		return err
	}
}

// TracingClient creates a client interceptor that adds OpenTelemetry tracing.
func TracingClient(tracer trace.Tracer) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// Start a new span
		ctx, span := tracer.Start(ctx, method,
			trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(
				attribute.String("rpc.system", "grpc"),
				attribute.String("rpc.service", extractServiceName(method)),
				attribute.String("rpc.method", extractMethodName(method)),
			),
		)
		defer span.End()

		// Inject trace context into metadata
		md, _ := metadata.FromOutgoingContext(ctx)
		md = md.Copy()
		otel.GetTextMapPropagator().Inject(ctx, &metadataCarrier{md: &md})
		ctx = metadata.NewOutgoingContext(ctx, md)

		// Call the handler
		err := invoker(ctx, method, req, reply, cc, opts...)

		// Set span status based on error
		if err != nil {
			span.RecordError(err)
			if st, ok := status.FromError(err); ok {
				span.SetStatus(codes.Error, st.Message())
				span.SetAttributes(attribute.Int("rpc.grpc.status_code", int(st.Code())))
			} else {
				span.SetStatus(codes.Error, err.Error())
			}
		} else {
			span.SetStatus(codes.Ok, "")
			span.SetAttributes(attribute.Int("rpc.grpc.status_code", int(grpcCodes.OK)))
		}

		return err
	}
}

// tracedServerStream wraps a grpc.ServerStream with a context containing the trace span.
type tracedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *tracedServerStream) Context() context.Context {
	return s.ctx
}

// metadataCarrier implements the TextMapCarrier interface for gRPC metadata.
type metadataCarrier struct {
	md *metadata.MD
}

func (c *metadataCarrier) Get(key string) string {
	values := c.md.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (c *metadataCarrier) Set(key, value string) {
	c.md.Set(key, value)
}

func (c *metadataCarrier) Keys() []string {
	keys := make([]string, 0, len(*c.md))
	for k := range *c.md {
		keys = append(keys, k)
	}
	return keys
}

// extractServiceName extracts the service name from a full method name.
// e.g., "/helloworld.Greeter/SayHello" -> "helloworld.Greeter"
func extractServiceName(fullMethod string) string {
	if len(fullMethod) == 0 {
		return ""
	}
	// Remove leading slash
	if fullMethod[0] == '/' {
		fullMethod = fullMethod[1:]
	}
	// Find the last slash
	for i := len(fullMethod) - 1; i >= 0; i-- {
		if fullMethod[i] == '/' {
			return fullMethod[:i]
		}
	}
	return fullMethod
}

// extractMethodName extracts the method name from a full method name.
// e.g., "/helloworld.Greeter/SayHello" -> "SayHello"
func extractMethodName(fullMethod string) string {
	if len(fullMethod) == 0 {
		return ""
	}
	// Find the last slash
	for i := len(fullMethod) - 1; i >= 0; i-- {
		if fullMethod[i] == '/' {
			return fullMethod[i+1:]
		}
	}
	return fullMethod
}
