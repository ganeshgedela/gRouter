package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	pb "grouter/templates/grpc-service/api/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const serverAddr = "localhost:9090"

func main() {
	fmt.Println("🚀 gRPC Client Example - Testing gRPC Service")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	// Set up connection to the server
	fmt.Printf("📡 Connecting to gRPC server at %s...\n", serverAddr)
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ Failed to connect: %v", err)
	}
	defer conn.Close()

	fmt.Println("✅ Connected to gRPC server")
	fmt.Println()

	// Create HelloService client
	client := pb.NewHelloServiceClient(conn)

	// Test 1: Simple greeting
	fmt.Println("📤 Test 1: Sending SayHello request...")
	testSayHello(client, "Alice")
	time.Sleep(500 * time.Millisecond)

	// Test 2: Another name
	fmt.Println("\n📤 Test 2: Sending SayHello request...")
	testSayHello(client, "Bob")
	time.Sleep(500 * time.Millisecond)

	// Test 3: With timeout
	fmt.Println("\n📤 Test 3: Testing with context timeout...")
	testSayHelloWithTimeout(client, "Charlie", 2*time.Second)

	fmt.Println("\n🎉 All gRPC tests completed successfully!")
	fmt.Println("👀 Check the server logs to see request processing")
}

func testSayHello(client pb.HelloServiceClient, name string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.HelloRequest{
		Name: name,
	}

	resp, err := client.SayHello(ctx, req)
	if err != nil {
		fmt.Printf("❌ Error calling SayHello: %v\n", err)
		return
	}

	fmt.Printf("✅ Request: {name: \"%s\"}\n", name)
	fmt.Printf("✅ Response: {message: \"%s\"}\n", resp.Message)
}

func testSayHelloWithTimeout(client pb.HelloServiceClient, name string, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req := &pb.HelloRequest{
		Name: name,
	}

	startTime := time.Now()
	resp, err := client.SayHello(ctx, req)
	duration := time.Since(startTime)

	if err != nil {
		fmt.Printf("❌ Error calling SayHello: %v\n", err)
		return
	}

	fmt.Printf("✅ Request: {name: \"%s\"} (timeout: %v)\n", name, timeout)
	fmt.Printf("✅ Response: {message: \"%s\"}\n", resp.Message)
	fmt.Printf("⏱️  Duration: %v\n", duration)
}
