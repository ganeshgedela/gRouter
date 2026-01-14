import { useState } from "react"
import {
    Search,
    User,
    Settings,
    Shield
} from "lucide-react"

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

const AUDIT_LOGS = [
    { id: 'evt-1029', timestamp: '2025-05-15 14:32:01', user: 'admin@polaris.io', action: 'UPDATE_POLICY', resource: 'acl-us-west-01', details: 'Added allow rule for port 443' },
    { id: 'evt-1028', timestamp: '2025-05-15 14:30:55', user: 'system', action: 'AUTO_SCALE', resource: 'gcp-eu-west-2', details: 'Triggered scale out due to high load' },
    { id: 'evt-1027', timestamp: '2025-05-15 12:15:22', user: 'j.doe@company.com', action: 'DELETE_ROUTE', resource: 'bgp-peer-1', details: 'Removed stale BGP neighbor' },
    { id: 'evt-1026', timestamp: '2025-05-15 11:05:00', user: 'admin@polaris.io', action: 'LOGIN_SUCCESS', resource: '-', details: 'Web console login from 192.168.1.5' },
    { id: 'evt-1025', timestamp: '2025-05-15 09:45:10', user: 'system', action: 'THREAT_BLOCK', resource: 'firewall-edge', details: 'Blocked DDoS pattern from 45.33.xx.xx' },
    { id: 'evt-1024', timestamp: '2025-05-14 18:20:05', user: 'audit-bot', action: 'CONFIG_BACKUP', resource: 'daily-backup', details: 'Snapshot created successfully' },
]

export function AuditLog() {
    const [filter, setFilter] = useState('')

    const filteredLogs = AUDIT_LOGS.filter(log =>
        log.user.toLowerCase().includes(filter.toLowerCase()) ||
        log.action.toLowerCase().includes(filter.toLowerCase()) ||
        log.resource.toLowerCase().includes(filter.toLowerCase())
    )

    return (
        <div className="space-y-4">
            <div className="flex items-center space-x-2">
                <div className="relative flex-1 max-w-sm">
                    <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
                    <Input
                        placeholder="Filter by user, action, or resource..."
                        className="pl-8"
                        value={filter}
                        onChange={(e) => setFilter(e.target.value)}
                    />
                </div>
            </div>

            <div className="rounded-md border bg-background/50">
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead className="w-[180px]">Timestamp</TableHead>
                            <TableHead className="w-[200px]">User</TableHead>
                            <TableHead className="w-[150px]">Action</TableHead>
                            <TableHead className="w-[180px]">Resource</TableHead>
                            <TableHead>Details</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {filteredLogs.map((log) => (
                            <TableRow key={log.id}>
                                <TableCell className="font-mono text-xs text-muted-foreground">
                                    {log.timestamp}
                                </TableCell>
                                <TableCell>
                                    <div className="flex items-center gap-2">
                                        <div className="h-6 w-6 rounded-full bg-primary/10 flex items-center justify-center">
                                            {log.user === 'system' ? <Settings className="h-3 w-3" /> : <User className="h-3 w-3" />}
                                        </div>
                                        <span className="text-sm font-medium">{log.user}</span>
                                    </div>
                                </TableCell>
                                <TableCell>
                                    <Badge variant="outline" className="font-mono text-[10px]">
                                        {log.action}
                                    </Badge>
                                </TableCell>
                                <TableCell className="font-mono text-xs">
                                    {log.resource}
                                </TableCell>
                                <TableCell className="text-sm text-muted-foreground">
                                    <div className="flex items-center gap-1">
                                        {log.action.includes('BLOCK') && <Shield className="h-3 w-3 text-red-500 mr-1" />}
                                        {log.details}
                                    </div>
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </div>
            <div className="text-xs text-muted-foreground text-center">
                Showing {filteredLogs.length} of {AUDIT_LOGS.length} events. Logs are immutable and retained for 365 days.
            </div>
        </div>
    )
}
