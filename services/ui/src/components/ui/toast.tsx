import { X, CheckCircle2, AlertCircle, Info } from "lucide-react"
import { cn } from "@/lib/utils"

export type ToastProps = React.HTMLAttributes<HTMLDivElement> & {
    variant?: "default" | "destructive" | "success"
    title?: React.ReactNode
    description?: React.ReactNode
    action?: React.ReactNode
    open?: boolean
    onOpenChange?: (open: boolean) => void
}

export type ToastActionElement = React.ReactElement

export function Toast({
    className,
    variant = "default",
    title,
    description,
    action,
    open,
    onOpenChange,
    ...props
}: ToastProps) {
    if (!open) return null

    const Icon = variant === "destructive" ? AlertCircle : variant === "success" ? CheckCircle2 : Info

    return (
        <div
            className={cn(
                "group pointer-events-auto relative flex w-full items-center justify-between gap-4 overflow-hidden rounded-md border p-4 shadow-lg transition-all data-[swipe=cancel]:translate-x-0 data-[swipe=end]:translate-x-[var(--radix-toast-swipe-end-x)] data-[swipe=move]:translate-x-[var(--radix-toast-swipe-move-x)] data-[swipe=move]:transition-none data-[state=open]:animate-in data-[state=closed]:animate-out data-[swipe=end]:animate-out data-[state=closed]:fade-out-80 data-[state=closed]:slide-out-to-right-full data-[state=open]:slide-in-from-top-full data-[state=open]:sm:slide-in-from-bottom-full border-zinc-200 bg-white text-zinc-950 dark:border-zinc-800 dark:bg-zinc-950 dark:text-zinc-50",
                variant === "destructive" &&
                "destructive group border-red-500 bg-red-500 text-zinc-50 dark:border-red-900 dark:bg-red-900 dark:text-zinc-50",
                variant === "success" &&
                "border-emerald-500/50 bg-emerald-50 text-emerald-900 dark:border-emerald-500/30 dark:bg-emerald-900/20 dark:text-emerald-50",
                className
            )}
            {...props}
        >
            <div className="flex gap-3 items-start">
                {variant !== "default" && <Icon className={cn("h-5 w-5 mt-0.5",
                    variant === "destructive" ? "text-white" :
                        variant === "success" ? "text-emerald-600 dark:text-emerald-400" :
                            "text-blue-500"
                )} />}
                <div className="grid gap-1">
                    {title && <div className="text-sm font-semibold opacity-90">{title}</div>}
                    {description && (
                        <div className="text-sm opacity-90">{description}</div>
                    )}
                </div>
            </div>

            {action}
            <button
                onClick={() => onOpenChange?.(false)}
                className={cn(
                    "absolute right-2 top-2 rounded-md p-1 pl-4 opacity-0 transition-opacity hover:text-foreground focus:opacity-100 focus:outline-none focus:ring-2 group-hover:opacity-100",
                    variant === "destructive" ? "text-red-50 hover:text-red-50" : "text-foreground/50"
                )}
            >
                <X className="h-4 w-4" />
            </button>
        </div>
    )
}
