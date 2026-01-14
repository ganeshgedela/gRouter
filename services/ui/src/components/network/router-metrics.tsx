import { useState, useEffect } from "react"
import { Activity, Cpu, Server, Network, Clock } from "lucide-react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Button } from "@/components/ui/button"
import { LineChart } from "@/components/metrics/line-chart"

interface RouterMetricsProps {
    routerId?: string
}

export function RouterMetrics({ routerId }: RouterMetricsProps) {
    const [period, setPeriod] = useState("24h")
    const [isLoading, setIsLoading] = useState(true)

    // Simulate loading effect for realism
    useEffect(() => {
        setIsLoading(true)
        const timer = setTimeout(() => setIsLoading(false), 600)
        return () => clearTimeout(timer)
    }, [period])

    // Mock Data Generators
    const generateSeries = (count: number, base: number, variance: number) => {
        return Array.from({ length: count }, (_, i) => {
            const noise = Math.random() * variance - (variance / 2)
            const trend = Math.sin(i / 5) * (variance / 2)
            return Math.max(0, Math.floor(base + noise + trend))
        })
    }

    const dataPoints = period === "1h" ? 12 : period === "24h" ? 24 : 14
    const intervalLabel = period === "1h" ? "5m" : period === "24h" ? "1h" : "1d"

    // Generate consistent data based on period
    const cpuSeries = generateSeries(dataPoints, 45, 15)
    const memSeries = generateSeries(dataPoints, 62, 5)
    const inboundSeries = generateSeries(dataPoints, 450, 150)
    const outboundSeries = generateSeries(dataPoints, 380, 120)
    const latencySeries = generateSeries(dataPoints, 12, 4)
    const jitterSeries = generateSeries(dataPoints, 2, 1)

    const chartData = Array.from({ length: dataPoints }, (_, i) => ({
        name: `T-${dataPoints - i}${intervalLabel}`,
        cpu: cpuSeries[i],
        memory: memSeries[i],
        inbound: inboundSeries[i],
        outbound: outboundSeries[i],
        latency: latencySeries[i],
        jitter: jitterSeries[i]
    }))

    if (isLoading) {
        return (
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4 h-[400px]">
                {[...Array(4)].map((_, i) => (
                    <div key={i} className="animate-pulse rounded-lg bg-muted/20 border" />
                ))}
            </div>
        )
    }

    return (
        <div className="space-y-6 animate-in fade-in duration-500">
            {/* Controls */}
            <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
                <div className="space-y-1">
                    <h3 className="text-lg font-medium leading-none">Telemetry & Diagnostics</h3>
                    <p className="text-sm text-muted-foreground">
                        Real-time performance metrics for <span className="font-mono text-primary">{routerId || 'Router'}</span>
                    </p>
                </div>
                <div className="flex items-center gap-2">
                    <Tabs value={period} onValueChange={setPeriod}>
                        <TabsList>
                            <TabsTrigger value="1h">1H</TabsTrigger>
                            <TabsTrigger value="24h">24H</TabsTrigger>
                            <TabsTrigger value="7d">7D</TabsTrigger>
                        </TabsList>
                    </Tabs>
                    <Button variant="outline" size="icon" onClick={() => setIsLoading(true)}>
                        <Activity className="h-4 w-4" />
                    </Button>
                </div>
            </div>

            {/* Quick Stats Row */}
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                <Card className="bg-muted/30">
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Avg CPU Load</CardTitle>
                        <Cpu className="h-4 w-4 text-violet-500" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">45%</div>
                        <p className="text-xs text-muted-foreground flex items-center gap-1">
                            <span className="text-emerald-500 font-medium">Stable</span> within limits
                        </p>
                    </CardContent>
                </Card>
                <Card className="bg-muted/30">
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Memory Usage</CardTitle>
                        <Server className="h-4 w-4 text-emerald-500" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">6.2 GB</div>
                        <p className="text-xs text-muted-foreground">
                            of 8.0 GB Total
                        </p>
                    </CardContent>
                </Card>
                <Card className="bg-muted/30">
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Peak Throughput</CardTitle>
                        <Network className="h-4 w-4 text-blue-500" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">850 Mbps</div>
                        <p className="text-xs text-muted-foreground text-orange-500">
                            Peak at 14:02 UTC
                        </p>
                    </CardContent>
                </Card>
                <Card className="bg-muted/30">
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Avg Latency</CardTitle>
                        <Clock className="h-4 w-4 text-orange-500" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">12ms</div>
                        <p className="text-xs text-muted-foreground">
                            ±2ms jitter
                        </p>
                    </CardContent>
                </Card>
            </div>

            {/* Main Charts */}
            <div className="grid gap-6 md:grid-cols-2">
                {/* System Resources */}
                <Card className="col-span-1">
                    <CardHeader>
                        <CardTitle>System Resources</CardTitle>
                        <CardDescription>CPU and Memory utilization trends.</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <LineChart
                            data={chartData}
                            height={300}
                            fillArea={true}
                            lines={[
                                { key: "cpu", color: "#8b5cf6", name: "CPU Load (%)" },
                                { key: "memory", color: "#10b981", name: "Memory Usage (%)" },
                            ]}
                        />
                    </CardContent>
                </Card>

                {/* Network Throughput */}
                <Card className="col-span-1">
                    <CardHeader>
                        <CardTitle>Network Interface Traffic</CardTitle>
                        <CardDescription>Ingress vs Egress throughput (Mbps).</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <LineChart
                            data={chartData}
                            height={300}
                            fillArea={false}
                            lines={[
                                { key: "inbound", color: "#3b82f6", name: "Inbound" },
                                { key: "outbound", color: "#6366f1", name: "Outbound" },
                            ]}
                        />
                    </CardContent>
                </Card>

                {/* Latency & Jitter */}
                <Card className="md:col-span-2">
                    <CardHeader>
                        <CardTitle>Latency & QoS</CardTitle>
                        <CardDescription>Round-trip time and jitter analysis.</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <LineChart
                            data={chartData}
                            height={250}
                            fillArea={true}
                            lines={[
                                { key: "latency", color: "#f97316", name: "Latency (ms)" },
                                { key: "jitter", color: "#ec4899", name: "Jitter (ms)" },
                            ]}
                        />
                    </CardContent>
                </Card>
            </div>
        </div>
    )
}
