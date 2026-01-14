import { TrendingUp, Activity, AlertTriangle, Clock, Signal } from "lucide-react"

interface HealthScoreProps {
    score: number // 0-100
    label: string
}

export function HealthScore({ score, label }: HealthScoreProps) {
    const getColor = (score: number) => {
        if (score >= 90) return "text-emerald-500"
        if (score >= 70) return "text-yellow-500"
        return "text-destructive"
    }

    const getGradient = (score: number) => {
        if (score >= 90) return "from-emerald-500 to-emerald-600"
        if (score >= 70) return "from-yellow-500 to-yellow-600"
        return "from-red-500 to-red-600"
    }

    const getStatusText = (score: number) => {
        if (score >= 90) return "Excellent"
        if (score >= 70) return "Degraded"
        return "Critical"
    }

    return (
        <div className="flex flex-col gap-6">
            <div className="flex items-center gap-6">
                {/* Gauge */}
                <div className="relative w-24 h-24 shrink-0">
                    <svg className="w-24 h-24 transform -rotate-90">
                        <circle
                            cx="48"
                            cy="48"
                            r="42"
                            stroke="currentColor"
                            strokeWidth="8"
                            fill="none"
                            className="text-muted/20"
                        />
                        <circle
                            cx="48"
                            cy="48"
                            r="42"
                            stroke="url(#health-gradient)"
                            strokeWidth="8"
                            fill="none"
                            strokeDasharray={`${(score / 100) * 264} 264`}
                            className="transition-all duration-1000 ease-out"
                            strokeLinecap="round"
                        />
                        <defs>
                            <linearGradient id="health-gradient" x1="0%" y1="0%" x2="100%" y2="100%">
                                <stop offset="0%" className={getGradient(score).split(' ')[0].replace('from-', 'stop-')} />
                                <stop offset="100%" className={getGradient(score).split(' ')[1].replace('to-', 'stop-')} />
                            </linearGradient>
                        </defs>
                    </svg>
                    <div className="absolute inset-0 flex items-center justify-center">
                        <span className={`text-2xl font-bold ${getColor(score)}`}>
                            {score}
                        </span>
                    </div>
                </div>

                {/* Main Status Text */}
                <div className="space-y-1">
                    <h3 className="text-2xl font-bold tracking-tight">{getStatusText(score)}</h3>
                    <p className="text-sm text-muted-foreground">{label}</p>
                    <div className={`flex items-center gap-1.5 text-xs font-medium ${getColor(score)} bg-current/10 w-fit px-2 py-0.5 rounded-full`}>
                        <TrendingUp className="h-3 w-3" />
                        <span>+2.4% vs last week</span>
                    </div>
                </div>
            </div>

            {/* Metrics Grid */}
            <div className="grid grid-cols-2 gap-3">
                <div className="flex flex-col gap-1 p-3 rounded-lg bg-muted/30 border border-border/50">
                    <div className="flex items-center gap-2 text-muted-foreground text-xs">
                        <Clock className="h-3.5 w-3.5" />
                        <span>Uptime</span>
                    </div>
                    <div className="font-mono font-semibold text-lg text-foreground">99.99%</div>
                </div>

                <div className="flex flex-col gap-1 p-3 rounded-lg bg-muted/30 border border-border/50">
                    <div className="flex items-center gap-2 text-muted-foreground text-xs">
                        <Activity className="h-3.5 w-3.5" />
                        <span>Latency</span>
                    </div>
                    <div className="font-mono font-semibold text-lg text-emerald-500">24ms</div>
                </div>

                <div className="flex flex-col gap-1 p-3 rounded-lg bg-muted/30 border border-border/50">
                    <div className="flex items-center gap-2 text-muted-foreground text-xs">
                        <Signal className="h-3.5 w-3.5" />
                        <span>Loss</span>
                    </div>
                    <div className="font-mono font-semibold text-lg text-foreground">0.00%</div>
                </div>

                <div className="flex flex-col gap-1 p-3 rounded-lg bg-muted/30 border border-border/50">
                    <div className="flex items-center gap-2 text-muted-foreground text-xs">
                        <AlertTriangle className="h-3.5 w-3.5" />
                        <span>Alerts</span>
                    </div>
                    <div className="font-mono font-semibold text-lg text-amber-500">2 Active</div>
                </div>
            </div>
        </div>
    )
}
