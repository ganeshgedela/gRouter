import { MoreHorizontal, Users, Router } from "lucide-react"

import { Button } from "@/components/ui/button"
import {
    Card,
    CardContent,
    CardFooter,
    CardHeader,
} from "@/components/ui/card"
import {
    Avatar,
    AvatarFallback,
    AvatarImage,
} from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"

interface TenantCardProps {
    tenant: {
        id: string
        name: string
        logo: string
        status: string
        plan: string
        region: string
        joined: string
    }
}

export function TenantCard({ tenant }: TenantCardProps) {
    return (
        <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <Badge variant={tenant.status === "Active" ? "default" : "destructive"}>
                    {tenant.status}
                </Badge>
                <Button variant="ghost" size="icon" className="h-8 w-8">
                    <MoreHorizontal className="h-4 w-4" />
                </Button>
            </CardHeader>
            <CardContent className="flex flex-col items-center space-y-4">
                <Avatar className="h-20 w-20">
                    <AvatarImage src={tenant.logo} alt={tenant.name} />
                    <AvatarFallback className="text-xl">{tenant.name.slice(0, 2).toUpperCase()}</AvatarFallback>
                </Avatar>
                <div className="text-center space-y-1">
                    <h3 className="text-xl font-bold">{tenant.name}</h3>
                    <p className="text-sm text-muted-foreground">{tenant.region}</p>
                </div>

                <div className="grid grid-cols-2 gap-4 w-full pt-2">
                    <div className="flex flex-col items-center justify-center p-2 bg-muted/50 rounded-md">
                        <Users className="h-4 w-4 mb-1 text-muted-foreground" />
                        <span className="text-lg font-bold">12</span>
                        <span className="text-[10px] uppercase text-muted-foreground font-medium">Users</span>
                    </div>
                    <div className="flex flex-col items-center justify-center p-2 bg-muted/50 rounded-md">
                        <Router className="h-4 w-4 mb-1 text-muted-foreground" />
                        <span className="text-lg font-bold">4</span>
                        <span className="text-[10px] uppercase text-muted-foreground font-medium">Routers</span>
                    </div>
                </div>
            </CardContent>
            <CardFooter>
                <Button className="w-full" variant="outline">Manage Tenant</Button>
            </CardFooter>
        </Card>
    )
}
