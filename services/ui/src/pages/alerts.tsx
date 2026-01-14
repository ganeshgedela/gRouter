import { useState } from "react"
import {
    AlertTriangle,
    CheckCircle2,
    XCircle,
    Clock,
    Filter,
    MoreHorizontal,
    Bell,
    ShieldAlert,
    Search
} from "lucide-react"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
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
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
    DropdownMenuCheckboxItem
} from "@/components/ui/dropdown-menu"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

// Mock Data for Alerts
const alertsData = [
    {
        id: "ALT-1024",
        severity: "Critical",
        title: "BGP Session Down: Peer 10.0.5.2",
        source: "router-us-west-1",
        time: "12 mins ago",
        status: "Active",
        service: "BGP",
    },
    {
        id: "ALT-1023",
        severity: "Warning",
        title: "High CPU Utilization (>85%)",
        source: "router-eu-central-1",
        time: "45 mins ago",
        status: "Active",
        service: "System",
    },
    {
        id: "ALT-1022",
        severity: "Info",
        title: "Configuration Change Detected",
        source: "router-us-east-2",
        time: "2 hours ago",
        status: "Resolved",
        service: "Audit",
    },
    {
        id: "ALT-1021",
        severity: "Critical",
        title: "Interface eth0 Link Down",
        source: "router-ap-se-1",
        time: "5 hours ago",
        status: "Resolved",
        service: "Hardware",
    },
    {
        id: "ALT-1020",
        severity: "Warning",
        title: "High Memory Usage (>75%)",
        source: "router-us-west-1",
        time: "1 day ago",
        status: "Resolved",
        service: "System",
    }
]

