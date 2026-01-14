/**
 * Custom React Query hooks for data fetching
 * Example hooks for common API operations
 */

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiClient } from '@/lib/api-client'

// Example: Router data types (update with your actual types)
interface Router {
    id: string
    hostname: string
    ip: string
    region: string
    status: 'online' | 'degraded' | 'offline'
    load: number
    uptime: string
    version: string
}

// Query Keys - centralized for cache invalidation
export const queryKeys = {
    routers: ['routers'] as const,
    router: (id: string) => ['router', id] as const,
    metrics: (range: string) => ['metrics', range] as const,
    tenants: ['tenants'] as const,
    alerts: ['alerts'] as const,
}

// ============================================================================
// ROUTERS
// ============================================================================

export function useRouters() {
    return useQuery({
        queryKey: queryKeys.routers,
        queryFn: async () => {
            // TODO: Replace with actual API endpoint when backend is ready
            // return apiClient.get<Router[]>('/api/v1/routers')

            // Mock data for now
            return Promise.resolve([
                {
                    id: '1',
                    hostname: 'router-us-west-1',
                    ip: '10.0.1.1',
                    region: 'us-west-1',
                    status: 'online' as const,
                    load: 45,
                    uptime: '99.9%',
                    version: '1.2.3',
                },
            ])
        },
    })
}

export function useRouter(id: string) {
    return useQuery({
        queryKey: queryKeys.router(id),
        queryFn: async () => {
            // TODO: Replace with actual API endpoint
            // return apiClient.get<Router>(`/api/v1/routers/${id}`)

            return Promise.resolve({
                id,
                hostname: `router-${id}`,
                ip: '10.0.1.1',
                region: 'us-west-1',
                status: 'online' as const,
                load: 45,
                uptime: '99.9%',
                version: '1.2.3',
            })
        },
        enabled: !!id, // Only run if id is provided
    })
}

// ============================================================================
// MUTATIONS
// ============================================================================

export function useCreateRouter() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: async (data: Partial<Router>) => {
            // TODO: Replace with actual API endpoint
            // return apiClient.post<Router>('/api/v1/routers', data)
            return Promise.resolve({ ...data, id: String(Date.now()) } as Router)
        },
        onSuccess: () => {
            // Invalidate and refetch routers list
            queryClient.invalidateQueries({ queryKey: queryKeys.routers })
        },
    })
}

export function useUpdateRouter() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: async ({ id, data }: { id: string; data: Partial<Router> }) => {
            // TODO: Replace with actual API endpoint
            // return apiClient.patch<Router>(`/api/v1/routers/${id}`, data)
            return Promise.resolve({ id, ...data } as Router)
        },
        onSuccess: (_, variables) => {
            // Invalidate specific router and routers list
            queryClient.invalidateQueries({ queryKey: queryKeys.router(variables.id) })
            queryClient.invalidateQueries({ queryKey: queryKeys.routers })
        },
    })
}

export function useDeleteRouter() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: async (id: string) => {
            // TODO: Replace with actual API endpoint
            // return apiClient.delete(`/api/v1/routers/${id}`)
            return Promise.resolve()
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: queryKeys.routers })
        },
    })
}

// ============================================================================
// METRICS
// ============================================================================

interface MetricsData {
    timestamp: string
    latency: number
    packetLoss: number
    bandwidth: number
}

export function useMetrics(timeRange: string = '2h') {
    return useQuery({
        queryKey: queryKeys.metrics(timeRange),
        queryFn: async () => {
            // TODO: Replace with actual API endpoint
            // return apiClient.get<MetricsData[]>('/api/v1/metrics', { range: timeRange })

            // Mock data
            return Promise.resolve([
                {
                    timestamp: new Date().toISOString(),
                    latency: 24,
                    packetLoss: 0.001,
                    bandwidth: 1000,
                },
            ])
        },
        refetchInterval: 30000, // Refetch every 30 seconds for real-time data
    })
}
