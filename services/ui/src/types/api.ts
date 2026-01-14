// Core API Types

export interface PaginationParams {
    page?: number
    limit?: number
}

export interface PaginationResponse {
    page: number
    limit: number
    total: number
    totalPages: number
    hasNext: boolean
    hasPrev: boolean
}

export interface ApiError {
    code: string
    message: string
    details?: Record<string, unknown>
}

// Authentication
export interface LoginRequest {
    email: string
    password: string
}

export interface AuthResponse {
    token: string
    user: User
    expiresAt: string
}

export interface User {
    id: string
    email: string
    name: string
    role: 'admin' | 'user' | 'viewer'
}

// Routers
export interface Router {
    id: string
    hostname: string
    ip: string
    region: string
    status: 'online' | 'degraded' | 'offline'
    load: number
    uptime: string
    version: string
    lastSeen: string
    metadata?: {
        datacenter?: string
        rack?: string
        [key: string]: unknown
    }
}

export interface RouterDetail extends Router {
    hardware: {
        model: string
        cpu: string
        cores: number
        memory: string
        storage: string
    }
    interfaces: NetworkInterface[]
    bgpPeers: number
    ospfNeighbors: number
    routes: {
        ipv4: number
        ipv6: number
    }
}

export interface NetworkInterface {
    name: string
    ip: string
    mac: string
    status: 'up' | 'down'
    speed: string
}

export interface CreateRouterRequest {
    hostname: string
    region: string
    datacenter: string
    config: {
        version: string
        bgpAsn: number
    }
}

export interface UpdateRouterRequest {
    config?: {
        bgpAsn?: number
        [key: string]: unknown
    }
}

export interface RoutersResponse {
    routers: Router[]
    pagination: PaginationResponse
}

export interface RoutersListParams extends PaginationParams {
    region?: string
    status?: 'online' | 'degraded' | 'offline'
}

// Tenants
export interface Tenant {
    id: string
    name: string
    domain: string
    status: 'active' | 'suspended'
    plan: 'free' | 'premium' | 'enterprise'
    routerCount: number
    bandwidth: string
    createdAt: string
}

export interface TenantDetail extends Tenant {
    routers: string[]
    contacts: TenantContact[]
    billing: {
        monthlySpend: number
        currency: string
    }
}

export interface TenantContact {
    name: string
    email: string
    role: string
}

export interface CreateTenantRequest {
    name: string
    domain: string
    plan: 'free' | 'premium' | 'enterprise'
    contacts: TenantContact[]
}

export interface TenantsListParams {
    search?: string
    status?: 'active' | 'suspended'
}

// Metrics
export interface MetricDataPoint {
    timestamp: string
    value: number
}

export interface RouterMetricsResponse {
    routerId: string
    metrics: {
        cpu?: MetricDataPoint[]
        memory?: MetricDataPoint[]
        network_in?: MetricDataPoint[]
        network_out?: MetricDataPoint[]
        latency?: MetricDataPoint[]
        packet_loss?: MetricDataPoint[]
    }
}

export interface MetricsQueryParams {
    from: string // ISO 8601
    to: string // ISO 8601
    metrics?: string[] // ['cpu', 'memory', 'network_in']
}

export interface AggregateMetrics {
    avgLatency: number
    totalBandwidth: number
    packetLoss: number
    activeRoutes: number
}

// Alerts
export interface Alert {
    id: string
    severity: 'critical' | 'warning' | 'info'
    title: string
    message: string
    routerId: string
    status: 'active' | 'acknowledged' | 'resolved'
    createdAt: string
    acknowledgedAt?: string
    resolvedAt?: string
}

export interface AlertsListParams {
    severity?: 'critical' | 'warning' | 'info'
    status?: 'active' | 'acknowledged' | 'resolved'
    routerId?: string
}

// Deployments
export interface Deployment {
    id: string
    version: string
    status: 'pending' | 'running' | 'success' | 'failed' | 'rollback'
    routerId: string
    user: string
    startedAt: string
    completedAt?: string
    duration?: string
    changes: string[]
}

export interface CreateDeploymentRequest {
    routerId: string
    version: string
    config?: {
        autoRollback?: boolean
        drainTraffic?: boolean
    }
}

// Workflows
export interface Workflow {
    id: string
    name: string
    category: string
    executions: number
    avgDuration: string
    successRate: number
}

export interface WorkflowExecution {
    id: string
    workflowId: string
    routerId: string
    status: 'pending' | 'running' | 'success' | 'failed'
    startedAt: string
    completedAt?: string
    user: string
}

export interface ExecuteWorkflowRequest {
    routerId: string
    parameters?: Record<string, unknown>
}

// WebSocket Events
export interface WebSocketMessage {
    event: 'metrics.update' | 'alert.created' | 'deployment.status'
    data: unknown
}

export interface MetricsUpdateEvent {
    routerId: string
    metrics: {
        cpu: number
        memory: number
        timestamp: string
    }
}

export interface AlertCreatedEvent {
    id: string
    severity: 'critical' | 'warning' | 'info'
    routerId: string
    message: string
}

export interface DeploymentStatusEvent {
    id: string
    status: 'success' | 'failed'
    routerId: string
}
