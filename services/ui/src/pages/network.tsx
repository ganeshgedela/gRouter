import { useState } from "react"
import {
    Search,
    Server,
    Network,
    Activity,
    Terminal,
    Power,
    FileText,
    RefreshCw,
    MoreHorizontal,
    Globe,
    Plus,
    Cloud,
    Layers
} from "lucide-react"
import { Link } from "react-router-dom"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
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
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { InteractiveTopology } from "@/components/topology/interactive-topology"
import { RouterTableSkeleton } from "@/components/skeletons/router-table-skeleton"
import { CardSkeleton } from "@/components/skeletons/card-skeleton"
import { ProvisionRouterWizard } from "@/components/wizards/provision-router-wizard"
import { CreateVPCWizard } from "@/components/wizards/create-vpc-wizard"
import { VPCDetailPanel } from "@/components/network/vpc-detail-panel"
import { SubnetCard } from "@/components/network/subnet-card"
import { BGPNeighborsTable } from "@/components/network/bgp-neighbors-table"
import { RouteTable } from "@/components/network/route-table"
import type { VPC, Subnet } from "@/types/network"
import { cn } from "@/lib/utils"

// Mock Data for Routers
const routersData = [
    {
        id: "RTR-US-W-01",
        hostname: "router-us-west-1",
        ip: "10.0.1.1",
        region: "US West (Oregon)",
        status: "Online",
        load: 45,
        uptime: "14d 2h",
        version: "v2.4.1",
    },
    {
        id: "RTR-US-E-02",
        hostname: "router-us-east-2",
        ip: "10.0.2.1",
        region: "US East (N. Virginia)",
        status: "Online",
        load: 62,
        uptime: "45d 12h",
        version: "v2.4.1",
    },
    {
        id: "RTR-EU-W-01",
        hostname: "router-eu-west-1",
        ip: "10.1.1.1",
        region: "EU West (London)",
        status: "Degraded",
        load: 92,
        uptime: "2d 5h",
        version: "v2.3.9",
    },
    {
        id: "RTR-AP-SE-03",
        hostname: "router-ap-se-3",
        ip: "10.2.1.5",
        region: "AP Southeast (Singapore)",
        status: "Offline",
        load: 0,
        uptime: "-",
        version: "v2.4.0",
    },
    {
        id: "RTR-SA-E-01",
        hostname: "router-sa-east-1",
        ip: "10.3.1.2",
        region: "SA East (Sao Paulo)",
        status: "Online",
        load: 23,
        uptime: "120d 1h",
        version: "v2.4.1",
    },
]

// Mock VPCs
const mockVPCs: VPC[] = [
    {
        id: 'vpc-001',
        name: 'production-vpc',
        description: 'Main production environment',
        cidr: '10.0.0.0/16',
        region: 'us-west-1',
        cloudProvider: 'AWS',
        routerIds: ['RTR-US-W-01', 'RTR-US-E-02'],
        status: 'active',
        createdAt: '2025-11-15T10:00:00Z',
        tags: { environment: 'production', team: 'platform' }
    },
    {
        id: 'vpc-002',
        name: 'staging-vpc',
        description: 'Staging and QA environment',
        cidr: '10.1.0.0/16',
        region: 'us-east-2',
        cloudProvider: 'AWS',
        routerIds: ['RTR-US-E-02'],
        status: 'active',
        createdAt: '2025-12-01T14:30:00Z',
        tags: { environment: 'staging', team: 'platform' }
    },
    {
        id: 'vpc-003',
        name: 'dev-vpc',
        description: 'Development and testing',
        cidr: '10.2.0.0/16',
        region: 'eu-west-1',
        cloudProvider: 'GCP',
        routerIds: ['RTR-EU-W-01'],
        status: 'active',
        createdAt: '2026-01-05T09:00:00Z',
        tags: { environment: 'development', team: 'engineering' }
    }
]

