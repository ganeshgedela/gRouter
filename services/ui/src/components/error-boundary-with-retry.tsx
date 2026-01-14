import React, { Component, ReactNode } from 'react'
import { ErrorState } from '@/components/ui/error-state'

interface Props {
    children: ReactNode
    onReset?: () => void
    fallback?: ReactNode
}

interface State {
    hasError: boolean
    error?: Error
}

export class ErrorBoundaryWithRetry extends Component<Props, State> {
    constructor(props: Props) {
        super(props)
        this.state = { hasError: false }
    }

    static getDerivedStateFromError(error: Error): State {
        return { hasError: true, error }
    }

    componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
        console.error('ErrorBoundary caught an error:', error, errorInfo)
    }

    handleReset = () => {
        this.setState({ hasError: false, error: undefined })
        if (this.props.onReset) {
            this.props.onReset()
        }
    }

    render() {
        if (this.state.hasError) {
            if (this.props.fallback) {
                return this.props.fallback
            }

            return (
                <ErrorState
                    title="Component Error"
                    message="This component encountered an unexpected error."
                    error={this.state.error}
                    onRetry={this.handleReset}
                />
            )
        }

        return this.props.children
    }
}
