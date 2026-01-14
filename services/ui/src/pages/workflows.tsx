import { useState } from "react"
import {
    Workflow,
    Play,
    Pause,
    Clock,
    CheckCircle2,
    AlertCircle,
    Code,
    Zap,
    Plus,
    TrendingUp
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import {
    Tabs,
    TabsContent,
    TabsList,
    TabsTrigger,
} from "@/components/ui/tabs"

const workflowTemplates = [
    {
        id: "WF-001",
        name: "Drain and Reboot",
        description: "Gracefully drain traffic before system reboot",
        category: "Maintenance",
        executions: 142,
        avgDuration: "4m 30s",
        successRate: 99.3,
        icon: TrendingUp
    },
    {
        id: "WF-002",
        name: "BGP Failover Test",
        description: "Simulate BGP peer failure for resilience testing",
        category: "Testing",
        executions: 45,
        avgDuration: "2m 15s",
        successRate: 100,
        icon: Zap
    },
    {
        id: "WF-003",
        name: "Auto-Scale Routes",
        description: "Automatically adjust routing based on traffic patterns",
        category: "Automation",
        executions: 1203,
        avgDuration: "30s",
        successRate: 98.1,
        icon: TrendingUp
    },
    {
        id: "WF-004",
        name: "Config Backup",
        description: "Backup router configuration to S3",
        category: "Backup",
        executions: 2456,
        avgDuration: "15s",
        successRate: 100,
        icon: CheckCircle2
    }
]

const recentExecutions = [
    {
        id: "EXE-2024-156",
        workflow: "Drain and Reboot",
        router: "router-us-west-1",
        status: "success",
        duration: "4m 22s",
        timestamp: "2024-01-09 18:30:00",
        user: "automation"
    },
    {
        id: "EXE-2024-155",
        workflow: "Config Backup",
        router: "All Routers",
        status: "success",
        duration: "12s",
        timestamp: "2024-01-09 18:00:00",
        user: "cron"
    },
    {
        id: "EXE-2024-154",
        workflow: "BGP Failover Test",
        router: "router-eu-west-1",
        status: "failed",
        duration: "1m 05s",
        timestamp: "2024-01-09 16:45:00",
        user: "admin@polaris.io"
    }
]

export default function WorkflowsPage() {
    const [selectedTemplate, setSelectedTemplate] = useState<string | null>(null)

    return (
        <div className="flex-1 space-y-4 p-8 pt-6">
            <div className="flex items-center justify-between">
                <div>
                    <h2 className="text-3xl font-bold tracking-tight">Workflow Automation</h2>
                    <p className="text-muted-foreground">
                        Pre-built templates and custom automation for operational tasks.
                    </p>
                </div>
                <div className="flex gap-2">
                    <Button variant="outline">
                        <Code className="mr-2 h-4 w-4" />
                        Custom Workflow
                    </Button>
                    <Button>
                        <Plus className="mr-2 h-4 w-4" />
                        Create from Template
                    </Button>
                </div>
            </div>

            {/* Stats */}
            <div className="grid gap-4 md:grid-cols-4">
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Active Workflows</CardTitle>
                        <Workflow className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">12</div>
                        <p className="text-xs text-muted-foreground">4 templates, 8 custom</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Total Executions</CardTitle>
                        <Play className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">3,846</div>
                        <p className="text-xs text-muted-foreground">Last 30 days</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Success Rate</CardTitle>
                        <CheckCircle2 className="h-4 w-4 text-emerald-500" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold text-emerald-500">99.2%</div>
                        <p className="text-xs text-muted-foreground">32 failed runs</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Avg Duration</CardTitle>
                        <Clock className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">1m 45s</div>
                        <p className="text-xs text-muted-foreground">-8s from last week</p>
                    </CardContent>
                </Card>
            </div>

            <Tabs defaultValue="templates" className="space-y-4">
                <TabsList>
                    <TabsTrigger value="templates">Templates</TabsTrigger>
                    <TabsTrigger value="history">Execution History</TabsTrigger>
                    <TabsTrigger value="custom">Custom Workflows</TabsTrigger>
                </TabsList>

                <TabsContent value="templates" className="space-y-4">
                    <div className="grid gap-4 md:grid-cols-2">
                        {workflowTemplates.map((template) => (
                            <Card key={template.id} className="hover:border-primary/50 transition-colors cursor-pointer">
                                <CardHeader>
                                    <div className="flex items-start justify-between">
                                        <div className="flex items-center gap-3">
                                            <div className="h-12 w-12 rounded-lg bg-primary/10 flex items-center justify-center">
                                                <template.icon className="h-6 w-6 text-primary" />
                                            </div>
                                            <div>
                                                <CardTitle className="text-lg">{template.name}</CardTitle>
                                                <Badge variant="outline" className="mt-1">{template.category}</Badge>
                                            </div>
                                        </div>
                                    </div>
                                    <CardDescription className="mt-3">{template.description}</CardDescription>
                                </CardHeader>
                                <CardContent>
                                    <div className="grid grid-cols-3 gap-4 text-sm">
                                        <div>
                                            <p className="text-muted-foreground text-xs">Executions</p>
                                            <p className="font-semibold">{template.executions.toLocaleString()}</p>
                                        </div>
                                        <div>
                                            <p className="text-muted-foreground text-xs">Avg Duration</p>
                                            <p className="font-semibold font-mono text-xs">{template.avgDuration}</p>
                                        </div>
                                        <div>
                                            <p className="text-muted-foreground text-xs">Success Rate</p>
                                            <p className="font-semibold text-emerald-500">{template.successRate}%</p>
                                        </div>
                                    </div>
                                    <div className="flex gap-2 mt-4">
                                        <Button size="sm" className="flex-1">
                                            <Play className="mr-2 h-3 w-3" />
                                            Run Now
                                        </Button>
                                        <Button size="sm" variant="outline">
                                            View Details
                                        </Button>
                                    </div>
                                </CardContent>
                            </Card>
                        ))}
                    </div>
                </TabsContent>

                <TabsContent value="history" className="space-y-4">
                    <Card>
                        <CardHeader>
                            <CardTitle>Recent Executions</CardTitle>
                            <CardDescription>Workflow runs from the last 24 hours</CardDescription>
                        </CardHeader>
                        <CardContent>
                            <div className="space-y-4">
                                {recentExecutions.map((exec) => (
                                    <div key={exec.id} className="flex items-center justify-between p-4 border rounded-lg">
                                        <div className="flex items-center gap-4">
                                            <div className={`h-10 w-10 rounded-lg flex items-center justify-center ${exec.status === 'success' ? 'bg-emerald-500/10' : 'bg-destructive/10'
                                                }`}>
                                                {exec.status === 'success' ? (
                                                    <CheckCircle2 className="h-5 w-5 text-emerald-500" />
                                                ) : (
                                                    <AlertCircle className="h-5 w-5 text-destructive" />
                                                )}
                                            </div>
                                            <div>
                                                <p className="font-semibold">{exec.workflow}</p>
                                                <p className="text-xs text-muted-foreground">
                                                    {exec.router} • {exec.user}
                                                </p>
                                            </div>
                                        </div>
                                        <div className="text-right">
                                            <p className="font-mono text-xs">{exec.duration}</p>
                                            <p className="text-xs text-muted-foreground">{exec.timestamp}</p>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>

                <TabsContent value="custom" className="space-y-4">
                    <div className="bg-muted/30 border border-dashed rounded-lg p-12 text-center">
                        <Code className="mx-auto h-12 w-12 text-muted-foreground opacity-50" />
                        <h3 className="mt-4 text-lg font-semibold">No Custom Workflows</h3>
                        <p className="text-sm text-muted-foreground mt-2 max-w-sm mx-auto">
                            Create custom workflows using YAML or our visual editor.
                        </p>
                        <Button className="mt-6" variant="outline">
                            <Plus className="mr-2 h-4 w-4" />
                            Create Custom Workflow
                        </Button>
                    </div>
                </TabsContent>
            </Tabs>
        </div>
    )
}
