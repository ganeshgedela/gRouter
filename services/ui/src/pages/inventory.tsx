import { useState } from "react"
import {
    HardDrive,
    Cpu,
    MemoryStick,
    Network,
    Upload,
    Filter,
    Search
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import {
    Tabs,
    TabsContent,
    TabsList,
    TabsTrigger,
} from "@/components/ui/tabs"

const hardwareInventory = [
    {
        id: "HW-001",
        hostname: "router-us-west-1",
        model: "HPE ProLiant DL360 Gen10",
        cpu: "Intel Xeon Gold 6230 @ 2.10GHz",
        cores: 40,
        memory: "256 GB DDR4",
        storage: "2x 960GB NVMe SSD",
        network: "4x 25GbE SFP28",
        firmware: "BIOS: U32 v2.44",
        location: "US-WEST-1 DC02",
        status: "online"
    },
    {
        id: "HW-002",
        hostname: "router-us-east-2",
        model: "Dell PowerEdge R640",
        cpu: "Intel Xeon Silver 4214 @ 2.20GHz",
        cores: 24,
        memory: "128 GB DDR4",
        storage: "4x 480GB SATA SSD",
        network: "2x 10GbE RJ45",
        firmware: "BIOS: 2.15.1",
        location: "US-EAST-2 DC01",
        status: "online"
    },
    {
        id: "HW-003",
        hostname: "router-eu-west-1",
        model: "HPE ProLiant DL380 Gen10",
        cpu: "AMD EPYC 7502 @ 2.50GHz",
        cores: 64,
        memory: "512 GB DDR4",
        storage: "8x 1.92TB NVMe SSD",
        network: "2x 100GbE QSFP28",
        firmware: "BIOS: U30 v2.62",
        location: "EU-WEST-1 DC01",
        status: "degraded"
    }
]

export default function InventoryPage() {
    const [searchTerm, setSearchTerm] = useState("")

    const filteredInventory = hardwareInventory.filter(item =>
        item.hostname.toLowerCase().includes(searchTerm.toLowerCase()) ||
        item.model.toLowerCase().includes(searchTerm.toLowerCase())
    )

    return (
        <div className="flex-1 space-y-4 p-8 pt-6">
            <div className="flex items-center justify-between">
                <div>
                    <h2 className="text-3xl font-bold tracking-tight">Hardware Inventory</h2>
                    <p className="text-muted-foreground">
                        Physical infrastructure specifications and firmware versions.
                    </p>
                </div>
                <div className="flex gap-2">
                    <Button variant="outline">
                        <Filter className="mr-2 h-4 w-4" />
                        Filters
                    </Button>
                    <Button>
                        <Upload className="mr-2 h-4 w-4" />
                        Export Inventory
                    </Button>
                </div>
            </div>

            {/* Summary Cards */}
            <div className="grid gap-4 md:grid-cols-4">
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Total Devices</CardTitle>
                        <HardDrive className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{hardwareInventory.length}</div>
                        <p className="text-xs text-muted-foreground">Across 4 datacenters</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Total Cores</CardTitle>
                        <Cpu className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">128</div>
                        <p className="text-xs text-muted-foreground">Physical CPU cores</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Total Memory</CardTitle>
                        <MemoryStick className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">896 GB</div>
                        <p className="text-xs text-muted-foreground">DDR4 ECC RAM</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Network Capacity</CardTitle>
                        <Network className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">334 Gbps</div>
                        <p className="text-xs text-muted-foreground">Aggregate bandwidth</p>
                    </CardContent>
                </Card>
            </div>

            <Tabs defaultValue="hardware" className="space-y-4">
                <TabsList>
                    <TabsTrigger value="hardware">Hardware Specs</TabsTrigger>
                    <TabsTrigger value="firmware">Firmware Matrix</TabsTrigger>
                    <TabsTrigger value="network">Network Interfaces</TabsTrigger>
                </TabsList>

                <TabsContent value="hardware" className="space-y-4">
                    <Card>
                        <CardHeader>
                            <div className="flex items-center justify-between">
                                <div>
                                    <CardTitle>Physical Infrastructure</CardTitle>
                                    <CardDescription>Detailed hardware specifications per router</CardDescription>
                                </div>
                                <div className="relative w-64">
                                    <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
                                    <Input
                                        placeholder="Search hostname or model..."
                                        className="pl-8"
                                        value={searchTerm}
                                        onChange={(e) => setSearchTerm(e.target.value)}
                                    />
                                </div>
                            </div>
                        </CardHeader>
                        <CardContent>
                            <div className="space-y-4">
                                {filteredInventory.map((item) => (
                                    <Card key={item.id} className="border-l-4 border-l-primary/20">
                                        <CardHeader>
                                            <div className="flex items-center justify-between">
                                                <div className="flex items-center gap-3">
                                                    <div className="h-12 w-12 rounded-lg bg-primary/10 flex items-center justify-center">
                                                        <HardDrive className="h-6 w-6 text-primary" />
                                                    </div>
                                                    <div>
                                                        <CardTitle className="text-base">{item.hostname}</CardTitle>
                                                        <p className="text-xs text-muted-foreground">{item.model}</p>
                                                    </div>
                                                </div>
                                                <Badge variant={item.status === "online" ? "default" : "destructive"}>
                                                    {item.status}
                                                </Badge>
                                            </div>
                                        </CardHeader>
                                        <CardContent>
                                            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                                                <div>
                                                    <p className="text-xs text-muted-foreground flex items-center gap-1">
                                                        <Cpu className="h-3 w-3" /> CPU
                                                    </p>
                                                    <p className="font-medium mt-1">{item.cpu}</p>
                                                    <p className="text-xs text-muted-foreground">{item.cores} cores</p>
                                                </div>
                                                <div>
                                                    <p className="text-xs text-muted-foreground flex items-center gap-1">
                                                        <MemoryStick className="h-3 w-3" /> Memory
                                                    </p>
                                                    <p className="font-medium mt-1">{item.memory}</p>
                                                </div>
                                                <div>
                                                    <p className="text-xs text-muted-foreground flex items-center gap-1">
                                                        <HardDrive className="h-3 w-3" /> Storage
                                                    </p>
                                                    <p className="font-medium mt-1">{item.storage}</p>
                                                </div>
                                                <div>
                                                    <p className="text-xs text-muted-foreground flex items-center gap-1">
                                                        <Network className="h-3 w-3" /> Network
                                                    </p>
                                                    <p className="font-medium mt-1">{item.network}</p>
                                                </div>
                                            </div>
                                            <div className="mt-4 pt-4 border-t flex items-center justify-between text-xs">
                                                <span className="text-muted-foreground">{item.location}</span>
                                                <span className="font-mono text-muted-foreground">{item.firmware}</span>
                                            </div>
                                        </CardContent>
                                    </Card>
                                ))}
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>

                <TabsContent value="firmware" className="space-y-4">
                    <Card>
                        <CardHeader>
                            <CardTitle>Firmware Version Matrix</CardTitle>
                            <CardDescription>BIOS and firmware versions across the fleet</CardDescription>
                        </CardHeader>
                        <CardContent>
                            <Table>
                                <TableHeader>
                                    <TableRow>
                                        <TableHead>Hostname</TableHead>
                                        <TableHead>BIOS Version</TableHead>
                                        <TableHead>NIC Firmware</TableHead>
                                        <TableHead>BMC Version</TableHead>
                                        <TableHead>Status</TableHead>
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    {hardwareInventory.map((item) => (
                                        <TableRow key={item.id}>
                                            <TableCell className="font-medium">{item.hostname}</TableCell>
                                            <TableCell className="font-mono text-xs">{item.firmware}</TableCell>
                                            <TableCell className="font-mono text-xs">v19.5.12</TableCell>
                                            <TableCell className="font-mono text-xs">v2.88</TableCell>
                                            <TableCell>
                                                <Badge variant="outline" className="bg-emerald-500/10 text-emerald-500">
                                                    Up to date
                                                </Badge>
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </CardContent>
                    </Card>
                </TabsContent>

                <TabsContent value="network" className="space-y-4">
                    <div className="bg-muted/30 border border-dashed rounded-lg p-12 text-center">
                        <Network className="mx-auto h-12 w-12 text-muted-foreground opacity-50" />
                        <h3 className="mt-4 text-lg font-semibold">Network Interface Details</h3>
                        <p className="text-sm text-muted-foreground mt-2 max-w-sm mx-auto">
                            Physical NIC topology, MAC addresses, and link status coming soon.
                        </p>
                    </div>
                </TabsContent>
            </Tabs>
        </div>
    )
}
