# Pulumi IaaS Architecture Design

## System Architecture Overview

```mermaid
graph TB
    subgraph "Multi-Cloud Deployment"
        subgraph "GCP Region"
            GKE[GKE Cluster]
            GCS[Cloud Storage]
            GCSQL[Cloud SQL]
            GLB[Global Load Balancer]
        end
        
        subgraph "AWS Region"
            EKS[EKS Cluster]
            S3[S3 Storage]
            RDS[RDS Database]
            ALB[Application LB]
        end
        
        subgraph "Azure Region"
            AKS[AKS Cluster]
            BLOB[Blob Storage]
            ADB[Azure Database]
            AzLB[Azure LB]
        end
        
        subgraph "On-Premises"
            K8S[Kubernetes]
            NFS[NFS Storage]
            PG[PostgreSQL]
            NGINX[NGINX]
        end
    end
    
    subgraph "Global Services"
        DNS[Global DNS]
        CDN[CDN]
        WAF[Web Application Firewall]
    end
    
    DNS --> GLB
    DNS --> ALB
    DNS --> AzLB
    DNS --> NGINX
    
    CDN --> GLB
    CDN --> ALB
```

## Cloud-Specific Architectures

### GCP Architecture

```mermaid
graph TB
    subgraph "Internet"
        Users[Users]
    end
    
    subgraph "GCP Project"
        subgraph "Networking"
            VPC[VPC Network]
            CloudNAT[Cloud NAT]
            CloudDNS[Cloud DNS]
        end
        
        subgraph "Compute - GKE"
            Ingress[Ingress Controller]
            
            subgraph "Workload Cluster"
                NS1[nats-service]
                NS2[rest-service]
                NS3[grpc-service]
                NS4[hybrid-service]
                NS5[messaging-rpc]
            end
            
            NATS[NATS Cluster]
        end
        
        subgraph "Data Layer"
            CloudSQL[Cloud SQL<br/>PostgreSQL]
            Memorystore[Memorystore<br/>Redis]
            GCS[Cloud Storage]
        end
        
        subgraph "Observability"
            CloudMonitoring[Cloud Monitoring]
            CloudLogging[Cloud Logging]
            CloudTrace[Cloud Trace]
        end
        
        subgraph "Security"
            KMS[Cloud KMS]
            SecretManager[Secret Manager]
            IAM[IAM & Service Accounts]
        end
    end
    
    Users --> Ingress
    Ingress --> NS1
    Ingress --> NS2
    Ingress --> NS3
    Ingress --> NS4
    Ingress --> NS5
    
    NS1 --> NATS
    NS2 --> CloudSQL
    NS3 --> CloudSQL
    NS4 --> NATS
    NS4 --> CloudSQL
    NS5 --> NATS
    NS5 --> CloudSQL
    
    NS2 --> Memorystore
    NS3 --> Memorystore
```

### AWS Architecture

```mermaid
graph TB
    subgraph "AWS Account"
        subgraph "VPC"
            subgraph "Public Subnets"
                ALB[Application LB]
                NAT[NAT Gateway]
            end
            
            subgraph "Private Subnets - EKS"
                subgraph "EKS Cluster"
                    Pod1[nats-service]
                    Pod2[rest-service]
                    Pod3[grpc-service]
                    Pod4[hybrid-service]
                    Pod5[messaging-rpc]
                end
                
                NATS[NATS StatefulSet]
            end
            
            subgraph "Data Subnets"
                RDS[(RDS PostgreSQL<br/>Multi-AZ)]
                ElastiCache[(ElastiCache<br/>Redis)]
            end
        end
        
        subgraph "Storage"
            S3[S3 Buckets]
            EBS[EBS Volumes]
        end
        
        subgraph "Observability"
            CloudWatch[CloudWatch]
            XRay[X-Ray]
        end
        
        subgraph "Security"
            SecretsManager[Secrets Manager]
            KMS_AWS[AWS KMS]
            SecurityGroups[Security Groups]
        end
    end
    
    Internet[Internet] --> Route53[Route53]
    Route53 --> ALB
    ALB --> Pod1
    ALB --> Pod2
    ALB --> Pod3
    ALB --> Pod4
    ALB --> Pod5
    
    Pod1 --> NATS
    Pod2 --> RDS
    Pod2 --> ElastiCache
    Pod4 --> NATS
    Pod4 --> RDS
```

### Azure Architecture

