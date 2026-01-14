import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog"
import { Badge } from "@/components/ui/badge"
import { Separator } from "@/components/ui/separator"
import { Stepper } from "@/components/ui/stepper"
import type { VPC, CloudProvider } from "@/types/network"
import { validateCIDR, calculateCIDRInfo, suggestSubnets } from "@/lib/cidr-utils"
import { AlertCircle, Check, CheckCircle2, Cloud, Globe, Network, Server } from "lucide-react"
import { cn } from "@/lib/utils"
import { Checkbox } from "@/components/ui/checkbox"

interface CreateVPCWizardProps {
    open: boolean
    onOpenChange: (open: boolean) => void
    routers?: { id: string; hostname: string; region: string; status: string }[]
    onCreate?: (vpc: Partial<VPC>) => void
}

const CLOUD_PROVIDERS: { value: CloudProvider; label: string; regions: string[] }[] = [
    {
        value: 'AWS',
        label: 'Amazon Web Services',
        regions: ['us-east-1', 'us-west-1', 'us-west-2', 'eu-west-1', 'ap-southeast-1']
    },
    {
        value: 'GCP',
        label: 'Google Cloud Platform',
        regions: ['us-central1', 'us-west1', 'europe-west1', 'asia-southeast1']
    },
    {
        value: 'Azure',
        label: 'Microsoft Azure',
        regions: ['eastus', 'westus', 'westeurope', 'southeastasia']
    },
    {
        value: 'On-Premise',
        label: 'On-Premise',
        regions: ['datacenter-1', 'datacenter-2', 'datacenter-3']
    }
]

