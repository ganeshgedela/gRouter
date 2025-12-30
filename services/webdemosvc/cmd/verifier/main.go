package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	urlPtr := flag.String("url", "http://localhost:8080", "URL of the webdemosvc")
	flag.Parse()

	baseURL := *urlPtr

	// Initialize simple logger
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()
	defer logger.Sync()

	logger.Info("Starting WebDemoSvc E2E Verifier", zap.String("url", baseURL))

	// 1. Check Health/Liveness
	if err := checkHealth(baseURL, logger); err != nil {
		logger.Fatal("Health check failed", zap.Error(err))
	}

	// 2. Start Service
	if err := triggerStart(baseURL, logger); err != nil {
		// It might be already started, which is fine, but let's log it
		logger.Warn("Start trigger returned error (maybe already started)", zap.Error(err))
	}

	// 3. Verify Hello Endpoint
	if err := checkHello(baseURL, logger); err != nil {
		logger.Fatal("Hello check failed", zap.Error(err))
	}

	// 4. Verify Echo Endpoint
	if err := checkEcho(baseURL, logger); err != nil {
		logger.Fatal("Echo check failed", zap.Error(err))
	}

	// 5. Stop Service
	if err := triggerStop(baseURL, logger); err != nil {
		logger.Fatal("Stop trigger failed", zap.Error(err))
	}

	logger.Info("Verification Successful!")
}

func checkHealth(baseURL string, logger *zap.Logger) error {
	logger.Info("Checking liveness...")
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(baseURL + "/health/live")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				logger.Info("Service is live")
				return nil
			}
		}
		time.Sleep(1 * time.Second)
		if i%5 == 0 {
			logger.Info("Waiting for service...", zap.Int("attempt", i+1))
		}
	}
	return fmt.Errorf("service not live after %d attempts", maxRetries)
}

func triggerStart(baseURL string, logger *zap.Logger) error {
	logger.Info("Triggering Start...")
	resp, err := http.Get(baseURL + "/start")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(body))
	}

	logger.Info("Start triggered successfully")
	// Give it a moment to register services
	time.Sleep(5 * time.Second)
	return nil
}

func checkHello(baseURL string, logger *zap.Logger) error {
	logger.Info("Checking /hello...")
	resp, err := http.Get(baseURL + "/hello")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	logger.Info("Hello endpoint verified")
	return nil
}

func checkEcho(baseURL string, logger *zap.Logger) error {
	logger.Info("Checking /echo...")
	resp, err := http.Get(baseURL + "/echo?msg=verifier")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	// Could check body content too, but status 200 is good enough for basic connectivity
	logger.Info("Echo endpoint verified")
	return nil
}

func triggerStop(baseURL string, logger *zap.Logger) error {
	logger.Info("Triggering Stop...")
	resp, err := http.Get(baseURL + "/stop")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	logger.Info("Stop triggered successfully")
	return nil
}
