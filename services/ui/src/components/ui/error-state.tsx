import { AlertTriangle, RefreshCw, Home } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

interface ErrorStateProps {
    title?: string
    message?: string
    error?: Error | null
    onRetry?: () => void
    showHomeButton?: boolean
}

export function ErrorState({
    title = "Something went wrong",
    message = "We encountered an error while loading this data. Please try again.",
    error,
    onRetry,
    showHomeButton = false
}: ErrorStateProps) {
    return (
        <div className="flex min-h-[400px] w-full items-center justify-center p-8">
            <Card className="max-w-md border-destructive/20">
                <CardHeader className="text-center">
                    <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-destructive/10">
                        <AlertTriangle className="h-8 w-8 text-destructive" />
                    </div>
                    <CardTitle className="text-xl">{title}</CardTitle>
                    <CardDescription className="text-base">{message}</CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                    {error && process.env.NODE_ENV === 'development' && (
                        <div className="rounded-md bg-muted p-3">
                            <p className="text-xs font-mono text-muted-foreground break-all">
                                {error.message}
                            </p>
                        </div>
                    )}
                    <div className="flex gap-2 justify-center">
                        {onRetry && (
                            <Button onClick={onRetry} variant="default">
                                <RefreshCw className="mr-2 h-4 w-4" />
                                Try Again
                            </Button>
                        )}
                        {showHomeButton && (
                            <Button onClick={() => window.location.href = '/'} variant="outline">
                                <Home className="mr-2 h-4 w-4" />
                                Go Home
                            </Button>
                        )}
                    </div>
                </CardContent>
            </Card>
        </div>
    )
}
