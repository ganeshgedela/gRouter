import { Component, type ErrorInfo, type ReactNode } from "react"
import { AlertCircle, RefreshCcw } from "lucide-react"
import { Button } from "@/components/ui/button"

interface Props {
    children?: ReactNode
}

interface State {
    hasError: boolean
    error?: Error
}

export class ErrorBoundary extends Component<Props, State> {
    public state: State = {
        hasError: false
    }

    public static getDerivedStateFromError(error: Error): State {
        return { hasError: true, error }
    }

    public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
        console.error("Uncaught error:", error, errorInfo)
    }

    public render() {
        if (this.state.hasError) {
            return (
                <div className="flex h-screen w-full flex-col items-center justify-center p-4 bg-muted/10">
                    <div className="flex max-w-md flex-col items-center text-center space-y-4">
                        <div className="flex h-20 w-20 items-center justify-center rounded-full bg-red-100 dark:bg-red-900/20">
                            <AlertCircle className="h-10 w-10 text-red-600 dark:text-red-500" />
                        </div>
                        <h1 className="text-2xl font-bold tracking-tight">Something went wrong</h1>
                        <p className="text-muted-foreground">
                            We encountered an unexpected error. The application state has been logged for investigation.
                        </p>
                        <div className="p-4 bg-muted rounded-md text-xs font-mono text-left w-full overflow-auto max-h-32 opacity-70">
                            {this.state.error?.message}
                        </div>
                        <Button
                            onClick={() => window.location.reload()}
                            className="gap-2"
                        >
                            <RefreshCcw className="h-4 w-4" />
                            Reload Application
                        </Button>
                        <button
                            onClick={() => this.setState({ hasError: false })}
                            className="text-sm text-muted-foreground hover:underline"
                        >
                            Try to recover
                        </button>
                    </div>
                </div>
            )
        }

        return this.props.children
    }
}
