package main

import (
	"context"
	"fmt"
	"time"

	"grouter/pkg/config"
	"grouter/pkg/logger"
	"grouter/pkg/messaging/nats"

	"go.uber.org/zap"
)

func main() {
	// 1. Initialize Logger
	log, err := logger.New(logger.Config{
		Level:      "info",
		Format:     "console",
		OutputPath: "stdout",
	})
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	log.Info("🚀 NATS Publisher Example - Testing NATS Worker Service")

	// 2. Load Configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load configuration", zap.Error(err))
	}

	// 3. Create NATS Messenger
	messenger := nats.NewMessenger(nil, nil, nil)
	if err := messenger.InitFromConfig(cfg, log); err != nil {
		log.Fatal("failed to initialize messenger", zap.Error(err))
	}

	if err := messenger.Start(); err != nil {
		log.Fatal("failed to start messenger", zap.Error(err))
	}
	defer messenger.Close()

	log.Info("✅ Connected to NATS server", zap.String("url", cfg.NATS.URL))
	ctx := context.Background()

	// 4. Test User Service - Publish to user.created
	log.Info("\n📤 Test 1: Publishing to user.created topic")
	userPayload := map[string]interface{}{
		"user_id": "user-123",
		"name":    "John Doe",
		"email":   "john.doe@example.com",
		"role":    "customer",
	}

	if err := messenger.Publish(ctx, "user.created", "user.created", userPayload, &nats.PublishOptions{}); err != nil {
		log.Error("failed to publish user.created", zap.Error(err))
	} else {
		log.Info("✅ Published to user.created", zap.Any("payload", userPayload))
	}

	time.Sleep(500 * time.Millisecond)

	// 5. Test Order Service - Publish to order.created
	log.Info("\n📤 Test 2: Publishing to order.created topic")
	orderPayload := map[string]interface{}{
		"order_id": "order-456",
		"user_id":  "user-123",
		"items": []map[string]interface{}{
			{"product": "Widget A", "quantity": 2, "price": 29.99},
			{"product": "Widget B", "quantity": 1, "price": 49.99},
		},
		"total": 109.97,
	}

	if err := messenger.Publish(ctx, "order.created", "order.created", orderPayload, &nats.PublishOptions{}); err != nil {
		log.Error("failed to publish order.created", zap.Error(err))
	} else {
		log.Info("✅ Published to order.created", zap.Any("payload", orderPayload))
	}

	time.Sleep(500 * time.Millisecond)

	// 6. Test Order Service - Publish to order.updated
	log.Info("\n📤 Test 3: Publishing to order.updated topic")
	updatePayload := map[string]interface{}{
		"order_id": "order-456",
		"status":   "shipped",
		"tracking": "TRACK123456",
	}

	if err := messenger.Publish(ctx, "order.updated", "order.updated", updatePayload, &nats.PublishOptions{}); err != nil {
		log.Error("failed to publish order.updated", zap.Error(err))
	} else {
		log.Info("✅ Published to order.updated", zap.Any("payload", updatePayload))
	}

	// 7. Wait for processing
	log.Info("\n⏳ Waiting for worker to process messages...")
	time.Sleep(1 * time.Second)

	fmt.Println("\n🎉 All test messages published successfully!")
	fmt.Println("👀 Check the worker service logs to see the messages being processed")
}
