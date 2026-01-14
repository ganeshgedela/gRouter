import {
    BarChart3,
    AlertTriangle,
    Globe,
    Activity
} from "lucide-react"

import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import { PredictiveTraffic } from "@/components/dashboard/predictive-traffic"
import { AlertFeed } from "@/components/dashboard/alert-feed"
import { MapView } from "@/components/dashboard/map-view"

export default function AnalyticsPage() {
    return (
        <div className="flex-1 space-y-4 p-8 pt-6">
            <div className="flex items-center justify-between space-y-2">
                <div>
                    <h2 className="text-3xl font-bold tracking-tight">Analytics & Reports</h2>
                    <p className="text-muted-foreground">
                        Deep dive into network performance, security incidents, and forecasting.
                    </p>
                </div>
            </div>

            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
                <Card className="col-span-4">
                    <CardHeader>
                        <CardTitle>Traffic Forecast (AI)</CardTitle>
                        <CardDescription>
                            Predictive bandwidth modeling for the next 4 hours.
                        </CardDescription>
                    </CardHeader>
                    <CardContent className="pl-2">
                        <PredictiveTraffic />
                    </CardContent>
                </Card>
                <Card className="col-span-3">
                    <CardHeader>
                        <CardTitle>Live Incident Feed</CardTitle>
                        <CardDescription>
                            Real-time anomalies and AI recommendations.
                        </CardDescription>
                    </CardHeader>
                    <CardContent>
                        <AlertFeed />
                    </CardContent>
                </Card>
            </div>

            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
                <Card className="col-span-4">
                    <CardHeader>
                        <CardTitle>Network Topology Heatmap</CardTitle>
                        <CardDescription>
                            Regional load distribution.
                        </CardDescription>
                    </CardHeader>
                    <CardContent>
                        <MapView />
                    </CardContent>
                </Card>
            </div>
        </div>
    )
}
