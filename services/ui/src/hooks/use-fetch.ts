import { useState, useEffect } from 'react'

interface UseFetchResult<T> {
    data: T | null
    error: Error | null
    isLoading: boolean
    refetch: () => void
}

export function useFetch<T>(
    fetcher: () => Promise<T>,
    dependencies: unknown[] = []
): UseFetchResult<T> {
    const [data, setData] = useState<T | null>(null)
    const [error, setError] = useState<Error | null>(null)
    const [isLoading, setIsLoading] = useState(true)
    const [refetchTrigger, setRefetchTrigger] = useState(0)

    const refetch = () => setRefetchTrigger(prev => prev + 1)

    useEffect(() => {
        let cancelled = false

        const fetchData = async () => {
            setIsLoading(true)
            setError(null)

            try {
                const result = await fetcher()
                if (!cancelled) {
                    setData(result)
                }
            } catch (err) {
                if (!cancelled) {
                    setError(err as Error)
                }
            } finally {
                if (!cancelled) {
                    setIsLoading(false)
                }
            }
        }

        fetchData()

        return () => {
            cancelled = true
        }
    }, [...dependencies, refetchTrigger])

    return { data, error, isLoading, refetch }
}

// Example usage hook for routers
import { routersApi } from '@/services/routers'
import type { RoutersListParams } from '@/types/api'

export function useRouters(params?: RoutersListParams) {
    return useFetch(() => routersApi.list(params), [JSON.stringify(params)])
}

export function useRouter(id: string) {
    return useFetch(() => routersApi.get(id), [id])
}
