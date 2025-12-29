package nats

import (
	"go.opentelemetry.io/otel"
)

const (
	// instrumentationName is the name of the instrumentation library
	instrumentationName = "grouter/pkg/messaging/nats"
	// systemName is the messaging system name
	systemName = "nats"
	// spanNamePublish is the span name for publish operations
	spanNamePublish = "nats.publish"
	// spanNameProcess is the span name for process operations
	spanNameProcess = "nats.process"
)

// tracer is the global tracer for this package
var tracer = otel.Tracer(instrumentationName)
