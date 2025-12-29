# Database Package Learning Guide

This document explains the architecture, design, and usage of the generic `pkg/database` package in `gRouter`.

## Overview
The `pkg/database` package provides a standardized, production-ready interface for interacting with databases. It abstracts the underlying driver (Postgres, SQLite) and ORM (GORM) complexities, offering a clean API for microservices.

## Architecture

The package is built on top of **GORM** and follows the **Repository Pattern**.

```mermaid
classDiagram
    class Database {
        +*gorm.DB DB
        +New(config, logger)
        +WithTransaction(ctx, fn)
        +HealthCheck(ctx)
    }

    class Repository~T~ {
        <<interface>>
        +Create(ctx, entity)
        +FindByID(ctx, id)
        +List(ctx, pagination)
        +Update(ctx, entity)
        +Delete(ctx, id)
    }

    class GORMRepository~T~ {
        -db *gorm.DB
    }

    class MetricsCollector {
        +Start(interval)
    }

    Database *-- MetricsCollector : monitors
    Repository <|.. GORMRepository : implements
    Database ..> Repository : creates
```

## Features & Sequence Diagrams

### 1. Initialization
The `New` factory function initializes the database connection, configures the connection pool, sets up the Zap logger adapter, and registers OpenTelemetry tracing.

```mermaid
sequenceDiagram
    participant App
    participant Factory as pkg/database
    participant GORM
    participant Driver as Postgres/SQLite
    participant Otel as OpenTelemetry

    App->>Factory: New(Config, Logger)
    Factory->>Driver: Open(DSN)
    Factory->>GORM: Open(Dialect, Config)
    GORM-->>Factory: *gorm.DB
    Factory->>Otel: Register Plugin
    Factory->>Factory: Configure Connection Pool (MaxOpen, MaxIdle)
    Factory-->>App: *Database
```

### 2. Generic CRUD Operations
The `Repository[T]` interface allows standard CRUD operations on any entity without writing boilerplate code. Supports **Pagination** and **Dynamic Filtering**.

```mermaid
sequenceDiagram
    participant Service
    participant Repo as Repository[User]
    participant DB as Database

    Service->>Repo: List(ctx, Pagination{Page:1, Filters:{"status":"active"}})
    Repo->>DB: Count(Where status='active')
    DB-->>Repo: total_count
    Repo->>DB: Find(Offset, Limit, Order, Where)
    DB-->>Repo: []User
    Repo-->>Service: []User, total
```

### 3. Transaction Management
Atomic operations are supported via `WithTransaction`.

```mermaid
sequenceDiagram
    participant Service
    participant DB as DatabaseWrapper
    participant Tx as GORM Transaction
    participant Repo as Repository[User]

    Service->>DB: WithTransaction(fn)
    DB->>Tx: Begin Transaction
    DB->>Service: invoke fn(txDB)
    
    Service->>Repo: NewRepository(txDB)
    Service->>Repo: Create(User)
    Repo->>Tx: INSERT ...
    
    alt Success
        Service-->>DB: return nil
        DB->>Tx: Commit
    else Failure
        Service-->>DB: return error
        DB->>Tx: Rollback
    end
```

### 4. Observability: Tracing Flow
The package uses `gorm.io/plugin/opentelemetry` to automatically trace queries.

```mermaid
sequenceDiagram
    participant Client
    participant Service
    participant Repo as Repository
    participant DB as GORM/DB
    participant OTel as OTel Plugin
    participant Jaeger

    Client->>Service: HTTP Request (Headers: trace-id)
    Note right of Client: B3/W3C Propagation
    Service->>Repo: List(ctx)
    Repo->>DB: Find(ctx)
    DB->>OTel: BeforeQuery()
    OTel->>OTel: Start Span "gorm.query"
    OTel->>OTel: Inject SQL tags
    DB->>DB: Execute SQL
    DB->>OTel: AfterQuery()
    OTel->>Jaeger: Export Span (trace-id from ctx)
    DB-->>Repo: Result
    Repo-->>Service: Result
```

## Build and Verification

### Unit Tests
Run the comprehensive test suite (includes SQLite in-memory tests).

**Using Go:**
```bash
go test -v ./pkg/database/...
```

**Using Bazel:**
```bash
bazel test //pkg/database/...
```

### Build
To build the library:

**Using Go:**
```bash
go build ./pkg/database/...
```

**Using Bazel:**
```bash
bazel build //pkg/database
```

## Usage Examples

### Basic Setup
```go
cfg := config.DatabaseConfig{
    Driver: "postgres",
    Host: "localhost",
    // ...
}
db, err := database.New(cfg, zapLogger)
```

### Using Generic Repository
```go
type User struct {
    ID   uint `gorm:"primarykey"`
    Name string
}

// Create Repo
repo := database.NewRepository[User](db.DB)

// CRUD
repo.Create(ctx, &User{Name: "Ganesh"})

// List with Filters
repo.List(ctx, database.Pagination{
    Page: 1,
    Filters: map[string]interface{}{"name": "Ganesh"},
})
```
