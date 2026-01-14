import { useParams, useNavigate } from "react-router-dom"
import {
    ArrowLeft,
    MoreHorizontal,
    ExternalLink,
    Users,
    Network,
    Activity,
    HardDrive
} from "lucide-react"

import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Progress } from "@/components/ui/progress"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Separator } from "@/components/ui/separator"
import { LineChart } from "@/components/metrics/line-chart"

// Mock Data
const tenantData = {
    id: "tenant-123",
    name: "Acme Corp",
    status: "Active",
    tier: "Enterprise",
    users: 12,
    routers: 3,
    usage: {
        bandwidth: 75, // %
        routes: 45, // %
        storage: 20 // %
    }
}

const bandwidthHistory = Array.from({ length: 20 }, (_, i) => ({
    name: `Day ${i + 1}`,
    value: 40 + Math.random() * 40
}))

export default function TenantDetailPage() {
    const { id } = useParams()
    const navigate = useNavigate()

    return (
        <div className="flex-1 space-y-6 p-8 pt-6">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div className="flex items-center space-x-4">
                    <Button variant="ghost" size="icon" onClick={() => navigate("/tenants")}>
                        <ArrowLeft className="h-4 w-4" />
                    </Button>
                    <div>
                        <h2 className="text-2xl font-bold tracking-tight">{tenantData.name}</h2>
                        <div className="flex items-center text-sm text-muted-foreground mt-1">
                            <span className="mr-2">ID: {id || tenantData.id}</span>
                            <Badge variant="outline" className="text-xs bg-emerald-500/10 text-emerald-500 border-emerald-500/20">
                                {tenantData.status}
                            </Badge>
                        </div>
                    </div>
                </div>
                <div className="flex items-center space-x-2">
                    <Button variant="outline">Edit Profile</Button>
                    <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                            <Button variant="outline" size="icon">
                                <MoreHorizontal className="h-4 w-4" />
                            </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                            <DropdownMenuLabel>Actions</DropdownMenuLabel>
                            <DropdownMenuItem>Suspend Tenant</DropdownMenuItem>
                            <DropdownMenuItem>Reset Keys</DropdownMenuItem>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem className="text-destructive">Delete Tenant</DropdownMenuItem>
                        </DropdownMenuContent>
                    </DropdownMenu>
                </div>
            </div>

            <Separator />

            {/* Overview Stats */}
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Assigned Users</CardTitle>
                        <Users className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{tenantData.users}</div>
                        <p className="text-xs text-muted-foreground">Across 3 groups</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Active Routers</CardTitle>
                        <Network className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{tenantData.routers}</div>
                        <p className="text-xs text-muted-foreground">All operational</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Service Tier</CardTitle>
                        <Activity className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{tenantData.tier}</div>
                        <p className="text-xs text-muted-foreground">Next billing: Jan 1</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Storage Used</CardTitle>
                        <HardDrive className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">1.2 TB</div>
                        <p className="text-xs text-muted-foreground">of 5 TB Quota</p>
                    </CardContent>
                </Card>
            </div>

            {/* Detailed Usage */}
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
                <Card className="col-span-4">
                    <CardHeader>
                        <CardTitle>Bandwidth Usage</CardTitle>
                        <CardDescription>Daily egress traffic over the last 30 days.</CardDescription>
                    </CardHeader>
                    <CardContent className="pl-2">
                        <LineChart
                            data={bandwidthHistory}
                            lines={[{ key: "value", color: "#8b5cf6", name: "Traffic (GB)" }]}
                            height={300}
                        />
                    </CardContent>
                </Card>
                <Card className="col-span-3">
                    <CardHeader>
                        <CardTitle>Resource Quotas</CardTitle>
                        <CardDescription>Current utilization vs hard limits.</CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-8">
                        <div className="space-y-2">
                            <div className="flex items-center justify-between text-sm">
                                <span className="font-medium">Total Bandwidth</span>
                                <span className="text-muted-foreground">{tenantData.usage.bandwidth}%</span>
                            </div>
                            <Progress value={tenantData.usage.bandwidth} className="h-2" />
                        </div>
                        <div className="space-y-2">
                            <div className="flex items-center justify-between text-sm">
                                <span className="font-medium">Route Entries</span>
                                <span className="text-muted-foreground">{tenantData.usage.routes}%</span>
                            </div>
                            <Progress value={tenantData.usage.routes} className="h-2 [&>div]:bg-blue-500" />
                        </div>
                        <div className="space-y-2">
                            <div className="flex items-center justify-between text-sm">
                                <span className="font-medium">Object Storage</span>
                                <span className="text-muted-foreground">{tenantData.usage.storage}%</span>
                            </div>
                            <Progress value={tenantData.usage.storage} className="h-2 [&>div]:bg-emerald-500" />
                        </div>

                        <div className="pt-4">
                            <Button variant="outline" className="w-full">
                                <ExternalLink className="mr-2 h-4 w-4" /> Request Quota Increase
                            </Button>
                        </div>
                    </CardContent>
                </Card>
            </div>
        </div>
    )
}