export function CreateVPCWizard({ open, onOpenChange, routers = [], onCreate }: CreateVPCWizardProps) {
    const [currentStep, setCurrentStep] = useState(0)

    // Step 1: Basic Info
    const [name, setName] = useState("")
    const [description, setDescription] = useState("")
    const [cloudProvider, setCloudProvider] = useState<CloudProvider>('AWS')
    const [region, setRegion] = useState("")

    // Step 2: Network Config
    const [cidr, setCidr] = useState("")
    const [enableDnsResolution, setEnableDnsResolution] = useState(true)
    const [enableDnsHostnames, setEnableDnsHostnames] = useState(true)

    // Step 3: Subnets
    const [subnetCount, setSubnetCount] = useState(2)
    const [suggestedSubnets, setSuggestedSubnets] = useState<string[]>([])

    // Step 4: Router Assignment
    const [selectedRouterIds, setSelectedRouterIds] = useState<string[]>([])

    const cidrInfo = validateCIDR(cidr) ? calculateCIDRInfo(cidr) : null
    const selectedProvider = CLOUD_PROVIDERS.find(p => p.value === cloudProvider)
    const availableRegions = selectedProvider?.regions || []

    const steps = [
        { title: "Basic Info", description: "VPC details" },
        { title: "Network", description: "CIDR configuration" },
        { title: "Subnets", description: "Subnet planning" },
        { title: "Routers", description: "Router assignment" },
        { title: "Review", description: "Confirm and create" }
    ]

    const canProceed = () => {
        switch (currentStep) {
            case 0:
                return name.trim() !== "" && region !== ""
            case 1:
                return validateCIDR(cidr)
            case 2:
                return true
            case 3:
                return selectedRouterIds.length > 0
            case 4:
                return true
            default:
                return false
        }
    }

    const handleNext = () => {
        if (currentStep === 1 && validateCIDR(cidr)) {
            // Generate subnet suggestions when moving from network config
            setSuggestedSubnets(suggestSubnets(cidr, subnetCount))
        }
        if (currentStep < steps.length - 1) {
            setCurrentStep(currentStep + 1)
        }
    }

    const handleBack = () => {
        if (currentStep > 0) {
            setCurrentStep(currentStep - 1)
        }
    }

    const handleCreate = () => {
        const newVPC: Partial<VPC> = {
            name,
            description,
            cidr,
            region,
            cloudProvider,
            routerIds: selectedRouterIds,
            status: 'pending',
            createdAt: new Date().toISOString(),
            tags: {}
        }
        onCreate?.(newVPC)
        handleClose()
    }

    const handleClose = () => {
        setCurrentStep(0)
        setName("")
        setDescription("")
        setCloudProvider('AWS')
        setRegion("")
        setCidr("")
        setSubnetCount(2)
        setSuggestedSubnets([])
        setSelectedRouterIds([])
        onOpenChange(false)
    }

    const toggleRouter = (routerId: string) => {
        setSelectedRouterIds(prev =>
            prev.includes(routerId)
                ? prev.filter(id => id !== routerId)
                : [...prev, routerId]
        )
    }

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="max-w-3xl max-h-[90vh] overflow-hidden flex flex-col">
                <DialogHeader>
                    <DialogTitle className="flex items-center gap-2">
                        <Network className="h-5 w-5" />
                        Create New VPC
                    </DialogTitle>
                    <DialogDescription>
                        Configure a new Virtual Private Cloud across your infrastructure
                    </DialogDescription>
                </DialogHeader>

                <Stepper currentStep={currentStep} steps={steps.map(s => ({ id: s.title.toLowerCase().replace(/\s+/g, '-'), title: s.title, description: s.description }))} />

                <div className="flex-1 overflow-y-auto px-1">
                    {/* Step 0: Basic Information */}
                    {currentStep === 0 && (
                        <div className="space-y-4 py-4">
                            <div className="space-y-2">
                                <Label htmlFor="name">VPC Name *</Label>
                                <Input
                                    id="name"
                                    placeholder="production-vpc"
                                    value={name}
                                    onChange={(e) => setName(e.target.value)}
                                />
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="description">Description</Label>
                                <Input
                                    id="description"
                                    placeholder="Main production environment"
                                    value={description}
                                    onChange={(e) => setDescription(e.target.value)}
                                />
                            </div>

                            <div className="space-y-2">
                                <Label>Cloud Provider *</Label>
                                <div className="grid grid-cols-2 gap-3">
                                    {CLOUD_PROVIDERS.map((provider) => (
                                        <button
                                            key={provider.value}
                                            type="button"
                                            onClick={() => {
                                                setCloudProvider(provider.value)
                                                setRegion("") // Reset region when provider changes
                                            }}
                                            className={cn(
                                                "flex items-center gap-3 p-3 border rounded-lg transition-colors text-left",
                                                cloudProvider === provider.value
                                                    ? "border-primary bg-primary/5"
                                                    : "border-border hover:border-primary/50"
                                            )}
                                        >
                                            <Cloud className="h-5 w-5" />
                                            <span className="font-medium text-sm">{provider.label}</span>
                                            {cloudProvider === provider.value && (
                                                <Check className="h-4 w-4 ml-auto text-primary" />
                                            )}
                                        </button>
                                    ))}
                                </div>
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="region">Region *</Label>
                                <select
                                    id="region"
                                    value={region}
                                    onChange={(e) => setRegion(e.target.value)}
                                    className="w-full px-3 py-2 border rounded-md bg-background"
                                >
                                    <option value="">Select a region</option>
                                    {availableRegions.map((r) => (
                                        <option key={r} value={r}>{r}</option>
                                    ))}
                                </select>
                            </div>
                        </div>
                    )}

                    {/* Step 1: Network Configuration */}
                    {currentStep === 1 && (
                        <div className="space-y-4 py-4">
                            <div className="space-y-2">
                                <Label htmlFor="cidr">CIDR Block *</Label>
                                <Input
                                    id="cidr"
                                    placeholder="10.0.0.0/16"
                                    value={cidr}
                                    onChange={(e) => setCidr(e.target.value)}
                                    className={cn(
                                        cidr && !validateCIDR(cidr) && "border-destructive"
                                    )}
                                />
                                {cidr && !validateCIDR(cidr) && (
                                    <p className="text-xs text-destructive flex items-center gap-1">
                                        <AlertCircle className="h-3 w-3" />
                                        Invalid CIDR format
                                    </p>
                                )}
                            </div>

                            {cidrInfo && (
                                <div className="p-4 bg-muted/50 rounded-lg space-y-2">
                                    <h4 className="font-semibold text-sm mb-3">CIDR Block Information</h4>
                                    <div className="grid grid-cols-2 gap-3 text-sm">
                                        <div>
                                            <p className="text-xs text-muted-foreground">Network Address</p>
                                            <p className="font-mono font-medium">{cidrInfo.networkAddress}</p>
                                        </div>
                                        <div>
                                            <p className="text-xs text-muted-foreground">Broadcast Address</p>
                                            <p className="font-mono font-medium">{cidrInfo.broadcastAddress}</p>
                                        </div>
                                        <div>
                                            <p className="text-xs text-muted-foreground">Subnet Mask</p>
                                            <p className="font-mono font-medium">{cidrInfo.subnetMask}</p>
                                        </div>
                                        <div>
                                            <p className="text-xs text-muted-foreground">Usable IPs</p>
                                            <p className="font-mono font-medium text-primary">{cidrInfo.usableIPs.toLocaleString()}</p>
                                        </div>
                                    </div>
                                    <div className="pt-2 border-t">
                                        <p className="text-xs text-muted-foreground">IP Range</p>
                                        <p className="font-mono text-xs">{cidrInfo.firstUsableIP} - {cidrInfo.lastUsableIP}</p>
                                    </div>
                                </div>
                            )}

                            <Separator />

                            <div className="space-y-3">
                                <h4 className="font-semibold text-sm">DNS Configuration</h4>
                                <div className="flex items-center space-x-2">
                                    <Checkbox
                                        id="dns-resolution"
                                        checked={enableDnsResolution}
                                        onCheckedChange={(checked) => setEnableDnsResolution(checked as boolean)}
                                    />
                                    <Label htmlFor="dns-resolution" className="text-sm font-normal cursor-pointer">
                                        Enable DNS resolution
                                    </Label>
                                </div>
                                <div className="flex items-center space-x-2">
                                    <Checkbox
                                        id="dns-hostnames"
                                        checked={enableDnsHostnames}
                                        onCheckedChange={(checked) => setEnableDnsHostnames(checked as boolean)}
                                    />
                                    <Label htmlFor="dns-hostnames" className="text-sm font-normal cursor-pointer">
                                        Enable DNS hostnames
                                    </Label>
                                </div>
                            </div>
                        </div>
                    )}

                    {/* Step 2: Subnet Planning */}
                    {currentStep === 2 && (
                        <div className="space-y-4 py-4">
                            <div className="space-y-2">
                                <Label htmlFor="subnet-count">Number of Subnets</Label>
                                <Input
                                    id="subnet-count"
                                    type="number"
                                    min="0"
                                    max="16"
                                    value={subnetCount}
                                    onChange={(e) => {
                                        const count = parseInt(e.target.value, 10)
                                        setSubnetCount(count)
                                        if (validateCIDR(cidr)) {
                                            setSuggestedSubnets(suggestSubnets(cidr, count))
                                        }
                                    }}
                                />
                                <p className="text-xs text-muted-foreground">
                                    Suggested subnet divisions for your VPC CIDR block
                                </p>
                            </div>

                            {suggestedSubnets.length > 0 && (
                                <div className="space-y-2">
                                    <h4 className="font-semibold text-sm">Suggested Subnets</h4>
                                    <div className="space-y-2 max-h-64 overflow-y-auto">
                                        {suggestedSubnets.map((subnet, idx) => {
                                            const subnetInfo = calculateCIDRInfo(subnet)
                                            return (
                                                <div key={idx} className="p-3 border rounded-lg bg-background">
                                                    <div className="flex items-center justify-between mb-1">
                                                        <p className="font-mono font-semibold text-sm">{subnet}</p>
                                                        <Badge variant="outline" className="text-xs">
                                                            Subnet {idx + 1}
                                                        </Badge>
                                                    </div>
                                                    {subnetInfo && (
                                                        <p className="text-xs text-muted-foreground">
                                                            {subnetInfo.usableIPs.toLocaleString()} usable IPs
                                                        </p>
                                                    )}
                                                </div>
                                            )
                                        })}
                                    </div>
                                    <p className="text-xs text-muted-foreground italic">
                                        These subnets can be created after VPC creation
                                    </p>
                                </div>
                            )}
                        </div>
                    )}

                    {/* Step 3: Router Assignment */}
                    {currentStep === 3 && (
                        <div className="space-y-4 py-4">
                            <div>
                                <h4 className="font-semibold mb-1">Select Routers</h4>
                                <p className="text-sm text-muted-foreground mb-4">
                                    Choose which routers will be associated with this VPC
                                </p>
                            </div>

                            <div className="space-y-2 max-h-80 overflow-y-auto">
                                {routers.length > 0 ? (
                                    routers.map((router) => (
                                        <button
                                            key={router.id}
                                            type="button"
                                            onClick={() => toggleRouter(router.id)}
                                            className={cn(
                                                "w-full flex items-center justify-between p-3 border rounded-lg transition-colors text-left",
                                                selectedRouterIds.includes(router.id)
                                                    ? "border-primary bg-primary/5"
                                                    : "border-border hover:border-primary/50"
                                            )}
                                        >
                                            <div className="flex items-center gap-3">
                                                <Server className="h-4 w-4 text-muted-foreground" />
                                                <div>
                                                    <p className="font-medium text-sm">{router.hostname}</p>
                                                    <p className="text-xs text-muted-foreground flex items-center gap-1">
                                                        <Globe className="h-3 w-3" />
                                                        {router.region}
                                                    </p>
                                                </div>
                                            </div>
                                            {selectedRouterIds.includes(router.id) && (
                                                <CheckCircle2 className="h-5 w-5 text-primary" />
                                            )}
                                        </button>
                                    ))
                                ) : (
                                    <div className="text-center py-8 border border-dashed rounded-lg">
                                        <p className="text-sm text-muted-foreground">No routers available</p>
                                    </div>
                                )}
                            </div>
                        </div>
                    )}

                    {/* Step 4: Review */}
                    {currentStep === 4 && (
                        <div className="space-y-4 py-4">
                            <h4 className="font-semibold">Review Configuration</h4>

                            <div className="space-y-3 p-4 bg-muted/30 rounded-lg">
                                <div>
                                    <p className="text-xs text-muted-foreground">VPC Name</p>
                                    <p className="font-semibold">{name}</p>
                                </div>
                                <Separator />
                                <div className="grid grid-cols-2 gap-4">
                                    <div>
                                        <p className="text-xs text-muted-foreground">Cloud Provider</p>
                                        <p className="font-medium">{cloudProvider}</p>
                                    </div>
                                    <div>
                                        <p className="text-xs text-muted-foreground">Region</p>
                                        <p className="font-medium">{region}</p>
                                    </div>
                                </div>
                                <Separator />
                                <div>
                                    <p className="text-xs text-muted-foreground">CIDR Block</p>
                                    <p className="font-mono font-semibold text-primary">{cidr}</p>
                                    {cidrInfo && (
                                        <p className="text-xs text-muted-foreground mt-1">
                                            {cidrInfo.usableIPs.toLocaleString()} usable IP addresses
                                        </p>
                                    )}
                                </div>
                                <Separator />
                                <div>
                                    <p className="text-xs text-muted-foreground mb-2">Associated Routers ({selectedRouterIds.length})</p>
                                    <div className="flex flex-wrap gap-2">
                                        {selectedRouterIds.map(id => {
                                            const router = routers.find(r => r.id === id)
                                            return router ? (
                                                <Badge key={id} variant="outline">
                                                    {router.hostname}
                                                </Badge>
                                            ) : null
                                        })}
                                    </div>
                                </div>
                                {description && (
                                    <>
                                        <Separator />
                                        <div>
                                            <p className="text-xs text-muted-foreground">Description</p>
                                            <p className="text-sm">{description}</p>
                                        </div>
                                    </>
                                )}
                            </div>
                        </div>
                    )}
                </div>

                <Separator />

                <div className="flex justify-between pt-4">
                    <Button
                        variant="outline"
                        onClick={currentStep === 0 ? handleClose : handleBack}
                    >
                        {currentStep === 0 ? 'Cancel' : 'Back'}
                    </Button>
                    <Button
                        onClick={currentStep === steps.length - 1 ? handleCreate : handleNext}
                        disabled={!canProceed()}
                    >
                        {currentStep === steps.length - 1 ? 'Create VPC' : 'Next'}
                    </Button>
                </div>
            </DialogContent>
        </Dialog>
    )
}
