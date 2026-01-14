/**
 * React Query Configuration
 * Centralized query client with default options
 */

import { QueryClient } from '@tanstack/react-query'

export const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            // Stale time: Consider data fresh for 1 minute
            staleTime: 1000 * 60,

            // Retry failed requests 3 times with exponential backoff
            retry: 3,
            retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000),

            // Refetch on window focus for real-time data
            refetchOnWindowFocus: true,

            // Refetch on reconnect
            refetchOnReconnect: true,

            // Don't refetch on mount if data is still fresh
            refetchOnMount: false,
        },
        mutations: {
            retry: 1,
        },
    },
})