```mermaid
graph TB
    subgraph "Azure Subscription"
        subgraph "Resource Group"
            subgraph "Virtual Network"
                subgraph "AKS Subnet"
                    AKS[AKS Cluster]
                    
                    subgraph "Services"
                        S1[nats-service]
                        S2[rest-service]
                        S3[grpc-service]
                        S4[hybrid-service]
                        S5[messaging-rpc]
                    end
                    
                    NATS_AZ[NATS Cluster]
                end
                
                subgraph "Data Subnet"
                    PG_AZ[(Azure Database<br/>PostgreSQL)]
                    Redis_AZ[(Azure Cache<br/>Redis)]
                end
            end
            
            LB[Azure Load Balancer]
            AppGW[Application Gateway]
        end
        
        subgraph "Storage"
            BlobStorage[Blob Storage]
            FileStorage[Azure Files]
        end
        
        subgraph "Monitoring"
            AppInsights[Application Insights]
            LogAnalytics[Log Analytics]
        end
        
        subgraph "Security"
            KeyVault[Key Vault]
            ManagedIdentity[Managed Identity]
        end
    end
    
    Internet_AZ[Internet] --> AppGW
    AppGW --> LB
    LB --> S1
    LB --> S2
    LB --> S3
    LB --> S4
    LB --> S5
```

## Network Topology

```mermaid
graph TB
    subgraph "Network Layers"
        subgraph "Edge Layer"
            CDN[CDN/Edge Cache]
            WAF[WAF]
            DDoS[DDoS Protection]
        end
        
        subgraph "Ingress Layer"
            LB[Load Balancer]
            Ingress[Ingress Controller]
            APIGateway[API Gateway]
        end
        
        subgraph "Service Mesh Layer"
            Envoy[Envoy Proxy]
            Istio[Istio Control Plane]
        end
        
        subgraph "Application Layer"
            Services[Microservices]
        end
        
        subgraph "Data Layer"
            Cache[Redis/Memcached]
            DB[(Database)]
            MessageBus[NATS]
        end
    end
    
    CDN --> WAF
    WAF --> DDoS
    DDoS --> LB
    LB --> Ingress
    Ingress --> APIGateway
    APIGateway --> Envoy
    Envoy --> Services
    Services --> Cache
    Services --> DB
    Services --> MessageBus
```

## Service Mesh Design

```mermaid
graph LR
    subgraph "Data Plane"
        subgraph "Pod 1"
            App1[Application]
            Sidecar1[Envoy Sidecar]
        end
        
        subgraph "Pod 2"
            App2[Application]
            Sidecar2[Envoy Sidecar]
        end
        
        subgraph "Pod 3"
            App3[Application]
            Sidecar3[Envoy Sidecar]
        end
    end
    
    subgraph "Control Plane"
        Pilot[Pilot<br/>Traffic Management]
        Citadel[Citadel<br/>Certificate Authority]
        Galley[Galley<br/>Configuration]
        Mixer[Mixer<br/>Telemetry]
    end
    
    App1 --> Sidecar1
    Sidecar1 <--> Sidecar2
    Sidecar2 --> App2
    Sidecar2 <--> Sidecar3
    Sidecar3 --> App3
    
    Pilot -.-> Sidecar1
    Pilot -.-> Sidecar2
    Pilot -.-> Sidecar3
    
    Citadel -.-> Sidecar1
    Citadel -.-> Sidecar2
    Citadel -.-> Sidecar3
```

## Data Flow Architecture

```mermaid
graph TD
    Client[Client Request]
    
    subgraph "Request Path"
        Client --> LB[Load Balancer]
        LB --> Ingress[Ingress]
        Ingress --> Service[Service Pod]
        
        Service --> Cache{Cache Hit?}
        Cache -->|Yes| Return1[Return Cached]
        Cache -->|No| DB[(Database)]
        DB --> Cache_Store[Update Cache]
        Cache_Store --> Return2[Return Data]
    end
    
    subgraph "Async Path"
        Service --> NATS[NATS Publish]
        NATS --> Worker1[Worker 1]
        NATS --> Worker2[Worker 2]
        Worker1 --> ProcessingDB[(Processing DB)]
        Worker2 --> ProcessingDB
    end
    
    subgraph "Event Flow"
        ProcessingDB --> EventStream[Event Stream]
        EventStream --> Analytics[Analytics Service]
        EventStream --> Audit[Audit Service]
    end
```

## Deployment Architecture

