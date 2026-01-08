# Pulumi Deployment Sequence Diagrams

## 1. Initial Infrastructure Provisioning

```mermaid
sequenceDiagram
    participant Dev as Developer
    participant Git as Git Repository
    participant Pulumi as Pulumi CLI
    participant Cloud as Cloud Provider
    participant K8s as Kubernetes
    participant DB as Database
    
    Dev->>Git: git push (Pulumi code)
    Dev->>Pulumi: pulumi up --stack prod
    
    Pulumi->>Pulumi: Load configuration
    Pulumi->>Pulumi: Generate execution plan
    Pulumi->>Dev: Display preview
    
    Dev->>Pulumi: Approve changes
    
    Note over Pulumi,Cloud: Provision Infrastructure
    
    Pulumi->>Cloud: Create VPC/Network
    Cloud-->>Pulumi: VPC created
    
    Pulumi->>Cloud: Create security groups
    Cloud-->>Pulumi: Security groups created
    
    Pulumi->>Cloud: Create Kubernetes cluster
    Cloud-->>Pulumi: Cluster creating...
    Cloud-->>Pulumi: Cluster ready
    
    Pulumi->>Cloud: Create managed database
    Cloud-->>Pulumi: Database provisioning...
    Cloud-->>Pulumi: Database ready
    
    Pulumi->>Cloud: Create storage buckets
    Cloud-->>Pulumi: Storage created
    
    Pulumi->>K8s: Deploy core services (NATS, monitoring)
    K8s-->>Pulumi: Services deployed
    
    Pulumi->>Pulumi: Save state
    Pulumi->>Dev: Deployment complete
    
    Dev->>K8s: kubectl get nodes
    K8s-->>Dev: Cluster status
```

## 2. Service Deployment Flow

```mermaid
sequenceDiagram
    participant Dev as Developer
    participant CI as CI/CD (GitHub Actions)
    participant Registry as Container Registry
    participant Pulumi as Pulumi
    participant K8s as Kubernetes Cluster
    participant Health as Health Check
    participant Monitor as Monitoring
    
    Dev->>CI: git push (service code)
    
    CI->>CI: Run tests
    CI->>CI: Build Docker image
    CI->>CI: Security scan
    
    alt Security scan fails
        CI->>Dev: ❌ Deployment blocked
    else Security scan passes
        CI->>Registry: Push image:v1.2.3
        Registry-->>CI: Image pushed
        
        CI->>Pulumi: Update image tag in config
        Pulumi->>Pulumi: pulumi preview
        
        CI->>Pulumi: pulumi up (auto-approve)
        
        Pulumi->>K8s: Update deployment
        K8s->>K8s: Rolling update
        
        loop Health checks
            K8s->>Health: Check pod health
            Health-->>K8s: Health status
        end
        
        alt Deployment successful
            K8s-->>Pulumi: Deployment complete
            Pulumi-->>CI: Success
            CI->>Monitor: Send deployment event
            Monitor-->>CI: Event recorded
            CI->>Dev: ✅ Deployment successful
        else Deployment failed
            K8s-->>Pulumi: Deployment failed
            Pulumi->>K8s: Rollback to previous version
            K8s-->>Pulumi: Rollback complete
            Pulumi-->>CI: Rollback executed
            CI->>Dev: ❌ Deployment failed, rolled back
        end
    end
```

## 3. CI/CD Pipeline with Pulumi

