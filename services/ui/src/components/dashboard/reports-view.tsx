import {
    FileText,
    Download,
    BarChart3,
    PieChart,
    Calendar,
    ArrowUpRight,
    Activity,
    AlertTriangle,
    Plus
} from "lucide-react"

import { Button } from "@/components/ui/button"
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import { useToast } from "@/components/ui/use-toast"


const REPORT_TEMPLATES = [
    {
        title: "Monthly Bandwidth Usage",
        description: "Aggregated ingress/egress data for all tenants.",
        date: "Generated Feb 01, 2024",
        size: "2.4 MB",
        icon: BarChart3
    },
    {
        title: "Router Uptime Summary",
        description: "Availability metrics for all active nodes.",
        date: "Generated Feb 01, 2024",
        size: "856 KB",
        icon: ActivityIcon
    },
    {
        title: "Security Incident Report",
        description: "Detailed log of all firewall triggers and blocks.",
        date: "Generated Jan 31, 2024",
        size: "1.2 MB",
        icon: AlertTriangleIcon
    },
    {
        title: "Tenant Billing Export",
        description: "Resource consumption breakdown for invoicing.",
        date: "Generated Feb 01, 2024",
        size: "4.1 MB",
        icon: PieChart
    },
]

function ActivityIcon(props: any) { return <div className="h-full w-full bg-emerald-500/10 rounded-md p-2 text-emerald-500"><Calendar {...props} /></div> }
function AlertTriangleIcon(props: any) { return <div className="h-full w-full bg-orange-500/10 rounded-md p-2 text-orange-500"><FileText {...props} /></div> }

export function ReportsView() {
    return (
        <div className="space-y-4">
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                {REPORT_TEMPLATES.map((report, i) => (
                    <Card key={i} className="group hover:border-primary/50 transition-colors">
                        <CardHeader className="flex flex-row items-start justify-between space-y-0 pb-2">
                            <CardTitle className="text-base font-medium">
                                {report.title}
                            </CardTitle>
                            <div className="h-8 w-8">
                                <report.icon className="h-full w-full p-1" />
                            </div>
                        </CardHeader>
                        <CardContent>
                            <CardDescription className="line-clamp-2 min-h-[40px]">
                                {report.description}
                            </CardDescription>
                            <div className="mt-4 flex items-center justify-between text-xs text-muted-foreground">
                                <span>{report.date}</span>
                                <span>{report.size}</span>
                            </div>
                            <Button className="w-full mt-4 gap-2" variant="outline">
                                <Download className="h-4 w-4" />
                                Download PDF
                            </Button>
                        </CardContent>
                    </Card>
                ))}
            </div>

            <Card>
                <CardHeader>
                    <CardTitle>Scheduled Reports</CardTitle>
                    <CardDescription>
                        Automated recurring report generation settings.
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <div className="space-y-4">
                        {[
                            { name: "Weekly Executive Summary", frequency: "Every Monday at 9:00 AM", recipients: "exec-team@polaris.net" },
                            { name: "Daily Operations Log", frequency: "Daily at 00:00 UTC", recipients: "ops-team@polaris.net" }
                        ].map((item, i) => (
                            <div key={i} className="flex items-center justify-between border-b last:border-0 pb-4 last:pb-0">
                                <div className="space-y-1">
                                    <p className="font-medium">{item.name}</p>
                                    <p className="text-sm text-muted-foreground">
                                        {item.frequency} • <span className="text-xs bg-muted px-2 py-1 rounded">{item.recipients}</span>
                                    </p>
                                </div>
                                <Button variant="ghost" size="icon">
                                    <ArrowUpRight className="h-4 w-4" />
                                </Button>
                            </div>
                        ))}
                    </div>
                </CardContent>
            </Card>
        </div>
    )
}
