import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { calculateUtilization, getUtilizationColor } from "@/lib/cidr-utils"
import type { Subnet } from "@/types/network"
import { Globe, Lock, Server, Unlock } from "lucide-react"
import { cn } from "@/lib/utils"

interface SubnetCardProps {
    subnet: Subnet
    routerHostname?: string
    onManage?: (subnet: Subnet) => void
    onAddRoute?: (subnet: Subnet) => void
}

export function SubnetCard({ subnet, routerHostname, onManage, onAddRoute }: SubnetCardProps) {
    const utilization = calculateUtilization(subnet.usedIPs, subnet.availableIPs)
    const utilizationColor = getUtilizationColor(utilization)

    return (
        <Card className="hover:border-primary/50 transition-colors">
            <CardHeader className="pb-3">
                <div className="flex items-start justify-between">
                    <div className="flex items-center gap-3">
                        <div className={cn(
                            "h-10 w-10 rounded-lg flex items-center justify-center",
                            subnet.type === 'public' ? "bg-blue-500/10" : "bg-purple-500/10"
                        )}>
                            {subnet.type === 'public' ? (
                                <Globe className={cn("h-5 w-5", "text-blue-500")} />
                            ) : (
                                <Lock className={cn("h-5 w-5", "text-purple-500")} />
                            )}
                        </div>
                        <div>
                            <CardTitle className="text-base">{subnet.name}</CardTitle>
                            <p className="text-xs text-muted-foreground font-mono mt-0.5">
                                {subnet.cidr} • {subnet.availabilityZone}
                            </p>
                        </div>
                    </div>
                    <Badge variant={subnet.status === 'active' ? 'default' : 'destructive'} className="uppercase text-[10px]">
                        {subnet.status}
                    </Badge>
                </div>
            </CardHeader>
            <CardContent className="space-y-3">
                <div className="grid grid-cols-2 gap-3 text-sm">
                    <div>
                        <p className="text-xs text-muted-foreground mb-1">Type</p>
                        <div className="flex items-center gap-1.5">
                            {subnet.type === 'public' ? (
                                <Unlock className="h-3 w-3 text-blue-500" />
                            ) : (
                                <Lock className="h-3 w-3 text-purple-500" />
                            )}
                            <span className="font-medium capitalize">{subnet.type}</span>
                        </div>
                    </div>
                    <div>
                        <p className="text-xs text-muted-foreground mb-1">Router</p>
                        <div className="flex items-center gap-1.5">
                            <Server className="h-3 w-3 text-muted-foreground" />
                            <span className="font-medium text-xs truncate">{routerHostname || subnet.routerId}</span>
                        </div>
                    </div>
                </div>

                <div>
                    <div className="flex justify-between items-center text-sm mb-1.5">
                        <span className="text-xs text-muted-foreground">IP Utilization</span>
                        <span className="font-mono text-xs">
                            {subnet.usedIPs.toLocaleString()} / {subnet.availableIPs.toLocaleString()}
                            <span className={cn(
                                "ml-1.5 font-semibold",
                                utilizationColor === 'green' && "text-emerald-500",
                                utilizationColor === 'yellow' && "text-yellow-500",
                                utilizationColor === 'red' && "text-red-500"
                            )}>
                                ({utilization}%)
                            </span>
                        </span>
                    </div>
                    <div className="h-2 w-full bg-secondary rounded-full overflow-hidden">
                        <div
                            className={cn(
                                "h-full transition-all",
                                utilizationColor === 'green' && "bg-emerald-500",
                                utilizationColor === 'yellow' && "bg-yellow-500",
                                utilizationColor === 'red' && "bg-red-500"
                            )}
                            style={{ width: `${utilization}%` }}
                        />
                    </div>
                </div>

                <div className="flex gap-2 pt-2">
                    <Button
                        size="sm"
                        variant="outline"
                        className="flex-1"
                        onClick={() => onManage?.(subnet)}
                    >
                        Manage
                    </Button>
                    <Button
                        size="sm"
                        className="flex-1"
                        onClick={() => onAddRoute?.(subnet)}
                    >
                        Add Route
                    </Button>
                </div>
            </CardContent>
        </Card>
    )
}