```mermaid
sequenceDiagram
    participant Dev as Developer
    participant Git as GitHub
    participant GHA as GitHub Actions
    participant Test as Test Suite
    participant Build as Docker Build
    participant Scan as Trivy Scanner
    participant Registry as Container Registry
    participant Pulumi as Pulumi
    participant GCP as GCP (Production)
    participant AWS as AWS (Staging)
    participant Notify as Notifications
    
    Dev->>Git: Push to main branch
    Git->>GHA: Trigger workflow
    
    GHA->>Test: Run unit tests
    Test-->>GHA: Tests passed ✓
    
    GHA->>Build: Build multi-arch image
    Build-->>GHA: Image built
    
    GHA->>Scan: Scan for vulnerabilities
    Scan-->>GHA: Scan results
    
    alt Critical vulnerabilities found
        GHA->>Notify: Send alert
        GHA->>Dev: ❌ Build failed
    else No critical issues
        GHA->>Registry: Push image
        Registry-->>GHA: Image stored
        
        par Deploy to Staging
            GHA->>Pulumi: Deploy to AWS (staging)
            Pulumi->>AWS: Update EKS deployment
            AWS-->>Pulumi: Deployed
            Pulumi-->>GHA: Staging ready
            
            GHA->>Test: Run E2E tests (staging)
            Test-->>GHA: E2E passed ✓
        end
        
        GHA->>Dev: Request production approval
        Dev->>GHA: Approve production deploy
        
        GHA->>Pulumi: Deploy to GCP (production)
        Pulumi->>GCP: Update GKE deployment
        Pulumi->>GCP: Deploy with canary (10%)
        GCP-->>Pulumi: Canary deployed
        
        Pulumi->>GCP: Monitor canary metrics
        
        alt Canary healthy
            Pulumi->>GCP: Scale canary to 50%
            Pulumi->>GCP: Scale canary to 100%
            GCP-->>Pulumi: Full deployment complete
            Pulumi-->>GHA: Production deployed ✓
            GHA->>Notify: Send success notification
        else Canary unhealthy
            Pulumi->>GCP: Rollback canary
            GCP-->>Pulumi: Rollback complete
            Pulumi-->>GHA: Deployment aborted
            GHA->>Notify: Send failure alert
        end
    end
```

## 4. Disaster Recovery Failover

```mermaid
sequenceDiagram
    participant Monitor as Monitoring
    participant Primary as Primary Region
    participant DNS as Global DNS
    participant DR as DR Region
    participant DB_Pri as Primary DB
    participant DB_DR as DR DB
    participant Alert as Alert Manager
    participant Ops as Operations Team
    
    loop Health monitoring
        Monitor->>Primary: Health check
        Primary-->>Monitor: Healthy
    end
    
    Note over Primary: Region failure occurs
    
    Monitor->>Primary: Health check
    Primary--xMonitor: No response (timeout)
    
    Monitor->>Monitor: Retry (3 attempts)
    Monitor->>Alert: Trigger critical alert
    
    Alert->>Ops: PagerDuty alert
    Alert->>Ops: Slack notification
    
    Ops->>Monitor: Acknowledge alert
    Ops->>DNS: Initiate failover procedure
    
    DNS->>DB_DR: Promote replica to primary
    DB_DR->>DB_DR: Stop replication
    DB_DR->>DB_DR: Enable writes
    DB_DR-->>DNS: Promotion complete
    
    DNS->>DNS: Update DNS records
    DNS->>DNS: Point to DR region
    
    DNS->>DR: Activate DR services
    DR->>DR: Scale up replicas
    DR-->>DNS: Services active
    
    Note over DNS: Traffic now routed to DR
    
    DNS->>Monitor: Update monitoring targets
    Monitor->>DR: Start monitoring DR region
    
    loop Wait for primary recovery
        Ops->>Primary: Check region status
        Primary-->>Ops: Still down
    end
    
    Note over Primary: Region recovered
    
    Primary-->>Ops: Region online
    
    Ops->>DB_DR: Setup replication back to primary
    DB_DR->>DB_Pri: Sync data
    DB_Pri-->>DB_DR: Replication established
    
    Ops->>DNS: Plan failback
    Ops->>DNS: Switch traffic back (gradual)
    
    DNS->>DNS: Update records (canary 10%)
    DNS->>DNS: Increase to 50%
    DNS->>DNS: Full cutover to primary
    
    DNS->>DR: Scale down DR services
    DR->>DR: Return to standby mode
    
    Ops->>Alert: Incident resolved
    Alert->>Ops: Incident closed
```

## 5. Auto-Scaling Event

