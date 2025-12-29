# Manager Module Architecture

The `pkg/manager` module is responsible for orchestrating the application lifecycle, managing internal services, and routing NATS messages.

## Components

*   **ServiceManager**: The main entry point. Initializes config, logger, NATS, and Web Server.
*   **ServiceRouter**: Routes incoming NATS messages to the correct registered service.
*   **ServiceStore**: A thread-safe registry of services.
*   **Service Interface**: The contract that all managed services must implement.

## Sequence Diagrams

### 1. Service Manager Initialization

```mermaid
sequenceDiagram
    participant App as Application
    participant Mgr as ServiceManager
    participant Config as Config Loader
    participant NATS as NATS Client
    participant Web as Web Server

    App->>Mgr: NewServiceManager()
    App->>Mgr: Init()
    activate Mgr
    
    Mgr->>Config: Load()
    Config-->>Mgr: Config
    
    Mgr->>Mgr: initLogger()
    
    Mgr->>NATS: NewNATSClient()
    Mgr->>NATS: Connect()
    NATS-->>Mgr: Connected
    
    Mgr->>Web: NewWebServer()
    
    deactivate Mgr
```

### 2. Service Registration

Services are registered with the manager so they can receive NATS messages routed by topic.

```mermaid
sequenceDiagram
    participant App as Application
    participant Mgr as ServiceManager
    participant Router as ServiceRouter
    participant Store as ServiceStore
    participant Svc as Service Instance

    App->>Mgr: RegisterNATSService(Svc)
    Mgr->>Svc: Name()
    Svc-->>Mgr: "my-service"
    
    Mgr->>Router: Register("my-service", Svc)
    Router->>Store: Add("my-service", Svc)
    Note right of Store: Stored in Map
```

### 3. Message Routing Flow

When a NATS message arrives (e.g., subject `app.my-service.do-work`), the manager routes it to the correct service.

```mermaid
sequenceDiagram
    participant NATS as NATS Server
    participant Mgr as ServiceManager
    participant Router as ServiceRouter
    participant Store as ServiceStore
    participant Svc as Service (my-service)

    NATS->>Mgr: onNATSMessage(Subject="app.my-service.op")
    activate Mgr
    
    Mgr->>Router: HandleMessage(Subject, Msg)
    activate Router
    
    Router->>Router: GetServiceFromTopic(Subject)
    Note right of Router: Extracts "my-service"
    
    Router->>Store: Get("my-service")
    Store-->>Router: Service Instance
    
    Router->>Svc: Handle(Context, Subject, Msg)
    activate Svc
    Svc-->>Router: Result / Error
    deactivate Svc
    
    Router-->>Mgr: Result / Error
    deactivate Router
    
    alt Error Occurred
        Mgr->>NATS: ReplyError(ErrMsg)
    end

    deactivate Mgr
```

### 4. WebService Registration

Services that expose HTTP endpoints register with the Web Server.

```mermaid
sequenceDiagram
    participant App as Application
    participant Mgr as ServiceManager
    participant Web as WebServer
    participant Gin as Gin Router
    participant Svc as Web Service

    App->>Mgr: RegisterWebService(Svc)
    Mgr->>Web: RegisterService(Svc)
    Web->>Svc: RegisterRoutes(RouterGroup)
    Svc->>Gin: GET("/path", Handler)
    Svc->>Gin: POST("/path", Handler)
    Note right of Gin: Routes Registered
```

### 5. Web Request Flow

How an incoming HTTP request is processed by the managed web server.

```mermaid
sequenceDiagram
    participant Client
    participant HTTP as HTTP Server
    participant Gin as Gin Engine
    participant MW as Middleware
    participant Handler as Service Handler

    Client->>HTTP: GET /api/v1/resource
    HTTP->>Gin: ServeHTTP()
    
    Gin->>MW: Logger / Metrics / Tracing
    activate MW
    Note right of MW: Pre-processing
    
    MW->>Handler: Handle Request
    activate Handler
    Handler-->>MW: JSON Response
    deactivate Handler
    
    Note right of MW: Post-processing
    MW-->>Gin: Response
    deactivate MW
    
    Gin-->>HTTP: Response
    HTTP-->>Client: 200 OK
```

