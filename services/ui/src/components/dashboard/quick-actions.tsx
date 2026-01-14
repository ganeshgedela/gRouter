import { Zap, Router, Workflow, Terminal, Plus } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

const quickActions = [
    {
        icon: Router,
        label: "Provision Router",
        description: "Deploy new node",
        action: () => console.log("Provision router"),
        color: "text-blue-500",
        bgColor: "bg-blue-500/10"
    },
    {
        icon: Workflow,
        label: "Run Workflow",
        description: "Execute automation",
        action: () => console.log("Run workflow"),
        color: "text-purple-500",
        bgColor: "bg-purple-500/10"
    },
    {
        icon: Terminal,
        label: "SSH Console",
        description: "Remote access",
        action: () => console.log("Open SSH"),
        color: "text-emerald-500",
        bgColor: "bg-emerald-500/10"
    },
    {
        icon: Zap,
        label: "Drain & Reboot",
        description: "Graceful restart",
        action: () => console.log("Drain & reboot"),
        color: "text-yellow-500",
        bgColor: "bg-yellow-500/10"
    }
]

export function QuickActions() {
    return (
        <Card>
            <CardHeader>
                <CardTitle className="flex items-center gap-2">
                    <Zap className="h-5 w-5 text-primary" />
                    Quick Actions
                </CardTitle>
                <CardDescription>
                    Common operations at your fingertips
                </CardDescription>
            </CardHeader>
            <CardContent className="grid gap-3">
                {quickActions.map((action, i) => (
                    <button
                        key={i}
                        onClick={action.action}
                        className="flex items-center gap-4 p-3 rounded-lg border hover:bg-accent/50 transition-colors text-left group"
                    >
                        <div className={`h-10 w-10 rounded-lg ${action.bgColor} flex items-center justify-center group-hover:scale-110 transition-transform`}>
                            <action.icon className={`h-5 w-5 ${action.color}`} />
                        </div>
                        <div className="flex-1">
                            <p className="font-medium text-sm">{action.label}</p>
                            <p className="text-xs text-muted-foreground">{action.description}</p>
                        </div>
                        <Plus className="h-4 w-4 text-muted-foreground group-hover:text-primary transition-colors" />
                    </button>
                ))}
            </CardContent>
        </Card>
    )
}
