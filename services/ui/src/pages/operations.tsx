import { useState } from "react"
import {
    History,
    Package,
    GitBranch,
    Clock,
    CheckCircle2,
    XCircle,
    AlertCircle,
    ChevronRight,
    ArrowUpCircle,
    ArrowDownCircle,
    RotateCcw
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import {
    Tabs,
    TabsContent,
    TabsList,
    TabsTrigger,
} from "@/components/ui/tabs"

// Mock deployment history
const deployments = [
    {
        id: "DEP-2024-001",
        version: "v2.4.1",
        status: "success",
        timestamp: "2024-01-09 14:23:00",
        duration: "2m 34s",
        router: "router-us-west-1",
        user: "admin@polaris.io",
        changes: ["Updated BGP config", "Increased memory limits"]
    },
    {
        id: "DEP-2024-002",
        version: "v2.4.1",
        status: "success",
        timestamp: "2024-01-09 14:18:00",
        duration: "3m 12s",
        router: "router-eu-west-1",
        user: "admin@polaris.io",
        changes: ["Updated BGP config", "Increased memory limits"]
    },
    {
        id: "DEP-2024-003",
        version: "v2.4.0",
        status: "failed",
        timestamp: "2024-01-08 09:45:00",
        duration: "1m 05s",
        router: "router-ap-se-3",
        user: "devops@polaris.io",
        changes: ["OSPF configuration update"]
    },
    {
        id: "DEP-2024-004",
        version: "v2.3.9",
        status: "rollback",
        timestamp: "2024-01-07 16:30:00",
        duration: "4m 22s",
        router: "router-us-east-2",
        user: "admin@polaris.io",
        changes: ["Reverted firewall rules"]
    }
]

const getStatusIcon = (status: string) => {
    switch (status) {
        case "success":
            return <CheckCircle2 className="h-4 w-4 text-emerald-500" />
        case "failed":
            return <XCircle className="h-4 w-4 text-destructive" />
        case "rollback":
            return <RotateCcw className="h-4 w-4 text-yellow-500" />
        default:
            return <AlertCircle className="h-4 w-4 text-muted-foreground" />
    }
}

const getStatusBadge = (status: string) => {
    const variants: Record<string, any> = {
        success: "default",
        failed: "destructive",
        rollback: "secondary"
    }
    return (
        <Badge variant={variants[status] || "outline"} className="capitalize">
            {status}
        </Badge>
    )
}

export default function OperationsPage() {
    const [selectedDeployment, setSelectedDeployment] = useState<string | null>(null)

    return (
        <div className="flex-1 space-y-4 p-8 pt-6">
            <div className="flex items-center justify-between">
                <div>
                    <h2 className="text-3xl font-bold tracking-tight">Operations</h2>
                    <p className="text-muted-foreground">
                        Deployment history, rollbacks, and version control.
                    </p>
                </div>
                <div className="flex gap-2">
                    <Button variant="outline">
                        <GitBranch className="mr-2 h-4 w-4" />
                        Version History
                    </Button>
                    <Button>
                        <Package className="mr-2 h-4 w-4" />
                        New Deployment
                    </Button>
                </div>
            </div>

            {/* Stats Cards */}
            <div className="grid gap-4 md:grid-cols-3">
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Total Deployments</CardTitle>
                        <History className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">142</div>
                        <p className="text-xs text-muted-foreground">Last 30 days</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Success Rate</CardTitle>
                        <CheckCircle2 className="h-4 w-4 text-emerald-500" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold text-emerald-500">98.6%</div>
                        <p className="text-xs text-muted-foreground">2 failed deployments</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Avg Duration</CardTitle>
                        <Clock className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">3m 22s</div>
                        <p className="text-xs text-muted-foreground">-12s from last week</p>
                    </CardContent>
                </Card>
            </div>

            <Tabs defaultValue="history" className="space-y-4">
                <TabsList>
                    <TabsTrigger value="history">Deployment History</TabsTrigger>
                    <TabsTrigger value="rollbacks">Rollback Queue</TabsTrigger>
                    <TabsTrigger value="versions">Version Matrix</TabsTrigger>
                </TabsList>

                <TabsContent value="history" className="space-y-4">
                    <Card>
                        <CardHeader>
                            <CardTitle>Recent Deployments</CardTitle>
                            <CardDescription>
                                Track configuration changes and version updates across your fleet.
                            </CardDescription>
                        </CardHeader>
                        <CardContent>
                            <Table>
                                <TableHeader>
                                    <TableRow>
                                        <TableHead>ID</TableHead>
                                        <TableHead>Status</TableHead>
                                        <TableHead>Version</TableHead>
                                        <TableHead>Router</TableHead>
                                        <TableHead>User</TableHead>
                                        <TableHead>Duration</TableHead>
                                        <TableHead>Timestamp</TableHead>
                                        <TableHead className="text-right">Actions</TableHead>
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    {deployments.map((dep) => (
                                        <TableRow key={dep.id}>
                                            <TableCell className="font-mono text-xs">{dep.id}</TableCell>
                                            <TableCell>
                                                <div className="flex items-center gap-2">
                                                    {getStatusIcon(dep.status)}
                                                    {getStatusBadge(dep.status)}
                                                </div>
                                            </TableCell>
                                            <TableCell>
                                                <Badge variant="outline" className="font-mono">
                                                    {dep.version}
                                                </Badge>
                                            </TableCell>
                                            <TableCell className="font-medium">{dep.router}</TableCell>
                                            <TableCell className="text-muted-foreground text-sm">{dep.user}</TableCell>
                                            <TableCell className="font-mono text-xs">{dep.duration}</TableCell>
                                            <TableCell className="text-muted-foreground text-sm">{dep.timestamp}</TableCell>
                                            <TableCell className="text-right">
                                                <Button
                                                    variant="ghost"
                                                    size="sm"
                                                    onClick={() => setSelectedDeployment(dep.id)}
                                                >
                                                    Details
                                                    <ChevronRight className="ml-1 h-4 w-4" />
                                                </Button>
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </CardContent>
                    </Card>
                </TabsContent>

                <TabsContent value="rollbacks" className="space-y-4">
                    <div className="bg-muted/30 border border-dashed rounded-lg p-12 text-center">
                        <ArrowDownCircle className="mx-auto h-12 w-12 text-muted-foreground opacity-50" />
                        <h3 className="mt-4 text-lg font-semibold">No Pending Rollbacks</h3>
                        <p className="text-sm text-muted-foreground mt-2 max-w-sm mx-auto">
                            All deployments are stable. Rollback actions will appear here when triggered.
                        </p>
                    </div>
                </TabsContent>

                <TabsContent value="versions" className="space-y-4">
                    <Card>
                        <CardHeader>
                            <CardTitle>Version Distribution</CardTitle>
                            <CardDescription>
                                Current version deployed across all routers.
                            </CardDescription>
                        </CardHeader>
                        <CardContent>
                            <div className="space-y-4">
                                <div className="flex items-center justify-between p-4 border rounded-lg">
                                    <div className="flex items-center gap-3">
                                        <div className="h-10 w-10 rounded-lg bg-primary/10 flex items-center justify-center">
                                            <ArrowUpCircle className="h-5 w-5 text-primary" />
                                        </div>
                                        <div>
                                            <p className="font-semibold font-mono">v2.4.1</p>
                                            <p className="text-xs text-muted-foreground">Latest Stable</p>
                                        </div>
                                    </div>
                                    <div className="text-right">
                                        <p className="text-2xl font-bold">3</p>
                                        <p className="text-xs text-muted-foreground">routers</p>
                                    </div>
                                </div>
                                <div className="flex items-center justify-between p-4 border rounded-lg">
                                    <div className="flex items-center gap-3">
                                        <div className="h-10 w-10 rounded-lg bg-muted flex items-center justify-center">
                                            <Package className="h-5 w-5 text-muted-foreground" />
                                        </div>
                                        <div>
                                            <p className="font-semibold font-mono">v2.4.0</p>
                                            <p className="text-xs text-muted-foreground">Previous</p>
                                        </div>
                                    </div>
                                    <div className="text-right">
                                        <p className="text-2xl font-bold">1</p>
                                        <p className="text-xs text-muted-foreground">router</p>
                                    </div>
                                </div>
                                <div className="flex items-center justify-between p-4 border rounded-lg">
                                    <div className="flex items-center gap-3">
                                        <div className="h-10 w-10 rounded-lg bg-yellow-500/10 flex items-center justify-center">
                                            <AlertCircle className="h-5 w-5 text-yellow-500" />
                                        </div>
                                        <div>
                                            <p className="font-semibold font-mono">v2.3.9</p>
                                            <p className="text-xs text-muted-foreground">Deprecated</p>
                                        </div>
                                    </div>
                                    <div className="text-right">
                                        <p className="text-2xl font-bold">1</p>
                                        <p className="text-xs text-muted-foreground">router</p>
                                    </div>
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>
            </Tabs>
        </div>
    )
}
