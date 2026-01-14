import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Skeleton } from "@/components/ui/skeleton"

export function MetricsSkeleton() {
    return (
        <div className="flex-1 space-y-4 p-8 pt-6">
            {/* Header Skeleton */}
            <div className="flex flex-col space-y-2 md:flex-row md:items-center md:justify-between md:space-y-0">
                <div>
                    <Skeleton className="h-8 w-48" />
                    <Skeleton className="h-4 w-72 mt-2" />
                </div>
                <div className="flex items-center gap-2">
                    <Skeleton className="h-9 w-32" />
                    <Skeleton className="h-9 w-32" />
                    <Skeleton className="h-4 w-[1px]" />
                    <Skeleton className="h-9 w-32" />
                    <Skeleton className="h-9 w-24" />
                </div>
            </div>

            {/* Key Indicators Skeleton */}
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                {[1, 2, 3, 4].map((i) => (
                    <Card key={i}>
                        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                            <Skeleton className="h-4 w-24" />
                            <Skeleton className="h-4 w-4" />
                        </CardHeader>
                        <CardContent>
                            <Skeleton className="h-8 w-20 mb-2" />
                            <Skeleton className="h-3 w-32" />
                        </CardContent>
                    </Card>
                ))}
            </div>

            {/* Charts Grid Skeleton */}
            <div className="grid gap-4 md:grid-cols-2">
                {/* Two charts side by side */}
                {[1, 2].map((i) => (
                    <Card key={i} className="col-span-1">
                        <CardHeader>
                            <Skeleton className="h-6 w-40" />
                            <Skeleton className="h-4 w-64 mt-2" />
                        </CardHeader>
                        <CardContent className="pl-0">
                            <Skeleton className="h-[300px] w-full" />
                        </CardContent>
                    </Card>
                ))}

                {/* Full width chart */}
                <Card className="col-span-2">
                    <CardHeader>
                        <Skeleton className="h-6 w-48" />
                        <Skeleton className="h-4 w-72 mt-2" />
                    </CardHeader>
                    <CardContent className="pl-0">
                        <Skeleton className="h-[350px] w-full" />
                    </CardContent>
                </Card>

                {/* Bottom two charts */}
                {[1, 2].map((i) => (
                    <Card key={`bottom-${i}`} className="col-span-2 md:col-span-1">
                        <CardHeader>
                            <Skeleton className="h-6 w-36" />
                            <Skeleton className="h-4 w-56 mt-2" />
                        </CardHeader>
                        <CardContent className="pl-0">
                            <Skeleton className="h-[250px] w-full" />
                        </CardContent>
                    </Card>
                ))}
            </div>
        </div>
    )
}
