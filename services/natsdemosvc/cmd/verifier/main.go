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
	flag.BoolVar(&runAsTest, "test", false, "Run verification steps and exit non-zero on failure")
	flag.Parse()

	nc, err := nats.Connect("nats://localhost:4222")
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
	log.Println("2. Sending Health Expecting Reply...")
	healthSubject := "gRouter.health" // Router routes 'health' -> HealthService
	replySubject := "verifier.reply.health"

	sub, err := nc.SubscribeSync(replySubject)
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	env := MessageEnvelope{
		ID:        "verify-health",
		Type:      "liveness",
		Reply:     replySubject,
		Timestamp: time.Now(),
	}
	data, _ := json.Marshal(env)
	if err := nc.Publish(healthSubject, data); err != nil {
		log.Fatalf("Failed to publish health: %v", err)
	}

	msg, err := sub.NextMsg(5 * time.Second)
	if err != nil {
		log.Printf("Health Check Failed: %v", err)
		if runAsTest {
			log.Fatal("Verification Failed: Health check timeout")
		}
	} else {
		log.Printf("Health Response Received: %s", string(msg.Data))
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
