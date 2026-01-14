import { useState } from "react"
import { CalendarDays, RefreshCw, Database, Activity, Signal, Router, Globe, Download, Clock } from "lucide-react"

import { Button } from "@/components/ui/button"
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import { LineChart } from "@/components/metrics/line-chart"
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuCheckboxItem,
    DropdownMenuTrigger,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuItem
} from "@/components/ui/dropdown-menu"
import { Badge } from "@/components/ui/badge"
import { Switch } from "@/components/ui/switch"
import { Label } from "@/components/ui/label"
import { MetricsSkeleton } from "@/components/skeletons/metrics-skeleton"
import { useToast } from "@/components/ui/use-toast"

// Mock Data Generators (Higher resolution for sparklines)
const generateTimeSeries = (points: number, base: number, volatility: number) => {
    return Array.from({ length: points }, (_, i) => {
        const time = new Date()
        time.setMinutes(time.getMinutes() - (points - i) * 5)
        return {
            name: time.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
            value: Math.max(0, Math.min(100, Math.round(base + Math.sin(i / 5) * volatility + (Math.random() - 0.5) * volatility / 2)))
        }
    })
}

// Data sets
const cpuData = generateTimeSeries(24, 45, 15).map(d => ({ name: d.name, cpu: d.value }))
const networkData = generateTimeSeries(24, 50, 20).map(d => ({
    name: d.name,
    ingress: d.value,
    egress: Math.max(0, Math.min(100, d.value + (Math.random() - 0.5) * 20))
}))
const latencyData = generateTimeSeries(24, 45, 10).map(d => ({
    name: d.name,
    p95: d.value,
    p50: Math.max(10, d.value - 15)
}))
const dropData = generateTimeSeries(24, 0.5, 0.5).map(d => ({
    name: d.name,
    drops: Math.max(0, d.value / 100)
}))
const routesData = generateTimeSeries(24, 40, 2).map((d, i) => ({
    name: d.name,
    ipv4: 8500 + i * 10,
    ipv6: 1200 + i * 2
}))

// Sparkline datasets (More points, smaller range)
const sparkLatency = generateTimeSeries(20, 25, 5).map(d => ({ name: '', val: d.value }))
const sparkLoss = generateTimeSeries(20, 5, 5).map(d => ({ name: '', val: d.value }))
const sparkRoutes = generateTimeSeries(20, 50, 2).map(d => ({ name: '', val: d.value }))
const sparkBGP = generateTimeSeries(20, 100, 0).map(d => ({ name: '', val: d.value }))

