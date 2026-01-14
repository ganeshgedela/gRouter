import { useState } from "react"
import { Filter, X } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from "@/components/ui/popover"
import { Checkbox } from "@/components/ui/checkbox"
import { Label } from "@/components/ui/label"

export interface DashboardFiltersState {
    regions: string[]
    statuses: string[]
    tenants: string[]
}

interface DashboardFiltersProps {
    value: DashboardFiltersState
    onChange: (filters: DashboardFiltersState) => void
    availableRegions?: string[]
    availableTenants?: string[]
}

const defaultRegions = ['us-west-1', 'us-east-2', 'eu-west-1', 'ap-southeast-1']
const defaultStatuses = ['online', 'degraded', 'offline']

export function DashboardFilters({
    value,
    onChange,
    availableRegions = defaultRegions,
    availableTenants = [],
}: DashboardFiltersProps) {
    const [isOpen, setIsOpen] = useState(false)

    const activeFilterCount =
        value.regions.length + value.statuses.length + value.tenants.length

    const toggleRegion = (region: string) => {
        const newRegions = value.regions.includes(region)
            ? value.regions.filter(r => r !== region)
            : [...value.regions, region]
        onChange({ ...value, regions: newRegions })
    }

    const toggleStatus = (status: string) => {
        const newStatuses = value.statuses.includes(status)
            ? value.statuses.filter(s => s !== status)
            : [...value.statuses, status]
        onChange({ ...value, statuses: newStatuses })
    }

    const clearAll = () => {
        onChange({
            regions: [],
            statuses: [],
            tenants: [],
        })
    }

    const hasActiveFilters = activeFilterCount > 0

    return (
        <Popover open={isOpen} onOpenChange={setIsOpen}>
            <PopoverTrigger asChild>
                <Button variant="outline" className="relative">
                    <Filter className="mr-2 h-4 w-4" />
                    Filters
                    {hasActiveFilters && (
                        <Badge className="ml-2 h-5 w-5 p-0 flex items-center justify-center" variant="default">
                            {activeFilterCount}
                        </Badge>
                    )}
                </Button>
            </PopoverTrigger>
            <PopoverContent className="w-[320px]" align="end">
                <div className="space-y-4">
                    {/* Header */}
                    <div className="flex items-center justify-between">
                        <h4 className="font-medium text-sm">Filters</h4>
                        {hasActiveFilters && (
                            <Button
                                variant="ghost"
                                size="sm"
                                className="h-auto p-0 text-xs"
                                onClick={clearAll}
                            >
                                Clear all
                            </Button>
                        )}
                    </div>

                    {/* Regions */}
                    <div className="space-y-2">
                        <Label className="text-xs text-muted-foreground">Regions</Label>
                        <div className="space-y-2">
                            {availableRegions.map((region) => (
                                <div key={region} className="flex items-center space-x-2">
                                    <Checkbox
                                        id={`region-${region}`}
                                        checked={value.regions.includes(region)}
                                        onCheckedChange={() => toggleRegion(region)}
                                    />
                                    <label
                                        htmlFor={`region-${region}`}
                                        className="text-sm cursor-pointer"
                                    >
                                        {region}
                                    </label>
                                </div>
                            ))}
                        </div>
                    </div>

                    {/* Status */}
                    <div className="space-y-2">
                        <Label className="text-xs text-muted-foreground">Status</Label>
                        <div className="space-y-2">
                            {defaultStatuses.map((status) => (
                                <div key={status} className="flex items-center space-x-2">
                                    <Checkbox
                                        id={`status-${status}`}
                                        checked={value.statuses.includes(status)}
                                        onCheckedChange={() => toggleStatus(status)}
                                    />
                                    <label
                                        htmlFor={`status-${status}`}
                                        className="text-sm cursor-pointer capitalize"
                                    >
                                        {status}
                                    </label>
                                </div>
                            ))}
                        </div>
                    </div>

                    {/* Active Filters Summary */}
                    {hasActiveFilters && (
                        <div className="pt-2 border-t">
                            <div className="flex flex-wrap gap-1">
                                {value.regions.map((region) => (
                                    <Badge key={region} variant="secondary" className="text-xs">
                                        {region}
                                        <button
                                            onClick={() => toggleRegion(region)}
                                            className="ml-1 hover:bg-muted-foreground/20 rounded-full p-0.5"
                                        >
                                            <X className="h-2 w-2" />
                                        </button>
                                    </Badge>
                                ))}
                                {value.statuses.map((status) => (
                                    <Badge key={status} variant="secondary" className="text-xs capitalize">
                                        {status}
                                        <button
                                            onClick={() => toggleStatus(status)}
                                            className="ml-1 hover:bg-muted-foreground/20 rounded-full p-0.5"
                                        >
                                            <X className="h-2 w-2" />
                                        </button>
                                    </Badge>
                                ))}
                            </div>
                        </div>
                    )}
                </div>
            </PopoverContent>
        </Popover>
    )
}
