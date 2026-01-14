"use client"

import {
    Area,
    AreaChart,
    CartesianGrid,
    ResponsiveContainer,
    Tooltip,
    XAxis,
    YAxis,
    ReferenceLine
} from "recharts"

const data = [
    { time: "00:00", actual: 300, predicted: 310 },
    { time: "02:00", actual: 350, predicted: 360 },
    { time: "04:00", actual: 200, predicted: 220 },
    { time: "06:00", actual: 450, predicted: 440 },
    { time: "08:00", actual: 600, predicted: 590 },
    { time: "10:00", actual: 750, predicted: 740 },
    { time: "12:00", actual: 900, predicted: 900 }, // Current
    { time: "14:00", predicted: 1100 },
    { time: "16:00", predicted: 1250 },
    { time: "18:00", predicted: 1100 },
    { time: "20:00", predicted: 950 },
    { time: "22:00", predicted: 800 },
    { time: "24:00", predicted: 600 },
]

export function PredictiveTraffic() {
    return (
        <div className="h-[350px] w-full">
            <ResponsiveContainer width="100%" height="100%">
                <AreaChart
                    data={data}
                    margin={{
                        top: 10,
                        right: 10,
                        left: 0,
                        bottom: 0,
                    }}
                >
                    <defs>
                        <linearGradient id="colorActual" x1="0" y1="0" x2="0" y2="1">
                            <stop offset="5%" stopColor="hsl(var(--primary))" stopOpacity={0.3} />
                            <stop offset="95%" stopColor="hsl(var(--primary))" stopOpacity={0} />
                        </linearGradient>
                        <linearGradient id="colorPredicted" x1="0" y1="0" x2="0" y2="1">
                            <stop offset="5%" stopColor="hsl(var(--muted-foreground))" stopOpacity={0.3} />
                            <stop offset="95%" stopColor="hsl(var(--muted-foreground))" stopOpacity={0} />
                        </linearGradient>
                    </defs>
                    <CartesianGrid strokeDasharray="3 3" className="stroke-muted" vertical={false} />
                    <XAxis
                        dataKey="time"
                        stroke="#888888"
                        fontSize={12}
                        tickLine={false}
                        axisLine={false}
                    />
                    <YAxis
                        stroke="#888888"
                        fontSize={12}
                        tickLine={false}
                        axisLine={false}
                        tickFormatter={(value) => `${value}Mbps`}
                    />
                    <Tooltip
                        contentStyle={{
                            backgroundColor: "hsl(var(--popover))",
                            border: "1px solid hsl(var(--border))",
                            borderRadius: "var(--radius)",
                        }}
                        itemStyle={{ color: "hsl(var(--foreground))" }}
                    />
                    <ReferenceLine x="12:00" stroke="hsl(var(--destructive))" strokeDasharray="3 3" label={{ position: 'top', value: 'Now', fill: 'hsl(var(--destructive))', fontSize: 12 }} />
                    <Area
                        type="monotone"
                        dataKey="predicted"
                        stroke="hsl(var(--muted-foreground))" // Grey for predicted
                        strokeDasharray="5 5"
                        fillOpacity={1}
                        fill="url(#colorPredicted)"
                        name="Predicted Traffic"
                    />
                    <Area
                        type="monotone"
                        dataKey="actual"
                        stroke="hsl(var(--primary))"
                        fillOpacity={1}
                        fill="url(#colorActual)"
                        name="Actual Traffic"
                    />
                </AreaChart>
            </ResponsiveContainer>
        </div>
    )
}
