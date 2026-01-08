package telemetry

import (
	"context"

	"grouter/pkg/config"
)

// Init initializes both Tracing and Metrics systems
func Init(tracingCfg config.TracingConfig) (func(context.Context) error, error) {
	// 1. Init Metrics
	InitMetrics()

	// 2. Init Tracer
	shutdown, err := InitTracer(tracingCfg)
	if err != nil {
		return nil, err
	}

	return shutdown, nil
}
