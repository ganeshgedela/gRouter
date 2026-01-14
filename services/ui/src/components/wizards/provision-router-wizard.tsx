import { useState } from "react"
import {
    Globe,
    Server,
    Network,
    Shield,
    CheckCircle2,
    Cpu,
    HardDrive,
    ArrowRight,
    ArrowLeft,
    Loader2
} from "lucide-react"

import { Button } from "@/components/ui/button"
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Stepper } from "@/components/ui/stepper"
import { cn } from "@/lib/utils"
import { toast } from "@/components/ui/use-toast"

interface ProvisionRouterWizardProps {
    open: boolean
    onOpenChange: (open: boolean) => void
}

const STEPS = [
    { id: 'region', title: 'Region' },
    { id: 'hardware', title: 'Hardware' },
    { id: 'network', title: 'Network' },
    { id: 'review', title: 'Review' },
]

const REGIONS = [
    { id: 'us-west', name: 'US West', loc: 'Oregon', flag: '🇺🇸' },
    { id: 'us-east', name: 'US East', loc: 'N. Virginia', flag: '🇺🇸' },
    { id: 'eu-west', name: 'EU West', loc: 'London', flag: '🇬🇧' },
    { id: 'ap-south', name: 'AP South', loc: 'Singapore', flag: '🇸🇬' },
]

const HARDWARE_SKUS = [
    {
        id: 'c1.small',
        name: 'Cloud Basic',
        cpu: '2 vCPU',
        ram: '4GB RAM',
        bw: '1 Gbps',
        price: '$45/mo'
    },
    {
        id: 'c1.medium',
        name: 'Cloud Standard',
        cpu: '4 vCPU',
        ram: '8GB RAM',
        bw: '5 Gbps',
        price: '$90/mo',
        recommended: true
    },
    {
        id: 'c1.large',
        name: 'Cloud Performance',
        cpu: '8 vCPU',
        ram: '16GB RAM',
        bw: '10 Gbps',
        price: '$180/mo'
    },
]

