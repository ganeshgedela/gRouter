import { useState } from "react"
import {
    Trash2,
    Save,
    Search,
    MoreHorizontal
} from "lucide-react"
import { AddPolicyDialog } from "@/components/policy/add-policy-dialog"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger
} from "@/components/ui/dropdown-menu"
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow
} from "@/components/ui/table"
import { cn } from "@/lib/utils"

type PolicyRule = {
    id: string
    priority: number
    name: string
    type: 'ingress' | 'egress'
    action: 'allow' | 'deny'
    protocol: 'tcp' | 'udp' | 'icmp' | 'any'
    source: string
    destination: string
    status: 'active' | 'inactive'
}

const INITIAL_RULES: PolicyRule[] = [
    { id: 'rule-100', priority: 1, name: 'Allow SSH Management', type: 'ingress', action: 'allow', protocol: 'tcp', source: '10.0.0.0/8', destination: '0.0.0.0/0', status: 'active' },
    { id: 'rule-200', priority: 2, name: 'Block Malicious IPs', type: 'ingress', action: 'deny', protocol: 'any', source: '192.168.1.100/32', destination: '0.0.0.0/0', status: 'active' },
    { id: 'rule-300', priority: 3, name: 'Web Traffic', type: 'ingress', action: 'allow', protocol: 'tcp', source: '0.0.0.0/0', destination: 'self', status: 'active' },
    { id: 'rule-999', priority: 999, name: 'Default Deny Info', type: 'ingress', action: 'deny', protocol: 'any', source: '0.0.0.0/0', destination: '0.0.0.0/0', status: 'inactive' },
]

export function PolicyManager() {
    const [rules, setRules] = useState<PolicyRule[]>(INITIAL_RULES)
    const [searchTerm, setSearchTerm] = useState('')

    const handleAdd = (newRule: PolicyRule) => {
        setRules([...rules, newRule].sort((a, b) => a.priority - b.priority))
    }

    const handleDelete = (id: string) => {
        setRules(rules.filter(r => r.id !== id))
    }

    const handleStatusToggle = (id: string) => {
        setRules(rules.map(r =>
            r.id === id ? { ...r, status: r.status === 'active' ? 'inactive' : 'active' } : r
        ))
    }

    const filteredRules = rules.filter(r =>
        r.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        r.source.includes(searchTerm) ||
        r.id.includes(searchTerm)
    )

    return (
        <Card className="flex flex-col h-[700px]">
            <CardHeader className="pb-3">
                <div className="flex items-center justify-between">
                    <div>
                        <CardTitle>Network Policies</CardTitle>
                        <CardDescription>Manage Firewall rules, ACLs, and Traffic Policies.</CardDescription>
                    </div>
                    <div className="flex items-center gap-2">
                        <div className="relative">
                            <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                            <Input
                                placeholder="Search rules..."
                                className="pl-8 w-[250px]"
                                value={searchTerm}
                                onChange={(e) => setSearchTerm(e.target.value)}
                            />
                        </div>
                        <AddPolicyDialog
                            onAdd={handleAdd}
                            existingPriorities={rules.map(r => r.priority)}
                        />
                        <Button variant="outline">
                            <Save className="mr-2 h-4 w-4" /> Deploy
                        </Button>
                    </div>
                </div>
            </CardHeader>
            <CardContent className="flex-1 p-0 overflow-auto">
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead className="w-[80px]">Order</TableHead>
                            <TableHead className="w-[200px]">Name</TableHead>
                            <TableHead className="w-[100px]">Type</TableHead>
                            <TableHead className="w-[100px]">Action</TableHead>
                            <TableHead>Traffic Match</TableHead>
                            <TableHead className="w-[100px]">Status</TableHead>
                            <TableHead className="w-[50px]"></TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {filteredRules.map((rule) => (
                            <TableRow key={rule.id} className="group">
                                <TableCell className="font-mono text-xs">
                                    <div className="flex items-center gap-1">
                                        <Badge variant="outline" className="w-8 justify-center p-0 h-6">
                                            {rule.priority}
                                        </Badge>
                                    </div>
                                </TableCell>
                                <TableCell>
                                    <div className="font-medium text-sm">{rule.name}</div>
                                    <div className="text-xs text-muted-foreground font-mono">{rule.id}</div>
                                </TableCell>
                                <TableCell>
                                    <Badge variant="secondary" className="text-xs uppercase">
                                        {rule.type}
                                    </Badge>
                                </TableCell>
                                <TableCell>
                                    <div className={cn(
                                        "inline-flex items-center rounded-md border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2",
                                        rule.action === 'allow'
                                            ? "border-emerald-500/50 text-emerald-700 bg-emerald-500/10 dark:text-emerald-400"
                                            : "border-red-500/50 text-red-700 bg-red-500/10 dark:text-red-400"
                                    )}>
                                        {rule.action.toUpperCase()}
                                    </div>
                                </TableCell>
                                <TableCell>
                                    <div className="grid grid-cols-[40px_1fr] gap-x-2 gap-y-1 text-xs">
                                        <span className="text-muted-foreground font-mono">PROT</span>
                                        <span className="font-mono">{rule.protocol.toUpperCase()}</span>

                                        <span className="text-muted-foreground font-mono">SRC</span>
                                        <span className="font-mono">{rule.source}</span>

                                        <span className="text-muted-foreground font-mono">DST</span>
                                        <span className="font-mono">{rule.destination}</span>
                                    </div>
                                </TableCell>
                                <TableCell>
                                    <div
                                        className="flex items-center gap-2 cursor-pointer"
                                        onClick={() => handleStatusToggle(rule.id)}
                                    >
                                        <div className={cn(
                                            "h-2.5 w-2.5 rounded-full transition-all",
                                            rule.status === 'active' ? "bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.5)]" : "bg-zinc-300 dark:bg-zinc-700"
                                        )} />
                                        <span className="text-xs text-muted-foreground capitalize">{rule.status}</span>
                                    </div>
                                </TableCell>
                                <TableCell>
                                    <DropdownMenu>
                                        <DropdownMenuTrigger asChild>
                                            <Button variant="ghost" className="h-8 w-8 p-0">
                                                <span className="sr-only">Open menu</span>
                                                <MoreHorizontal className="h-4 w-4" />
                                            </Button>
                                        </DropdownMenuTrigger>
                                        <DropdownMenuContent align="end">
                                            <DropdownMenuItem onClick={() => handleStatusToggle(rule.id)}>
                                                {rule.status === 'active' ? 'Disable Rule' : 'Enable Rule'}
                                            </DropdownMenuItem>
                                            <DropdownMenuItem>Edit Details</DropdownMenuItem>
                                            <DropdownMenuItem
                                                className="text-red-600 focus:text-red-600"
                                                onClick={() => handleDelete(rule.id)}
                                            >
                                                <Trash2 className="mr-2 h-4 w-4" /> Delete Rule
                                            </DropdownMenuItem>
                                        </DropdownMenuContent>
                                    </DropdownMenu>
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </CardContent>
        </Card>
    )
}
