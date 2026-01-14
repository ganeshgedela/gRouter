import { useState } from "react"
import {
    AlertOctagon,
    ShieldAlert,
    Flame,
    Activity,
    Lock,
    RotateCcw
} from "lucide-react"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Switch } from "@/components/ui/switch"
import { Label } from "@/components/ui/label"
import { Separator } from "@/components/ui/separator"

export function DDoSControls() {
    const [underAttackMode, setUnderAttackMode] = useState(false)
    const [rateLimitICMP, setRateLimitICMP] = useState(true)
    const [rateLimitUDP, setRateLimitUDP] = useState(true)

    return (
        <div className="space-y-6">
            {/* Panic Mode Header */}
            <Card className={`border-2 transition-all ${underAttackMode ? 'border-red-500 bg-red-500/10' : 'border-muted'}`}>
                <CardHeader className="flex flex-row items-center justify-between">
                    <div>
                        <CardTitle className="text-2xl flex items-center gap-2">
                            {underAttackMode ? <Flame className="h-6 w-6 text-red-500 animate-pulse" /> : <ShieldAlert className="h-6 w-6" />}
                            Under Attack Mode
                        </CardTitle>
                        <CardDescription className={underAttackMode ? "text-red-400" : ""}>
                            Emergency override for active volumetric attacks. Enforces strict challenges and rate limiting.
                        </CardDescription>
                    </div>
                    <Switch
                        checked={underAttackMode}
                        onChange={(e) => setUnderAttackMode(e.target.checked)}
                        className="scale-150 data-[state=checked]:bg-red-500"
                    />
                </CardHeader>
                {underAttackMode && (
                    <CardContent>
                        <div className="bg-red-500/20 border border-red-500/50 rounded-lg p-4 flex items-center gap-4 text-red-200">
                            <AlertOctagon className="h-8 w-8" />
                            <div>
                                <div className="font-bold">SYSTEM IN DEFENSIVE POSTURE</div>
                                <div className="text-sm">Javascript challenges enabled for all incoming traffic. Rate limits tightened by 500%. Latency may increase.</div>
                            </div>
                        </div>
                    </CardContent>
                )}
            </Card>

            <div className="grid gap-6 md:grid-cols-2">
                {/* Traffic Shaping */}
                <Card>
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <Activity className="h-5 w-5" /> Traffic Shaping
                        </CardTitle>
                        <CardDescription>Protocol-specific flood mitigation</CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-6">
                        <div className="flex items-center justify-between">
                            <div className="space-y-0.5">
                                <Label className="text-base">ICMP Rate Limiting</Label>
                                <div className="text-xs text-muted-foreground">Drop ICMP exceeeding 100 pps</div>
                            </div>
                            <Switch checked={rateLimitICMP} onChange={(e) => setRateLimitICMP(e.target.checked)} />
                        </div>
                        <Separator />
                        <div className="flex items-center justify-between">
                            <div className="space-y-0.5">
                                <Label className="text-base">UDP Flood Protection</Label>
                                <div className="text-xs text-muted-foreground">Block non-DNS UDP &gt; 10Mbps</div>
                            </div>
                            <Switch checked={rateLimitUDP} onChange={(e) => setRateLimitUDP(e.target.checked)} />
                        </div>
                        <Separator />
                        <div className="flex items-center justify-between">
                            <div className="space-y-0.5">
                                <Label className="text-base">SYN Cookie Validation</Label>
                                <div className="text-xs text-muted-foreground">Mitigate SYN floods at kernel level</div>
                            </div>
                            <Switch defaultChecked disabled />
                        </div>
                    </CardContent>
                </Card>

                {/* Geo-Blocking */}
                <Card>
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <Lock className="h-5 w-5" /> Geo-Fencing
                        </CardTitle>
                        <CardDescription>Regional traffic policies</CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <div className="flex flex-wrap gap-2">
                            <Badge variant="outline" className="border-red-500/50 bg-red-500/10 text-red-500 cursor-pointer hover:bg-red-500/20">
                                BLOCK North Korea (KP)
                            </Badge>
                            <Badge variant="outline" className="border-red-500/50 bg-red-500/10 text-red-500 cursor-pointer hover:bg-red-500/20">
                                BLOCK Iran (IR)
                            </Badge>
                            <Badge variant="outline" className="border-muted bg-muted/50 cursor-pointer hover:bg-muted">
                                + Add Region to Blocklist
                            </Badge>
                        </div>
                        <div className="p-4 bg-muted/30 rounded-lg text-sm text-muted-foreground">
                            <p>Global Whitelist is active. Only traffic from 42 designated friendly nations is passed without deep inspection.</p>
                        </div>
                    </CardContent>
                    <CardFooter>
                        <Button variant="outline" size="sm" className="w-full">
                            <RotateCcw className="mr-2 h-4 w-4" /> Reset to Defaults
                        </Button>
                    </CardFooter>
                </Card>
            </div>
        </div>
    )
}
