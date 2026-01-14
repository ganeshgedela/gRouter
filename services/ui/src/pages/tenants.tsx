import { useState } from "react"
import { Link } from "react-router-dom"
import { Search, LayoutGrid, List, Filter } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
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
    Avatar,
    AvatarFallback,
    AvatarImage,
} from "@/components/ui/avatar"
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuCheckboxItem,
    DropdownMenuTrigger
} from "@/components/ui/dropdown-menu"
import { TenantCard } from "@/components/tenants/tenant-card"
import { AddTenantDialog } from "@/components/tenants/add-tenant-dialog"

// Mock Data
const tenantsData = [
    {
        id: "TEN-001",
        name: "Acme Corp",
        logo: "https://ui.shadcn.com/avatars/01.png",
        status: "Active",
        plan: "Enterprise",
        region: "US East (N. Virginia)",
        joined: "Oct 24, 2024",
    },
    {
        id: "TEN-002",
        name: "Globex Corporation",
        logo: "https://ui.shadcn.com/avatars/02.png",
        status: "Active",
        plan: "Pro",
        region: "EU West (London)",
        joined: "Nov 12, 2024",
    },
    {
        id: "TEN-003",
        name: "Soylent Corp",
        logo: "https://ui.shadcn.com/avatars/03.png",
        status: "Suspended",
        plan: "Basic",
        region: "US West (Oregon)",
        joined: "Dec 05, 2024",
    },
    {
        id: "TEN-004",
        name: "Umbrella Corp",
        logo: "https://ui.shadcn.com/avatars/04.png",
        status: "Active",
        plan: "Enterprise",
        region: "AP Southeast (Singapore)",
        joined: "Jan 15, 2025",
    },
    {
        id: "TEN-005",
        name: "Stark Industries",
        logo: "https://ui.shadcn.com/avatars/05.png",
        status: "Active",
        plan: "Enterprise",
        region: "US East (N. Virginia)",
        joined: "Feb 02, 2025",
    },
]

export default function Tenants() {
    const [searchTerm, setSearchTerm] = useState("")
    const [viewMode, setViewMode] = useState<"list" | "grid">("grid")
    const [statusFilter, setStatusFilter] = useState<string[]>(["Active", "Suspended"])

    const filteredTenants = tenantsData.filter((tenant) => {
        const matchesSearch = tenant.name.toLowerCase().includes(searchTerm.toLowerCase())
        const matchesStatus = statusFilter.includes(tenant.status)
        return matchesSearch && matchesStatus
    })

    const toggleStatusFilter = (status: string) => {
        setStatusFilter(prev =>
            prev.includes(status)
                ? prev.filter(s => s !== status)
                : [...prev, status]
        )
    }

    return (
        <div className="flex-1 space-y-4 p-8 pt-6">
            <div className="flex items-center justify-between space-y-2">
                <div>
                    <h2 className="text-3xl font-bold tracking-tight">Tenants</h2>
                    <p className="text-muted-foreground">
                        Manage your organization's tenants and their access permissions.
                    </p>
                </div>
                <div className="flex items-center space-x-2">
                    <AddTenantDialog />
                </div>
            </div>
            <div className="space-y-4">
                <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                        <div className="relative w-64">
                            <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
                            <Input
                                placeholder="Search tenants..."
                                className="pl-8"
                                value={searchTerm}
                                onChange={(e) => setSearchTerm(e.target.value)}
                            />
                        </div>
                        {/* Status Filter */}
                        <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                                <Button variant="outline" size="sm" className="h-10 border-dashed">
                                    <Filter className="mr-2 h-4 w-4" />
                                    Filter
                                </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="start">
                                <DropdownMenuCheckboxItem
                                    checked={statusFilter.includes("Active")}
                                    onCheckedChange={() => toggleStatusFilter("Active")}
                                >
                                    Active
                                </DropdownMenuCheckboxItem>
                                <DropdownMenuCheckboxItem
                                    checked={statusFilter.includes("Suspended")}
                                    onCheckedChange={() => toggleStatusFilter("Suspended")}
                                >
                                    Suspended
                                </DropdownMenuCheckboxItem>
                            </DropdownMenuContent>
                        </DropdownMenu>
                    </div>

                    {/* View Toggle */}
                    <div className="flex items-center space-x-2 bg-muted/50 p-1 rounded-md border">
                        <Button
                            variant={viewMode === "list" ? "secondary" : "ghost"}
                            size="sm"
                            className="h-8 w-8 p-0"
                            onClick={() => setViewMode("list")}
                        >
                            <List className="h-4 w-4" />
                        </Button>
                        <Button
                            variant={viewMode === "grid" ? "secondary" : "ghost"}
                            size="sm"
                            className="h-8 w-8 p-0"
                            onClick={() => setViewMode("grid")}
                        >
                            <LayoutGrid className="h-4 w-4" />
                        </Button>
                    </div>
                </div>

                {viewMode === "list" ? (
                    <div className="rounded-md border">
                        <Table>
                            <TableHeader>
                                <TableRow>
                                    <TableHead className="w-[250px]">Tenant</TableHead>
                                    <TableHead>Status</TableHead>
                                    <TableHead>Plan</TableHead>
                                    <TableHead>Region</TableHead>
                                    <TableHead>Joined</TableHead>
                                    <TableHead className="text-right">Actions</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {filteredTenants.length > 0 ? (
                                    filteredTenants.map((tenant) => (
                                        <TableRow key={tenant.id}>
                                            <TableCell className="font-medium">
                                                <div className="flex items-center space-x-3">
                                                    <Avatar>
                                                        <AvatarImage src={tenant.logo} />
                                                        <AvatarFallback>{tenant.name.slice(0, 2).toUpperCase()}</AvatarFallback>
                                                    </Avatar>
                                                    <div className="flex flex-col">
                                                        <Link to={`/tenants/${tenant.id}`} className="hover:underline hover:text-primary transition-colors">
                                                            {tenant.name}
                                                        </Link>
                                                        <span className="text-xs text-muted-foreground">{tenant.id}</span>
                                                    </div>
                                                </div>
                                            </TableCell>
                                            <TableCell>
                                                <Badge
                                                    variant={
                                                        tenant.status === "Active" ? "default" : "destructive"
                                                    }
                                                >
                                                    {tenant.status}
                                                </Badge>
                                            </TableCell>
                                            <TableCell>{tenant.plan}</TableCell>
                                            <TableCell>{tenant.region}</TableCell>
                                            <TableCell>{tenant.joined}</TableCell>
                                            <TableCell className="text-right">
                                                <Button variant="ghost" size="sm">
                                                    Edit
                                                </Button>
                                            </TableCell>
                                        </TableRow>
                                    ))
                                ) : (
                                    <TableRow>
                                        <TableCell
                                            colSpan={6}
                                            className="h-24 text-center"
                                        >
                                            No results.
                                        </TableCell>
                                    </TableRow>
                                )}
                            </TableBody>
                        </Table>
                    </div>
                ) : (
                    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                        {filteredTenants.length > 0 ? (
                            filteredTenants.map((tenant) => (
                                <TenantCard key={tenant.id} tenant={tenant} />
                            ))
                        ) : (
                            <div className="col-span-full h-24 text-center flex items-center justify-center text-muted-foreground border border-dashed rounded-md">
                                No results found.
                            </div>
                        )}
                    </div>
                )}
            </div>
        </div>
    )
}
