import { Globe } from "lucide-react"

export default function EdgePage() {
    return (
        <div className="flex flex-col items-center justify-center h-full p-8 text-center space-y-4">
            <div className="p-6 bg-muted/20 rounded-full">
                <Globe className="h-12 w-12 text-muted-foreground" />
            </div>
            <h2 className="text-2xl font-bold tracking-tight">Edge & Multi-Cloud</h2>
            <p className="text-muted-foreground max-w-md">
                Manage edge nodes and cross-cloud connectivity policies.
            </p>
        </div>
    )
}
