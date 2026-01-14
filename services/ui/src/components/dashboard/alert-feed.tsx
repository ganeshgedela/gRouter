import {
    AlertTriangle,
    CheckCircle2,
    XCircle,
    Bell,
    Bot
} from "lucide-react"

import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"

const ALERTS = [
    {
        id: 1,
        title: "High Latency Detected",
        description: "Route optimization recommended for neighbor 10.0.5.2",
        time: "Just now",
        severity: "warning",
        type: "ai"
    },
    {
        id: 2,
        title: "BGP Session Established",
        description: "Peering with ISP-A (AS65001) is active.",
        time: "15 min ago",
        severity: "success",
        type: "system"
    },
    {
        id: 3,
        title: "Packet Loss Spike",
        description: "Interface eth2 experienced 0.5% loss.",
        time: "1 hour ago",
        severity: "critical",
        type: "system"
    }
]

export function AlertFeed() {
    return (
        <div className="space-y-4">
            {ALERTS.map((alert) => (
                <div key={alert.id} className="flex items-start space-x-4 p-4 rounded-lg border bg-card/50 hover:bg-card transition-colors">
                    <div className={`
                        p-2 rounded-full ring-1 ring-inset
                        ${alert.severity === 'critical' ? 'bg-red-500/10 text-red-500 ring-red-500/20' :
                            alert.severity === 'warning' ? 'bg-amber-500/10 text-amber-500 ring-amber-500/20' :
                                'bg-emerald-500/10 text-emerald-500 ring-emerald-500/20'}
                    `}>
                        {alert.type === 'ai' ? <Bot className="h-4 w-4" /> :
                            alert.severity === 'critical' ? <XCircle className="h-4 w-4" /> :
                                alert.severity === 'warning' ? <AlertTriangle className="h-4 w-4" /> :
                                    <CheckCircle2 className="h-4 w-4" />}
                    </div>
                    <div className="flex-1 space-y-1">
                        <div className="flex items-center justify-between">
                            <p className="text-sm font-medium leading-none">{alert.title}</p>
                            <span className="text-xs text-muted-foreground">{alert.time}</span>
                        </div>
                        <p className="text-sm text-muted-foreground">
                            {alert.description}
                        </p>
                        {alert.type === 'ai' && (
                            <div className="pt-2">
                                <Button size="sm" variant="outline" className="h-7 text-xs border-primary/20 text-primary hover:bg-primary/10">
                                    Apply Fix
                                </Button>
                            </div>
                        )}
                    </div>
                </div>
            ))}
            <Button variant="ghost" className="w-full text-xs text-muted-foreground">
                View All Alerts
            </Button>
        </div>
    )
}
