import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Separator } from "@/components/ui/separator"
import type { VPC, Subnet } from "@/types/network"
import { X, Server, Globe, Calendar, Tag, Trash2, AlertTriangle } from "lucide-react"
import { cn } from "@/lib/utils"
import { SubnetCard } from "./subnet-card"
import { calculateCIDRInfo } from "@/lib/cidr-utils"
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
} from "@/components/ui/alert-dialog"

interface VPCDetailPanelProps {
    vpc: VPC | null
    subnets: Subnet[]
    routers: { id: string; hostname: string; status: string }[]
    onClose: () => void
    onDelete?: (vpc: VPC) => void
}

export function VPCDetailPanel({ vpc, subnets, routers, onClose, onDelete }: VPCDetailPanelProps) {
    const [showDeleteDialog, setShowDeleteDialog] = useState(false)

    if (!vpc) return null

    const vpcSubnets = subnets.filter(s => s.vpcId === vpc.id)
    const vpcRouters = routers.filter(r => vpc.routerIds.includes(r.id))
    const cidrInfo = calculateCIDRInfo(vpc.cidr)

    const handleDelete = () => {
        onDelete?.(vpc)
        setShowDeleteDialog(false)
        onClose()
    }

    return (
        <>
            <div className={cn(
                "fixed top-0 right-0 h-full w-[500px] bg-background border-l shadow-2xl z-50",
                "transform transition-transform duration-300 ease-in-out",
                vpc ? "translate-x-0" : "translate-x-full"
            )}>
                <ScrollArea className="h-full">
                    <div className="p-6 space-y-6">
                        {/* Header */}
                        <div className="flex items-start justify-between">
                            <div className="flex-1">
                                <div className="flex items-center gap-2 mb-2">
                                    <h3 className="font-bold text-xl">{vpc.name}</h3>
                                    <Badge variant={vpc.status === 'active' ? 'default' : 'destructive'}>
                                        {vpc.status}
                                    </Badge>
                                </div>
                                <p className="text-sm text-muted-foreground">{vpc.description}</p>
                            </div>
                            <Button variant="ghost" size="icon" onClick={onClose}>
                                <X className="h-4 w-4" />
                            </Button>
                        </div>

                        <Separator />

                        {/* VPC Metadata */}
                        <div className="space-y-4">
                            <div>
                                <p className="text-xs text-muted-foreground mb-1">VPC ID</p>
                                <p className="font-mono text-sm">{vpc.id}</p>
                            </div>

                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <p className="text-xs text-muted-foreground mb-1">Cloud Provider</p>
                                    <Badge variant="outline" className="font-medium">
                                        {vpc.cloudProvider}
                                    </Badge>
                                </div>
                                <div>
                                    <p className="text-xs text-muted-foreground mb-1">Region</p>
                                    <div className="flex items-center gap-1.5">
                                        <Globe className="h-3 w-3 text-muted-foreground" />
                                        <span className="text-sm font-medium">{vpc.region}</span>
                                    </div>
                                </div>
                            </div>

                            <div>
                                <p className="text-xs text-muted-foreground mb-1">CIDR Block</p>
                                <p className="font-mono text-sm font-semibold">{vpc.cidr}</p>
                                {cidrInfo && (
                                    <div className="mt-2 p-3 bg-muted/50 rounded-md text-xs space-y-1">
                                        <div className="flex justify-between">
                                            <span className="text-muted-foreground">Network:</span>
                                            <span className="font-mono">{cidrInfo.networkAddress}</span>
                                        </div>
                                        <div className="flex justify-between">
                                            <span className="text-muted-foreground">Broadcast:</span>
                                            <span className="font-mono">{cidrInfo.broadcastAddress}</span>
                                        </div>
                                        <div className="flex justify-between">
                                            <span className="text-muted-foreground">Usable IPs:</span>
                                            <span className="font-mono font-semibold">{cidrInfo.usableIPs.toLocaleString()}</span>
                                        </div>
                                    </div>
                                )}
                            </div>

                            <div>
                                <p className="text-xs text-muted-foreground mb-1 flex items-center gap-1">
                                    <Calendar className="h-3 w-3" />
                                    Created
                                </p>
                                <p className="text-sm">{new Date(vpc.createdAt).toLocaleString()}</p>
                            </div>
                        </div>

                        <Separator />

                        {/* Associated Routers */}
                        <div>
                            <h4 className="font-semibold mb-3 flex items-center gap-2">
                                <Server className="h-4 w-4" />
                                Associated Routers ({vpcRouters.length})
                            </h4>
                            <div className="space-y-2">
                                {vpcRouters.length > 0 ? (
                                    vpcRouters.map(router => (
                                        <div key={router.id} className="flex items-center justify-between p-2 rounded border bg-background/50">
                                            <div className="flex items-center gap-2">
                                                <div className={cn(
                                                    "h-2 w-2 rounded-full",
                                                    router.status === 'Online' ? "bg-emerald-500" : "bg-red-500"
                                                )} />
                                                <span className="text-sm font-medium">{router.hostname}</span>
                                            </div>
                                            <span className="text-xs text-muted-foreground font-mono">{router.id}</span>
                                        </div>
                                    ))
                                ) : (
                                    <p className="text-sm text-muted-foreground text-center py-4">No routers associated</p>
                                )}
                            </div>
                        </div>

                        <Separator />

                        {/* Subnets */}
                        <div>
                            <h4 className="font-semibold mb-3">Subnets ({vpcSubnets.length})</h4>
                            <div className="space-y-3">
                                {vpcSubnets.length > 0 ? (
                                    vpcSubnets.map(subnet => {
                                        const router = routers.find(r => r.id === subnet.routerId)
                                        return (
                                            <SubnetCard
                                                key={subnet.id}
                                                subnet={subnet}
                                                routerHostname={router?.hostname}
                                            />
                                        )
                                    })
                                ) : (
                                    <div className="text-center py-8 border border-dashed rounded-lg">
                                        <p className="text-sm text-muted-foreground">No subnets created yet</p>
                                        <Button size="sm" variant="outline" className="mt-3">
                                            Create Subnet
                                        </Button>
                                    </div>
                                )}
                            </div>
                        </div>

                        {/* Tags */}
                        {Object.keys(vpc.tags).length > 0 && (
                            <>
                                <Separator />
                                <div>
                                    <h4 className="font-semibold mb-3 flex items-center gap-2">
                                        <Tag className="h-4 w-4" />
                                        Tags
                                    </h4>
                                    <div className="flex flex-wrap gap-2">
                                        {Object.entries(vpc.tags).map(([key, value]) => (
                                            <Badge key={key} variant="outline" className="font-mono text-xs">
                                                {key}: {value}
                                            </Badge>
                                        ))}
                                    </div>
                                </div>
                            </>
                        )}

                        <Separator />

                        {/* Actions */}
                        <div className="space-y-2">
                            <Button className="w-full" variant="outline">
                                Edit VPC
                            </Button>
                            <Button
                                className="w-full"
                                variant="destructive"
                                onClick={() => setShowDeleteDialog(true)}
                            >
                                <Trash2 className="mr-2 h-4 w-4" />
                                Delete VPC
                            </Button>
                        </div>
                    </div>
                </ScrollArea>
            </div>

            {/* Delete Confirmation Dialog */}
            <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle className="flex items-center gap-2">
                            <AlertTriangle className="h-5 w-5 text-destructive" />
                            Delete VPC?
                        </AlertDialogTitle>
                        <AlertDialogDescription className="space-y-2">
                            <p>
                                Are you sure you want to delete <strong>{vpc.name}</strong>?
                            </p>
                            {vpcSubnets.length > 0 && (
                                <p className="text-destructive font-semibold">
                                    Warning: This will also delete {vpcSubnets.length} associated subnet{vpcSubnets.length !== 1 ? 's' : ''}.
                                </p>
                            )}
                            <p className="text-xs">This action cannot be undone.</p>
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction onClick={handleDelete} className="bg-destructive text-destructive-foreground hover:bg-destructive/90">
                            Delete
                        </AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>

            {/* Overlay */}
            {vpc && (
                <div
                    className="fixed inset-0 bg-background/80 backdrop-blur-sm z-40"
                    onClick={onClose}
                />
            )}
        </>
    )
}
