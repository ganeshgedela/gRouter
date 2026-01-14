import { useState, useEffect } from "react"
import {
    Activity,
    AlertTriangle,
    Download,
    RefreshCw,
    Router,
    Users,
    FileDown,
    FileSpreadsheet,
    ImageDown,
    TrendingUp,
    Database,
    Server,
    HardDrive
} from "lucide-react"

import { Button } from "@/components/ui/button"
import { Switch } from "@/components/ui/switch"
import { Progress } from "@/components/ui/progress"
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
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
} from "@/components/ui/dropdown-menu"
import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
} from "@/components/ui/tooltip"

import { Overview } from "@/components/dashboard/overview"
import { RecentActivity } from "@/components/dashboard/recent-activity"
import { RouterHealth } from "@/components/dashboard/router-health"
import { MapView } from "@/components/dashboard/map-view"
import { ReportsView } from "@/components/dashboard/reports-view"
import { PredictiveTraffic } from "@/components/dashboard/predictive-traffic"
import { AlertFeed } from "@/components/dashboard/alert-feed"
import { DashboardSkeleton } from "@/components/skeletons/dashboard-skeleton"
import { HealthScore } from "@/components/dashboard/health-score"
import { QuickActions } from "@/components/dashboard/quick-actions"
import { RecentEvents } from "@/components/dashboard/recent-events"
import { ShortcutHelp } from "@/components/dashboard/shortcut-help"
import { useDashboardShortcuts } from "@/hooks/use-keyboard-shortcuts"
import { useToast } from "@/components/ui/use-toast"

const stats = [
    {
        title: "Total Routers",
        value: "128",
        change: "+4 from last week",
        icon: Router,
        trend: "up"
    },
    {
        title: "Active Tenants",
        value: "24",
        change: "+2 new this month",
        icon: Users,
        trend: "up"
    },
    {
        title: "Critical Alerts",
        value: "3",
        change: "-2 from yesterday",
        icon: AlertTriangle,
        color: "text-destructive",
        trend: "down", // actually good that it's down, but structurally 'down'
        descriptionColor: "text-emerald-500" // positive context
    },
    {
        title: "System Health",
        value: "98.2%",
        change: "All systems operational",
        icon: Activity,
        color: "text-emerald-500",
        trend: "neutral"
    }
]

const topConsumers = [
    { name: "Tenant-Alpha-DB", usage: 78, value: "2.4 GB/s", icon: Database },
    { name: "Content-CDN-West", usage: 64, value: "1.8 GB/s", icon: Server },
    { name: "Backup-Service-01", usage: 42, value: "1.1 GB/s", icon: HardDrive },
]

