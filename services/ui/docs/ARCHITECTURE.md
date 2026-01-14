# Polaris - Cloud-Native Router Orchestrator Platform
## Complete Architecture & Design Document

**Version:** 1.0  
**Last Updated:** January 2026  
**Status:** Design Complete, Implementation Pending

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [System Architecture Overview](#system-architecture-overview)
3. [Technology Stack](#technology-stack)
4. [Core Components](#core-components)
5. [Authentication Architecture](#authentication-architecture)
6. [Microservices Layer](#microservices-layer)
7. [Data Architecture](#data-architecture)
8. [Multi-Tenancy Strategy](#multi-tenancy-strategy)
9. [Multi-Cloud Deployment](#multi-cloud-deployment)
10. [Infrastructure as Code (Pulumi)](#infrastructure-as-code)
11. [Security Architecture](#security-architecture)
12. [Data Flows & Use Cases](#data-flows--use-cases)
13. [Performance & Scalability](#performance--scalability)
14. [Observability](#observability)
15. [DevOps & CI/CD](#devops--cicd)
16. [Implementation Roadmap](#implementation-roadmap)

---

## Executive Summary

**Polaris** is a production-ready, cloud-native platform for orchestrating network routers across multi-cloud and on-premise environments. Built on event-driven microservices architecture with Go, NATS, and Kubernetes, it provides enterprise-grade multi-tenancy, real-time monitoring, GitOps-based configuration management, and AI-powered operations assistance.

### Key Capabilities
- ✅ **Multi-Tenant SaaS** - Schema-per-tenant isolation with Keycloak realms
- ✅ **Multi-Cloud** - AWS, GCP, Azure, On-Premise support
- ✅ **Real-Time Operations** - WebSocket updates, 100k metrics/sec ingestion
- ✅ **Event-Driven** - NATS JetStream for async messaging
- ✅ **GitOps Ready** - Git-backed configuration management
- ✅ **AI-Powered** - LLM-based troubleshooting assistant
- ✅ **Production UI** - Complete React dashboard with 8 functional pages

### Performance Targets
| Metric | Target |
|--------|--------|
| API Latency (p95) | < 100ms |
| Metrics Ingestion | 100k points/sec |
| Concurrent Users | 10,000+ |
| Managed Routers | 100,000+ |
| Uptime SLA | 99.95% |

---

## System Architecture Overview

### High-Level Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLIENT LAYER                            │
│                                                                 │
│  React UI (Vite + TypeScript + Tailwind)                       │
│  Pages: Dashboard, Network, Security, Analytics, Chat,         │
│         Alerts, Config, Workflows, Settings                     │
│                                                                 │
│  Authentication: Keycloak JS Adapter                            │
│  Real-Time: WebSocket connection                                │
└──────────────────┬──────────────────────────────────────────────┘
                   │ HTTPS/WSS
                   ▼
┌─────────────────────────────────────────────────────────────────┐
│                         EDGE LAYER                              │
│                                                                 │
│  ┌───────────────┐  ┌───────────────┐  ┌──────────────────┐   │
│  │   Ingress     │  │   Keycloak    │  │   WebSocket      │   │
│  │  Controller   │  │   (Auth)      │  │   Server         │   │
│  │               │  │               │  │                  │   │
│  │ - TLS Term    │  │ - OAuth2/OIDC │  │ - Live Updates   │   │
│  │ - Rate Limit  │  │ - Multi-realm │  │ - Alerts Push    │   │
│  │ - DDoS        │  │ - SSO/2FA     │  │ - Topology       │   │
│  └───────┬───────┘  └───────────────┘  └──────────────────┘   │
└──────────┼──────────────────────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────────────────┐
│                       API GATEWAY LAYER                         │
│                                                                 │
│  API Gateway (Go + Fiber)                                       │
│  - JWT Validation (Keycloak public key)                         │
│  - Tenant Extraction (from JWT realm)                           │
│  - RBAC Authorization                                           │
│  - Request Routing                                              │
│  - Per-Tenant Rate Limiting                                     │
└──────────────────┬──────────────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────────────┐
│                    MICROSERVICES LAYER (Go)                     │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │   Router     │  │    Config    │  │   Metrics    │         │
│  │   Manager    │  │   Service    │  │  Collector   │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │    Alert     │  │   Workflow   │  │    Policy    │         │
│  │   Engine     │  │ Orchestrator │  │   Engine     │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │   Tenant     │  │    Audit     │  │  AI Agent    │         │
│  │   Manager    │  │   Logger     │  │   (LLM)      │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
└──────────────────┬──────────────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────────────┐
│                      MESSAGE BUS LAYER                          │
│                                                                 │
│  NATS JetStream                                                 │
│  - Event Streaming (10M+ msgs/sec)                              │
│  - At-least-once delivery                                      │
│  - Topics: router.*, config.*, metrics.*, alert.*              │
│  - Protocol Buffers serialization                              │
└──────────────────┬──────────────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────────────┐
│                        DATA LAYER                               │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │  PostgreSQL  │  │   InfluxDB   │  │    Redis     │         │
│  │   (Primary)  │  │ (Time-Series)│  │   (Cache)    │         │
│  │              │  │              │  │              │         │
│  │ Schema-per-  │  │ Tenant Tags  │  │ Key Prefixes │         │
│  │ Tenant       │  │              │  │              │         │
│  │              │  │              │  │              │         │
│  │ - Routers    │  │ - Metrics    │  │ - Sessions   │         │
│  │ - Configs    │  │ - Telemetry  │  │ - Metrics    │         │
│  │ - Users      │  │ - Events     │  │ - Rate Limits│         │
│  │ - Policies   │  │              │  │ - Alerts     │         │
│  │ - Workflows  │  │ Downsampling │  │ TTL: 30s-24h │         │
│  │ - Audit Logs │  │ 10s→1m→1h    │  │              │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
└─────────────────────────────────────────────────────────────────┘
                   ▲
                   │ NATS/TLS
                   │
┌─────────────────────────────────────────────────────────────────┐
│                    ROUTER AGENTS LAYER                          │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │Router Agent  │  │Router Agent  │  │Router Agent  │  (100k+)│
│  │   (AWS)      │  │   (GCP)      │  │  (Azure)     │         │
│  │              │  │              │  │              │         │
│  │ Publishes:   │  │ Publishes:   │  │ Publishes:   │         │
│  │ - Metrics    │  │ - Metrics    │  │ - Metrics    │         │
│  │ - Events     │  │ - Events     │  │ - Events     │         │
│  │              │  │              │  │              │         │
│  │ Subscribes:  │  │ Subscribes:  │  │ Subscribes:  │         │
│  │ - Config     │  │ - Config     │  │ - Config     │         │
│  │ - Commands   │  │ - Commands   │  │ - Commands   │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
└─────────────────────────────────────────────────────────────────┘
```

---

## Technology Stack

### Frontend
| Component | Technology | Version |
|-----------|------------|---------|
| Framework | React | 18.x |
| Build Tool | Vite | 5.x |
| Language | TypeScript | 5.x |
| Styling | Tailwind CSS | 3.x |
| UI Components | shadcn/ui | Latest |
| Charts | Recharts | 2.x |
| Maps | Leaflet + react-leaflet | Latest |
| Auth Client | Keycloak JS | Latest |

### Backend
| Component | Technology | Version |
|-----------|------------|---------|
| Language | Go | 1.21+ |
| HTTP Framework | Fiber | 2.x |
| Message Bus | NATS JetStream | 2.10+ |
| Primary Database | PostgreSQL | 15+ |
| Time-Series DB | InfluxDB | 2.x |
| Cache | Redis | 7+ |
| Authentication | Keycloak | 23+ |
| API Protocol | REST + gRPC | - |
| Serialization | Protocol Buffers | 3 |

### Infrastructure
| Component | Technology |
|-----------|------------|
| Container Orchestration | Kubernetes 1.28+ |
| Service Mesh | Istio 1.20+ |
| Ingress | NGINX/Traefik |
| IaC | **Pulumi (Go)** |
| Helm | Kubernetes packaging |
| CI/CD | GitHub Actions + ArgoCD |

### Observability
| Component | Technology |
|-----------|------------|
| Metrics | Prometheus + Grafana |
| Distributed Tracing | Jaeger / Tempo |
| Logging | Loki + Promtail |
| APM | OpenTelemetry |

---

## Core Components

### 1. Client Layer (React UI)

**Status:** ✅ Complete

**Pages:**
- **Dashboard** - Network health overview, key metrics
- **Network** - Router inventory, topology map, VPCs, subnets
- **Security** - Threat map, policies, audit logs
- **Analytics** - Traffic forecasting, reports
- **Alerts** - Incident management, active alerts
- **Config** - Configuration editor with diff viewer
- **Workflows** - Automation templates, execution history
- **Settings** - User preferences, API keys
- **Chat** - AI-powered assistant for troubleshooting

**Key Features:**
- Dark/Light theme support
- Real-time updates via WebSocket
- Responsive design (mobile-ready)
- Keycloak integration for SSO
- Session management
- Toast notifications

### 2. Edge Layer

#### Ingress Controller (NGINX/Traefik)
**Responsibilities:**
- TLS termination (HTTPS → HTTP)
- Global rate limiting (per-IP)
- DDoS protection
- Load balancing across API Gateway pods
- Request routing based on host/path

**Routes:**
```
auth.polaris.io  → Keycloak (port 8080)
api.polaris.io   → API Gateway (port 8080)
ws.polaris.io    → WebSocket Server (port 8080)
app.polaris.io   → React UI (static files)
```

#### Keycloak (Authentication Server)
**Deployment:** Standalone service (direct UI access)
**Features:**
- OAuth2/OIDC flows
- Multi-realm architecture (1 realm = 1 tenant)
- SSO/SAML integration
- Social login (Google, GitHub)
- 2FA/MFA support
- LDAP/Active Directory federation
- User self-service portal

**JWT Token Structure:**
```json
{
  "iss": "https://auth.polaris.io/realms/acme-corp",
  "sub": "user-uuid",
  "exp": 1736792400,
  "realm_access": {
    "roles": ["admin", "operator"]
  },
  "email": "john@acme.com"
}
```

#### WebSocket Server (Go)
**Purpose:** Push real-time updates to UI
**Subscriptions:**
- Metrics updates (router CPU, bandwidth, etc.)
- Topology changes (BGP peer down, link status)
- Alerts (critical events)
- Workflow status updates

---

## Authentication Architecture

### Flow: User Login

```
1. User visits app.polaris.io/dashboard
2. UI detects no JWT → redirects to Keycloak
   https://auth.polaris.io/realms/acme-corp/protocol/openid-connect/auth

3. User sees Keycloak login page
4. User enters credentials
5. Keycloak validates → redirects back with auth code
   https://app.polaris.io/callback?code=xyz123

6. UI exchanges code for tokens:
   POST https://auth.polaris.io/realms/acme-corp/protocol/openid-connect/token
   Returns: { access_token, refresh_token, expires_in: 900 }

7. UI stores tokens in localStorage
8. UI navigates to /dashboard
```

### Flow: API Request with JWT

```
1. UI makes API call:
   GET https://api.polaris.io/api/v1/routers
   Header: Authorization: Bearer <JWT>

2. Ingress Controller forwards to API Gateway (no auth check)

3. API Gateway:
   - Extracts JWT from Authorization header
   - Validates signature using Keycloak public key (cached)
   - Checks expiration
   - Extracts tenant_id from "iss" claim (realm name)
   - Validates user roles

4. API Gateway calls Router Service:
   GET http://router-service:8080/routers
   Headers:
     X-Tenant-ID: acme-corp
     X-User-ID: user-uuid
     X-User-Roles: admin,operator

5. Router Service queries database:
   SELECT * FROM tenant_acme_corp.routers

6. Response flows back to UI
```

### API Gateway - JWT Middleware (Go)

```go
func JWTMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
        
        // Validate JWT signature
        claims, err := validateJWT(token)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
        }
        
        // Extract tenant from issuer
        // "https://auth.polaris.io/realms/acme-corp" → "acme-corp"
        tenantID := extractTenantFromIssuer(claims.Issuer)
        
        // Check expiration
        if time.Now().Unix() > claims.ExpiresAt {
            return c.Status(401).JSON(fiber.Map{"error": "token expired"})
        }
        
        // RBAC: Check required roles
        if !hasRequiredRole(claims.RealmAccess.Roles, c.Path()) {
            return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
        }
        
        // Inject context
        c.Locals("tenant_id", tenantID)
        c.Locals("user_id", claims.Subject)
        c.Locals("roles", claims.RealmAccess.Roles)
        
        return c.Next()
    }
}
```

---

## Microservices Layer

### 1. API Gateway
**Language:** Go (Fiber framework)
**Port:** 8080
**Responsibilities:**
- JWT validation
- Tenant context injection
- Request routing to services
- Per-tenant rate limiting
- API versioning
- REST & GraphQL (optional)

**Endpoints:**
```
POST   /api/v1/auth/login
GET    /api/v1/routers
POST   /api/v1/routers
GET    /api/v1/routers/:id
POST   /api/v1/routers/:id/config
GET    /api/v1/metrics/routers/:id
POST   /api/v1/workflows/:id/execute
GET    /api/v1/alerts
POST   /api/v1/policies
GET    /api/v1/tenants/:id
```

### 2. Router Manager Service
**Purpose:** CRUD operations for routers, inventory management

**Database Schema (PostgreSQL):**
```sql
CREATE TABLE tenant_{tenant_id}.routers (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    region VARCHAR(100),
    cloud_provider VARCHAR(50), -- AWS, GCP, Azure, On-Prem
    status VARCHAR(50), -- online, offline, alert
    asn VARCHAR(20),
    ip_address INET,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

**NATS Topics:**
- Publishes: `router.created`, `router.updated`, `router.deleted`
- Subscribes: `router.health.check`

### 3. Config Management Service
**Purpose:** Git-backed configuration storage, versioning, deployment

**Features:**
- GitOps integration (watch GitHub repo)
- Configuration versioning
- Diff generation (current vs desired)
- Rollback support
- Template rendering (Jinja-like)

**Database Schema:**
```sql
CREATE TABLE tenant_{tenant_id}.router_configs (
    id UUID PRIMARY KEY,
    router_id UUID NOT NULL,
    version INT NOT NULL,
    config_yaml TEXT NOT NULL,
    applied_at TIMESTAMPTZ,
    created_by UUID,
    git_commit_sha VARCHAR(40),
    FOREIGN KEY (router_id) REFERENCES routers(id)
);
```

**NATS Topics:**
- Publishes: `config.deploy.request.{router_id}`
- Subscribes: `config.deploy.result.{router_id}`

### 4. Metrics Collector Service
**Purpose:** Ingest metrics from routers, batch write to InfluxDB

**Flow:**
```
1. Subscribe to NATS: router.metrics.*
2. Receive metrics (CPU, memory, network)
3. Accumulate batch (1000 points)
4. Write to InfluxDB
5. Update Redis cache (current metrics, TTL 30s)
6. Publish to NATS: metrics.updated
```

**InfluxDB Schema:**
```
Measurement: router_cpu
Tags: tenant_id, router_id, region, cloud
Fields: usage_percent, cores
Timestamp: nanosecond precision

Measurement: router_network
Tags: tenant_id, router_id, interface
Fields: bytes_in, bytes_out, packets_in, packets_out, errors
```

### 5. Alert Engine Service
**Purpose:** Rule evaluation, multi-channel notifications

**Rule Example:**
```go
type AlertRule struct {
    ID          string
    TenantID    string
    Condition   string // "cpu_usage > 80"
    Duration    time.Duration // 5min
    Severity    string // critical, warning, info
    Channels    []string // email, slack, pagerduty
}
```

**Deduplication:**
- Redis: `alert:state:{router_id}:{rule_id}` = OK | ALERT
- Only trigger on state change

**NATS Topics:**
- Subscribes: `metrics.updated`
- Publishes: `alert.triggered`, `alert.resolved`

### 6. Workflow Orchestrator Service
**Purpose:** Execute automation workflows (Temporal.io-based)

**Workflow Example (YAML):**
```yaml
name: drain-and-reboot
steps:
  - name: set-bgp-weight
    action: router_command
    params:
      command: "set bgp neighbor weight 0"
    timeout: 5m
  
  - name: wait-for-drain
    action: wait_until
    condition: "active_connections < 10"
    timeout: 10m
  
  - name: reboot
    action: exec
    params:
      command: "sudo reboot"
```

**NATS Topics:**
- Publishes: `router.{router_id}.command`
- Subscribes: `command.result`

### 7. Policy Engine Service
**Purpose:** Firewall rules, BGP policies, traffic shaping

**Database Schema:**
```sql
CREATE TABLE tenant_{tenant_id}.policies (
    id UUID PRIMARY KEY,
    name VARCHAR(255),
    type VARCHAR(50), -- firewall, bgp, traffic_shaping
    rules JSONB,
    scope JSONB, -- which routers
    status VARCHAR(20) -- active, inactive
);
```

**NATS Topics:**
- Publishes: `router.{router_id}.policy`
- Subscribes: `policy.applied`

### 8. Tenant Manager Service
**Purpose:** Onboarding, quota management, billing

**Onboarding Flow:**
```
1. Create Keycloak realm
2. Create PostgreSQL schema (tenant_{tenant_id})
3. Create InfluxDB bucket
4. Set Redis quotas
5. Send invite email to admin
```

### 9. Audit Logger Service
**Purpose:** Log all events for compliance

**NATS Topics:**
- Subscribes: `*.created`, `*.updated`, `*.deleted`, `config.*`, `workflow.*`

**Database Schema:**
```sql
CREATE TABLE tenant_{tenant_id}.audit_logs (
    id UUID PRIMARY KEY,
    event_type VARCHAR(100),
    user_id UUID,
    resource_type VARCHAR(100),
    resource_id UUID,
    details JSONB,
    timestamp TIMESTAMPTZ DEFAULT NOW()
);
```

### 10. AI Agent Service
**Purpose:** LLM-powered troubleshooting assistant

**Tech Stack:**
- LLM: OpenAI GPT-4 / Anthropic Claude / self-hosted Llama 2
- Vector DB: Pinecone/Weaviate (for RAG)
- Context: Queries other services for router data, metrics, alerts

**Capabilities:**
- Conversational Q&A ("Why is router-123 slow?")
- Config generation ("Create BGP peer config for...")
- Troubleshooting ("Users in EU reporting latency")
- Automation execution ("Drain and reboot all routers in us-east with CPU > 85%")
- Predictive insights ("Will we have capacity issues next month?")

---

## Data Architecture

### PostgreSQL (Primary Database)

**Multi-Tenancy:** Schema-per-Tenant

```sql
-- Each tenant gets a schema
CREATE SCHEMA tenant_acme_corp;
CREATE SCHEMA tenant_globex_corp;

-- Schema structure (same for all tenants)
CREATE TABLE tenant_{tenant_id}.routers (...);
CREATE TABLE tenant_{tenant_id}.router_configs (...);
CREATE TABLE tenant_{tenant_id}.users (...);
CREATE TABLE tenant_{tenant_id}.policies (...);
CREATE TABLE tenant_{tenant_id}.workflows (...);
CREATE TABLE tenant_{tenant_id}.alerts (...);
CREATE TABLE tenant_{tenant_id}.audit_logs (...);
```

**Connection Pooling:** pgxpool (5000 max connections)
**Replication:** Primary + 3 Read Replicas (for queries)
**Backup:** Automated daily snapshots, 30-day retention

### InfluxDB (Time-Series Database)

**Multi-Tenancy:** Tag-based filtering

**Measurements:**
- `router_cpu` → usage_percent, cores
- `router_memory` → used_bytes, total_bytes, usage_percent
- `router_network` → bytes_in, bytes_out, packets_in, packets_out
- `bgp_peers` → state, uptime, routes_received
- `link_latency` → latency_ms, jitter_ms, packet_loss_percent

**Tags (all measurements):**
- `tenant_id`
- `router_id`
- `region`
- `cloud` (AWS, GCP, Azure, On-Prem)

**Retention Policies:**
- Raw data: 7 days
- 1-minute aggregates: 30 days
- 1-hour aggregates: 1 year

**Continuous Queries:** Auto-downsample for long-term storage

### Redis (Cache & Session Store)

**Multi-Tenancy:** Key prefixing

**Key Patterns:**
```
tenant:{tenant_id}:session:{user_id}            # TTL: 24h
tenant:{tenant_id}:metrics:current:{router_id}  # TTL: 30s
tenant:{tenant_id}:ratelimit:{ip}               # TTL: 1min
alert:state:{router_id}:{rule_id}               # TTL: 1h
config:cache:{router_id}                        # TTL: 5min
ai:conversation:{user_id}:history               # TTL: 7d
```

**Deployment:** Redis Cluster (3 primary + 3 replicas, sharding enabled)

---

## Multi-Tenancy Strategy

### Isolation Model: Schema-per-Tenant

**PostgreSQL:**
```
Tenant A → Schema: tenant_acme
Tenant B → Schema: tenant_globex
```

**InfluxDB:**
```
All data in same bucket, filtered by tenant_id tag
WHERE tenant_id = 'acme'
```

**Redis:**
```
tenant:acme:*
tenant:globex:*
```

**Keycloak:**
```
Realm: acme-corp (users, roles, clients)
Realm: globex-corp (users, roles, clients)
```

### Tenant Context Flow

```
1. User logs in via Keycloak realm "acme-corp"
2. JWT issued with iss: "https://auth.polaris.io/realms/acme-corp"
3. API Gateway extracts tenant_id = "acme-corp" from JWT
4. API Gateway injects X-Tenant-ID header
5. All services use tenant_id for database queries
6. PostgreSQL: SET search_path TO tenant_acme_corp
7. InfluxDB: Query with WHERE tenant_id = 'acme-corp'
8. Redis: Keys prefixed with tenant:acme:
```

### Resource Quotas (per tenant)

```go
type TenantQuota struct {
    MaxRouters      int
    MaxBandwidth    int64  // Gbps
    MaxUsers        int
    MaxAPIRate      int    // requests/sec
    MaxStorage      int64  // GB
    MaxWorkflows    int
}
```

---

## Multi-Cloud Deployment

### Deployment Topology

```
┌─────────────────────────────────────┐
│     AWS Region (us-east-1)          │
│  - EKS Cluster                      │
│  - RDS PostgreSQL (Multi-AZ)        │
│  - ElastiCache Redis                │
│  - InfluxDB (EC2)                   │
│  - Router Agents (EC2)              │
└──────────────┬──────────────────────┘
               │ Cross-Region Replication
┌──────────────┴──────────────────────┐
│    GCP Region (europe-west1)        │
│  - GKE Cluster                      │
│  - Cloud SQL PostgreSQL             │
│  - Memorystore Redis                │
│  - InfluxDB (Compute Engine)        │
│  - Router Agents (Compute Engine)   │
└──────────────┬──────────────────────┘
               │
┌──────────────┴──────────────────────┐
│     Azure Region (eastus)           │
│  - AKS Cluster                      │
│  - Azure Database for PostgreSQL    │
│  - Azure Cache for Redis            │
│  - InfluxDB (VM)                    │
│  - Router Agents (VMs)              │
└─────────────────────────────────────┘
```

### Cloud Provider Abstraction

```go
type CloudProvider interface {
    CreateRouter(ctx context.Context, spec RouterSpec) (*Router, error)
    DeleteRouter(ctx context.Context, id string) error
    GetMetrics(ctx context.Context, id string) (*Metrics, error)
    ConfigureNetwork(ctx context.Context, config NetworkConfig) error
}

// Implementations
type AWSProvider struct { ... }
type GCPProvider struct { ... }
type AzureProvider struct { ... }
type OnPremiseProvider struct { ... }
```

### Cross-Cluster Communication

- **Service Mesh:** Istio multi-cluster federation
- **NATS:** Leaf nodes with gateway connections
- **Database:** PostgreSQL logical replication
- **Disaster Recovery:** Automated failover (RTO: 5min, RPO: 1min)

---

## Infrastructure as Code (Pulumi)

### Why Pulumi (Go)
- ✅ Same language as backend services
- ✅ Full Go type system and IDE support
- ✅ Native testing with Go testing framework
- ✅ Code reuse between infrastructure and application
- ✅ Multi-cloud support (AWS, GCP, Azure)

### Project Structure

```
infrastructure/
├── main.go              # Entry point
├── Pulumi.yaml          # Project config
├── Pulumi.dev.yaml      # Dev environment
├── Pulumi.prod.yaml     # Prod environment
│
├── aws/
│   ├── vpc.go           # VPC, subnets, security groups
│   ├── eks.go           # EKS cluster
│   ├── rds.go           # PostgreSQL RDS
│   └── elasticache.go   # Redis
│
├── gcp/
│   ├── network.go       # VPC network
│   ├── gke.go           # GKE cluster
│   └── cloudsql.go      # Cloud SQL PostgreSQL
│
└── kubernetes/
    ├── namespaces.go    # Create namespaces
    ├── services.go      # Deploy microservices
    └── ingress.go       # Ingress controllers
```

### Example: AWS Infrastructure (Pulumi Go)

```go
func deployAWS(ctx *pulumi.Context) error {
    // Create VPC
    vpc, _ := ec2.NewVpc(ctx, "polaris-vpc", &ec2.VpcArgs{
        CidrBlock: pulumi.String("10.0.0.0/16"),
    })
    
    // Create EKS cluster
    cluster, _ := eks.NewCluster(ctx, "polaris-eks", &eks.ClusterArgs{
        VpcId: vpc.ID(),
        Version: pulumi.String("1.28"),
    })
    
    // Create RDS PostgreSQL
    db, _ := rds.NewInstance(ctx, "polaris-postgres", &rds.InstanceArgs{
        Engine:        pulumi.String("postgres"),
        InstanceClass: pulumi.String("db.r6g.xlarge"),
        MultiAz:       pulumi.Bool(true),
    })
    
    // Export outputs
    ctx.Export("kubeconfig", cluster.KubeconfigJson)
    ctx.Export("dbEndpoint", db.Endpoint)
    
    return nil
}
```

### Deployment Commands

```bash
# Initialize stack
pulumi stack init prod

# Preview changes
pulumi preview

# Deploy infrastructure
pulumi up

# View outputs
pulumi stack output kubeconfig > ~/.kube/polaris-prod

# Destroy (when needed)
pulumi destroy
```

---

## Security Architecture

### Defense in Depth

**1. Network Layer**
- VPC isolation (private subnets for databases)
- Security groups / firewall rules
- TLS everywhere (mTLS for inter-service)
- WAF (Web Application Firewall) at edge

**2. Application Layer**
- JWT authentication (15min expiry)
- RBAC enforcement (admin, operator, viewer)
- API rate limiting (per-tenant)
- Input validation
- SQL injection prevention (parameterized queries)
- XSS protection

**3. Data Layer**
- Encryption at rest (AES-256)
- Encryption in transit (TLS 1.3)
- Database access control (least privilege)
- Automated backup encryption

**4. Secrets Management**
- HashiCorp Vault or Kubernetes Secrets
- Automated secret rotation (90 days)
- No secrets in code/config
- Keycloak for API keys

**5. Audit & Compliance**
- Every API call logged
- Config change audit trail
- Compliance reports (CIS, NIST, ISO 27001, SOC 2)
- GDPR-compliant data handling

### RBAC Roles

| Role | Permissions |
|------|-------------|
| **Admin** | Full access (create/update/delete all resources) |
| **Operator** | Manage routers, deploy configs, execute workflows |
| **Viewer** | Read-only access to dashboards, metrics |
| **Auditor** | Access audit logs, generate reports |

---

## Data Flows & Use Cases

### Use Case 1: User Views Dashboard

```
1. UI sends: GET /api/v1/dashboard (JWT in header)
2. API Gateway validates JWT → extracts tenant_id = "acme"
3. API Gateway calls Router Service
4. Router Service queries:
   - PostgreSQL: SELECT * FROM tenant_acme.routers
   - Redis: GET tenant:acme:metrics:current:*
5. Router Service aggregates data
6. API Gateway returns JSON response
7. UI renders dashboard
```

### Use Case 2: Router Publishes Metrics

```
1. Router Agent collects metrics (CPU: 65%, Memory: 78%)
2. Serialize to Protocol Buffers
3. Publish to NATS: router.metrics.router-123
4. Metrics Collector receives message
5. Batch accumulation (wait for 1000 points)
6. Write batch to InfluxDB
7. Update Redis: tenant:acme:metrics:current:router-123 (TTL 30s)
8. Publish to NATS: metrics.updated
9. WebSocket Server receives, pushes to UI
10. UI updates chart in real-time
```

### Use Case 3: GitOps Config Deployment

```
1. Admin pushes config to GitHub (main branch)
2. GitHub webhook triggers Config Service
3. Config Service pulls latest commit
4. Parse YAML config
5. Query PostgreSQL for current config (version 42)
6. Generate diff (old vs new)
7. Insert new config (version 43, git_sha)
8. Publish to NATS: config.deploy.request.router-123
9. Router Agent receives message
10. Validate config syntax
11. Apply config via SSH/NETCONF
12. Publish to NATS: config.deploy.result (success)
13. Config Service marks as applied
14. Audit Logger logs event
15. Slack notification: "Config v43 deployed to router-123"
```

### Use Case 4: Alert Triggered

```
1. Metrics Collector writes CPU: 92% to InfluxDB
2. Alert Engine queries InfluxDB every 15s
3. Evaluate rule: "cpu_usage > 80 for 5min"
4. Condition MET (sustained 5 minutes)
5. Check Redis: alert:state:router-123:cpu_high = OK
6. State changed: OK → ALERT
7. Update Redis: alert:state:router-123:cpu_high = ALERT
8. Insert into PostgreSQL: alerts table
9. Send notifications:
   - PagerDuty: Create incident
   - Email: ops@acme.com
   - Slack: #alerts channel
10. Publish to NATS: alert.triggered
11. WebSocket Server pushes to UI
12. UI displays toast: "⚠️ High CPU on router-123"
```

### Use Case 5: AI Agent Troubleshooting

```
User: "Why is router-us-west-1 slow?"

AI Agent:
1. Query Router Service: GetRouter("us-west-1")
2. Query Metrics Service: GetMetrics("us-west-1", last_hour)
3. Query Alert Service: GetActiveAlerts("us-west-1")
4. Query Audit Service: GetRecentChanges("us-west-1")
5. Analyze data:
   - CPU: 92% (high)
   - BGP routes: increased 40% in last hour
   - Config change: added OSPF neighbors 25min ago
   - Active connections: 8,245 (normal: 5k)
6. Generate response using LLM

Response:
"Router us-west-1 is slow due to:
1. High CPU (92%) caused by increased BGP route table
2. Recent config change added OSPF neighbors (25min ago)
3. Above-normal connections (8.2k vs 5k baseline)

Recommended actions:
- Scale to c5.xlarge instance ✅
- Review BGP route filtering
- Check for traffic anomalies

Shall I execute the scaling workflow?"
```

---

## Performance & Scalability

### Performance Targets

| Metric | Target | Strategy |
|--------|--------|----------|
| **API Latency (p95)** | < 100ms | Redis caching, read replicas, connection pooling |
| **UI Load Time** | < 2s | CDN, code splitting, lazy loading |
| **Metrics Ingestion** | 100k points/sec | InfluxDB batching, NATS streaming |
| **Concurrent Users** | 10,000+ | Horizontal scaling, session caching |
| **Supported Routers** | 100,000+ | NATS horizontal scaling, sharding |
| **Database Queries** | < 50ms | Indexes, materialized views, connection pooling |

### Horizontal Scaling

**All services are stateless:**
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: router-service-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: router-service
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

### Database Scaling

**PostgreSQL:**
- Primary (write) + 5-10 Read Replicas (queries)
- Connection pooling: pgxpool (5000 max)
- Partitioning: audit_logs table by month

**InfluxDB:**
- 3-node cluster (minimum)
- Replication factor: 2
- Shard duration: 7 days

**Redis:**
- Cluster mode with 6 nodes (3 primary + 3 replicas)
- Hash slot sharding
- Read replicas for read-heavy workloads

### NATS Scaling

- 3-node cluster (minimum for HA)
- Add nodes for more throughput (10M+ msgs/sec per node)
- Leaf nodes for geo-distribution

---

## Observability

### Metrics (Prometheus)

**Service Metrics:**
```
# API Gateway
api_requests_total{method, status, tenant}
api_request_duration_seconds{method, endpoint}
api_active_connections

# Router Service
router_operations_total{operation, status}
router_query_duration_seconds

# Config Service
config_deploys_total{status}
config_deploy_duration_seconds

# Metrics Collector
metrics_ingested_total{source}
metrics_batch_size
```

**Business Metrics:**
```
routers_online{tenant, region, cloud}
alerts_active{severity, tenant}
workflows_executed_total{template, status}
```

### Distributed Tracing (Jaeger)

**Trace Example:**
```
Trace ID: abc123
Span 1: API Gateway (10ms)
  ├─ Span 2: Router Service (30ms)
  │   ├─ Span 3: PostgreSQL Query (15ms)
  │   └─ Span 4: Redis Get (5ms)
  └─ Span 5: Metrics Service (20ms)
      └─ Span 6: InfluxDB Query (18ms)
Total: 60ms
```

**Context Propagation:**
```go
// API Gateway injects trace context into NATS message
natsMsg := &nats.Msg{
    Subject: "router.metrics.123",
    Data:    payload,
    Header: map[string][]string{
        "traceparent": {traceID},
    },
}
```

### Logging (Loki + Promtail)

**Structured Logging (JSON):**
```json
{
  "timestamp": "2026-01-13T15:30:00Z",
  "level": "info",
  "service": "router-service",
  "tenant_id": "acme",
  "user_id": "user-123",
  "trace_id": "abc123",
  "message": "Router created",
  "router_id": "rtr-456"
}
```

### Dashboards (Grafana)

**System Health Dashboard:**
- Service uptime
- Request rate (RPS)
- Error rate (5xx)
- Latency (p50, p95, p99)
- Database connection pool status

**Tenant Dashboard:**
- Routers online/offline
- Active alerts
- Top routers by CPU/bandwidth
- Recent config deployments

**Cost Dashboard:**
- Cloud spend per tenant
- Compute costs
- Data transfer costs

---

## DevOps & CI/CD

### CI/CD Pipeline

```
Git Push (main branch)
  ↓
GitHub Actions
  ├─ Lint (golangci-lint)
  ├─ Unit Tests
  ├─ Build Docker Images
  └─ Push to Registry (Docker Hub / ECR)
  ↓
Update Helm Chart (values.yaml)
  ↓
Commit to GitOps Repo
  ↓
ArgoCD Detects Change
  ↓
Deploy to Kubernetes
  ├─ Rolling Update
  ├─ Health Checks
  └─ Smoke Tests
  ↓
Canary Deployment (10% → 50% → 100%)
  ↓
Full Rollout
```

### GitHub Actions Workflow

```yaml
name: Build and Deploy

on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run tests
        run: go test -v ./...
      
      - name: Build Docker image
        run: docker build -t polaris/router-service:${{ github.sha }} .
      
      - name: Push to registry
        run: docker push polaris/router-service:${{ github.sha }}
  
  deploy-staging:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - name: Update staging
        run: |
          kubectl set image deployment/router-service \
            router-service=polaris/router-service:${{ github.sha }} \
            -n staging
  
  deploy-production:
    needs: deploy-staging
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
      - name: Update ArgoCD app
        run: argocd app sync polaris-prod
```

### ArgoCD Application

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: polaris-prod
spec:
  project: default
  source:
    repoURL: https://github.com/polaris/gitops
    path: environments/production
    targetRevision: HEAD
  destination:
    server: https://kubernetes.default.svc
    namespace: polaris
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
```

---

## Implementation Roadmap

### Phase 1: Foundation (Weeks 1-4) ✅ COMPLETE
- [x] UI Development (all 8 pages)
- [x] Architecture design
- [x] Technology selection
- [x] Documentation

### Phase 2: Core Backend (Weeks 5-10) ⏭️ NEXT
- [ ] Set up Go monorepo
- [ ] API Gateway + JWT middleware
- [ ] Router Manager Service
- [ ] PostgreSQL schema setup
- [ ] Docker Compose for local dev (NATS, PostgreSQL, Redis, Keycloak)
- [ ] Basic CRUD operations

### Phase 3: Metrics & Monitoring (Weeks 11-14)
- [ ] Metrics Collector Service
- [ ] InfluxDB integration
- [ ] WebSocket Server for real-time updates
- [ ] Connect UI to backend APIs
- [ ] Prometheus metrics instrumentation

### Phase 4: Configuration Management (Weeks 15-18)
- [ ] Config Service with GitOps
- [ ] Router Agent (Go binary)
- [ ] NATS-based config deployment
- [ ] Rollback functionality

### Phase 5: Orchestration & Alerts (Weeks 19-22)
- [ ] Alert Engine Service
- [ ] Workflow Orchestrator (Temporal.io)
- [ ] Policy Engine
- [ ] PagerDuty/Slack integrations

### Phase 6: Multi-Cloud & IaC (Weeks 23-26)
- [ ] Pulumi infrastructure code
- [ ] Deploy to AWS (dev environment)
- [ ] Deploy to GCP (staging)
- [ ] Cross-cloud testing

### Phase 7: AI & Advanced Features (Weeks 27-30)
- [ ] AI Agent Service
- [ ] LLM integration (GPT-4/Claude)
- [ ] RAG for context-aware responses
- [ ] Tenant Manager Service

### Phase 8: Production Hardening (Weeks 31-34)
- [ ] Security audit
- [ ] Performance optimization
- [ ] Load testing (100k routers simulation)
- [ ] Disaster recovery testing
- [ ] Documentation finalization

### Phase 9: Beta Launch (Week 35+)
- [ ] Onboard first 5 beta customers
- [ ] Collect feedback
- [ ] Bug fixes & iterations

### Phase 10: GA (General Availability)
- [ ] Public launch
- [ ] Marketing & sales enablement
- [ ] 24/7 support setup

---

## Conclusion

**Polaris** is designed as a production-ready, enterprise-grade platform from day one. The architecture supports:

- ✅ **Massive Scale:** 100k+ routers, 10k+ concurrent users
- ✅ **Multi-Tenancy:** Strong isolation with schema-per-tenant
- ✅ **Multi-Cloud:** AWS, GCP, Azure, On-Premise
- ✅ **Real-Time:** WebSocket updates, sub-second metrics
- ✅ **Event-Driven:** NATS for async, decoupled services
- ✅ **GitOps:** Git as source of truth for configs
- ✅ **AI-Powered:** LLM assistant for operations
- ✅ **Observable:** Full tracing, metrics, logging
- ✅ **Secure:** mTLS, JWT, RBAC, audit logging
- ✅ **Resilient:** Auto-scaling, failover, circuit breakers

**Technology Highlights:**
- Go for all backend services & infrastructure (Pulumi)
- NATS JetStream for message bus
- PostgreSQL (schema-per-tenant) + InfluxDB + Redis
- Kubernetes + Istio for orchestration
- Keycloak for authentication
- React for UI

This architecture is ready for implementation. Let's build it! 🚀

---

**Document Version:** 1.0  
**Author:** Polaris Architecture Team  
**Last Updated:** January 13, 2026  
**Status:** Approved for Implementation
