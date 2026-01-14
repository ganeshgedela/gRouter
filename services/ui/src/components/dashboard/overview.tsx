import { Area, AreaChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts"

const data = [
    { name: "00:00", ingress: 400, egress: 240 },
    { name: "04:00", ingress: 300, egress: 139 },
    { name: "08:00", ingress: 200, egress: 980 },
    { name: "12:00", ingress: 278, egress: 390 },
    { name: "16:00", ingress: 189, egress: 480 },
    { name: "20:00", ingress: 239, egress: 380 },
    { name: "23:59", ingress: 349, egress: 430 },
]

export function Overview() {
    return (
        <ResponsiveContainer width="100%" height={350}>
            <AreaChart data={data}>
                <defs>
                    <linearGradient id="colorIngress" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="5%" stopColor="#8884d8" stopOpacity={0.8} />
                        <stop offset="95%" stopColor="#8884d8" stopOpacity={0} />
                    </linearGradient>
                    <linearGradient id="colorEgress" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="5%" stopColor="#82ca9d" stopOpacity={0.8} />
                        <stop offset="95%" stopColor="#82ca9d" stopOpacity={0} />
                    </linearGradient>
                </defs>
                <XAxis
                    dataKey="name"
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
                    tickFormatter={(value) => `${value}G`}
                />
                <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#333" />
                <Tooltip
                    contentStyle={{ backgroundColor: "hsl(var(--popover))", borderRadius: "8px", border: "1px solid hsl(var(--border))" }}
                    labelStyle={{ color: "hsl(var(--popover-foreground))", fontWeight: "bold" }}
                />
                <Area
                    type="monotone"
                    dataKey="ingress"
                    stroke="#8884d8"
                    fillOpacity={1}
                    fill="url(#colorIngress)"
                />
                <Area
                    type="monotone"
                    dataKey="egress"
                    stroke="#82ca9d"
                    fillOpacity={1}
                    fill="url(#colorEgress)"
                />
            </AreaChart>
        </ResponsiveContainer>
    )
}