export default function AlertsPage() {
    const [selectedSeverity, setSelectedSeverity] = useState<string[]>(["Critical", "Warning", "Info"])

    const toggleSeverity = (severity: string) => {
        setSelectedSeverity(prev =>
            prev.includes(severity)
                ? prev.filter(s => s !== severity)
                : [...prev, severity]
        )
    }

    const filteredAlerts = alertsData.filter(alert => selectedSeverity.includes(alert.severity))

    return (
        <div className="flex-1 space-y-4 p-8 pt-6">
            <div className="flex items-center justify-between space-y-2">
                <div>
                    <h2 className="text-3xl font-bold tracking-tight">Alerts & Incidents</h2>
                    <p className="text-muted-foreground">
                        Monitor system anomalies and manage incident response.
                    </p>
                </div>
                <div className="flex items-center space-x-2">
                    <Button variant="outline">
                        <Clock className="mr-2 h-4 w-4" /> History
                    </Button>
                    <Button>
                        <Bell className="mr-2 h-4 w-4" /> Configure Rules
                    </Button>
                </div>
            </div>

            {/* Stats Header */}
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Active Criticals</CardTitle>
                        <ShieldAlert className="h-4 w-4 text-destructive" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold text-destructive">1</div>
                        <p className="text-xs text-muted-foreground">Requires immediate attention</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Active Warnings</CardTitle>
                        <AlertTriangle className="h-4 w-4 text-amber-500" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold text-amber-500">1</div>
                        <p className="text-xs text-muted-foreground">Potential degradation</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Mean Time To Resolve</CardTitle>
                        <Clock className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">45m</div>
                        <p className="text-xs text-muted-foreground">-12m vs last week</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Total Incidents (24h)</CardTitle>
                        <Bell className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">12</div>
                        <p className="text-xs text-muted-foreground">Across 4 regions</p>
                    </CardContent>
                </Card>
            </div>

            <Tabs defaultValue="active" className="space-y-4">
                <TabsList>
                    <TabsTrigger value="active">Active Incidents</TabsTrigger>
                    <TabsTrigger value="history">Alert History</TabsTrigger>
                </TabsList>

                {/* ACTIVE INCIDENTS TAB */}
                <TabsContent value="active" className="space-y-4">
                    <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-2">
                            <div className="relative w-64">
                                <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
                                <Input placeholder="Search alerts..." className="pl-8" />
                            </div>
                            <DropdownMenu>
                                <DropdownMenuTrigger asChild>
                                    <Button variant="outline" size="sm" className="h-10 border-dashed">
                                        <Filter className="mr-2 h-4 w-4" />
                                        Severity
                                    </Button>
                                </DropdownMenuTrigger>
                                <DropdownMenuContent align="start">
                                    <DropdownMenuLabel>Filter by Severity</DropdownMenuLabel>
                                    <DropdownMenuSeparator />
                                    <DropdownMenuCheckboxItem checked={selectedSeverity.includes("Critical")} onCheckedChange={() => toggleSeverity("Critical")}>
                                        Critical
                                    </DropdownMenuCheckboxItem>
                                    <DropdownMenuCheckboxItem checked={selectedSeverity.includes("Warning")} onCheckedChange={() => toggleSeverity("Warning")}>
                                        Warning
                                    </DropdownMenuCheckboxItem>
                                    <DropdownMenuCheckboxItem checked={selectedSeverity.includes("Info")} onCheckedChange={() => toggleSeverity("Info")}>
                                        Info
                                    </DropdownMenuCheckboxItem>
                                </DropdownMenuContent>
                            </DropdownMenu>
                        </div>
                    </div>

                    <div className="rounded-md border">
                        <Table>
                            <TableHeader>
                                <TableRow>
                                    <TableHead className="w-[100px]">Severity</TableHead>
                                    <TableHead className="w-[300px]">Alert Title</TableHead>
                                    <TableHead>Source</TableHead>
                                    <TableHead>Service</TableHead>
                                    <TableHead>Time</TableHead>
                                    <TableHead>Status</TableHead>
                                    <TableHead className="text-right">Actions</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {filteredAlerts.length > 0 ? (
                                    filteredAlerts.map((alert) => (
                                        <TableRow key={alert.id}>
                                            <TableCell>
                                                <Badge
                                                    className={`
                                                        ${alert.severity === 'Critical' ? 'bg-destructive/10 text-destructive hover:bg-destructive/20 border-destructive/20' :
                                                            alert.severity === 'Warning' ? 'bg-amber-500/10 text-amber-500 hover:bg-amber-500/20 border-amber-500/20' :
                                                                'bg-blue-500/10 text-blue-500 hover:bg-blue-500/20 border-blue-500/20'}
                                                    `}
                                                    variant="outline"
                                                >
                                                    {alert.severity}
                                                </Badge>
                                            </TableCell>
                                            <TableCell className="font-medium">
                                                {alert.title}
                                                <div className="text-xs text-muted-foreground mt-0.5">{alert.id}</div>
                                            </TableCell>
                                            <TableCell>{alert.source}</TableCell>
                                            <TableCell>{alert.service}</TableCell>
                                            <TableCell className="text-muted-foreground">{alert.time}</TableCell>
                                            <TableCell>
                                                {alert.status === 'Active' ? (
                                                    <div className="flex items-center text-destructive text-sm font-medium">
                                                        <XCircle className="mr-1 h-3 w-3" /> Active
                                                    </div>
                                                ) : (
                                                    <div className="flex items-center text-emerald-500 text-sm font-medium">
                                                        <CheckCircle2 className="mr-1 h-3 w-3" /> Resolved
                                                    </div>
                                                )}
                                            </TableCell>
                                            <TableCell className="text-right">
                                                <DropdownMenu>
                                                    <DropdownMenuTrigger asChild>
                                                        <Button variant="ghost" className="h-8 w-8 p-0">
                                                            <MoreHorizontal className="h-4 w-4" />
                                                        </Button>
                                                    </DropdownMenuTrigger>
                                                    <DropdownMenuContent align="end">
                                                        <DropdownMenuItem>View Details</DropdownMenuItem>
                                                        <DropdownMenuItem>Acknowledge</DropdownMenuItem>
                                                        <DropdownMenuItem>Silence (1h)</DropdownMenuItem>
                                                        <DropdownMenuSeparator />
                                                        <DropdownMenuItem className="text-emerald-600">Resolve Incident</DropdownMenuItem>
                                                    </DropdownMenuContent>
                                                </DropdownMenu>
                                            </TableCell>
                                        </TableRow>
                                    ))
                                ) : (
                                    <TableRow>
                                        <TableCell colSpan={7} className="h-24 text-center">
                                            No alerts found.
                                        </TableCell>
                                    </TableRow>
                                )}
                            </TableBody>
                        </Table>
                    </div>
                </TabsContent>

                {/* HISTORY TAB */}
                <TabsContent value="history" className="space-y-4">
                    <div className="flex flex-col items-center justify-center min-h-[300px] border border-dashed rounded-md bg-muted/10">
                        <Clock className="h-10 w-10 text-muted-foreground mb-4 opacity-50" />
                        <h3 className="text-lg font-semibold">Archived Incidents</h3>
                        <p className="text-sm text-muted-foreground max-w-sm text-center">
                            Historical alert data is retained for 90 days. Use the Filters to search older logs.
                        </p>
                        <Button variant="outline" className="mt-4">Load Archived Data</Button>
                    </div>
                </TabsContent>
            </Tabs>
        </div>
    )
}