export default function Dashboard() {
    const [isLoading, setIsLoading] = useState(false)
    const [showShortcutHelp, setShowShortcutHelp] = useState(false)
    const [activeTab, setActiveTab] = useState('overview')
    const [lastRefresh, setLastRefresh] = useState(new Date())
    const [isLive, setIsLive] = useState(false)
    const { toast } = useToast()

    const hour = new Date().getHours();
    const greeting = hour < 12 ? "Good morning" : hour < 18 ? "Good afternoon" : "Good evening";

    // Live mode effect
    useEffect(() => {
        if (!isLive) return
        const interval = setInterval(() => {
            setLastRefresh(new Date())
            // Ideally trigger a refetch here
        }, 5000)
        return () => clearInterval(interval)
    }, [isLive])

    // Handle refresh
    const handleRefresh = () => {
        setIsLoading(true)
        setLastRefresh(new Date())
        // Simulate fetch
        setTimeout(() => {
            setIsLoading(false)
            toast({
                title: "Dashboard Refreshed",
                description: "Latest metrics have been pulled successfully.",
            })
        }, 1500)
    }

    // Handle exports
    const handleExport = (format: 'pdf' | 'csv' | 'png') => {
        toast({
            title: "Export Started",
            description: `Generating ${format.toUpperCase()} report...`,
        })
        // TODO: Implement actual export logic
    }

    // Keyboard shortcuts
    useDashboardShortcuts({
        onRefresh: handleRefresh,
        onToggleFilters: () => { }, // Filters are already inline
        onSwitchTab: (tab) => {
            const tabs = ['overview', 'analytics', 'reports']
            setActiveTab(tabs[tab])
        },
        onShowHelp: () => setShowShortcutHelp(true),
    })

    if (isLoading) {
        return <DashboardSkeleton />
    }

    return (
        <div className="flex-1 space-y-4 p-8 pt-6 relative overflow-hidden animate-in fade-in duration-500">
            {/* Hero Background Gradient */}
            <div className="absolute top-0 left-0 w-full h-[300px] bg-gradient-to-br from-primary/5 via-background to-background -z-10" />

            <div className="flex flex-col space-y-4 md:flex-row md:items-start md:justify-between md:space-y-0 pb-6">
                <div className="space-y-1">
                    <h2 className="text-3xl font-bold tracking-tight bg-gradient-to-r from-foreground to-foreground/70 bg-clip-text text-transparent">
                        {greeting}, Administrator
                    </h2>
                    <p className="text-muted-foreground text-lg">
                        Polaris is orchestrating <span className="font-semibold text-foreground">128</span> active nodes across <span className="font-semibold text-foreground">4</span> regions.
                    </p>
                </div>

                {/* Enhanced Top-Right Controls */}
                <div className="flex flex-wrap items-center gap-4">
                    <div className="flex items-center gap-2 bg-muted/50 p-1.5 rounded-lg border">
                        <Switch
                            id="live-mode"
                            checked={isLive}
                            onCheckedChange={setIsLive}
                            className="scale-90"
                        />
                        <label
                            htmlFor="live-mode"
                            className="text-xs font-medium cursor-pointer select-none flex items-center gap-1.5"
                        >
                            <span className={`relative flex h-2 w-2 mr-0.5`}>
                                {isLive && (
                                    <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                                )}
                                <span className={`relative inline-flex rounded-full h-2 w-2 ${isLive ? 'bg-emerald-500' : 'bg-slate-400'}`}></span>
                            </span>
                            Live Updates
                        </label>
                    </div>

                    <TooltipProvider>
                        {/* Export Menu */}
                        <DropdownMenu>
                            <Tooltip>
                                <TooltipTrigger asChild>
                                    <DropdownMenuTrigger asChild>
                                        <Button variant="outline" size="icon" className="hover:bg-muted/50 transition-colors">
                                            <Download className="h-4 w-4" />
                                        </Button>
                                    </DropdownMenuTrigger>
                                </TooltipTrigger>
                                <TooltipContent>Export Dashboard</TooltipContent>
                            </Tooltip>
                            <DropdownMenuContent align="end" className="w-48">
                                <DropdownMenuLabel>Export As</DropdownMenuLabel>
                                <DropdownMenuSeparator />
                                <DropdownMenuItem onClick={() => handleExport('pdf')} className="cursor-pointer">
                                    <FileDown className="mr-2 h-4 w-4 text-red-500" />
                                    Export as PDF
                                </DropdownMenuItem>
                                <DropdownMenuItem onClick={() => handleExport('csv')} className="cursor-pointer">
                                    <FileSpreadsheet className="mr-2 h-4 w-4 text-green-500" />
                                    Export as CSV
                                </DropdownMenuItem>
                                <DropdownMenuItem onClick={() => handleExport('png')} className="cursor-pointer">
                                    <ImageDown className="mr-2 h-4 w-4 text-blue-500" />
                                    Export as PNG
                                </DropdownMenuItem>
                            </DropdownMenuContent>
                        </DropdownMenu>

                        {/* Refresh Button */}
                        <Tooltip>
                            <TooltipTrigger asChild>
                                <Button size="icon" onClick={handleRefresh} variant="outline" className="hover:bg-muted/50 transition-colors">
                                    <RefreshCw className={`h-4 w-4 ${isLoading ? "animate-spin" : ""}`} />
                                </Button>
                            </TooltipTrigger>
                            <TooltipContent>
                                <div className="text-xs">
                                    <div>Refresh Dashboard</div>
                                    <div className="text-muted-foreground mt-0.5">
                                        Last updated: {lastRefresh.toLocaleTimeString()}
                                    </div>
                                </div>
                            </TooltipContent>
                        </Tooltip>
                    </TooltipProvider>
                </div>
            </div>

            {/* Soft Shadow Line */}
            <div className="h-px w-full bg-gradient-to-r from-transparent via-border to-transparent -mt-6 mb-6" />

            <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-4">
                <TabsList className="bg-background/50 backdrop-blur-sm border">
                    <TabsTrigger value="overview">Overview</TabsTrigger>
                    <TabsTrigger value="analytics">Analytics</TabsTrigger>
                    <TabsTrigger value="reports">Reports</TabsTrigger>
                </TabsList>

                {/* OVERVIEW TAB */}
                <TabsContent value="overview" className="space-y-4">
                    {/* Stats Row */}
                    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                        {stats.map((stat, i) => (
                            <Card key={i} className="bg-card/50 backdrop-blur-sm hover:shadow-md transition-all duration-300 border-primary/10">
                                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                    <CardTitle className="text-sm font-medium">
                                        {stat.title}
                                    </CardTitle>
                                    <stat.icon className={`h-4 w-4 ${stat.color || "text-muted-foreground"}`} />
                                </CardHeader>
                                <CardContent>
                                    <div className={`text-2xl font-bold ${stat.color}`}>{stat.value}</div>
                                    <p className={`text-xs ${stat.descriptionColor || "text-muted-foreground"}`}>
                                        {stat.change}
                                    </p>
                                </CardContent>
                            </Card>
                        ))}
                    </div>

                    {/* Main Dashboard Grid */}
                    <div className="grid gap-4 grid-cols-1 md:grid-cols-7">
                        {/* Left Column - Health + Quick Actions */}
                        <div className="col-span-1 md:col-span-2 space-y-4">
                            <Card className="h-fit bg-card/50 backdrop-blur-sm">
                                <CardHeader>
                                    <CardTitle>Network Health</CardTitle>
                                </CardHeader>
                                <CardContent>
                                    <HealthScore score={92} label="Current Status" />
                                </CardContent>
                            </Card>
                            <QuickActions />
                        </div>

                        {/* Center - Traffic Chart & Top Consumers */}
                        <div className="col-span-1 md:col-span-3 space-y-4">
                            <Card className="bg-card/50 backdrop-blur-sm border-primary/10">
                                <CardHeader>
                                    <CardTitle className="flex items-center gap-2">
                                        <TrendingUp className="h-4 w-4 text-primary" />
                                        Traffic Overview
                                    </CardTitle>
                                    <CardDescription>
                                        Global bandwidth consumption (Ingress vs Egress)
                                    </CardDescription>
                                </CardHeader>
                                <CardContent className="pl-2">
                                    <Overview />
                                </CardContent>
                            </Card>

                            {/* Top Consumers Widget */}
                            <Card className="bg-card/50 backdrop-blur-sm border-primary/10">
                                <CardHeader className="pb-2">
                                    <CardTitle className="text-sm font-medium">Top Bandwidth Consumers</CardTitle>
                                </CardHeader>
                                <CardContent className="space-y-4">
                                    {topConsumers.map((consumer, i) => (
                                        <div key={i} className="space-y-1">
                                            <div className="flex items-center justify-between text-xs">
                                                <div className="flex items-center gap-2">
                                                    <consumer.icon className="h-3 w-3 text-muted-foreground" />
                                                    <span className="font-medium">{consumer.name}</span>
                                                </div>
                                                <span className="font-mono text-muted-foreground">{consumer.value}</span>
                                            </div>
                                            <Progress value={consumer.usage} className="h-1.5" />
                                        </div>
                                    ))}
                                </CardContent>
                            </Card>
                        </div>

                        {/* Right Column - Recent Events */}
                        <div className="col-span-1 md:col-span-2">
                            <RecentEvents />
                        </div>
                    </div>

                    {/* Router Health and Map */}
                    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
                        <Card className="col-span-4 bg-card/50 backdrop-blur-sm">
                            <CardHeader>
                                <CardTitle>Router Health Distribution</CardTitle>
                                <CardDescription>
                                    Status breakdown across all regions
                                </CardDescription>
                            </CardHeader>
                            <CardContent className="pl-2">
                                <RouterHealth />
                            </CardContent>
                        </Card>
                        <Card className="col-span-3 bg-card/50 backdrop-blur-sm">
                            <CardHeader>
                                <CardTitle>Recent Activity</CardTitle>
                                <CardDescription>
                                    Latest events from your network
                                </CardDescription>
                            </CardHeader>
                            <CardContent>
                                <RecentActivity />
                            </CardContent>
                        </Card>
                    </div>

                    {/* Global Map */}
                    <Card className="col-span-full bg-card/50 backdrop-blur-sm">
                        <CardHeader>
                            <CardTitle>Global Network Map</CardTitle>
                            <CardDescription>
                                Real-time visualization of your infrastructure
                            </CardDescription>
                        </CardHeader>
                        <CardContent>
                            <MapView />
                        </CardContent>
                    </Card>
                </TabsContent>

                {/* ANALYTICS TAB */}
                <TabsContent value="analytics" className="space-y-4">
                    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
                        <Card className="col-span-4 bg-card/50 backdrop-blur-sm">
                            <CardHeader>
                                <CardTitle>Predictive Traffic Analysis</CardTitle>
                                <CardDescription>
                                    AI-powered forecasting for the next 7 days
                                </CardDescription>
                            </CardHeader>
                            <CardContent className="pl-2">
                                <PredictiveTraffic />
                            </CardContent>
                        </Card>
                        <Card className="col-span-3 bg-card/50 backdrop-blur-sm">
                            <CardHeader>
                                <CardTitle>Alert Feed</CardTitle>
                                <CardDescription>
                                    Real-time system notifications
                                </CardDescription>
                            </CardHeader>
                            <CardContent>
                                <AlertFeed />
                            </CardContent>
                        </Card>
                    </div>
                </TabsContent>

                {/* REPORTS TAB */}
                <TabsContent value="reports" className="space-y-4">
                    <ReportsView />
                </TabsContent>
            </Tabs>

            {/* Keyboard Shortcuts Help */}
            <ShortcutHelp open={showShortcutHelp} onOpenChange={setShowShortcutHelp} />
        </div>
    )
}
