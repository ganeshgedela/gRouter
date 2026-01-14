import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog"
import { Plus } from "lucide-react"
import { useToast } from "@/components/ui/use-toast"


export function AddTenantDialog() {
    const [open, setOpen] = useState(false)
    const { toast } = useToast()
    const [loading, setLoading] = useState(false)

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        setLoading(true)
        // Simulate API
        await new Promise(resolve => setTimeout(resolve, 1000))
        setLoading(false)
        setOpen(false)
        toast({
            title: "Tenant Created",
            description: "The new tenant has been successfully provisioned.",
        })
    }

    const selectClass = "flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 col-span-3 bg-background"

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                <Button>
                    <Plus className="mr-2 h-4 w-4" /> Add Tenant
                </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-[425px]">
                <form onSubmit={handleSubmit}>
                    <DialogHeader>
                        <DialogTitle>Add New Tenant</DialogTitle>
                        <DialogDescription>
                            Create a new tenant environment. Usage quotas will be set to default.
                        </DialogDescription>
                    </DialogHeader>
                    <div className="grid gap-4 py-4">
                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label htmlFor="name" className="text-right">
                                Name
                            </Label>
                            <Input id="name" placeholder="Acme Corp" className="col-span-3" required />
                        </div>
                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label htmlFor="region" className="text-right">
                                Region
                            </Label>
                            <select
                                id="region"
                                className={selectClass}
                                required
                            >
                                <option value="us-east">US East (N. Virginia)</option>
                                <option value="eu-west">EU West (London)</option>
                                <option value="ap-south">AP Southeast (Singapore)</option>
                            </select>
                        </div>
                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label htmlFor="plan" className="text-right">
                                Plan
                            </Label>
                            <select
                                id="plan"
                                className={selectClass}
                            >
                                <option value="basic">Basic Tier</option>
                                <option value="pro">Pro Tier</option>
                                <option value="enterprise">Enterprise Tier</option>
                            </select>
                        </div>
                    </div>
                    <DialogFooter>
                        <Button type="button" variant="outline" onClick={() => setOpen(false)}>Cancel</Button>
                        <Button type="submit" disabled={loading}>
                            {loading ? "Creating..." : "Create Tenant"}
                        </Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    )
}