// Mock Subnets
const mockSubnets: Subnet[] = [
    {
        id: 'subnet-001',
        name: 'prod-public-us-west-1a',
        vpcId: 'vpc-001',
        cidr: '10.0.1.0/24',
        availabilityZone: 'us-west-1a',
        routerId: 'RTR-US-W-01',
        type: 'public',
        availableIPs: 251,
        usedIPs: 178,
        status: 'active',
        createdAt: '2025-11-15T10:05:00Z',
        tags: { tier: 'web' }
    },
    {
        id: 'subnet-002',
        name: 'prod-private-us-west-1a',
        vpcId: 'vpc-001',
        cidr: '10.0.2.0/24',
        availabilityZone: 'us-west-1a',
        routerId: 'RTR-US-W-01',
        type: 'private',
        availableIPs: 251,
        usedIPs: 89,
        status: 'active',
        createdAt: '2025-11-15T10:06:00Z',
        tags: { tier: 'application' }
    },
    {
        id: 'subnet-003',
        name: 'prod-public-us-east-2a',
        vpcId: 'vpc-001',
        cidr: '10.0.10.0/24',
        availabilityZone: 'us-east-2a',
        routerId: 'RTR-US-E-02',
        type: 'public',
        availableIPs: 251,
        usedIPs: 156,
        status: 'active',
        createdAt: '2025-11-20T11:00:00Z',
        tags: { tier: 'web' }
    },
    {
        id: 'subnet-004',
        name: 'staging-public-us-east-2a',
        vpcId: 'vpc-002',
        cidr: '10.1.1.0/24',
        availabilityZone: 'us-east-2a',
        routerId: 'RTR-US-E-02',
        type: 'public',
        availableIPs: 251,
        usedIPs: 45,
        status: 'active',
        createdAt: '2025-12-01T14:35:00Z',
        tags: { tier: 'web' }
    },
    {
        id: 'subnet-005',
        name: 'dev-private-eu-west-1a',
        vpcId: 'vpc-003',
        cidr: '10.2.1.0/24',
        availabilityZone: 'eu-west-1a',
        routerId: 'RTR-EU-W-01',
        type: 'private',
        availableIPs: 251,
        usedIPs: 12,
        status: 'active',
        createdAt: '2026-01-05T09:10:00Z',
        tags: { tier: 'application' }
    }
]

