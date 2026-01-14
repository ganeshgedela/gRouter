import {
    LineChart as RechartsLineChart,
    Line,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    ResponsiveContainer,
    AreaChart,
    Area
} from 'recharts'
import { Card } from "@/components/ui/card"

interface DataPoint {
    name: string
    [key: string]: string | number
}

interface CustomLineChartProps {
    data: DataPoint[]
    lines: {
        key: string
        color: string
        name: string
    }[]
    height?: number
    className?: string
    showAxes?: boolean
    showGrid?: boolean
    showTooltip?: boolean
    fillArea?: boolean
}

const CustomTooltip = ({ active, payload, label }: any) => {
    if (active && payload && payload.length) {
        return (
            <Card className="bg-background/95 backdrop-blur border-border p-3 shadow-xl">
                <div className="text-xs font-semibold mb-2 border-b pb-1">{label}</div>
                <div className="space-y-1">
                    {payload.map((entry: any, index: number) => (
                        <div key={index} className="flex items-center gap-3 text-xs">
                            <div className="flex items-center gap-1.5">
                                <div
                                    className="w-2 h-2 rounded-full shadow-[0_0_4px_currentColor]"
                                    style={{ color: entry.color, backgroundColor: entry.color }}
                                />
                                <span className="text-muted-foreground">{entry.name}</span>
                            </div>
                            <span className="font-mono font-medium">{entry.value}</span>
                        </div>
                    ))}
                </div>
            </Card>
        )
    }
    return null
}

export function LineChart({
    data,
    lines,
    height = 300,
    className,
    showAxes = true,
    showGrid = true,
    showTooltip = true,
    fillArea = false
}: CustomLineChartProps) {
    if (fillArea) {
        return (
            <div className={`w-full ${className}`} style={{ height }}>
                <ResponsiveContainer width="100%" height="100%">
                    <AreaChart data={data}>
                        <defs>
                            {lines.map((line) => (
                                <linearGradient key={line.key} id={`gradient-${line.key}`} x1="0" y1="0" x2="0" y2="1">
                                    <stop offset="5%" stopColor={line.color} stopOpacity={0.3} />
                                    <stop offset="95%" stopColor={line.color} stopOpacity={0} />
                                </linearGradient>
                            ))}
                        </defs>
                        {showGrid && (
                            <CartesianGrid
                                strokeDasharray="3 3"
                                stroke="hsl(var(--muted-foreground))"
                                strokeOpacity={0.1}
                                vertical={false}
                            />
                        )}
                        {showAxes && (
                            <XAxis
                                dataKey="name"
                                stroke="hsl(var(--muted-foreground))"
                                fontSize={10}
                                tickLine={false}
                                axisLine={false}
                                minTickGap={30}
                            />
                        )}
                        {showAxes && (
                            <YAxis
                                stroke="hsl(var(--muted-foreground))"
                                fontSize={10}
                                tickLine={false}
                                axisLine={false}
                                width={30}
                            />
                        )}
                        {showTooltip && (
                            <Tooltip content={<CustomTooltip />} cursor={{ stroke: 'hsl(var(--muted-foreground))', strokeWidth: 1, strokeDasharray: '4 4' }} />
                        )}
                        {lines.map((line) => (
                            <Area
                                key={line.key}
                                type="monotone"
                                dataKey={line.key}
                                name={line.name}
                                stroke={line.color}
                                fill={`url(#gradient-${line.key})`}
                                strokeWidth={2}
                                activeDot={{ r: 4, strokeWidth: 0, fill: line.color }}
                            />
                        ))}
                    </AreaChart>
                </ResponsiveContainer>
            </div>
        )
    }

    return (
        <div className={`w-full ${className}`} style={{ height }}>
            <ResponsiveContainer width="100%" height="100%">
                <RechartsLineChart data={data}>
                    {showGrid && (
                        <CartesianGrid
                            strokeDasharray="3 3"
                            stroke="hsl(var(--muted-foreground))"
                            strokeOpacity={0.1}
                            vertical={false}
                        />
                    )}
                    {showAxes && (
                        <XAxis
                            dataKey="name"
                            stroke="hsl(var(--muted-foreground))"
                            fontSize={10}
                            tickLine={false}
                            axisLine={false}
                            minTickGap={30}
                        />
                    )}
                    {showAxes && (
                        <YAxis
                            stroke="hsl(var(--muted-foreground))"
                            fontSize={10}
                            tickLine={false}
                            axisLine={false}
                            width={30}
                        />
                    )}
                    {showTooltip && (
                        <Tooltip content={<CustomTooltip />} cursor={{ stroke: 'hsl(var(--muted-foreground))', strokeWidth: 1, strokeDasharray: '4 4' }} />
                    )}
                    {lines.map((line) => (
                        <Line
                            key={line.key}
                            type="monotone"
                            dataKey={line.key}
                            name={line.name}
                            stroke={line.color}
                            strokeWidth={2}
                            dot={false}
                            activeDot={{ r: 4, strokeWidth: 0, fill: line.color }}
                        />
                    ))}
                </RechartsLineChart>
            </ResponsiveContainer>
        </div>
    )
}
