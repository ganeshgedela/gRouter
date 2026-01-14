import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"

const routes = [
    {
        destination: "0.0.0.0/0",
        gateway: "192.168.1.1",
        interface: "eth0",
        protocol: "BGP",
        metric: 20,
        age: "14d 2h"
    },
    {
        destination: "10.0.0.0/16",
        gateway: "0.0.0.0",
        interface: "eth1",
        protocol: "Connected",
        metric: 0,
        age: "45d 12h"
    },
    {
        destination: "10.1.0.0/16",
        gateway: "10.0.5.2",
        interface: "eth2",
        protocol: "OSPF",
        metric: 110,
        age: "5d 4h"
    },
    {
        destination: "192.168.50.0/24",
        gateway: "192.168.1.5",
        interface: "eth0",
        protocol: "BGP",
        metric: 20,
        age: "2h 30m"
    },
    {
        destination: "172.16.0.0/12",
        gateway: "10.0.0.254",
        interface: "eth1",
        protocol: "Static",
        metric: 1,
        age: "120d"
    }
]

export function RouteTable() {
    return (
        <div className="rounded-md border">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead>Destination</TableHead>
                        <TableHead>Gateway</TableHead>
                        <TableHead>Interface</TableHead>
                        <TableHead>Protocol</TableHead>
                        <TableHead>Metric</TableHead>
                        <TableHead className="text-right">Age</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {routes.map((route, i) => (
                        <TableRow key={i}>
                            <TableCell className="font-mono font-medium">{route.destination}</TableCell>
                            <TableCell className="font-mono text-muted-foreground">{route.gateway}</TableCell>
                            <TableCell>{route.interface}</TableCell>
                            <TableCell>
                                <Badge variant="outline" className={
                                    route.protocol === "BGP" ? "border-blue-500 text-blue-500" :
                                        route.protocol === "Connected" ? "border-emerald-500 text-emerald-500" :
                                            route.protocol === "Static" ? "border-purple-500 text-purple-500" :
                                                "border-orange-500 text-orange-500"
                                }>
                                    {route.protocol}
                                </Badge>
                            </TableCell>
                            <TableCell>{route.metric}</TableCell>
                            <TableCell className="text-right text-muted-foreground">{route.age}</TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </div>
    )
}