export function ProvisionRouterWizard({ open, onOpenChange }: ProvisionRouterWizardProps) {
    const [step, setStep] = useState(0)
    const [isSubmitting, setIsSubmitting] = useState(false)
    const [formData, setFormData] = useState({
        region: '',
        sku: '',
        hostname: '',
        vpcId: '',
    })

    const handleNext = () => {
        if (step < STEPS.length - 1) {
            setStep(step + 1)
        } else {
            handleSubmit()
        }
    }

    const handleBack = () => {
        if (step > 0) {
            setStep(step - 1)
        }
    }

    const handleSubmit = async () => {
        setIsSubmitting(true)
        // Simulate API call
        await new Promise(resolve => setTimeout(resolve, 2000))
        setIsSubmitting(false)
        toast({
            title: "Router Provisioned Successfully",
            description: `${formData.hostname} is initializing in ${formData.region}.`,
        })
        onOpenChange(false)
        setStep(0)
        setFormData({ region: '', sku: '', hostname: '', vpcId: '' })
    }

    const canProceed = () => {
        switch (step) {
            case 0: return !!formData.region
            case 1: return !!formData.sku
            case 2: return !!formData.hostname && !!formData.vpcId
            default: return true
        }
    }

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-[700px] p-0 overflow-hidden gap-0">
                <div className="p-6 bg-muted/30 border-b">
                    <DialogHeader className="mb-4">
                        <DialogTitle>Provision New Router</DialogTitle>
                        <DialogDescription>
                            Deploy a new cloud-native routing instance in minutes.
                        </DialogDescription>
                    </DialogHeader>
                    <Stepper steps={STEPS} currentStep={step} className="w-full" />
                </div>

                <div className="p-6 min-h-[350px]">
                    {/* STEP 1: REGION */}
                    {step === 0 && (
                        <div className="grid grid-cols-2 gap-4">
                            {REGIONS.map(region => (
                                <div
                                    key={region.id}
                                    className={cn(
                                        "cursor-pointer rounded-xl border-2 p-4 transition-all hover:bg-muted/50",
                                        formData.region === region.id ? "border-primary bg-primary/5 ring-1 ring-primary" : "border-muted-foreground/20"
                                    )}
                                    onClick={() => setFormData({ ...formData, region: region.id })}
                                >
                                    <div className="flex items-start justify-between">
                                        <div className="space-y-1">
                                            <p className="font-medium leading-none flex items-center gap-2">
                                                {region.name}
                                                <span className="text-xl">{region.flag}</span>
                                            </p>
                                            <p className="text-sm text-muted-foreground">{region.loc}</p>
                                        </div>
                                        {formData.region === region.id && (
                                            <CheckCircle2 className="h-5 w-5 text-primary" />
                                        )}
                                    </div>
                                </div>
                            ))}
                        </div>
                    )}

                    {/* STEP 2: HARDWARE */}
                    {step === 1 && (
                        <div className="space-y-4">
                            {HARDWARE_SKUS.map(sku => (
                                <div
                                    key={sku.id}
                                    className={cn(
                                        "relative cursor-pointer rounded-xl border-2 p-4 transition-all hover:bg-muted/50",
                                        formData.sku === sku.id ? "border-primary bg-primary/5 ring-1 ring-primary" : "border-muted-foreground/20"
                                    )}
                                    onClick={() => setFormData({ ...formData, sku: sku.id })}
                                >
                                    {sku.recommended && (
                                        <div className="absolute -top-3 left-4 bg-primary text-primary-foreground text-xs px-2 py-1 rounded-full">
                                            Recommended
                                        </div>
                                    )}
                                    <div className="flex items-center justify-between">
                                        <div className="flex items-center gap-4">
                                            <div className="h-10 w-10 rounded-lg bg-secondary flex items-center justify-center">
                                                <Server className="h-6 w-6 text-muted-foreground" />
                                            </div>
                                            <div>
                                                <p className="font-semibold">{sku.name}</p>
                                                <div className="flex gap-3 text-sm text-muted-foreground mt-1">
                                                    <span className="flex items-center gap-1"><Cpu className="h-3 w-3" /> {sku.cpu}</span>
                                                    <span className="flex items-center gap-1"><HardDrive className="h-3 w-3" /> {sku.ram}</span>
                                                    <span className="flex items-center gap-1"><Network className="h-3 w-3" /> {sku.bw}</span>
                                                </div>
                                            </div>
                                        </div>
                                        <div className="text-right">
                                            <p className="font-bold">{sku.price}</p>
                                            <span className="text-xs text-muted-foreground">monthly</span>
                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>
                    )}

                    {/* STEP 3: NETWORK */}
                    {step === 2 && (
                        <div className="space-y-6 max-w-md mx-auto py-4">
                            <div className="space-y-2">
                                <Label>Router Hostname</Label>
                                <Input
                                    placeholder="e.g. router-us-west-prod-01"
                                    value={formData.hostname}
                                    onChange={(e) => setFormData({ ...formData, hostname: e.target.value })}
                                />
                                <p className="text-xs text-muted-foreground">Must be unique within the region.</p>
                            </div>

                            <div className="space-y-2">
                                <Label>VPC ID</Label>
                                <Input
                                    placeholder="vpc-xxxxxxxx"
                                    value={formData.vpcId}
                                    onChange={(e) => setFormData({ ...formData, vpcId: e.target.value })}
                                />
                                <p className="text-xs text-muted-foreground">Target VPC for interface attachment.</p>
                            </div>

                            <div className="rounded-lg bg-yellow-500/10 p-3 border border-yellow-500/20 flex gap-2">
                                <Shield className="h-5 w-5 text-yellow-600 shrink-0" />
                                <p className="text-xs text-yellow-600">
                                    This router will be deployed in a <strong>private subnet</strong> by default.
                                    Public access must be configured via Elastic IP later.
                                </p>
                            </div>
                        </div>
                    )}

                    {/* STEP 4: REVIEW */}
                    {step === 3 && (
                        <div className="space-y-6">
                            <div className="rounded-lg border bg-muted/20 p-4 space-y-4">
                                <div className="flex justify-between border-b pb-2">
                                    <span className="text-muted-foreground">Region</span>
                                    <span className="font-medium flex items-center gap-2">
                                        {REGIONS.find(r => r.id === formData.region)?.flag}
                                        {REGIONS.find(r => r.id === formData.region)?.name}
                                    </span>
                                </div>
                                <div className="flex justify-between border-b pb-2">
                                    <span className="text-muted-foreground">Hardware</span>
                                    <span className="font-medium">{HARDWARE_SKUS.find(s => s.id === formData.sku)?.name}</span>
                                </div>
                                <div className="flex justify-between border-b pb-2">
                                    <span className="text-muted-foreground">Hostname</span>
                                    <span className="font-medium">{formData.hostname}</span>
                                </div>
                                <div className="flex justify-between border-b pb-2">
                                    <span className="text-muted-foreground">VPC Target</span>
                                    <span className="font-medium font-mono text-xs">{formData.vpcId}</span>
                                </div>
                                <div className="flex justify-between items-center pt-2">
                                    <span className="text-muted-foreground">Estimated Cost</span>
                                    <span className="text-xl font-bold">{HARDWARE_SKUS.find(s => s.id === formData.sku)?.price}</span>
                                </div>
                            </div>
                            <div className="text-center text-sm text-muted-foreground">
                                By clicking Confirm, you authorize the provisioning of these resources.
                            </div>
                        </div>
                    )}
                </div>

                <DialogFooter className="p-6 border-t bg-muted/10">
                    <Button variant="outline" onClick={handleBack} disabled={step === 0 || isSubmitting}>
                        <ArrowLeft className="mr-2 h-4 w-4" /> Back
                    </Button>
                    <Button onClick={handleNext} disabled={!canProceed() || isSubmitting}>
                        {isSubmitting ? (
                            <>
                                <Loader2 className="mr-2 h-4 w-4 animate-spin" /> Provisioning...
                            </>
                        ) : step === STEPS.length - 1 ? (
                            <>Confirm Deployment</>
                        ) : (
                            <>Next Step <ArrowRight className="ml-2 h-4 w-4" /></>
                        )}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    )
}
