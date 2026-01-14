import { AlertTriangle, Info, CheckCircle2, XCircle } from "lucide-react"
import { cn } from "@/lib/utils"

interface EmptyStateProps {
    icon?: React.ComponentType<{ className?: string }>
    title: string
    description: string
    action?: React.ReactNode
    variant?: "default" | "info" | "warning" | "success" | "error"
}

const variantStyles = {
    default: "bg-muted/20",
    info: "bg-blue-500/10",
    warning: "bg-yellow-500/10",
    success: "bg-emerald-500/10",
    error: "bg-destructive/10"
}

const variantIcons = {
    default: null,
    info: Info,
    warning: AlertTriangle,
    success: CheckCircle2,
    error: XCircle
}

const variantIconColors = {
    default: "text-muted-foreground",
    info: "text-blue-500",
    warning: "text-yellow-500",
    success: "text-emerald-500",
    error: "text-destructive"
}

export function EmptyState({
    icon: CustomIcon,
    title,
    description,
    action,
    variant = "default"
}: EmptyStateProps) {
    const Icon = CustomIcon || variantIcons[variant]
    const iconColor = variantIconColors[variant]

    return (
        <div className={cn(
            "flex flex-col items-center justify-center min-h-[400px] p-12 text-center rounded-lg border border-dashed",
            variantStyles[variant]
        )}>
            {Icon && (
                <div className="mb-4">
                    <Icon className={cn("h-12 w-12", iconColor, "opacity-50")} />
                </div>
            )}
            <h3 className="text-lg font-semibold mb-2">{title}</h3>
            <p className="text-sm text-muted-foreground max-w-sm mb-6">
                {description}
            </p>
            {action && <div>{action}</div>}
        </div>
    )
}
