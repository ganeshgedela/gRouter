import { useState, useEffect } from "react"
import { useNavigate } from "react-router-dom"
import {
    Shield,
    AlertOctagon,
    Globe,
    Activity,
    Lock,
    MapPin,
    Radio
} from "lucide-react"

import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { ScrollArea } from "@/components/ui/scroll-area"
import { LineChart } from "@/components/metrics/line-chart"
import { ThreatMap } from "@/components/security/threat-map"
import { SystemPanicDialog } from "@/components/security/system-panic-dialog"

// --- Mock Data ---
const ATTACK_VECTORS = [
    { name: '00:00', ddos: 120, bot: 45, malware: 12 },
    { name: '04:00', ddos: 135, bot: 55, malware: 15 },
    { name: '08:00', ddos: 250, bot: 90, malware: 22 },
    { name: '12:00', ddos: 450, bot: 120, malware: 35 },
    { name: '16:00', ddos: 320, bot: 85, malware: 28 },
    { name: '20:00', ddos: 180, bot: 60, malware: 18 },
    { name: '24:00', ddos: 140, bot: 50, malware: 15 },
]

const LIVE_THREATS = [
    { id: 1, src: '192.168.45.2', country: 'RU', type: 'DDoS Flood', target: 'US-West-Core', time: 'Just now' },
    { id: 2, src: '10.5.2.11', country: 'CN', type: 'SSH Brute', target: 'EU-West-GW', time: '2s ago' },
    { id: 3, src: '172.16.99.5', country: 'BR', type: 'SQL Injection', target: 'App-Server-01', time: '5s ago' },
    { id: 4, src: '192.168.1.100', country: 'US', type: 'Port Scan', target: 'US-East-Core', time: '12s ago' },
    { id: 5, src: '45.33.22.11', country: 'IN', type: 'Malware C2', target: 'User-VLAN-5', time: '25s ago' },
    { id: 6, src: '89.11.22.33', country: 'DE', type: 'DDoS Flood', target: 'US-West-Core', time: '32s ago' },
]

export function SecurityDashboard() {
    const navigate = useNavigate()
    const [activeAttacks, setActiveAttacks] = useState(1243)
    const [blockedRequests, setBlockedRequests] = useState(854021)
    const [panicDialogOpen, setPanicDialogOpen] = useState(false)
    const [isPanicMode, setIsPanicMode] = useState(false)

    // Simulate live data updates
    useEffect(() => {
        const interval = setInterval(() => {
            setActiveAttacks(prev => Math.max(0, prev + Math.floor(Math.random() * 20) - 10))
            setBlockedRequests(prev => prev + Math.floor(Math.random() * 50))
        }, 2000)
        return () => clearInterval(interval)
    }, [])

    return (
        <div className="space-y-4">
            {/* Top Metrics Row */}
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                <Card className="border-red-500/50 bg-red-500/10">
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Active Threats</CardTitle>
                        <AlertOctagon className="h-4 w-4 text-red-500 animate-pulse" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold text-red-500">{activeAttacks.toLocaleString()}</div>
                        <p className="text-xs text-muted-foreground">+12% from last hour</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Blocked Requests</CardTitle>
                        <Shield className="h-4 w-4 text-emerald-500" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{blockedRequests.toLocaleString()}</div>
                        <p className="text-xs text-muted-foreground">Last 24 hours</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Global Threat Level</CardTitle>
                        <Activity className="h-4 w-4 text-yellow-500" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold text-yellow-500">ELEVATED</div>
                        <p className="text-xs text-muted-foreground">Due to regional DDoS activity</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Active Policies</CardTitle>
                        <Lock className="h-4 w-4 text-primary" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">142</div>
                        <p className="text-xs text-muted-foreground">12 rules updated today</p>
                    </CardContent>
                </Card>
            </div>

            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
                {/* Threat Map */}
                <Card className="col-span-4 h-[500px] flex flex-col">
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <Globe className="h-5 w-5" /> Live Threat Map
                        </CardTitle>
                        <CardDescription>Real-time visualization of attack origins.</CardDescription>
                    </CardHeader>
                    <CardContent className="flex-1 p-0 flex flex-col">
                        <ThreatMap />
                    </CardContent>
                </Card>

                {/* Right Column: Feed and Stats */}
                <div className="col-span-3 space-y-4">
                    {/* Attack Vectors Chart */}
                    <Card>
                        <CardHeader>
                            <CardTitle>Attack Vector Analysis</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <LineChart
                                data={ATTACK_VECTORS}
                                height={180}
                                lines={[
                                    { key: 'ddos', color: '#ef4444', name: 'DDoS' },
                                    { key: 'bot', color: '#f59e0b', name: 'Botnet' },
                                    { key: 'malware', color: '#a855f7', name: 'Malware' },
                                ]}
                            />
                        </CardContent>
                    </Card>

                    {/* Live Feed */}
                    <Card className="flex-1 flex flex-col">
                        <CardHeader className="mb-2">
                            <div className="flex items-center justify-between">
                                <CardTitle className="flex items-center gap-2">
                                    <Radio className="h-4 w-4 text-red-500 animate-pulse" /> Live Intercepts
                                </CardTitle>
                                <Badge variant="outline" className="font-mono text-xs">Real-time</Badge>
                            </div>
                        </CardHeader>
                        <CardContent className="flex-1 p-0">
                            <ScrollArea className="h-[250px]">
                                <div className="space-y-1 p-4 pt-0">
                                    {LIVE_THREATS.map(threat => (
                                        <div key={threat.id} className="flex items-center justify-between p-2 rounded border bg-muted/40 font-mono text-xs hover:bg-muted/80 transition-colors cursor-pointer group">
                                            <div className="flex items-center gap-3">
                                                <Badge variant="outline" className="w-8 justify-center p-0 h-5 border-zinc-700 bg-zinc-900">{threat.country}</Badge>
                                                <div className="space-y-0.5">
                                                    <div className="text-red-400 group-hover:text-red-300 font-semibold">{threat.src}</div>
                                                    <div className="text-muted-foreground">{threat.type}</div>
                                                </div>
                                            </div>
                                            <div className="text-right">
                                                <div className="text-muted-foreground">{threat.time}</div>
                                                <MapPin className="h-3 w-3 inline ml-1 opacity-50" />
                                            </div>
                                        </div>
                                    ))}
                                    <div className="p-2 text-center text-xs text-muted-foreground animate-pulse">
                                        Scanning for new threats...
                                    </div>
                                </div>
                            </ScrollArea>
                        </CardContent>
                    </Card>
                </div>
            </div>

            <div className="flex justify-end gap-2">
                <Button
                    variant="outline"
                    className="text-red-500 hover:text-red-500 hover:bg-red-500/10 border-red-500/20"
                    onClick={() => setPanicDialogOpen(true)}
                    disabled={isPanicMode}
                >
                    <AlertOctagon className="mr-2 h-4 w-4" />
                    {isPanicMode ? 'Panic Mode Active' : 'System Panic'}
                </Button>
                <Button onClick={() => navigate('/policies')}>
                    <Shield className="mr-2 h-4 w-4" /> Configure Firewall
                </Button>
            </div>

            <SystemPanicDialog
                open={panicDialogOpen}
                onOpenChange={setPanicDialogOpen}
                onConfirm={() => setIsPanicMode(true)}
            />
        </div>
    )
}
