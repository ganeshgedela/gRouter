import * as React from "react"
import { cn } from "@/lib/utils"

const PopoverContext = React.createContext<{
    open: boolean
    onOpenChange: (open: boolean) => void
} | undefined>(undefined)

const Popover = ({
    children,
    open,
    onOpenChange
}: {
    children: React.ReactNode
    open?: boolean
    onOpenChange?: (open: boolean) => void
}) => {
    const [internalOpen, setInternalOpen] = React.useState(open || false)
    const isControlled = open !== undefined

    const handleOpenChange = (newOpen: boolean) => {
        if (!isControlled) setInternalOpen(newOpen)
        onOpenChange?.(newOpen)
    }

    const currentOpen = isControlled ? open : internalOpen

    return (
        <PopoverContext.Provider value={{ open: !!currentOpen, onOpenChange: handleOpenChange }}>
            <div className="relative inline-block">{children}</div>
        </PopoverContext.Provider>
    )
}

const PopoverTrigger = React.forwardRef<
    HTMLButtonElement,
    React.ButtonHTMLAttributes<HTMLButtonElement> & { asChild?: boolean }
>(({ className, children, onClick, ...props }, ref) => {
    const context = React.useContext(PopoverContext)
    return (
        <button
            ref={ref}
            type="button"
            className={className} // Caller handles trigger styling
            onClick={(e) => {
                onClick?.(e)
                context?.onOpenChange(!context.open)
            }}
            {...props}
        >
            {children}
        </button>
    )
})
PopoverTrigger.displayName = "PopoverTrigger"

const PopoverContent = React.forwardRef<
    HTMLDivElement,
    React.HTMLAttributes<HTMLDivElement> & { align?: 'start' | 'center' | 'end', sideOffset?: number }
>(({ className, align = "center", sideOffset = 4, ...props }, ref) => {
    const context = React.useContext(PopoverContext)

    React.useEffect(() => {
        // We'll use a transparent scrim for closing, simpler than ref checking for now
    }, [context?.open])

    if (!context?.open) return null

    return (
        <>
            <div
                className="fixed inset-0 z-40 bg-transparent"
                onClick={() => context.onOpenChange(false)}
            />
            <div
                ref={ref}
                className={cn(
                    "absolute z-50 w-72 rounded-md border bg-popover p-4 text-popover-foreground shadow-md outline-none animate-in fade-in-0 zoom-in-95",
                    // Positioning logic is simplified: default to bottom
                    "top-[calc(100%+0.5rem)]",
                    align === "start" && "left-0",
                    align === "center" && "left-1/2 -translate-x-1/2",
                    align === "end" && "right-0",
                    className
                )}
                {...props}
            />
        </>
    )
})
PopoverContent.displayName = "PopoverContent"

export { Popover, PopoverTrigger, PopoverContent }
