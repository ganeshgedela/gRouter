import { useState } from "react"
import { useParams, useNavigate } from "react-router-dom"
import { ArrowLeft, Server, Activity, Terminal, Shield, Network } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { WebTerminal } from "@/components/terminal/web-terminal"
import { LineChart } from "@/components/metrics/line-chart"
import { ConfigEditor } from "@/components/config/config-editor"
import { RouterMetrics } from "@/components/network/router-metrics"

// Mock Data
const routerData = {
    id: "rtr-us-west-01",
    name: "Edge Router - Seattle",
    status: "Active",
    ip: "10.45.2.1",
    uptime: "14d 2h 12m",
    version: "v2.4.1"
}

export default function RouterDetailPage() {
    const { id } = useParams()
    const navigate = useNavigate()
    const [activeTab, setActiveTab] = useState("overview")

    return (
        <div className="flex-1 space-y-6 p-8 pt-6">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div className="flex items-center space-x-4">
                    <Button variant="ghost" size="icon" onClick={() => navigate("/routers")}>
                        <ArrowLeft className="h-4 w-4" />
                    </Button>
                    <div>
                        <h2 className="text-2xl font-bold tracking-tight">{routerData.name}</h2>
                        <div className="flex items-center text-sm text-muted-foreground mt-1">
                            <span className="mr-2">ID: {routerData.id}</span>
                            <Badge className="bg-emerald-500/10 text-emerald-500 hover:bg-emerald-500/20 shadow-none border-0">
                                {routerData.status}
                            </Badge>
                        </div>
                    </div>
                </div>
                <div className="flex items-center space-x-2">
                    <Button disabled variant="outline">
                        <Activity className="mr-2 h-4 w-4 text-orange-500" /> Maintenance Mode
                    </Button>
                    <Button>
                        Reboot Router
                    </Button>
                </div>
            </div>

            <Tabs defaultValue="overview" onValueChange={setActiveTab} className="space-y-4">
                <TabsList>
                    <TabsTrigger value="overview">Overview</TabsTrigger>
                    <TabsTrigger value="config">Configuration</TabsTrigger>
                    <TabsTrigger value="metrics">Metrics</TabsTrigger>
                    <TabsTrigger value="terminal" className="flex items-center gap-2">
                        <Terminal className="h-3 w-3" /> Console
                    </TabsTrigger>
                </TabsList>

                <TabsContent value="overview" className="space-y-4">
                    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                        <Card>
                            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                <CardTitle className="text-sm font-medium">Uptime</CardTitle>
                                <Activity className="h-4 w-4 text-muted-foreground" />
                            </CardHeader>
                            <CardContent>
                                <div className="text-2xl font-bold">{routerData.uptime}</div>
                                <p className="text-xs text-muted-foreground">Since last patch</p>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                <CardTitle className="text-sm font-medium">Throughput</CardTitle>
                                <Network className="h-4 w-4 text-muted-foreground" />
                            </CardHeader>
                            <CardContent>
                                <div className="text-2xl font-bold">4.2 Gbps</div>
                                <p className="text-xs text-muted-foreground">+12% peak usage</p>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                <CardTitle className="text-sm font-medium">CPU Load</CardTitle>
                                <Server className="h-4 w-4 text-muted-foreground" />
                            </CardHeader>
                            <CardContent>
                                <div className="text-2xl font-bold">45%</div>
                                <p className="text-xs text-muted-foreground">Normal load</p>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                <CardTitle className="text-sm font-medium">Security Alerts</CardTitle>
                                <Shield className="h-4 w-4 text-muted-foreground" />
                            </CardHeader>
                            <CardContent>
                                <div className="text-2xl font-bold">0</div>
                                <p className="text-xs text-muted-foreground">No active threats</p>
                            </CardContent>
                        </Card>
                    </div>

                    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
                        <Card className="col-span-4">
                            <CardHeader>
                                <CardTitle>Traffic Analysis</CardTitle>
                                <CardDescription>Ingress vs Egress traffic over time.</CardDescription>
                            </CardHeader>
                            <CardContent className="pl-2">
                                <LineChart
                                    data={Array.from({ length: 12 }, (_, i) => ({ name: `${i}h`, value: Math.floor(Math.random() * 100) }))}
                                    lines={[{ key: "value", color: "#3b82f6", name: "Traffic" }]}
                                    height={300}
                                />
                            </CardContent>
                        </Card>
                        <Card className="col-span-3">
                            <CardHeader>
                                <CardTitle>System Properties</CardTitle>
                                <CardDescription>Hardware and software details.</CardDescription>
                            </CardHeader>
                            <CardContent className="space-y-4">
                                <div className="grid grid-cols-2 gap-4 text-sm">
                                    <div className="text-muted-foreground">Hostname</div>
                                    <div className="font-medium">{routerData.id}</div>

                                    <div className="text-muted-foreground">IP Address</div>
                                    <div className="font-medium">{routerData.ip}</div>

                                    <div className="text-muted-foreground">Firmware</div>
                                    <div className="font-medium">{routerData.version}</div>

                                    <div className="text-muted-foreground">Kernel</div>
                                    <div className="font-medium">Linux 6.1.0-generic</div>

                                    <div className="text-muted-foreground">Architecture</div>
                                    <div className="font-medium">amd64</div>
                                </div>
                            </CardContent>
                        </Card>
                    </div>
                </TabsContent>

                <TabsContent value="config" className="space-y-4">
                    <ConfigEditor routerId={routerData.id} initialConfig="" />
                </TabsContent>

                <TabsContent value="metrics" className="pt-2">
                    <RouterMetrics routerId={routerData.id} />
                </TabsContent>

                <TabsContent value="terminal" className="space-y-4">
                    <Card className="border-0 shadow-none bg-transparent">
                        <WebTerminal routerId={routerData.id} />
                        <div className="mt-2 text-xs text-muted-foreground text-center">
                            Connected via secure WebSocket (TLS 1.3) • Session ID: #88219A
                        </div>
                    </Card>
                </TabsContent>
            </Tabs>
        </div>
    )
}
