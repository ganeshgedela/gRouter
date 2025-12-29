# Application Lifecycle Flow

This document illustrates the event-driven lifecycle of the application, specifically how `BootstrapService` and `StopService` control the main application loop via NATS messages.

## Sequence Diagram

The following diagram shows the separation between the **Main Application Loop** (blocking and waiting) and the **NATS Event Handlers** (receiving messages asynchronously).

```mermaid
sequenceDiagram
    participant Ext as External (User/System)
    participant NATS as NATS Server
    participant Main as App.Start (Main Loop)
    participant Boot as BootstrapService
    participant Stop as StopService
    participant Svcs as ServiceManager

    Note over Main: App.Start() -> Initializes<br/>Registers Bootstrap & Stop<br/>Waits on startChan

    %% START FLOW
    Ext->>NATS: Publish "natsdemosvc.start"
    NATS->>Boot: Handle(msg) (NATS Goroutine)
    activate Boot
    Boot->>Main: startChan <- struct{}{} (Signal)
    deactivate Boot
    
    activate Main
    Note over Main: Unblocks from startChan
    Main->>Svcs: RegisterServices() (Load config & modules)
    Svcs-->>Main: Services Ready
    Note over Main: Waits on stopChan
    deactivate Main

    %% STOP FLOW
    Ext->>NATS: Publish "natsdemosvc.stop"
    NATS->>Stop: Handle(msg) (NATS Goroutine)
    activate Stop
    Stop->>Main: stopChan <- struct{}{} (Signal)
    deactivate Stop

    activate Main
    Note over Main: Unblocks from stopChan
    Main->>Svcs: UnregisterServices() (Teardown modules)
    Svcs-->>Main: Services Stopped
    Note over Main: Loops back to wait on startChan
    deactivate Main
```

## Key Concepts

1.  **Channel Synchronization**: `startChan` and `stopChan` act as synchronization bridges between the concurrent NATS message handlers and the main application state machine.
2.  **State Preservation**: The application process does not exit when "stopped". It merely unloads the heavy business logic services and goes into a dormant state, listening only for the "start" signal again.
3.  **Concurrency Safety**: By forcing service registration/unregistration to happen on the Main Loop (triggered by channels), we avoid race conditions that could occur if we tried to modify the Service Manager directly from the NATS handler goroutines.