```mermaid
sequenceDiagram
    participant Traffic as User Traffic
    participant LB as Load Balancer
    participant Metrics as Metrics Server
    participant HPA as Horizontal Pod Autoscaler
    participant K8s as Kubernetes
    participant Cloud as Cloud Provider
    participant Pods as Service Pods
    
    Traffic->>LB: Increased requests
    LB->>Pods: Forward traffic
    
    Pods->>Metrics: Report CPU usage (85%)
    Metrics->>HPA: CPU threshold exceeded
    
    HPA->>HPA: Calculate desired replicas
    Note over HPA: Current: 3 pods<br/>Target: 6 pods<br/>Metric: 85% CPU
    
    HPA->>K8s: Scale deployment to 6 replicas
    
    K8s->>K8s: Check cluster capacity
    
    alt Sufficient capacity
        K8s->>Pods: Create 3 new pods
        Pods->>Pods: Initialize containers
        Pods->>Pods: Run health checks
        Pods-->>K8s: Pods ready
        K8s->>LB: Add new pods to pool
        LB->>Pods: Distribute traffic (6 pods)
        
        Note over Pods: Load balanced across 6 pods
        Pods->>Metrics: Report CPU usage (45%)
        Metrics->>HPA: Within target range
    else Insufficient capacity
        K8s->>Cloud: Trigger cluster autoscaler
        Cloud->>Cloud: Provision new node
        Cloud-->>K8s: Node ready
        K8s->>Pods: Create pods on new node
        Pods-->>K8s: Pods ready
        K8s->>LB: Add new pods
    end
    
    Note over Traffic: Traffic decreases
    
    Traffic->>LB: Reduced requests
    Pods->>Metrics: Report CPU usage (20%)
    Metrics->>HPA: Below target (scale down)
    
    HPA->>K8s: Scale deployment to 3 replicas
    K8s->>Pods: Gracefully terminate 3 pods
    Pods->>Pods: Drain connections
    Pods->>K8s: Shutdown complete
    K8s->>LB: Remove pods from pool
    
    K8s->>Cloud: Node underutilized
    Cloud->>Cloud: Scale down cluster (after delay)
```

## 6. Secret Rotation Flow

```mermaid
sequenceDiagram
    participant Scheduler as Scheduled Job
    participant Vault as Secret Store
    participant Pulumi as Pulumi
    participant K8s as Kubernetes
    participant Pods as Service Pods
    participant DB as Database
    
    Note over Scheduler: Monthly rotation trigger
    
    Scheduler->>Vault: Initiate secret rotation
    Vault->>Vault: Generate new credentials
    Vault->>DB: Update database user password
    DB-->>Vault: Password updated
    
    Vault->>Vault: Store new secret (version 2)
    Vault->>Vault: Keep old secret (version 1)
    
    Vault->>Pulumi: Notify of new secret version
    Pulumi->>K8s: Update secret (version 2)
    K8s-->>Pulumi: Secret updated
    
    Pulumi->>K8s: Rolling restart pods
    
    loop For each pod
        K8s->>Pods: Terminate pod (old secret)
        Pods->>Pods: Graceful shutdown
        K8s->>Pods: Create new pod
        Pods->>K8s: Load secret (version 2)
        Pods->>DB: Connect with new password
        DB-->>Pods: Connection successful
        Pods-->>K8s: Pod ready
    end
    
    Note over Vault: Grace period for old secret
    
    Vault->>Vault: Monitor for old secret usage
    
    alt All pods using new secret
        Vault->>Vault: Revoke old secret (version 1)
        Vault->>Scheduler: Rotation complete
    else Some pods still using old secret
        Vault->>Scheduler: Alert - manual intervention needed
    end
```

## 7. Database Migration

```mermaid
sequenceDiagram
    participant Dev as Developer
    participant CI as CI/CD
    participant Backup as Backup Service
    participant DB as Database
    participant Migration as Migration Tool
    participant App as Application
    participant Monitor as Monitoring
    
    Dev->>CI: Push migration script
    CI->>CI: Run migration tests (local DB)
    CI-->>Dev: Tests passed
    
    Dev->>CI: Approve production migration
    
    CI->>App: Enable maintenance mode
    App-->>CI: Maintenance mode active
    
    CI->>Backup: Create database backup
    Backup->>DB: Backup current state
    DB-->>Backup: Backup complete
    Backup-->>CI: Backup stored
    
    CI->>Migration: Run migration (dry-run)
    Migration->>DB: Validate migration
    DB-->>Migration: Validation successful
    
    CI->>Migration: Execute migration
    Migration->>DB: Begin transaction
    Migration->>DB: Run DDL statements
    Migration->>DB: Update schema version
    
    alt Migration successful
        DB-->>Migration: Transaction committed
        Migration-->>CI: Migration complete
        
        CI->>App: Deploy new app version
        App->>DB: Test new schema
        DB-->>App: Schema compatible
        
        CI->>App: Disable maintenance mode
        App-->>CI: Service restored
        
        CI->>Monitor: Record migration success
        Monitor-->>CI: Logged
        CI->>Dev: ✅ Migration successful
    else Migration failed
        Migration->>DB: Rollback transaction
        DB-->>Migration: Transaction rolled back
        Migration-->>CI: Migration failed
        
        CI->>Backup: Restore from backup
        Backup->>DB: Restore previous state
        DB-->>Backup: Restore complete
        
        CI->>App: Disable maintenance mode
        CI->>Monitor: Record migration failure
        CI->>Dev: ❌ Migration failed, rolled back
    end
```