```mermaid
graph TB
    subgraph "CI/CD Pipeline"
        Git[Git Repository]
        CI[CI System<br/>GitHub Actions]
        Registry[Container Registry]
        
        Git --> CI
        CI --> Build[Build & Test]
        Build --> Scan[Security Scan]
        Scan --> Registry
    end
    
    subgraph "Pulumi Deployment"
        Registry --> Pulumi[Pulumi]
        Pulumi --> Preview[pulumi preview]
        Preview --> Approve{Manual<br/>Approval}
        Approve -->|Yes| Deploy[pulumi up]
        Approve -->|No| Reject[Reject]
        
        Deploy --> K8s[Kubernetes]
    end
    
    subgraph "Validation"
        K8s --> HealthCheck[Health Checks]
        HealthCheck --> E2E[E2E Tests]
        E2E --> Monitoring[Monitoring]
    end
    
    Monitoring --> Alert{Issues?}
    Alert -->|Yes| Rollback[Automatic Rollback]
    Alert -->|No| Success[Deployment Success]
```

## Security Architecture

```mermaid
graph TB
    subgraph "Security Layers"
        subgraph "Network Security"
            Firewall[Cloud Firewall]
            NSG[Network Security Groups]
            PrivateLink[Private Link/Endpoint]
        end
        
        subgraph "Identity & Access"
            IAM_SEC[IAM/RBAC]
            SA[Service Accounts]
            OIDC[OIDC Provider]
        end
        
        subgraph "Data Security"
            Encryption[Encryption at Rest]
            TLS[TLS in Transit]
            KMS_SEC[Key Management]
        end
        
        subgraph "Application Security"
            Secrets[Secret Management]
            ImageScan[Image Scanning]
            PodSecurity[Pod Security Policies]
        end
        
        subgraph "Monitoring & Audit"
            AuditLog[Audit Logs]
            IDS[Intrusion Detection]
            SIEM[SIEM Integration]
        end
    end
```

## High Availability Design

```mermaid
graph TB
    subgraph "Multi-Region Setup"
        subgraph "Region 1 - Primary"
            App1[Applications]
            DB1[(Primary DB)]
            NATS1[NATS Cluster]
        end
        
        subgraph "Region 2 - Secondary"
            App2[Applications]
            DB2[(Replica DB)]
            NATS2[NATS Cluster]
        end
        
        subgraph "Region 3 - DR"
            App3[Applications - Standby]
            DB3[(Replica DB)]
            NATS3[NATS Cluster]
        end
    end
    
    GLB[Global Load Balancer]
    
    GLB --> App1
    GLB -.->|Failover| App2
    GLB -.->|DR| App3
    
    DB1 -->|Replication| DB2
    DB1 -->|Replication| DB3
    
    NATS1 <-->|Cluster| NATS2
    NATS2 <-->|Cluster| NATS3
```

## Observability Stack

```mermaid
graph TB
    subgraph "Services"
        S1[Service 1]
        S2[Service 2]
        S3[Service 3]
    end
    
    subgraph "Metrics Collection"
        S1 --> |Prometheus| Metrics[Metrics Store]
        S2 --> |Prometheus| Metrics
        S3 --> |Prometheus| Metrics
        
        Metrics --> Grafana[Grafana Dashboard]
    end
    
    subgraph "Logging"
        S1 --> Fluentd[Fluentd/Fluent Bit]
        S2 --> Fluentd
        S3 --> Fluentd
        
        Fluentd --> Loki[Loki/Elasticsearch]
        Loki --> Grafana
    end
    
    subgraph "Tracing"
        S1 --> Jaeger[Jaeger Collector]
        S2 --> Jaeger
        S3 --> Jaeger
        
        Jaeger --> Tempo[Tempo/Storage]
        Tempo --> Grafana
    end
    
    subgraph "Alerting"
        Grafana --> AlertManager[Alert Manager]
        AlertManager --> Slack[Slack]
        AlertManager --> PagerDuty[PagerDuty]
        AlertManager --> Email[Email]
    end
```

## Disaster Recovery Design

```mermaid
graph TB
    subgraph "Normal Operations"
        Primary[Primary Region]
        App[Applications]
        DB_Primary[(Primary Database)]
        
        Primary --> App
        App --> DB_Primary
    end
    
    subgraph "Backup & Replication"
        DB_Primary -->|Continuous<br/>Replication| DB_Standby[(Standby Database)]
        DB_Primary -->|Daily Backup| Backup[Backup Storage]
        
        App -->|Config Backup| GitOps[GitOps Repository]
        App -->|State Backup| StateBackup[Pulumi State Backup]
    end
    
    subgraph "Disaster Scenario"
        Failure{Primary<br/>Region Fail}
        
        Failure -->|Detected| Failover[Automatic Failover]
        Failover --> DR_Region[DR Region]
        DR_Region --> DR_App[DR Applications]
        DR_App --> DB_Standby
        
        DB_Standby -->|Promote| DB_New_Primary[(New Primary)]
    end
    
    subgraph "Recovery"
        DB_New_Primary -->|Restore| Primary_Recovered[Recovered Primary]
        DR_App -->|Continue| Normal[Return to Normal]
    end
```

