import { useState } from "react"
import { Calendar, Clock, RefreshCw } from "lucide-react"
import { Button } from "@/components/ui/button"
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from "@/components/ui/popover"
import { Badge } from "@/components/ui/badge"

export type TimeRange = '1h' | '6h' | '24h' | '7d' | '30d' | 'custom'
export type RefreshInterval = 'off' | '30s' | '1m' | '5m'

interface TimeRangeSelectorProps {
    value: TimeRange
    onChange: (range: TimeRange) => void
    refreshInterval?: RefreshInterval
    onRefreshIntervalChange?: (interval: RefreshInterval) => void
}

const timeRanges: { value: TimeRange; label: string }[] = [
    { value: '1h', label: 'Last Hour' },
    { value: '6h', label: 'Last 6 Hours' },
    { value: '24h', label: 'Last 24 Hours' },
    { value: '7d', label: 'Last 7 Days' },
    { value: '30d', label: 'Last 30 Days' },
    { value: 'custom', label: 'Custom Range' },
]

const refreshIntervals: { value: RefreshInterval; label: string }[] = [
    { value: 'off', label: 'Off' },
    { value: '30s', label: '30 seconds' },
    { value: '1m', label: '1 minute' },
    { value: '5m', label: '5 minutes' },
]

export function TimeRangeSelector({
    value,
    onChange,
    refreshInterval = 'off',
    onRefreshIntervalChange,
}: TimeRangeSelectorProps) {
    const [isOpen, setIsOpen] = useState(false)

    const selectedRange = timeRanges.find(r => r.value === value)

    return (
        <div className="flex items-center gap-2">
            {/* Time Range Selector */}
            <Popover open={isOpen} onOpenChange={setIsOpen}>
                <PopoverTrigger asChild>
                    <Button variant="outline" className="w-[180px] justify-start">
                        <Clock className="mr-2 h-4 w-4" />
                        {selectedRange?.label}
                    </Button>
                </PopoverTrigger>
                <PopoverContent className="w-[200px] p-0" align="start">
                    <div className="p-2">
                        <div className="text-xs font-medium text-muted-foreground mb-2 px-2">
                            Quick Select
                        </div>
                        <div className="space-y-1">
                            {timeRanges.map((range) => (
                                <button
                                    key={range.value}
                                    onClick={() => {
                                        onChange(range.value)
                                        setIsOpen(false)
                                    }}
                                    className={`w-full text-left px-2 py-1.5 rounded text-sm hover:bg-accent transition-colors ${value === range.value ? 'bg-accent font-medium' : ''
                                        }`}
                                >
                                    {range.label}
                                </button>
                            ))}
                        </div>
                    </div>
                </PopoverContent>
            </Popover>

            {/* Refresh Interval */}
            {onRefreshIntervalChange && (
                <Popover>
                    <PopoverTrigger asChild>
                        <Button variant="outline" size="icon" className="relative">
                            <RefreshCw className="h-4 w-4" />
                            {refreshInterval !== 'off' && (
                                <span className="absolute -top-1 -right-1 h-3 w-3 rounded-full bg-emerald-500 animate-pulse" />
                            )}
                        </Button>
                    </PopoverTrigger>
                    <PopoverContent className="w-[160px] p-0" align="end">
                        <div className="p-2">
                            <div className="text-xs font-medium text-muted-foreground mb-2 px-2">
                                Auto Refresh
                            </div>
                            <div className="space-y-1">
                                {refreshIntervals.map((interval) => (
                                    <button
                                        key={interval.value}
                                        onClick={() => {
                                            onRefreshIntervalChange(interval.value)
                                        }}
                                        className={`w-full text-left px-2 py-1.5 rounded text-sm hover:bg-accent transition-colors ${refreshInterval === interval.value ? 'bg-accent font-medium' : ''
                                            }`}
                                    >
                                        {interval.label}
                                    </button>
                                ))}
                            </div>
                        </div>
                    </PopoverContent>
                </Popover>
            )}

            {/* Active indicator */}
            {refreshInterval !== 'off' && (
                <Badge variant="outline" className="text-xs">
                    <span className="mr-1 h-1.5 w-1.5 rounded-full bg-emerald-500 animate-pulse" />
                    Live
                </Badge>
            )}
        </div>
    )
}
