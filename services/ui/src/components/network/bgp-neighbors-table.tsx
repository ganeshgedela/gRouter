import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import { ArrowUpRight, ArrowDownLeft, Clock } from "lucide-react"

const bgpNeighbors = [
    {
        neighborIp: "192.168.1.5",
        asn: "65001",
        state: "Established",
        uptime: "14d 2h 15m",
        receivedRoutes: 450,
        sentRoutes: 12,
        desc: "US-West-ISP-A"
    },
    {
        neighborIp: "192.168.2.9",
        asn: "65002",
        state: "Established",
        uptime: "5d 10h 30m",
        receivedRoutes: 120,
        sentRoutes: 12,
        desc: "US-East-ISP-B"
    },
    {
        neighborIp: "10.0.5.2",
        asn: "64512",
        state: "Connect",
        uptime: "0s",
        receivedRoutes: 0,
        sentRoutes: 0,
        desc: "Internal-Mesh-GW"
    }
]

export function BGPNeighborsTable() {
    return (
        <div className="rounded-md border">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead>Neighbor</TableHead>
                        <TableHead>ASN</TableHead>
                        <TableHead>State</TableHead>
                        <TableHead>Routes (Rx/Tx)</TableHead>
                        <TableHead>Uptime</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {bgpNeighbors.map((peer, i) => (
                        <TableRow key={i}>
                            <TableCell>
                                <div className="flex flex-col">
                                    <span className="font-medium font-mono">{peer.neighborIp}</span>
                                    <span className="text-xs text-muted-foreground">{peer.desc}</span>
                                </div>
                            </TableCell>
                            <TableCell className="font-mono">{peer.asn}</TableCell>
                            <TableCell>
                                <Badge variant={peer.state === "Established" ? "default" : "destructive"}>
                                    {peer.state}
                                </Badge>
                            </TableCell>
                            <TableCell>
                                <div className="flex items-center gap-4 text-xs">
                                    <span className="flex items-center gap-1 text-emerald-500">
                                        <ArrowDownLeft className="h-3 w-3" /> {peer.receivedRoutes}
                                    </span>
                                    <span className="flex items-center gap-1 text-blue-500">
                                        <ArrowUpRight className="h-3 w-3" /> {peer.sentRoutes}
                                    </span>
                                </div>
                            </TableCell>
                            <TableCell>
                                <div className="flex items-center gap-2 text-muted-foreground">
                                    <Clock className="h-3 w-3" />
                                    {peer.uptime}
                                </div>
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </div>
    )
}
