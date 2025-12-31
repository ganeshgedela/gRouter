# gRouter Architecture Documentation

This document provides a high-level overview of the gRouter architectural design, component interactions, and message flows.

## 1. System Architecture Overview

gRouter is a generic microservices framework designed for high-performance messaging, web services, and observability. It follows a modular design where a central `ServiceManager` orchestrates various subsystems.

```mermaid
graph TD
    subgraph "Infrastructure Layer"
        NATS["NATS Cluster"]
        Keycloak["Keycloak (IAM)"]
        Jaeger["Jaeger (Tracing)"]
        Prometheus["Prometheus (Metrics)"]
        OpenSearch["OpenSearch (Logs)"]
    end

    subgraph "gRouter Core (pkg/)"
        Manager["ServiceManager (pkg/manager)"]
        WebSrv["Web Server (pkg/web)"]
        Messenger["Messenger (pkg/messaging/nats)"]
        Telemetry["Telemetry (pkg/telemetry)"]
        Config["Config (pkg/config)"]
    end

    subgraph "Service Layer"
        WebDemo["WebDemoSvc"]
        NATSDemo["NATSDemoSvc"]
    end

    WebDemo --> Manager
    NATSDemo --> Manager
    
    Manager --> WebSrv
    Manager --> Messenger
    Manager --> Telemetry
    
    Messenger --> NATS
    WebSrv --> Keycloak
    Telemetry --> Jaeger
    Telemetry --> Prometheus
```

## 2. Interface Sequence Diagrams

### 2.1 Authentication & Request Flow (OIDC)

The following diagram illustrates how an authenticated HTTP request is processed by `webdemosvc` using Keycloak.

```mermaid
sequenceDiagram
    participant User
    participant WebSvc as WebDemoSvc
    participant Keycloak as Keycloak (OIDC)
    participant Service as Internal Logic

    User->>Keycloak: Get Access Token (testuser/password)
    Keycloak-->>User: Bearer JWT
    User->>WebSvc: HTTP GET /start (Authorization: Bearer JWT)
    WebSvc->>WebSvc: AuthMiddleware: Decode & Validate
    WebSvc->>Keycloak: Fetch JWKS (Cachable)
    Keycloak-->>WebSvc: Public Keys
    WebSvc->>WebSvc: Verify Sign/Issuer/Audience
    alt Authorized
        WebSvc->>Service: Execute Start Logic
        Service-->>WebSvc: Success
        WebSvc-->>User: 200 OK {"status": "starting"}
    else Unauthorized
        WebSvc-->>User: 401 Unauthorized
    end
```

### 2.2 NATS Messaging Flow (Request-Reply)

The `natsdemosvc` uses a request-reply pattern for service activation and communication.

```mermaid
sequenceDiagram
    participant Requester as NATS Client (nats-box)
    participant NATS as NATS Server
    participant App as NATSDemoSvc
    participant Manager as ServiceManager
    participant Module as NATSDemoModule

    Requester->>NATS: PUB natsdemosvc.start (payload)
    NATS->>App: Deliver Message
    App->>Manager: onNATSMessage (routing)
    Manager->>Module: Handle (start)
    Module-->>Manager: Done
    Manager->>Manager: Register Internal Services
    Manager-->>NATS: Reply (optional)
```

### 2.3 Tracing & Observability Flow (OTLP)

Trace spans are generated within the application and exported via OTLP to Jaeger.

```mermaid
sequenceDiagram
    participant App as Application Service
    participant SDK as OTEL SDK / Tracer
    participant Exporter as OTLP HTTP Exporter
    participant Jaeger as Jaeger (All-in-One)
    participant UI as Jaeger UI

    App->>SDK: Start Span (operation_name)
    Note over App,SDK: Middleware/Interceptors
    App->>App: Business Logic
    App->>SDK: End Span
    SDK->>SDK: Batch Spans
    SDK->>Exporter: Export(batch)
    Exporter->>Jaeger: POST /v1/traces (Protobuf/JSON)
    Jaeger->>Jaeger: Process & Store Spans
    
    UI->>Jaeger: Query Services/Operations
    Jaeger-->>UI: Return Trace Data
```

### 2.4 Logging & Log Aggregation Flow (Fluentd)

Logs are captured from application output, processed by Fluentd, and indexed in OpenSearch.

```mermaid
sequenceDiagram
    participant App as Application Pod
    participant Node as Node (Log Files)
    participant Fluentd as Fluentd (DaemonSet)
    participant OS as OpenSearch (Backend)
    participant Dash as OpenSearch Dashboards

    App->>App: logger.Info("Message")
    App->>Node: Write to stdout (Container Log)
    Fluentd->>Node: Tail Log File (/var/log/pods/...)
    Fluentd->>Fluentd: Parse JSON & Enrich Metadata
    Fluentd->>OS: POST /fluentd/_doc (Bulk API)
    OS->>OS: Index Log Entry
    
    Dash->>OS: Search/Discover Indices
    OS-->>Dash: Return Log Hits
```

### 2.5 Metrics & Monitoring Flow (Prometheus)

Metrics are exposed by the application and periodically scraped by Prometheus for storage and visualization.

```mermaid
sequenceDiagram
    participant App as Application Service
    participant Registry as Prometheus Registry
    participant Prom as Prometheus (Server)
    participant Grafana as Grafana Dashboard

    App->>Registry: Observe(metric_name, value)
    Registry->>Registry: Update Counters/Histograms
    
    loop Every 15s (Scrape Interval)
        Prom->>App: GET /metrics
        App-->>Prom: Text Format (Prometheus Exposition)
    end
    
    Prom->>Prom: Store Time Series Data
    
    Grafana->>Prom: Query (PromQL)
    Prom-->>Grafana: Return Matrix/Vector Data
    Grafana->>Grafana: Render Visual Charts
```

## 3. Component Message Flows

### 3.1 Observability Pipeline

Logging and Tracing flows ensure full visibility across the distributed system.

```mermaid
flowchart LR
    App["Application Pod"] -- OTLP/HTTP --> Jaeger["Jaeger Exporter"]
    App -- stdout --> Fluentd["Fluentd DaemonSet"]
    Fluentd -- opensearch --> OS["OpenSearch"]
    App -- HTTP /metrics --> Prom["Prometheus"]
    Prom -- scrape --> Grafana["Grafana Dashboard"]
```

### 3.2 Service Lifecycle

All gRouter services follow a stateful transition pattern:

1.  **Bootstrap State**:
    - Only `start`, `stop`, and `health` endpoints/topics are active.
    - The application waits for a "Start Signal".
2.  **Running State**:
    - Upon receiving a signal, the `ServiceManager` registers all dynamic modules.
    - The `WebServer` or `Messenger` reloads to apply new routes/subscriptions.
3.  **Graceful Stop**:
    - Services are unregistered.
    - Internal state is cleaned up.
    - Application returns to Bootstrap or exits.

## 4. Component Summary

| Component | Responsibility |
| :--- | :--- |
| **ServiceManager** | Core orchestrator. Handles Init, Start, Stop, and Service Registration. |
| **Messenger** | Wrapper around NATS. Provides Middleware support (Logging, Tracing) and Request-Reply abstractions. |
| **Web Server** | Gin-based HTTP server. Implements OIDC Auth, CORS, Security Headers, and Prometheus metrics. |
| **Telemetry** | Centralized OpenTelemetry initialization for distributed tracing (Jaeger) and metrics. |
| **Config** | Viper-based configuration loader with validation and environment variable overrides. |