export default function NetworkPage() {
    const [isLoading, setIsLoading] = useState(false)
    const [isProvisionWizardOpen, setIsProvisionWizardOpen] = useState(false)
    const [isCreateVPCWizardOpen, setIsCreateVPCWizardOpen] = useState(false)
    const [selectedVPC, setSelectedVPC] = useState<VPC | null>(null)
    const [vpcs, setVpcs] = useState<VPC[]>(mockVPCs)
    const [subnets, setSubnets] = useState<Subnet[]>(mockSubnets)
    const [searchQuery, setSearchQuery] = useState("")

    const filteredRouters = routersData.filter(router =>
        router.hostname.toLowerCase().includes(searchQuery.toLowerCase()) ||
        router.ip.includes(searchQuery)
    )

    const handleRefresh = () => {
        setIsLoading(true)
        setTimeout(() => setIsLoading(false), 2000)
    }

    const handleCreateVPC = (newVPC: Partial<VPC>) => {
        const vpc: VPC = {
            id: `vpc-${Date.now()}`,
            name: newVPC.name || '',
            description: newVPC.description || '',
            cidr: newVPC.cidr || '',
            region: newVPC.region || '',
            cloudProvider: newVPC.cloudProvider || 'AWS',
            routerIds: newVPC.routerIds || [],
            status: 'active',
            createdAt: new Date().toISOString(),
            tags: newVPC.tags || {}
        }
        setVpcs([...vpcs, vpc])
    }

    const handleDeleteVPC = (vpc: VPC) => {
        setVpcs(vpcs.filter(v => v.id !== vpc.id))
        setSubnets(subnets.filter(s => s.vpcId !== vpc.id))
    }

    // Calculate statistics
    const totalIPs = subnets.reduce((sum, subnet) => sum + subnet.availableIPs, 0)
    const usedIPs = subnets.reduce((sum, subnet) => sum + subnet.usedIPs, 0)
    const ipUtilization = totalIPs > 0 ? Math.round((usedIPs / totalIPs) * 100) : 0
    const cloudProviders = new Set(vpcs.map(v => v.cloudProvider)).size

    return (
        <div className="flex-1 space-y-4 p-8 pt-6">
            <div className="flex items-center justify-between space-y-2">
                <div>
                    <h2 className="text-3xl font-bold tracking-tight">Network Infrastructure</h2>
                    <p className="text-muted-foreground">
                        Manage physical endpoints, topology, and routing policies.
                    </p>
                </div>
                <div className="flex items-center space-x-2">
                    <Button onClick={() => setIsProvisionWizardOpen(true)}>
                        <Network className="mr-2 h-4 w-4" /> Provision Router
                    </Button>
                </div>
            </div>

            <ProvisionRouterWizard
                open={isProvisionWizardOpen}
                onOpenChange={setIsProvisionWizardOpen}
            />

            <CreateVPCWizard
                open={isCreateVPCWizardOpen}
                onOpenChange={setIsCreateVPCWizardOpen}
                routers={routersData}
                onCreate={handleCreateVPC}
            />

            <VPCDetailPanel
                vpc={selectedVPC}
                subnets={subnets}
                routers={routersData}
                onClose={() => setSelectedVPC(null)}
                onDelete={handleDeleteVPC}
            />

            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                {isLoading ? (
                    <>
                        <CardSkeleton />
                        <CardSkeleton />
                    </>
                ) : (
                    <>
                        <Card>
                            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                <CardTitle className="text-sm font-medium">Total Routers</CardTitle>
                                <Server className="h-4 w-4 text-muted-foreground" />
                            </CardHeader>
                            <CardContent>
                                <div className="text-2xl font-bold">{routersData.length}</div>
                                <p className="text-xs text-muted-foreground">Across 4 regions</p>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                <CardTitle className="text-sm font-medium">Critical Issues</CardTitle>
                                <Activity className="h-4 w-4 text-muted-foreground" />
                            </CardHeader>
                            <CardContent>
                                <div className="text-2xl font-bold text-destructive">1</div>
                                <p className="text-xs text-muted-foreground">1 Offline, 1 Degraded</p>
                            </CardContent>
                        </Card>
                    </>
                )}
            </div>

            <Tabs defaultValue="inventory" className="space-y-4">
                <TabsList>
                    <TabsTrigger value="inventory">Inventory</TabsTrigger>
                    <TabsTrigger value="topology">Topology</TabsTrigger>
                    <TabsTrigger value="routing">Routing</TabsTrigger>
                    <TabsTrigger value="subnets">Subnets & VPCs</TabsTrigger>
                </TabsList>

                {/* ROUTING TAB */}
                <TabsContent value="routing" className="space-y-4">
                    <div className="grid gap-4 md:grid-cols-2">
                        <Card>
                            <CardHeader>
                                <CardTitle>BGP Neighbors</CardTitle>
                                <CardDescription>Active peering sessions and state.</CardDescription>
                            </CardHeader>
                            <CardContent>
                                <BGPNeighborsTable />
                            </CardContent>
                        </Card>
                        <Card>
                            <CardHeader>
                                <CardTitle>Global Routing Table</CardTitle>
                                <CardDescription>Active routes across the fabric.</CardDescription>
                            </CardHeader>
                            <CardContent>
                                <RouteTable />
                            </CardContent>
                        </Card>
                    </div>
                </TabsContent>



                {/* INVENTORY TAB */}
                <TabsContent value="inventory" className="space-y-4">
                    <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-2">
                            <div className="relative w-64">
                                <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
                                <Input
                                    placeholder="Filter by hostname or IP..."
                                    className="pl-8"
                                    value={searchQuery}
                                    onChange={(e) => setSearchQuery(e.target.value)}
                                />
                            </div>
                        </div>
                        <Button variant="outline" size="sm" onClick={handleRefresh} disabled={isLoading}>
                            <RefreshCw className={`mr-2 h-4 w-4 ${isLoading ? 'animate-spin' : ''}`} /> Refresh Status
                        </Button>
                    </div>

                    {isLoading ? (
                        <RouterTableSkeleton />
                    ) : (
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
                                    {filteredRouters.map((router) => (
                                        <TableRow key={router.id}>
                                            <TableCell className="font-medium">
                                                <div className="flex items-center">
                                                    <Server className="mr-2 h-4 w-4 text-muted-foreground" />
                                                    <Link to={`/routers/${router.id}`} className="hover:underline hover:text-primary transition-colors">
                                                        {router.hostname}
                                                    </Link>
                                                </div>
                                                <span className="text-xs text-muted-foreground ml-6">{router.version}</span>
                                            </TableCell>
                                            <TableCell>
                                                <div className="flex items-center">
                                                    <span className={`flex h-2 w-2 rounded-full mr-2 ${router.status === 'Online' ? 'bg-emerald-500' :
                                                        router.status === 'Degraded' ? 'bg-yellow-500' : 'bg-red-500'
                                                        }`} />
                                                    {router.status}
                                                </div>
                                            </TableCell>
                                            <TableCell>
                                                <div className="flex items-center text-muted-foreground">
                                                    <Globe className="mr-2 h-3 w-3" />
                                                    {router.region}
                                                </div>
                                            </TableCell>
                                            <TableCell className="font-mono text-xs">{router.ip}</TableCell>
                                            <TableCell>
                                                <div className="flex items-center gap-2">
                                                    <div className="w-16 h-2 bg-secondary rounded-full overflow-hidden">
                                                        <div
                                                            className={`h-full ${router.load > 90 ? 'bg-red-500' : router.load > 70 ? 'bg-yellow-500' : 'bg-green-500'}`}
                                                            style={{ width: `${router.load}%` }}
                                                        />
                                                    </div>
                                                    <span className="text-xs text-muted-foreground">{router.load}%</span>
                                                </div>
                                            </TableCell>
                                            <TableCell>{router.uptime}</TableCell>
                                            <TableCell className="text-right">
                                                <DropdownMenu>
                                                    <DropdownMenuTrigger asChild>
                                                        <Button variant="ghost" className="h-8 w-8 p-0">
                                                            <MoreHorizontal className="h-4 w-4" />
                                                        </Button>
                                                    </DropdownMenuTrigger>
                                                    <DropdownMenuContent align="end">
                                                        <DropdownMenuLabel>Management</DropdownMenuLabel>
                                                        <DropdownMenuItem>
                                                            <Terminal className="mr-2 h-4 w-4" /> SSH Console
                                                        </DropdownMenuItem>
                                                        <DropdownMenuItem>
                                                            <FileText className="mr-2 h-4 w-4" /> View Logs
                                                        </DropdownMenuItem>
                                                        <DropdownMenuSeparator />
                                                        <DropdownMenuItem className="text-destructive">
                                                            <Power className="mr-2 h-4 w-4" /> Reboot System
                                                        </DropdownMenuItem>
                                                    </DropdownMenuContent>
                                                </DropdownMenu>
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </div>
                    )}
                </TabsContent>

                {/* TOPOLOGY TAB */}
                <TabsContent value="topology" className="space-y-4">
                    <Card className="h-[600px] flex flex-col">
                        <CardHeader>
                            <CardTitle>Global Mesh Topology</CardTitle>
                            <CardDescription>
                                Visualizing active BGP peers and inter-region links.
                            </CardDescription>
                        </CardHeader>
                        <CardContent className="flex-1 p-0 relative bg-zinc-950/50">
                            {/* Full interactive topo */}
                            <div className="absolute inset-0">
                                <InteractiveTopology />
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>

                {/* SUBNETS & VPCs TAB */}
                <TabsContent value="subnets" className="space-y-4">
                    {/* Summary Cards */}
                    <div className="grid gap-4 md:grid-cols-4">
                        <Card>
                            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                <CardTitle className="text-sm font-medium">Total VPCs</CardTitle>
                                <Layers className="h-4 w-4 text-muted-foreground" />
                            </CardHeader>
                            <CardContent>
                                <div className="text-2xl font-bold">{vpcs.length}</div>
                                <p className="text-xs text-muted-foreground">Across {cloudProviders} cloud provider{cloudProviders !== 1 ? 's' : ''}</p>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                <CardTitle className="text-sm font-medium">Total Subnets</CardTitle>
                                <Network className="h-4 w-4 text-muted-foreground" />
                            </CardHeader>
                            <CardContent>
                                <div className="text-2xl font-bold">{subnets.length}</div>
                                <p className="text-xs text-muted-foreground">
                                    {subnets.filter(s => s.type === 'public').length} public, {subnets.filter(s => s.type === 'private').length} private
                                </p>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                <CardTitle className="text-sm font-medium">IP Utilization</CardTitle>
                                <Activity className="h-4 w-4 text-muted-foreground" />
                            </CardHeader>
                            <CardContent>
                                <div className={cn(
                                    "text-2xl font-bold",
                                    ipUtilization > 90 ? "text-red-500" : ipUtilization > 70 ? "text-yellow-500" : "text-emerald-500"
                                )}>{ipUtilization}%</div>
                                <p className="text-xs text-muted-foreground">{usedIPs.toLocaleString()} / {totalIPs.toLocaleString()} IPs used</p>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                <CardTitle className="text-sm font-medium">Multi-Cloud</CardTitle>
                                <Cloud className="h-4 w-4 text-muted-foreground" />
                            </CardHeader>
                            <CardContent>
                                <div className="text-2xl font-bold">{cloudProviders}</div>
                                <p className="text-xs text-muted-foreground">Active cloud providers</p>
                            </CardContent>
                        </Card>
                    </div>

                    {/* VPCs Section */}
                    <Card>
                        <CardHeader>
                            <div className="flex items-center justify-between">
                                <div>
                                    <CardTitle>Virtual Private Clouds</CardTitle>
                                    <CardDescription>Manage logical network isolation across cloud providers</CardDescription>
                                </div>
                                <Button onClick={() => setIsCreateVPCWizardOpen(true)}>
                                    <Plus className="mr-2 h-4 w-4" />
                                    Create VPC
                                </Button>
                            </div>
                        </CardHeader>
                        <CardContent>
                            {vpcs.length > 0 ? (
                                <div className="space-y-4">
                                    {vpcs.map((vpc) => {
                                        const vpcSubnets = subnets.filter(s => s.vpcId === vpc.id)
                                        const vpcRouters = routersData.filter(r => vpc.routerIds.includes(r.id))

                                        return (
                                            <Card
                                                key={vpc.id}
                                                className="border-l-4 border-l-primary hover:bg-muted/30 transition-colors cursor-pointer"
                                                onClick={() => setSelectedVPC(vpc)}
                                            >
                                                <CardHeader>
                                                    <div className="flex items-start justify-between">
                                                        <div className="flex items-center gap-3">
                                                            <div className="h-12 w-12 rounded-lg bg-primary/10 flex items-center justify-center">
                                                                <Layers className="h-6 w-6 text-primary" />
                                                            </div>
                                                            <div>
                                                                <div className="flex items-center gap-2">
                                                                    <CardTitle className="text-base">{vpc.name}</CardTitle>
                                                                    <Badge variant={vpc.status === 'active' ? 'default' : 'destructive'}>
                                                                        {vpc.status}
                                                                    </Badge>
                                                                </div>
                                                                <p className="text-xs text-muted-foreground mt-1">{vpc.description}</p>
                                                            </div>
                                                        </div>
                                                        <Badge variant="outline">{vpc.cloudProvider}</Badge>
                                                    </div>
                                                </CardHeader>
                                                <CardContent>
                                                    <div className="grid grid-cols-4 gap-4 text-sm">
                                                        <div>
                                                            <p className="text-xs text-muted-foreground mb-1">CIDR Block</p>
                                                            <p className="font-mono font-semibold">{vpc.cidr}</p>
                                                        </div>
                                                        <div>
                                                            <p className="text-xs text-muted-foreground mb-1">Region</p>
                                                            <div className="flex items-center gap-1">
                                                                <Globe className="h-3 w-3 text-muted-foreground" />
                                                                <span className="font-medium">{vpc.region}</span>
                                                            </div>
                                                        </div>
                                                        <div>
                                                            <p className="text-xs text-muted-foreground mb-1">Subnets</p>
                                                            <p className="font-semibold">{vpcSubnets.length}</p>
                                                        </div>
                                                        <div>
                                                            <p className="text-xs text-muted-foreground mb-1">Routers</p>
                                                            <p className="font-semibold">{vpcRouters.length}</p>
                                                        </div>
                                                    </div>

                                                    {vpcSubnets.length > 0 && (
                                                        <div className="mt-4 pt-4 border-t">
                                                            <p className="text-xs text-muted-foreground mb-2">Subnets</p>
                                                            <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
                                                                {vpcSubnets.slice(0, 4).map(subnet => (
                                                                    <div key={subnet.id} className="flex items-center justify-between text-xs p-2 rounded bg-muted/30">
                                                                        <span className="font-mono">{subnet.cidr}</span>
                                                                        <Badge variant="outline" className="text-[10px]">{subnet.type}</Badge>
                                                                    </div>
                                                                ))}
                                                            </div>
                                                            {vpcSubnets.length > 4 && (
                                                                <p className="text-xs text-muted-foreground text-center mt-2">
                                                                    +{vpcSubnets.length - 4} more subnets
                                                                </p>
                                                            )}
                                                        </div>
                                                    )}
                                                </CardContent>
                                            </Card>
                                        )
                                    })}
                                </div>
                            ) : (
                                <div className="bg-muted/30 border border-dashed rounded-lg p-12 text-center">
                                    <Layers className="mx-auto h-12 w-12 text-muted-foreground opacity-50" />
                                    <h3 className="mt-4 text-lg font-semibold">No VPCs Created</h3>
                                    <p className="text-sm text-muted-foreground mt-2 max-w-sm mx-auto">
                                        Create your first Virtual Private Cloud to start managing logical networks and CIDR blocks.
                                    </p>
                                    <Button variant="outline" className="mt-6" onClick={() => setIsCreateVPCWizardOpen(true)}>
                                        <Plus className="mr-2 h-4 w-4" />
                                        Create VPC
                                    </Button>
                                </div>
                            )}
                        </CardContent>
                    </Card>

                    {/* All Subnets View */}
                    {subnets.length > 0 && (
                        <Card>
                            <CardHeader>
                                <CardTitle>All Subnets</CardTitle>
                                <CardDescription>Complete list of subnet allocations</CardDescription>
                            </CardHeader>
                            <CardContent>
                                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                                    {subnets.map(subnet => {
                                        const router = routersData.find(r => r.id === subnet.routerId)
                                        return (
                                            <SubnetCard
                                                key={subnet.id}
                                                subnet={subnet}
                                                routerHostname={router?.hostname}
                                            />
                                        )
                                    })}
                                </div>
                            </CardContent>
                        </Card>
                    )}
                </TabsContent>
            </Tabs>
        </div>
    )
}