## 8. Monitoring Alert Flow

```mermaid
sequenceDiagram
    participant Service as Microservice
    participant Metrics as Metrics Collector
    participant Prometheus as Prometheus
    participant AlertMgr as Alert Manager
    participant OnCall as On-Call Engineer
    participant Runbook as Runbook System
    participant K8s as Kubernetes
    
    loop Every 15s
        Service->>Metrics: Export metrics
        Metrics->>Prometheus: Scrape metrics
    end
    
    Note over Service: Error rate increases
    
    Service->>Metrics: error_rate=15% (threshold: 5%)
    Metrics->>Prometheus: Store metric
    
    Prometheus->>Prometheus: Evaluate alert rules
    Prometheus->>Prometheus: Alert: HighErrorRate
    
    Prometheus->>AlertMgr: Fire alert
    AlertMgr->>AlertMgr: Group & deduplicate
    AlertMgr->>AlertMgr: Wait for group_wait
    
    AlertMgr->>OnCall: Send PagerDuty alert
    AlertMgr->>OnCall: Send Slack notification
    
    OnCall->>AlertMgr: Acknowledge alert
    AlertMgr->>AlertMgr: Stop repeat notifications
    
    OnCall->>Runbook: Query troubleshooting steps
    Runbook-->>OnCall: Return runbook
    
    OnCall->>Prometheus: Check related metrics
    Prometheus-->>OnCall: CPU, memory, latency data
    
    OnCall->>K8s: Check pod logs
    K8s-->>OnCall: Error logs
    
    OnCall->>K8s: Identify root cause (DB connection issue)
    
    OnCall->>K8s: Restart affected pods
    K8s->>Service: Rolling restart
    Service->>Service: Reconnect to database
    Service-->>K8s: Pods healthy
    
    Service->>Metrics: error_rate=0.5%
    Metrics->>Prometheus: Store metric
    
    Prometheus->>Prometheus: Evaluate rules
    Prometheus->>AlertMgr: Resolve alert
    
    AlertMgr->>OnCall: Send resolution notification
    OnCall->>AlertMgr: Close incident
    AlertMgr->>Runbook: Update incident log
```

## 9. Multi-Cloud Sync

```mermaid
sequenceDiagram
    participant GCP as GCP (Primary)
    participant Sync as Sync Service
    participant AWS as AWS (Secondary)
    participant Azure as Azure (Tertiary)
    participant State as State Store
    
    Note over GCP: Service deployed
    
    GCP->>Sync: Deployment event
    Sync->>State: Record deployment (GCP)
    State-->>Sync: State saved
    
    Sync->>AWS: Check sync policy
    
    alt Auto-sync enabled
        Sync->>AWS: Trigger deployment
        AWS->>AWS: Deploy service
        AWS-->>Sync: Deployment complete
        Sync->>State: Update state (AWS)
    else Manual sync
        Sync->>State: Mark for manual sync
    end
    
    Sync->>Azure: Check sync policy
    
    alt Selective sync (prod only)
        Note over Sync: Skip Azure (non-prod env)
    else Full sync
        Sync->>Azure: Trigger deployment
        Azure->>Azure: Deploy service
        Azure-->>Sync: Deployment complete
        Sync->>State: Update state (Azure)
    end
    
    Sync->>Sync: Verify consistency
    Sync->>State: All regions synced
```

## Flow Descriptions

### Key Processes

1. **Infrastructure Provisioning**: Complete setup from VPC to K8s cluster
2. **Service Deployment**: CI/CD pipeline with health checks and rollback
3. **Disaster Recovery**: Automated failover with data replication
4. **Auto-Scaling**: Dynamic scaling based on metrics
5. **Secret Rotation**: Zero-downtime credential updates
6. **Database Migration**: Safe schema changes with rollback
7. **Monitoring**: Alert-driven incident response
8. **Multi-Cloud Sync**: Cross-region deployment coordination

### Common Patterns

- **Retry Logic**: All API calls include retries with exponential backoff
- **Health Checks**: Continuous validation before and after changes
- **Rollback**: Automatic rollback on failure detection
- **State Management**: Consistent state tracking across operations
- **Notifications**: Multi-channel alerts for critical events
