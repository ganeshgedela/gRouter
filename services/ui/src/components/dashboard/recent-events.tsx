import { useState } from "react"
import { Clock, Activity, AlertTriangle, CheckCircle2, Info, Filter, ChevronRight } from "lucide-react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"

type EventType = "success" | "activity" | "warning" | "critical"
type EventFilter = "all" | EventType

const recentEvents = [
    {
        type: "success" as EventType,
        title: "Deployment Complete",
        description: "router-us-west-1 upgraded to v2.4.1",
        timestamp: "2 minutes ago",
        icon: CheckCircle2,
        color: "text-emerald-500",
        bgColor: "bg-emerald-500/10"
    },
    {
        type: "activity" as EventType,
        title: "Traffic Spike Detected",
        description: "EU region experiencing 40% increase",
        timestamp: "15 minutes ago",
        icon: Activity,
        color: "text-blue-500",
        bgColor: "bg-blue-500/10"
    },
    {
        type: "warning" as EventType,
        title: "Memory Usage High",
        description: "router-ap-se-3 at 85% capacity",
        timestamp: "1 hour ago",
        icon: AlertTriangle,
        color: "text-yellow-500",
        bgColor: "bg-yellow-500/10"
    },
    {
        type: "critical" as EventType,
        title: "Router Offline",
        description: "router-sa-east-1 unreachable",
        timestamp: "2 hours ago",
        icon: AlertTriangle,
        color: "text-red-500",
        bgColor: "bg-red-500/10"
    },
    {
        type: "success" as EventType,
        title: "Backup Completed",
        description: "Configuration backup successful",
        timestamp: "3 hours ago",
        icon: CheckCircle2,
        color: "text-emerald-500",
        bgColor: "bg-emerald-500/10"
    },
    {
        type: "activity" as EventType,
        title: "New Tenant Added",
        description: "Organization 'Acme Corp' created",
        timestamp: "5 hours ago",
        icon: Info,
        color: "text-blue-500",
        bgColor: "bg-blue-500/10"
    }
]

export function RecentEvents() {
    const [filter, setFilter] = useState<EventFilter>("all")
    const [showAll, setShowAll] = useState(false)

    const filteredEvents = recentEvents.filter(event =>
        filter === "all" ? true : event.type === filter
    )

    const displayedEvents = showAll ? filteredEvents : filteredEvents.slice(0, 5)

    const criticalCount = recentEvents.filter(e => e.type === "critical").length
    const warningCount = recentEvents.filter(e => e.type === "warning").length

    return (
        <Card className="h-full flex flex-col">
            <CardHeader className="pb-3">
                <div className="flex items-center justify-between">
                    <div>
                        <CardTitle className="text-lg">Recent Events</CardTitle>
                        <CardDescription className="text-xs">Live activity from your network</CardDescription>
                    </div>
                    <Badge variant="outline" className="font-mono text-xs gap-1.5">
                        <span className="flex h-2 w-2 rounded-full bg-emerald-500 animate-pulse" />
                        Live
                    </Badge>
                </div>

                {/* Quick Stats */}
                <div className="flex items-center gap-3 pt-3 text-xs">
                    <div className="flex items-center gap-1">
                        <span className="font-semibold text-foreground">{recentEvents.length}</span>
                        <span className="text-muted-foreground">events</span>
                    </div>
                    <div className="h-3 w-px bg-border" />
                    {criticalCount > 0 && (
                        <>
                            <div className="flex items-center gap-1">
                                <span className="font-semibold text-red-500">{criticalCount}</span>
                                <span className="text-muted-foreground">critical</span>
                            </div>
                            <div className="h-3 w-px bg-border" />
                        </>
                    )}
                    <div className="flex items-center gap-1">
                        <span className="font-semibold text-yellow-500">{warningCount}</span>
                        <span className="text-muted-foreground">warnings</span>
                    </div>
                </div>

                {/* Filter Buttons */}
                <div className="flex items-center gap-1 pt-3 flex-wrap">
                    <Button
                        variant={filter === "all" ? "default" : "outline"}
                        size="sm"
                        onClick={() => setFilter("all")}
                        className="h-7 text-xs"
                    >
                        All
                    </Button>
                    <Button
                        variant={filter === "critical" ? "default" : "outline"}
                        size="sm"
                        onClick={() => setFilter("critical")}
                        className="h-7 text-xs gap-1"
                    >
                        <span className="h-1.5 w-1.5 rounded-full bg-red-500" />
                        Critical
                    </Button>
                    <Button
                        variant={filter === "warning" ? "default" : "outline"}
                        size="sm"
                        onClick={() => setFilter("warning")}
                        className="h-7 text-xs gap-1"
                    >
                        <span className="h-1.5 w-1.5 rounded-full bg-yellow-500" />
                        Warnings
                    </Button>
                    <Button
                        variant={filter === "success" ? "default" : "outline"}
                        size="sm"
                        onClick={() => setFilter("success")}
                        className="h-7 text-xs gap-1"
                    >
                        <span className="h-1.5 w-1.5 rounded-full bg-emerald-500" />
                        Success
                    </Button>
                </div>
            </CardHeader>

            <CardContent className="flex-1 overflow-auto space-y-3 pt-0">
                {displayedEvents.length > 0 ? (
                    <>
                        {displayedEvents.map((event, i) => (
                            <div
                                key={i}
                                className={cn(
                                    "flex gap-3 items-start p-3 rounded-lg border transition-all cursor-pointer group",
                                    "hover:shadow-md hover:border-primary/30 hover:bg-muted/30"
                                )}
                            >
                                <div className={cn(
                                    "h-9 w-9 rounded-lg flex items-center justify-center shrink-0",
                                    event.bgColor
                                )}>
                                    <event.icon className={cn("h-4 w-4", event.color)} />
                                </div>
                                <div className="flex-1 space-y-1 min-w-0">
                                    <p className="text-sm font-medium leading-none">{event.title}</p>
                                    <p className="text-xs text-muted-foreground line-clamp-1">
                                        {event.description}
                                    </p>
                                    <div className="flex items-center gap-1 text-xs text-muted-foreground">
                                        <Clock className="h-3 w-3" />
                                        <span>{event.timestamp}</span>
                                    </div>
                                </div>
                                <ChevronRight className="h-4 w-4 text-muted-foreground opacity-0 group-hover:opacity-100 transition-opacity shrink-0" />
                            </div>
                        ))}

                        {/* Show More/Less Button */}
                        {filteredEvents.length > 5 && (
                            <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => setShowAll(!showAll)}
                                className="w-full text-xs"
                            >
                                {showAll ? "Show Less" : `View All ${filteredEvents.length} Events`}
                            </Button>
                        )}
                    </>
                ) : (
                    <div className="flex flex-col items-center justify-center py-8 text-center">
                        <Filter className="h-8 w-8 text-muted-foreground/50 mb-2" />
                        <p className="text-sm text-muted-foreground">No {filter} events</p>
                    </div>
                )}
            </CardContent>
        </Card>
    )
}
