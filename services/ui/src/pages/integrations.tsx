import { Link2 } from "lucide-react"

export default function IntegrationsPage() {
    return (
        <div className="flex flex-col items-center justify-center h-full p-8 text-center space-y-4">
            <div className="p-6 bg-muted/20 rounded-full">
                <Link2 className="h-12 w-12 text-muted-foreground" />
            </div>
            <h2 className="text-2xl font-bold tracking-tight">Integrations</h2>
            <p className="text-muted-foreground max-w-md">
                Manage connections to external providers, webhooks, and third-party monitoring tools.
            </p>
        </div>
    )
}