export default function MetricsPage() {
    const { toast } = useToast()
    const [isRefreshing, setIsRefreshing] = useState(false)
    const [selectedRegion, setSelectedRegion] = useState("All Regions")
    const [selectedRouter, setSelectedRouter] = useState("All Routers")
    const [timeRange, setTimeRange] = useState("2h")
    const [autoRefresh, setAutoRefresh] = useState(true)

    const handleRefresh = () => {
        setIsRefreshing(true)
        setTimeout(() => setIsRefreshing(false), 800)
    }

    const handleExport = (format: string) => {
        toast({
            title: "Export Started",
            description: `Generating ${format} report for ${timeRange}...`,
        })
    }

    if (isRefreshing) {
        return <MetricsSkeleton />
    }

    return (
        <div className="flex-1 space-y-4 p-8 pt-6">
            <div className="flex flex-col space-y-2 md:flex-row md:items-center md:justify-between md:space-y-0">
                <div>
                    <h2 className="text-3xl font-bold tracking-tight">System Metrics</h2>
                    <p className="text-muted-foreground">
                        Advanced telemetry for Control Plane and Data Plane.
                    </p>
                </div>

                {/* Visual Toolbar */}
                <div className="flex items-center space-x-2 bg-muted/40 p-1 rounded-lg border">
                    {/* Filters */}
                    <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="sm" className="h-8">
                                <Globe className="mr-2 h-4 w-4" />
                                {selectedRegion}
                            </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="start">
                            <DropdownMenuCheckboxItem checked={selectedRegion === "All Regions"} onCheckedChange={() => setSelectedRegion("All Regions")}>All Regions</DropdownMenuCheckboxItem>
                            <DropdownMenuCheckboxItem checked={selectedRegion === "US West"} onCheckedChange={() => setSelectedRegion("US West")}>US West</DropdownMenuCheckboxItem>
                            <DropdownMenuCheckboxItem checked={selectedRegion === "EU Central"} onCheckedChange={() => setSelectedRegion("EU Central")}>EU Central</DropdownMenuCheckboxItem>
                        </DropdownMenuContent>
                    </DropdownMenu>

                    <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="sm" className="h-8">
                                <Router className="mr-2 h-4 w-4" />
                                {selectedRouter}
                            </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="start">
                            <DropdownMenuCheckboxItem checked={selectedRouter === "All Routers"} onCheckedChange={() => setSelectedRouter("All Routers")}>All Routers</DropdownMenuCheckboxItem>
                            <DropdownMenuCheckboxItem checked={selectedRouter === "router-us-west-1"} onCheckedChange={() => setSelectedRouter("router-us-west-1")}>router-us-west-1</DropdownMenuCheckboxItem>
                        </DropdownMenuContent>
                    </DropdownMenu>

                    <div className="h-4 w-[1px] bg-border mx-1" />

                    {/* Time Range Selector */}
                    <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="sm" className="h-8">
                                <Clock className="mr-2 h-4 w-4" />
                                {timeRange === '1h' ? 'Last Hour' :
                                    timeRange === '2h' ? 'Last 2 Hours' :
                                        timeRange === '24h' ? 'Last 24 Hours' : 'Last 7 Days'}
                            </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                            <DropdownMenuLabel>Time Range</DropdownMenuLabel>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem onClick={() => setTimeRange("1h")}>Last Hour</DropdownMenuItem>
                            <DropdownMenuItem onClick={() => setTimeRange("2h")}>Last 2 Hours</DropdownMenuItem>
                            <DropdownMenuItem onClick={() => setTimeRange("6h")}>Last 6 Hours</DropdownMenuItem>
                            <DropdownMenuItem onClick={() => setTimeRange("24h")}>Last 24 Hours</DropdownMenuItem>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem onClick={() => setTimeRange("7d")}>Last 7 Days</DropdownMenuItem>
                        </DropdownMenuContent>
                    </DropdownMenu>

                    <div className="h-4 w-[1px] bg-border mx-1" />

                    {/* Auto Refresh Toggle */}
                    <div className="flex items-center space-x-2 px-2">
                        <Switch
                            id="auto-refresh"
                            checked={autoRefresh}
                            onChange={(e) => setAutoRefresh(e.target.checked)}
                            className="scale-75"
                        />
                        <Label htmlFor="auto-refresh" className="text-xs font-normal text-muted-foreground cursor-pointer">Auto</Label>
                    </div>

                    <div className="h-4 w-[1px] bg-border mx-1" />

                    {/* Export */}
                    <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="icon" className="h-8 w-8">
                                <Download className="h-4 w-4" />
                            </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                            <DropdownMenuItem onClick={() => handleExport('CSV')}>Export CSV</DropdownMenuItem>
                            <DropdownMenuItem onClick={() => handleExport('PDF')}>Export PDF Report</DropdownMenuItem>
                        </DropdownMenuContent>
                    </DropdownMenu>

                    <Button size="icon" variant="ghost" className="h-8 w-8" onClick={handleRefresh} disabled={isRefreshing}>
                        <RefreshCw className={`h-4 w-4 ${isRefreshing ? "animate-spin" : ""}`} />
                    </Button>
                </div>
            </div>



            {/* Key Indicators Row with Sparklines */}
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Avg Latency</CardTitle>
                        <Activity className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="flex items-end justify-between">
                            <div>
                                <div className="text-2xl font-bold">24 ms</div>
                                <p className="text-xs text-muted-foreground flex items-center mt-1">
                                    <Badge variant="secondary" className="mr-2 bg-emerald-500/10 text-emerald-500 hover:bg-emerald-500/20 text-[10px] px-1 py-0">Healthy</Badge>
                                    -2ms vs avg
                                </p>
                            </div>
                            <div className="w-[80px] h-[30px]">
                                <LineChart data={sparkLatency} lines={[{ key: 'val', color: '#10b981', name: '' }]} height={30} showAxes={false} showGrid={false} showTooltip={false} />
                            </div>
                        </div>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Packet Loss</CardTitle>
                        <Signal className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="flex items-end justify-between">
                            <div>
                                <div className="text-2xl font-bold">0.001%</div>
                                <p className="text-xs text-muted-foreground mt-1">
                                    0 drops in last 5m
                                </p>
                            </div>
                            <div className="w-[80px] h-[30px]">
                                <LineChart data={sparkLoss} lines={[{ key: 'val', color: '#f59e0b', name: '' }]} height={30} showAxes={false} showGrid={false} showTooltip={false} />
                            </div>
                        </div>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Active Routes</CardTitle>
                        <Database className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="flex items-end justify-between">
                            <div>
                                <div className="text-2xl font-bold">12,450</div>
                                <p className="text-xs text-muted-foreground mt-1 text-emerald-500">
                                    +45 new routes
                                </p>
                            </div>
                            <div className="w-[80px] h-[30px]">
                                <LineChart data={sparkRoutes} lines={[{ key: 'val', color: '#8b5cf6', name: '' }]} height={30} showAxes={false} showGrid={false} showTooltip={false} />
                            </div>
                        </div>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">BGP Sessions</CardTitle>
                        <Globe className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="flex items-end justify-between">
                            <div>
                                <div className="text-2xl font-bold text-emerald-500">24/24</div>
                                <p className="text-xs text-muted-foreground mt-1">
                                    All peers established
                                </p>
                            </div>
                            <div className="w-[80px] h-[30px]">
                                <LineChart data={sparkBGP} lines={[{ key: 'val', color: '#10b981', name: '' }]} height={30} showAxes={false} showGrid={false} showTooltip={false} />
                            </div>
                        </div>
                    </CardContent>
                </Card>
            </div>

            {/* Main Charts Grid - Now with Gradients */}
            <div className="grid gap-4 md:grid-cols-2">
                {/* Latency Chart */}
                <Card className="col-span-1">
                    <CardHeader>
                        <CardTitle>Network Latency (RTT)</CardTitle>
                        <CardDescription>P50 vs P95 round trip times across the mesh.</CardDescription>
                    </CardHeader>
                    <CardContent className="pl-0">
                        <LineChart
                            data={latencyData}
                            lines={[
                                { key: "p95", color: "#f59e0b", name: "P95 (ms)" },
                                { key: "p50", color: "#3b82f6", name: "P50 (ms)" }
                            ]}
                            height={300}
                            fillArea={true}
                        />
                    </CardContent>
                </Card>

                {/* Packet Loss / Errors */}
                <Card className="col-span-1">
                    <CardHeader>
                        <CardTitle>Hardware Interface Drops</CardTitle>
                        <CardDescription>IF-MIB Discards per second.</CardDescription>
                    </CardHeader>
                    <CardContent className="pl-0">
                        <LineChart
                            data={dropData}
                            lines={[{ key: "drops", color: "#ef4444", name: "Drops/s" }]}
                            height={300}
                            fillArea={true}
                        />
                    </CardContent>
                </Card>

                {/* Traffic - Expanded */}
                <Card className="col-span-2">
                    <CardHeader>
                        <CardTitle>Throughput Volume</CardTitle>
                        <CardDescription>Aggregate Ingress/Egress bandwidth usage.</CardDescription>
                    </CardHeader>
                    <CardContent className="pl-0">
                        <LineChart
                            data={networkData}
                            lines={[
                                { key: "ingress", color: "#10b981", name: "Ingress (Gbps)" },
                                { key: "egress", color: "#f43f5e", name: "Egress (Gbps)" }
                            ]}
                            height={350}
                            fillArea={true}
                        />
                    </CardContent>
                </Card>

                {/* Control Plane Stats */}
                <Card className="col-span-2 md:col-span-1">
                    <CardHeader>
                        <CardTitle>Route Table Growth</CardTitle>
                        <CardDescription>Total RIB entries over time.</CardDescription>
                    </CardHeader>
                    <CardContent className="pl-0">
                        <LineChart
                            data={routesData}
                            lines={[
                                { key: "ipv4", color: "#8b5cf6", name: "IPv4 Routes" },
                                { key: "ipv6", color: "#ec4899", name: "IPv6 Routes" }
                            ]}
                            height={250}
                            fillArea={true}
                        />
                    </CardContent>
                </Card>

                <Card className="col-span-2 md:col-span-1">
                    <CardHeader>
                        <CardTitle>System Load (Control Plane)</CardTitle>
                        <CardDescription>CPU usage for routing daemons (BIRD/FRR).</CardDescription>
                    </CardHeader>
                    <CardContent className="pl-0">
                        <LineChart
                            data={cpuData}
                            lines={[{ key: "cpu", color: "#6366f1", name: "Load %" }]}
                            height={250}
                            fillArea={true}
                        />
                    </CardContent>
                </Card>
            </div>
        </div>
    )
}
