import { Keyboard } from "lucide-react"
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog"
import { Badge } from "@/components/ui/badge"

interface ShortcutHelpProps {
    open: boolean
    onOpenChange: (open: boolean) => void
}

const shortcuts = [
    { keys: ['R'], description: 'Refresh data' },
    { keys: ['F'], description: 'Toggle filters' },
    { keys: ['1'], description: 'Overview tab' },
    { keys: ['2'], description: 'Analytics tab' },
    { keys: ['3'], description: 'Reports tab' },
    { keys: ['Ctrl', 'K'], description: 'Command palette' },
    { keys: ['?'], description: 'Show shortcuts' },
    { keys: ['/'], description: 'Focus search' },
]

export function ShortcutHelp({ open, onOpenChange }: ShortcutHelpProps) {
    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-md">
                <DialogHeader>
                    <DialogTitle className="flex items-center gap-2">
                        <Keyboard className="h-5 w-5" />
                        Keyboard Shortcuts
                    </DialogTitle>
                    <DialogDescription>
                        Navigate faster with keyboard shortcuts
                    </DialogDescription>
                </DialogHeader>
                <div className="space-y-2">
                    {shortcuts.map((shortcut, i) => (
                        <div
                            key={i}
                            className="flex items-center justify-between py-2 px-3 rounded-lg hover:bg-accent/50 transition-colors"
                        >
                            <span className="text-sm">{shortcut.description}</span>
                            <div className="flex gap-1">
                                {shortcut.keys.map((key, j) => (
                                    <Badge
                                        key={j}
                                        variant="secondary"
                                        className="font-mono text-xs px-2"
                                    >
                                        {key}
                                    </Badge>
                                ))}
                            </div>
                        </div>
                    ))}
                </div>
            </DialogContent>
        </Dialog>
    )
}
