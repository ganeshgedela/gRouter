package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8080"

func main() {
	fmt.Println("🚀 REST API Client Example - Testing REST API Service")
	fmt.Println("=" + string(make([]byte, 60)) + "=\n")

	// Wait a moment for server to be ready
	time.Sleep(500 * time.Millisecond)

	// Test 1: Health Check
	fmt.Println("📤 Test 1: Health Check")
	testHealthCheck()
	time.Sleep(500 * time.Millisecond)

	// Test 2: Liveness Probe
	fmt.Println("\n📤 Test 2: Liveness Probe")
	testLiveness()
	time.Sleep(500 * time.Millisecond)

	// Test 3: Readiness Probe
	fmt.Println("\n📤 Test 3: Readiness Probe")
	testReadiness()
	time.Sleep(500 * time.Millisecond)

	// Test 4: Create User
	fmt.Println("\n📤 Test 4: Create User (Example)")
	testCreateUser()
	time.Sleep(500 * time.Millisecond)

	// Test 5: Get Users
	fmt.Println("\n📤 Test 5: Get Users (Example)")
	testGetUsers()
	time.Sleep(500 * time.Millisecond)

	// Test 6: Create Order
	fmt.Println("\n📤 Test 6: Create Order (Example)")
	testCreateOrder()
	time.Sleep(500 * time.Millisecond)

	// Test 7: Get Orders
	fmt.Println("\n📤 Test 7: Get Orders (Example)")
	testGetOrders()

	fmt.Println("\n\n🎉 All API tests completed!")
	fmt.Println("👀 Check the API server logs to see request processing")
}

func testHealthCheck() {
	resp, err := http.Get(baseURL + "/api/v1/health")
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	printResponse("GET /api/v1/health", resp)
}

func testLiveness() {
	resp, err := http.Get(baseURL + "/health/live")
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	printResponse("GET /health/live", resp)
}

func testReadiness() {
	resp, err := http.Get(baseURL + "/health/ready")
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	printResponse("GET /health/ready", resp)
}

func testCreateUser() {
	user := map[string]interface{}{
		"name":  "John Doe",
		"email": "john.doe@example.com",
		"role":  "customer",
	}

	jsonData, _ := json.Marshal(user)
	resp, err := http.Post(baseURL+"/api/v1/users", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	printResponse("POST /api/v1/users", resp)
}

func testGetUsers() {
	resp, err := http.Get(baseURL + "/api/v1/users")
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	printResponse("GET /api/v1/users", resp)
}

func testCreateOrder() {
	order := map[string]interface{}{
		"user_id": "user-123",
		"items": []map[string]interface{}{
			{"product": "Widget A", "quantity": 2, "price": 29.99},
			{"product": "Widget B", "quantity": 1, "price": 49.99},
		},
		"total": 109.97,
	}

	jsonData, _ := json.Marshal(order)
	resp, err := http.Post(baseURL+"/api/v1/orders", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	printResponse("POST /api/v1/orders", resp)
}

func testGetOrders() {
	resp, err := http.Get(baseURL + "/api/v1/orders")
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	printResponse("GET /api/v1/orders", resp)
}

func printResponse(endpoint string, resp *http.Response) {
	body, _ := io.ReadAll(resp.Body)

	statusEmoji := "✅"
	if resp.StatusCode >= 400 {
		statusEmoji = "❌"
	} else if resp.StatusCode >= 300 {
		statusEmoji = "⚠️"
	}

	fmt.Printf("%s %s\n", statusEmoji, endpoint)
	fmt.Printf("   Status: %d %s\n", resp.StatusCode, resp.Status)

	// Try to pretty print JSON
	var jsonBody interface{}
	if err := json.Unmarshal(body, &jsonBody); err == nil {
		prettyJSON, _ := json.MarshalIndent(jsonBody, "   ", "  ")
		fmt.Printf("   Response: %s\n", string(prettyJSON))
	} else {
		fmt.Printf("   Response: %s\n", string(body))
	}
}
