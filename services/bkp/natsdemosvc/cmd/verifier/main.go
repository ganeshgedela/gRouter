package main

import (
	"encoding/json"
	"flag"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type MessageEnvelope struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Source    string            `json:"source"`
	Timestamp time.Time         `json:"timestamp"`
	Reply     string            `json:"reply,omitempty"`
	Data      json.RawMessage   `json:"data,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

func main() {
	var runAsTest bool
	var natsURL string
	flag.BoolVar(&runAsTest, "test", false, "Run verification steps and exit non-zero on failure")
	flag.StringVar(&natsURL, "url", "nats://localhost:4222", "NATS server URL")
	flag.Parse()

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// 1. Send Start
	log.Println("1. Sending Start Signal...")
	if err := publishMessage(nc, "natsdemosvc.start", "start"); err != nil {
		log.Printf("Failed: %v", err)
	}
	time.Sleep(2 * time.Second)

	// 2. Health Check (Proof of Life)
	log.Println("2. Sending Health Live Expecting Reply...")
	healthSubject := "natsdemosvc.health.live"
	replySubject := "verifier.reply.health.live"

	sub, err := nc.SubscribeSync(replySubject)
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	env := MessageEnvelope{
		ID:        "verify-health-live",
		Type:      "health.live",
		Reply:     replySubject,
		Timestamp: time.Now(),
	}
	data, _ := json.Marshal(env)
	if err := nc.Publish(healthSubject, data); err != nil {
		log.Fatalf("Failed to publish health: %v", err)
	}

	msg, err := sub.NextMsg(5 * time.Second)
	if err != nil {
		log.Printf("Health Live Check Failed: %v", err)
		if runAsTest {
			log.Fatal("Verification Failed: Health check timeout")
		}
	} else {
		log.Printf("Health Live Response Received: %s", string(msg.Data))
	}

	// 2.1 Health Ready Check
	log.Println("2.1 Sending Health Ready Expecting Reply...")
	healthReadySubject := "natsdemosvc.health.ready"
	replyReadySubject := "verifier.reply.health.ready"

	subReady, err := nc.SubscribeSync(replyReadySubject)
	if err != nil {
		log.Fatalf("Failed to subscribe ready: %v", err)
	}

	envReady := MessageEnvelope{
		ID:        "verify-health-ready",
		Type:      "health.ready",
		Reply:     replyReadySubject,
		Timestamp: time.Now(),
	}
	dataReady, _ := json.Marshal(envReady)
	if err := nc.Publish(healthReadySubject, dataReady); err != nil {
		log.Fatalf("Failed to publish health ready: %v", err)
	}

	msgReady, err := subReady.NextMsg(5 * time.Second)
	if err != nil {
		log.Printf("Health Ready Check Failed: %v", err)
	} else {
		log.Printf("Health Ready Response Received: %s", string(msgReady.Data))
	}

	// 3. Send Stop
	log.Println("3. Sending Stop Signal...")
	if err := publishMessage(nc, "natsdemosvc.stop", "stop"); err != nil {
		log.Printf("Failed: %v", err)
	}

	log.Println("Verification Complete")
}

func publishMessage(nc *nats.Conn, subject, msgType string) error {
	env := MessageEnvelope{
		ID:        "verify-" + msgType,
		Type:      msgType,
		Source:    "verifier",
		Timestamp: time.Now(),
		Data:      json.RawMessage("{}"),
	}
	data, _ := json.Marshal(env)
	return nc.Publish(subject, data)
}
