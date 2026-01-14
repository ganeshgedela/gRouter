import { apiClient } from '@/lib/api-client'
import type {
    Router,
    RouterDetail,
    RoutersResponse,
    RoutersListParams,
    CreateRouterRequest,
    UpdateRouterRequest,
} from '@/types/api'

export const routersApi = {
    // List all routers
    list: async (params?: RoutersListParams) => {
        return apiClient.get<RoutersResponse>('/routers', params)
    },

    // Get router by ID
    get: async (id: string) => {
        return apiClient.get<RouterDetail>(`/routers/${id}`)
    },

    // Create new router
    create: async (data: CreateRouterRequest) => {
        return apiClient.post<{ id: string; status: string; message: string }>(
            '/routers',
            data
        )
    },

    // Update router
    update: async (id: string, data: UpdateRouterRequest) => {
        return apiClient.patch<RouterDetail>(`/routers/${id}`, data)
    },

    // Delete router
    delete: async (id: string) => {
        return apiClient.delete<{ message: string }>(`/routers/${id}`)
    },
}
