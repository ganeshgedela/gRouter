import { useState } from "react"
import { AlertOctagon, Flame, Shield } from "lucide-react"
import { Button } from "@/components/ui/button"
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog"
import { useToast } from "@/components/ui/use-toast"

interface SystemPanicDialogProps {
    open: boolean
    onOpenChange: (open: boolean) => void
    onConfirm: () => void
}

export function SystemPanicDialog({ open, onOpenChange, onConfirm }: SystemPanicDialogProps) {
    const { toast } = useToast()
    const [countdown, setCountdown] = useState<number | null>(null)

    const handleConfirm = () => {
        // Start 3-second countdown
        setCountdown(3)

        const countdownInterval = setInterval(() => {
            setCountdown(prev => {
                if (prev === null || prev <= 1) {
                    clearInterval(countdownInterval)
                    onConfirm()
                    onOpenChange(false)

                    toast({
                        title: "🚨 SYSTEM PANIC ACTIVATED",
                        description: "All traffic blocked. Under Attack Mode enabled. Security team notified.",
                        variant: "destructive",
                    })

                    return null
                }
                return prev - 1
            })
        }, 1000)
    }

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-[500px] border-red-500/50 bg-red-500/5">
                <DialogHeader>
                    <DialogTitle className="flex items-center gap-2 text-red-500 text-xl">
                        <Flame className="h-6 w-6 animate-pulse" />
                        System Panic Mode
                    </DialogTitle>
                    <DialogDescription className="text-base">
                        This emergency action will immediately:
                    </DialogDescription>
                </DialogHeader>

                <div className="space-y-3 py-4">
                    <div className="bg-red-500/10 border border-red-500/30 rounded-lg p-4">
                        <ul className="space-y-2 text-sm">
                            <li className="flex items-start gap-2">
                                <AlertOctagon className="h-4 w-4 text-red-500 mt-0.5 flex-shrink-0" />
                                <span><strong>Block ALL incoming traffic</strong> except from trusted IPs</span>
                            </li>
                            <li className="flex items-start gap-2">
                                <AlertOctagon className="h-4 w-4 text-red-500 mt-0.5 flex-shrink-0" />
                                <span><strong>Enable Under Attack Mode</strong> with maximum rate limiting</span>
                            </li>
                            <li className="flex items-start gap-2">
                                <AlertOctagon className="h-4 w-4 text-red-500 mt-0.5 flex-shrink-0" />
                                <span><strong>Disable non-essential services</strong> to preserve resources</span>
                            </li>
                            <li className="flex items-start gap-2">
                                <AlertOctagon className="h-4 w-4 text-red-500 mt-0.5 flex-shrink-0" />
                                <span><strong>Alert security team</strong> and create incident ticket</span>
                            </li>
                        </ul>
                    </div>

                    <div className="bg-yellow-500/10 border border-yellow-500/30 rounded-lg p-3 text-sm text-yellow-200">
                        <strong>⚠️ Warning:</strong> This will cause service disruption. Only use during active security incidents.
                    </div>

                    {countdown !== null && (
                        <div className="bg-red-500/20 border border-red-500 rounded-lg p-4 text-center">
                            <div className="text-2xl font-bold text-red-500 font-mono">
                                Activating in {countdown}...
                            </div>
                        </div>
                    )}
                </div>

                <DialogFooter className="gap-2">
                    <Button
                        type="button"
                        variant="outline"
                        onClick={() => {
                            setCountdown(null)
                            onOpenChange(false)
                        }}
                        disabled={countdown !== null}
                    >
                        <Shield className="mr-2 h-4 w-4" />
                        Cancel
                    </Button>
                    <Button
                        type="button"
                        variant="destructive"
                        onClick={handleConfirm}
                        disabled={countdown !== null}
                        className="bg-red-600 hover:bg-red-700"
                    >
                        <Flame className="mr-2 h-4 w-4" />
                        {countdown !== null ? 'Activating...' : 'Activate Panic Mode'}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    )
}
