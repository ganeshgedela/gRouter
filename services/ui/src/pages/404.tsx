import { Home, Search } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"

export default function NotFoundPage() {
    return (
        <div className="flex min-h-screen flex-col items-center justify-center p-8 text-center">
            <div className="space-y-6 max-w-md">
                {/* 404 Visual */}
                <div className="relative">
                    <h1 className="text-[180px] font-black leading-none text-primary/10">
                        404
                    </h1>
                    <div className="absolute inset-0 flex items-center justify-center">
                        <div className="flex h-24 w-24 items-center justify-center rounded-full bg-primary/10 backdrop-blur-sm">
                            <Search className="h-12 w-12 text-primary" />
                        </div>
                    </div>
                </div>

                {/* Message */}
                <div className="space-y-2">
                    <h2 className="text-3xl font-bold tracking-tight">Page Not Found</h2>
                    <p className="text-muted-foreground">
                        The page you're looking for doesn't exist or has been moved.
                    </p>
                </div>

                {/* Quick Search */}
                <div className="relative">
                    <Search className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                    <Input
                        placeholder="Search for routers, tenants, metrics..."
                        className="pl-10"
                    />
                </div>

                {/* Actions */}
                <div className="flex gap-2 justify-center pt-4">
                    <Button onClick={() => window.location.href = '/'}>
                        <Home className="mr-2 h-4 w-4" />
                        Back to Dashboard
                    </Button>
                    <Button variant="outline" onClick={() => window.history.back()}>
                        Go Back
                    </Button>
                </div>
            </div>
        </div>
    )
}
