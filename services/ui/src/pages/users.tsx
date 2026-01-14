import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Input } from "@/components/ui/input"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { AddUserDialog } from "@/components/users/add-user-dialog"
import { Search, MoreHorizontal, Shield } from "lucide-react"
import { useState } from "react"

const users = [
    {
        id: "usr_1",
        name: "Alice Johnson",
        email: "alice@polaris.network",
        role: "Admin",
        status: "Active",
        lastActive: "2 mins ago",
    },
    {
        id: "usr_2",
        name: "Bob Smith",
        email: "bob@polaris.network",
        role: "Editor",
        status: "Active",
        lastActive: "1 hour ago",
    },
    {
        id: "usr_3",
        name: "Charlie Davis",
        email: "charlie@polaris.network",
        role: "Viewer",
        status: "Offline",
        lastActive: "2 days ago",
    },
]

export default function Users() {
    const [searchTerm, setSearchTerm] = useState("")

    const filteredUsers = users.filter(user =>
        user.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        user.email.toLowerCase().includes(searchTerm.toLowerCase())
    )

    return (
        <div className="flex-1 space-y-4 p-8 pt-6">
            <div className="flex items-center justify-between space-y-2">
                <div>
                    <h2 className="text-3xl font-bold tracking-tight">Users</h2>
                    <p className="text-muted-foreground">
                        Manage access and permissions for the platform.
                    </p>
                </div>
                <div className="flex items-center space-x-2">
                    <AddUserDialog />
                </div>
            </div>

            <div className="flex items-center space-x-2 py-4">
                <div className="relative w-64">
                    <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
                    <Input
                        placeholder="Search users..."
                        className="pl-8"
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                    />
                </div>
            </div>
            <div className="rounded-md border">
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead>User</TableHead>
                            <TableHead>Role</TableHead>
                            <TableHead>Status</TableHead>
                            <TableHead>Last Active</TableHead>
                            <TableHead className="text-right">Actions</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {filteredUsers.map((user) => (
                            <TableRow key={user.id}>
                                <TableCell className="font-medium">
                                    <div className="flex items-center gap-3">
                                        <Avatar className="h-9 w-9">
                                            <AvatarImage src={`https://ui.shadcn.com/avatars/0${user.id.split('_')[1]}.png`} alt={user.name} />
                                            <AvatarFallback>{user.name.slice(0, 2).toUpperCase()}</AvatarFallback>
                                        </Avatar>
                                        <div>
                                            <div className="font-semibold">{user.name}</div>
                                            <div className="text-xs text-muted-foreground">{user.email}</div>
                                        </div>
                                    </div>
                                </TableCell>
                                <TableCell>
                                    <Badge variant="outline" className="flex w-fit items-center gap-1">
                                        <Shield className="h-3 w-3" />
                                        {user.role}
                                    </Badge>
                                </TableCell>
                                <TableCell>
                                    <Badge
                                        variant={user.status === "Active" ? "default" : "secondary"}
                                    >
                                        {user.status}
                                    </Badge>
                                </TableCell>
                                <TableCell className="text-muted-foreground">{user.lastActive}</TableCell>
                                <TableCell className="text-right">
                                    <DropdownMenu>
                                        <DropdownMenuTrigger asChild>
                                            <Button variant="ghost" className="h-8 w-8 p-0">
                                                <span className="sr-only">Open menu</span>
                                                <MoreHorizontal className="h-4 w-4" />
                                            </Button>
                                        </DropdownMenuTrigger>
                                        <DropdownMenuContent align="end">
                                            <DropdownMenuLabel>Actions</DropdownMenuLabel>
                                            <DropdownMenuItem>View Profile</DropdownMenuItem>
                                            <DropdownMenuItem>Edit Permissions</DropdownMenuItem>
                                            <DropdownMenuSeparator />
                                            <DropdownMenuItem className="text-destructive">
                                                Deactivate
                                            </DropdownMenuItem>
                                        </DropdownMenuContent>
                                    </DropdownMenu>
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </div>
        </div>
    )
}
