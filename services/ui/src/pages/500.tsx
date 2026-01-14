import { ServerCrash, RefreshCw, Home } from "lucide-react"
import { Button } from "@/components/ui/button"

export default function ServerErrorPage() {
    return (
        <div className="flex min-h-screen flex-col items-center justify-center p-8 text-center">
            <div className="space-y-6 max-w-md">
                {/* 500 Visual */}
                <div className="relative">
                    <h1 className="text-[180px] font-black leading-none text-destructive/10">
                        500
                    </h1>
                    <div className="absolute inset-0 flex items-center justify-center">
                        <div className="flex h-24 w-24 items-center justify-center rounded-full bg-destructive/10 backdrop-blur-sm">
                            <ServerCrash className="h-12 w-12 text-destructive" />
                        </div>
                    </div>
                </div>

                {/* Message */}
                <div className="space-y-2">
                    <h2 className="text-3xl font-bold tracking-tight">Server Error</h2>
                    <p className="text-muted-foreground">
                        Something went wrong on our end. Our team has been notified and is working on a fix.
                    </p>
                </div>

                {/* Actions */}
                <div className="flex gap-2 justify-center pt-4">
                    <Button onClick={() => window.location.reload()}>
                        <RefreshCw className="mr-2 h-4 w-4" />
                        Reload Page
                    </Button>
                    <Button variant="outline" onClick={() => window.location.href = '/'}>
                        <Home className="mr-2 h-4 w-4" />
                        Go Home
                    </Button>
                </div>

                {/* Additional Info */}
                <div className="pt-8 border-t">
                    <p className="text-xs text-muted-foreground">
                        If this problem persists, contact support with error ID: <span className="font-mono">ERR-{Date.now()}</span>
                    </p>
                </div>
            </div>
        </div>
    )
}
