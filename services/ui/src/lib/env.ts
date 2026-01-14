/**
 * Environment Configuration
 * Typed access to environment variables
 */

interface EnvConfig {
    apiUrl: string
    wsUrl: string
    enableMocks: boolean
    enableDevtools: boolean
    sentryDsn?: string
    analyticsId?: string
    isDevelopment: boolean
    isProduction: boolean
}

function getEnv(key: string, defaultValue?: string): string {
    const value = import.meta.env[key]
    if (value === undefined && defaultValue === undefined) {
        console.warn(`Environment variable ${key} is not defined`)
        return ''
    }
    return value || defaultValue || ''
}

function getBoolEnv(key: string, defaultValue: boolean = false): boolean {
    const value = import.meta.env[key]
    if (value === undefined) return defaultValue
    return value === 'true' || value === '1'
}

export const env: EnvConfig = {
    apiUrl: getEnv('VITE_API_URL', 'http://localhost:8080'),
    wsUrl: getEnv('VITE_WS_URL', 'ws://localhost:8080/ws'),
    enableMocks: getBoolEnv('VITE_ENABLE_MOCKS', true),
    enableDevtools: getBoolEnv('VITE_ENABLE_DEVTOOLS', true),
    sentryDsn: getEnv('VITE_SENTRY_DSN'),
    analyticsId: getEnv('VITE_ANALYTICS_ID'),
    isDevelopment: import.meta.env.DEV,
    isProduction: import.meta.env.PROD,
}

// Log configuration in development
if (env.isDevelopment) {
    console.log('🔧 Environment Configuration:', {
        apiUrl: env.apiUrl,
        wsUrl: env.wsUrl,
        enableMocks: env.enableMocks,
        mode: import.meta.env.MODE,
    })
}
