import { FlaskConical } from "lucide-react"

export default function LabPage() {
    return (
        <div className="flex flex-col items-center justify-center h-full p-8 text-center space-y-4">
            <div className="p-6 bg-muted/20 rounded-full">
                <FlaskConical className="h-12 w-12 text-muted-foreground" />
            </div>
            <h2 className="text-2xl font-bold tracking-tight">Lab & Sandbox</h2>
            <p className="text-muted-foreground max-w-md">
                Isolated environment for testing configurations and simulating network failure scenarios.
            </p>
        </div>
    )
}
