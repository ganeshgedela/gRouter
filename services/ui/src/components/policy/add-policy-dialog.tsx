import { useState } from "react"
import { Plus } from "lucide-react"
import { Button } from "@/components/ui/button"
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"
import { useToast } from "@/components/ui/use-toast"

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

interface AddPolicyDialogProps {
    onAdd: (rule: PolicyRule) => void
    existingPriorities: number[]
}

export function AddPolicyDialog({ onAdd, existingPriorities }: AddPolicyDialogProps) {
    const [open, setOpen] = useState(false)
    const { toast } = useToast()

    const [formData, setFormData] = useState({
        name: '',
        type: 'ingress' as 'ingress' | 'egress',
        action: 'allow' as 'allow' | 'deny',
        protocol: 'tcp' as 'tcp' | 'udp' | 'icmp' | 'any',
        source: '',
        destination: '',
        priority: '',
    })

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()

        // Validation
        if (!formData.name || !formData.source || !formData.destination || !formData.priority) {
            toast({
                title: "Validation Error",
                description: "Please fill in all required fields.",
                variant: "destructive"
            })
            return
        }

        const priority = parseInt(formData.priority)
        if (isNaN(priority) || priority < 1 || priority > 999) {
            toast({
                title: "Invalid Priority",
                description: "Priority must be between 1 and 999.",
                variant: "destructive"
            })
            return
        }

        if (existingPriorities.includes(priority)) {
            toast({
                title: "Priority Conflict",
                description: `Priority ${priority} is already in use. Please choose a different priority.`,
                variant: "destructive"
            })
            return
        }

        const newRule: PolicyRule = {
            id: `rule-${priority}`,
            priority,
            name: formData.name,
            type: formData.type,
            action: formData.action,
            protocol: formData.protocol,
            source: formData.source,
            destination: formData.destination,
            status: 'active'
        }

        onAdd(newRule)

        toast({
            title: "Rule Created",
            description: `Policy rule "${formData.name}" has been created successfully.`,
        })

        // Reset form and close dialog
        setFormData({
            name: '',
            type: 'ingress',
            action: 'allow',
            protocol: 'tcp',
            source: '',
            destination: '',
            priority: '',
        })
        setOpen(false)
    }

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                <Button>
                    <Plus className="mr-2 h-4 w-4" /> Add Rule
                </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-[600px]">
                <form onSubmit={handleSubmit}>
                    <DialogHeader>
                        <DialogTitle>Create Policy Rule</DialogTitle>
                        <DialogDescription>
                            Define a new firewall or access control policy rule.
                        </DialogDescription>
                    </DialogHeader>
                    <div className="grid gap-4 py-4">
                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label htmlFor="name" className="text-right">
                                Rule Name *
                            </Label>
                            <Input
                                id="name"
                                value={formData.name}
                                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                                className="col-span-3"
                                placeholder="e.g., Allow HTTPS Traffic"
                            />
                        </div>

                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label htmlFor="priority" className="text-right">
                                Priority *
                            </Label>
                            <Input
                                id="priority"
                                type="number"
                                min="1"
                                max="999"
                                value={formData.priority}
                                onChange={(e) => setFormData({ ...formData, priority: e.target.value })}
                                className="col-span-3"
                                placeholder="1-999 (lower = higher priority)"
                            />
                        </div>

                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label htmlFor="type" className="text-right">
                                Type *
                            </Label>
                            <Select
                                value={formData.type}
                                onValueChange={(value: 'ingress' | 'egress') => setFormData({ ...formData, type: value })}
                            >
                                <SelectTrigger className="col-span-3">
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="ingress">Ingress (Inbound)</SelectItem>
                                    <SelectItem value="egress">Egress (Outbound)</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>

                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label htmlFor="action" className="text-right">
                                Action *
                            </Label>
                            <Select
                                value={formData.action}
                                onValueChange={(value: 'allow' | 'deny') => setFormData({ ...formData, action: value })}
                            >
                                <SelectTrigger className="col-span-3">
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="allow">Allow</SelectItem>
                                    <SelectItem value="deny">Deny</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>

                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label htmlFor="protocol" className="text-right">
                                Protocol *
                            </Label>
                            <Select
                                value={formData.protocol}
                                onValueChange={(value: 'tcp' | 'udp' | 'icmp' | 'any') => setFormData({ ...formData, protocol: value })}
                            >
                                <SelectTrigger className="col-span-3">
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="tcp">TCP</SelectItem>
                                    <SelectItem value="udp">UDP</SelectItem>
                                    <SelectItem value="icmp">ICMP</SelectItem>
                                    <SelectItem value="any">Any</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>

                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label htmlFor="source" className="text-right">
                                Source *
                            </Label>
                            <Input
                                id="source"
                                value={formData.source}
                                onChange={(e) => setFormData({ ...formData, source: e.target.value })}
                                className="col-span-3 font-mono"
                                placeholder="e.g., 10.0.0.0/8 or 0.0.0.0/0"
                            />
                        </div>

                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label htmlFor="destination" className="text-right">
                                Destination *
                            </Label>
                            <Input
                                id="destination"
                                value={formData.destination}
                                onChange={(e) => setFormData({ ...formData, destination: e.target.value })}
                                className="col-span-3 font-mono"
                                placeholder="e.g., 192.168.1.0/24 or self"
                            />
                        </div>
                    </div>
                    <DialogFooter>
                        <Button type="button" variant="outline" onClick={() => setOpen(false)}>
                            Cancel
                        </Button>
                        <Button type="submit">Create Rule</Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    )
}
