package telemetry

import (
	"context"

	"grouter/pkg/config"
)

// Init initializes both Tracing and Metrics systems
func Init(cfg config.Config) (func(context.Context) error, error) {
	// 1. Init Metrics
	InitMetrics(cfg.Metrics)

	// 2. Init Tracer
	shutdown, err := InitTracer(cfg.Tracing)
	if err != nil {
		return nil, err
	}

	return shutdown, nil
}
