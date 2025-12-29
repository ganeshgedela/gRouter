# Service Manager Package (`pkg/manager`)

The `manager` package is the core orchestration layer of the gRouter application. It is responsible for managing the lifecycle of internal services, handling external NATS messages, and routing them to the appropriate service handlers.

## Design Overview

The Service Manager follows a centralized orchestration pattern where a single `ServiceManager` instance controls the application's message loop and delegates work to registered services.

### Core Components

1.  **ServiceManager (`manager.go`)**:
    -   Initializes the application (Config, Logger, NATS).
    -   Manages the NATS subscription to the application's root topic (`<app_name>.>`).
    -   Handles the graceful shutdown of all components.

2.  **ServiceRouter (`router.go`)**:
    -   Inspects incoming NATS subjects.
    -   Extracts the target service name based on a convention (e.g., `grouter.<service>.<action>`).
    -   Dispatches the message to the registered `Service` implementation.

3.  **ServiceStore (`store.go`)**:
    -   A thread-safe registry for holding service instances.
    -   Allows dynamic registration and lookup of services by name.

4.  **Service Interface (`types.go`)**:
    -   Defines the contract that all internal services (e.g., IPSec, Firewall) must implement.

## Implementation Sketch

### Message Flow

1.  **Ingress**: A message arrives on NATS topic `grouter.ipsec.create`.
2.  **Listener**: The `ServiceManager`'s wildcard subscription (`grouter.>`) picks up the message.
3.  **Routing**: `ServiceManager` passes the message to `ServiceRouter`.
4.  **Resolution**: `ServiceRouter` parses the topic, identifies `ipsec` as the target service, and looks it up in `ServiceStore`.
5.  **Execution**: The `ipsec` service's `Handle` method is invoked.
6.  **Response**: If the `Handle` method returns a response and the message has a `Reply` subject, `ServiceManager` publishes the response.

### Service Interface

To integrate a new service, implement the `Service` interface:

```go
type Service interface {
    // Name returns the unique identifier for the service (e.g., "ipsec")
    Name() string

    // Lifecycle methods
    Ready(ctx context.Context) error
    Start(ctx context.Context) error
    Stop(ctx context.Context) error

    // Message handling
    Handle(ctx context.Context, topic string, msg *messaging.MessageEnvelope) (*messaging.MessageEnvelope, error)
}
```

## Implementation Details

-   **Topic Format**: `<manager_name>.<service_name>.<operation>`
-   **Concurrency**: Each message is processed in its own goroutine (managed by NATS client), but the `ServiceManager` imposes a timeout context for every handler execution.
-   **Error Handling**: Errors returned by services are automatically wrapped and sent back to the caller if a `Reply` subject is present.
