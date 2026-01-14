import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import { Skeleton } from "@/components/ui/skeleton"

export function RouterTableSkeleton({ rows = 5 }: { rows?: number }) {
    return (
        <div className="rounded-md border">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead className="w-[200px]">Hostname</TableHead>
                        <TableHead>Status</TableHead>
                        <TableHead>Region</TableHead>
                        <TableHead>IP Address</TableHead>
                        <TableHead>System Load</TableHead>
                        <TableHead>Uptime</TableHead>
                        <TableHead className="text-right">Actions</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {Array.from({ length: rows }).map((_, i) => (
                        <TableRow key={i}>
                            <TableCell>
                                <div className="space-y-2">
                                    <Skeleton className="h-4 w-32" />
                                    <Skeleton className="h-3 w-16" />
                                </div>
                            </TableCell>
                            <TableCell>
                                <div className="flex items-center gap-2">
                                    <Skeleton className="h-2 w-2 rounded-full" />
                                    <Skeleton className="h-4 w-16" />
                                </div>
                            </TableCell>
                            <TableCell>
                                <Skeleton className="h-4 w-40" />
                            </TableCell>
                            <TableCell>
                                <Skeleton className="h-4 w-24" />
                            </TableCell>
                            <TableCell>
                                <div className="flex items-center gap-2">
                                    <Skeleton className="h-2 w-16 rounded-full" />
                                    <Skeleton className="h-3 w-8" />
                                </div>
                            </TableCell>
                            <TableCell>
                                <Skeleton className="h-4 w-16" />
                            </TableCell>
                            <TableCell className="text-right">
                                <Skeleton className="h-8 w-8 ml-auto" />
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </div>
    )
}
