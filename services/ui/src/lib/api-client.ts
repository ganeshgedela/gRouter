import type { ApiError } from '@/types/api'
import { env } from './env'

const API_BASE_URL = env.apiUrl

interface RequestConfig extends RequestInit {
    timeout?: number
    retry?: number
    retryDelay?: number
}

class ApiClient {
    private baseUrl: string
    private token: string | null = null

    constructor(baseUrl: string) {
        this.baseUrl = baseUrl
        this.loadToken()
    }

    private loadToken() {
        this.token = localStorage.getItem('auth_token')
    }

    setToken(token: string) {
        this.token = token
        localStorage.setItem('auth_token', token)
    }

    clearToken() {
        this.token = null
        localStorage.removeItem('auth_token')
    }

    private async fetchWithTimeout(
        url: string,
        config: RequestConfig = {}
    ): Promise<Response> {
        const { timeout = 30000, ...fetchConfig } = config

        const controller = new AbortController()
        const timeoutId = setTimeout(() => controller.abort(), timeout)

        try {
            const response = await fetch(url, {
                ...fetchConfig,
                signal: controller.signal,
            })
            clearTimeout(timeoutId)
            return response
        } catch (error) {
            clearTimeout(timeoutId)
            if (error instanceof Error && error.name === 'AbortError') {
                throw new Error('Request timeout')
            }
            throw error
        }
    }

    private async request<T>(
        endpoint: string,
        config: RequestConfig = {}
    ): Promise<T> {
        const { retry = 3, retryDelay = 1000, ...fetchConfig } = config
        const url = `${this.baseUrl}${endpoint}`
        const headers: Record<string, string> = {
            'Content-Type': 'application/json',
            ...(fetchConfig.headers as Record<string, string>),
        }

        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`
        }

        let lastError: Error | null = null

        // Exponential backoff retry logic
        for (let attempt = 0; attempt <= retry; attempt++) {
            try {
                const response = await this.fetchWithTimeout(url, {
                    ...fetchConfig,
                    headers,
                })

                // Handle 401 Unauthorized
                if (response.status === 401) {
                    this.clearToken()
                    window.location.href = '/login'
                    throw new Error('Unauthorized')
                }

                // Parse response
                const data = await response.json()

                // Handle non-2xx responses
                if (!response.ok) {
                    const error = data as { error: ApiError }
                    const apiError = new Error(error.error?.message || 'API request failed')

                    // Don't retry client errors (4xx)
                    if (response.status >= 400 && response.status < 500) {
                        throw apiError
                    }

                    throw apiError
                }

                return data as T
            } catch (error) {
                lastError = error as Error

                // Don't retry on client errors
                if (error instanceof Error && error.message === 'Unauthorized') {
                    throw error
                }

                // Wait before retrying (exponential backoff)
                if (attempt < retry) {
                    await new Promise((resolve) =>
                        setTimeout(resolve, retryDelay * Math.pow(2, attempt))
                    )
                }
            }
        }

        // All retries failed
        throw lastError
    }

    // HTTP Methods
    async get<T>(endpoint: string, params?: Record<string, unknown>): Promise<T> {
        const queryString = params
            ? '?' + new URLSearchParams(params as Record<string, string>).toString()
            : ''
        return this.request<T>(`${endpoint}${queryString}`, { method: 'GET' })
    }

    async post<T>(endpoint: string, data?: unknown): Promise<T> {
        return this.request<T>(endpoint, {
            method: 'POST',
            body: JSON.stringify(data),
        })
    }

    async patch<T>(endpoint: string, data?: unknown): Promise<T> {
        return this.request<T>(endpoint, {
            method: 'PATCH',
            body: JSON.stringify(data),
        })
    }

    async delete<T>(endpoint: string): Promise<T> {
        return this.request<T>(endpoint, { method: 'DELETE' })
    }
}

// Singleton instance
export const apiClient = new ApiClient(API_BASE_URL)

// Helper for retry logic
export async function withRetry<T>(
    fn: () => Promise<T>,
    maxRetries = 3,
    delay = 1000
): Promise<T> {
    let lastError: Error

    for (let i = 0; i < maxRetries; i++) {
        try {
            return await fn()
        } catch (error) {
            lastError = error as Error
            if (i < maxRetries - 1) {
                await new Promise(resolve => setTimeout(resolve, delay * (i + 1)))
            }
        }
    }

    throw lastError!
}