## Cost Optimization Architecture

```mermaid
graph TB
    subgraph "Workload Classification"
        Critical[Critical Workloads]
        Standard[Standard Workloads]
        Batch[Batch Processing]
        Dev[Development]
    end
    
    subgraph "Compute Strategy"
        Critical --> OnDemand[On-Demand<br/>Reserved Instances]
        Standard --> Savings[Savings Plans<br/>Committed Use]
        Batch --> Spot[Spot/Preemptible<br/>Instances]
        Dev --> Burstable[Burstable<br/>Instances]
    end
    
    subgraph "Scaling Strategy"
        OnDemand --> HPA[Horizontal Pod<br/>Autoscaler]
        Savings --> VPA[Vertical Pod<br/>Autoscaler]
        Spot --> Keda[KEDA Event-driven<br/>Autoscaling]
        Burstable --> Schedule[Scheduled<br/>Scaling]
    end
    
    subgraph "Storage Optimization"
        HPA --> HotStorage[Hot Storage<br/>SSD]
        VPA --> WarmStorage[Warm Storage<br/>Standard]
        Keda --> ColdStorage[Cold Storage<br/>Archive]
        Schedule --> Lifecycle[Lifecycle<br/>Policies]
    end
```

## Service Dependencies

```mermaid
graph LR
    subgraph "gRouter Services"
        REST[rest-service]
        GRPC[grpc-service]
        NATS_SVC[nats-service]
        HYBRID[hybrid-service]
        MSG_RPC[messaging-rpc-service]
    end
    
    subgraph "Infrastructure"
        K8S[Kubernetes]
        DB[(PostgreSQL)]
        NATS_INFRA[NATS Cluster]
        REDIS[(Redis)]
        DNS[Service Discovery]
    end
    
    subgraph "External"
        Monitoring[Monitoring]
        Logging[Logging]
        Secrets[Secret Store]
    end
    
    REST --> DB
    REST --> REDIS
    REST --> DNS
    
    GRPC --> DB
    GRPC --> DNS
    
    NATS_SVC --> NATS_INFRA
    NATS_SVC --> DNS
    
    HYBRID --> NATS_INFRA
    HYBRID --> DB
    HYBRID --> REDIS
    
    MSG_RPC --> NATS_INFRA
    MSG_RPC --> DB
    
    K8S --> Monitoring
    K8S --> Logging
    K8S --> Secrets
```

## Technology Stack

| Layer | Technology | Purpose |
|-------|-----------|---------|
| **IaC** | Pulumi (TypeScript) | Infrastructure as Code |
| **Orchestration** | Kubernetes | Container orchestration |
| **Compute** | GKE/EKS/AKS | Managed Kubernetes |
| **Database** | PostgreSQL | Primary data store |
| **Cache** | Redis | Application cache |
| **Messaging** | NATS | Async messaging |
| **Load Balancing** | Cloud LB / Ingress | Traffic distribution |
| **Service Mesh** | Istio (optional) | Service-to-service communication |
| **Monitoring** | Prometheus + Grafana | Metrics & dashboards |
| **Logging** | Loki / ELK | Log aggregation |
| **Tracing** | Jaeger / Tempo | Distributed tracing |
| **Secrets** | Cloud KMS / Vault | Secret management |
| **Storage** | S3 / GCS / Blob | Object storage |
| **CI/CD** | GitHub Actions | Automation |
| **Registry** | GCR / ECR / ACR | Container images |

## Design Principles

1. **Cloud Native**: Leverage managed services where possible
2. **High Availability**: Multi-AZ deployments, health checks, auto-recovery
3. **Scalability**: Horizontal scaling, load balancing, caching
4. **Security**: Defense in depth, least privilege, encryption
5. **Observability**: Comprehensive monitoring, logging, tracing
6. **Cost Optimization**: Right-sizing, auto-scaling, reserved capacity
7. **Disaster Recovery**: Backups, replication, failover automation
8. **DevOps**: Infrastructure as Code, CI/CD, GitOps
