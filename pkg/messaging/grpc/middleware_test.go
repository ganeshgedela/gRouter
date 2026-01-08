package grpc

import (
	"context"
	"testing"

	"google.golang.org/grpc"
)

// TestChainUnaryServer tests the unary server interceptor chaining.
func TestChainUnaryServer(t *testing.T) {
	var callOrder []int

	interceptor1 := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		callOrder = append(callOrder, 1)
		resp, err := handler(ctx, req)
		callOrder = append(callOrder, 4)
		return resp, err
	}

	interceptor2 := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		callOrder = append(callOrder, 2)
		resp, err := handler(ctx, req)
		callOrder = append(callOrder, 3)
		return resp, err
	}

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}

	chained := ChainUnaryServer(interceptor1, interceptor2)

	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}
	_, err := chained(context.Background(), nil, info, handler)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := []int{1, 2, 3, 4}
	if len(callOrder) != len(expected) {
		t.Errorf("expected call order %v, got %v", expected, callOrder)
	}

	for i, v := range expected {
		if callOrder[i] != v {
			t.Errorf("at index %d: expected %d, got %d", i, v, callOrder[i])
		}
	}
}

// TestChainStreamServer tests the stream server interceptor chaining.
func TestChainStreamServer(t *testing.T) {
	var callOrder []int

	interceptor1 := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		callOrder = append(callOrder, 1)
		err := handler(srv, ss)
		callOrder = append(callOrder, 4)
		return err
	}

	interceptor2 := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		callOrder = append(callOrder, 2)
		err := handler(srv, ss)
		callOrder = append(callOrder, 3)
		return err
	}

	handler := func(srv interface{}, stream grpc.ServerStream) error {
		return nil
	}

	chained := ChainStreamServer(interceptor1, interceptor2)

	info := &grpc.StreamServerInfo{FullMethod: "/test.Service/StreamMethod"}
	err := chained(nil, nil, info, handler)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := []int{1, 2, 3, 4}
	if len(callOrder) != len(expected) {
		t.Errorf("expected call order %v, got %v", expected, callOrder)
	}
}

// TestChainUnaryClient tests the unary client interceptor chaining.
func TestChainUnaryClient(t *testing.T) {
	var callOrder []int

	interceptor1 := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		callOrder = append(callOrder, 1)
		err := invoker(ctx, method, req, reply, cc, opts...)
		callOrder = append(callOrder, 4)
		return err
	}

	interceptor2 := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		callOrder = append(callOrder, 2)
		err := invoker(ctx, method, req, reply, cc, opts...)
		callOrder = append(callOrder, 3)
		return err
	}

	invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return nil
	}

	chained := ChainUnaryClient(interceptor1, interceptor2)

	err := chained(context.Background(), "/test", nil, nil, nil, invoker)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := []int{1, 2, 3, 4}
	if len(callOrder) != len(expected) {
		t.Errorf("expected call order %v, got %v", expected, callOrder)
	}
}

// TestChainStreamClient tests the stream client interceptor chaining.
func TestChainStreamClient(t *testing.T) {
	var callOrder []int

	interceptor1 := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		callOrder = append(callOrder, 1)
		stream, err := streamer(ctx, desc, cc, method, opts...)
		callOrder = append(callOrder, 4)
		return stream, err
	}

	interceptor2 := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		callOrder = append(callOrder, 2)
		stream, err := streamer(ctx, desc, cc, method, opts...)
		callOrder = append(callOrder, 3)
		return stream, err
	}

	streamer := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return nil, nil
	}

	chained := ChainStreamClient(interceptor1, interceptor2)

	_, err := chained(context.Background(), nil, nil, "/test", streamer)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := []int{1, 2, 3, 4}
	if len(callOrder) != len(expected) {
		t.Errorf("expected call order %v, got %v", expected, callOrder)
	}
}

// TestEmptyChains tests chaining with no interceptors.
func TestEmptyChains(t *testing.T) {
	t.Run("empty unary server chain", func(t *testing.T) {
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return "response", nil
		}

		chained := ChainUnaryServer()
		info := &grpc.UnaryServerInfo{}

		resp, err := chained(context.Background(), nil, info, handler)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp != "response" {
			t.Errorf("expected 'response', got %v", resp)
		}
	})

	t.Run("empty unary client chain", func(t *testing.T) {
		invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			return nil
		}

		chained := ChainUnaryClient()
		err := chained(context.Background(), "/test", nil, nil, nil, invoker)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
