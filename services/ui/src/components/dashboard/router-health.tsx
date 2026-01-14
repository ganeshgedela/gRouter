export function RouterHealth() {
    const data = [
        { name: "Healthy", value: 400, color: "#10b981" },
        { name: "Degraded", value: 30, color: "#f59e0b" },
        { name: "Down", value: 10, color: "#ef4444" },
    ]

    const total = data.reduce((acc, curr) => acc + curr.value, 0)
    let cumulativePercent = 0

    const getCoordinatesForPercent = (percent: number) => {
        const x = Math.cos(2 * Math.PI * percent)
        const y = Math.sin(2 * Math.PI * percent)
        return [x, y]
    }

    const slices = data.map(slice => {
        const startPercent = cumulativePercent
        const slicePercent = slice.value / total
        cumulativePercent += slicePercent
        const endPercent = cumulativePercent

        // SVG coordinates start from 3 o'clock position (0 radians)
        // We want to start from 12 o'clock (-PI/2), so we rotate -90deg in CSS or adjust coordinates
        // Adjusting coordinates: x = cos(theta), y = sin(theta) where theta is standard angle
        // -0.25 offset to start at top

        const [startX, startY] = getCoordinatesForPercent(startPercent - 0.25)
        const [endX, endY] = getCoordinatesForPercent(endPercent - 0.25)

        const largeArcFlag = slicePercent > 0.5 ? 1 : 0

        const pathData = [
            `M ${startX} ${startY}`,
            `A 1 1 0 ${largeArcFlag} 1 ${endX} ${endY}`,
            `L 0 0`,
        ].join(' ')

        return { pathData, color: slice.color }
    })

    return (
        <div className="h-[300px] w-full flex items-center justify-center relative">
            <svg viewBox="-1 -1 2 2" className="h-full w-full -rotate-90">
                {slices.map((slice, i) => (
                    <path key={i} d={slice.pathData} fill={slice.color} stroke="var(--background)" strokeWidth="0.05" />
                ))}
                {/* Inner Circle for Donut Effect */}
                <circle cx="0" cy="0" r="0.7" fill="hsl(var(--card))" />
            </svg>

            <div className="absolute inset-0 flex flex-col items-center justify-center pointer-events-none">
                <span className="text-3xl font-bold">{total}</span>
                <span className="text-sm text-muted-foreground">Total Routers</span>
            </div>
        </div>
    )
}
