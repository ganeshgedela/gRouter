import * as React from "react"
import { createPortal } from "react-dom"
import { X } from "lucide-react"
import { cn } from "@/lib/utils"

const Dialog = ({
    open,
    onOpenChange,
    children
}: {
    open?: boolean
    onOpenChange?: (open: boolean) => void
    children: React.ReactNode
}) => {
    // Simple state if uncontrolled, but usually controlled
    const [isOpen, setIsOpen] = React.useState(open || false)

    React.useEffect(() => {
        if (open !== undefined) setIsOpen(open)
    }, [open])

    const handleOpenChange = (newOpen: boolean) => {
        setIsOpen(newOpen)
        onOpenChange?.(newOpen)
    }

    // Provide context for children to close dialog
    return (
        <DialogContext.Provider value={{ open: isOpen, onOpenChange: handleOpenChange }}>
            {children}
        </DialogContext.Provider>
    )
}

interface DialogContextType {
    open: boolean
    onOpenChange: (open: boolean) => void
}

const DialogContext = React.createContext<DialogContextType | undefined>(undefined)

const DialogTrigger = React.forwardRef<
    HTMLButtonElement,
    React.HTMLAttributes<HTMLButtonElement> & { asChild?: boolean }
>(({ className, onClick, children, ...props }, ref) => {
    const context = React.useContext(DialogContext)
    return (
        <button
            ref={ref}
            className={className}
            onClick={(e) => {
                onClick?.(e as any)
                context?.onOpenChange(true)
            }}
            {...props}
        >
            {children}
        </button>
    )
})
DialogTrigger.displayName = "DialogTrigger"

const DialogPortal = ({ children }: { children: React.ReactNode }) => {
    // In a real app we'd portal to body, but for simplicity we render inline if open
    // Actually let's use portal for z-index correctness
    const context = React.useContext(DialogContext)
    if (!context?.open) return null

    return createPortal(
        <div className="fixed inset-0 z-50 flex items-start justify-center sm:items-center">
            {children}
        </div>,
        document.body
    )
}

const DialogOverlay = React.forwardRef<
    HTMLDivElement,
    React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => {
    const context = React.useContext(DialogContext)
    return (
        <div
            ref={ref}
            className={cn(
                "fixed inset-0 z-50 bg-background/80 backdrop-blur-sm transition-all duration-100 data-[state=closed]:animate-out data-[state=closed]:fade-out data-[state=open]:fade-in",
                className
            )}
            onClick={() => context?.onOpenChange(false)}
            {...props}
        />
    )
})
DialogOverlay.displayName = "DialogOverlay"

const DialogContent = React.forwardRef<
    HTMLDivElement,
    React.HTMLAttributes<HTMLDivElement>
>(({ className, children, ...props }, ref) => {
    const context = React.useContext(DialogContext)
    return (
        <DialogPortal>
            <DialogOverlay />
            <div
                ref={ref}
                className={cn(
                    "fixed left-[50%] top-[50%] z-50 grid w-full max-w-lg translate-x-[-50%] translate-y-[-50%] gap-4 border bg-background p-6 shadow-lg duration-200 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[state=closed]:slide-out-to-left-1/2 data-[state=closed]:slide-out-to-top-[48%] data-[state=open]:slide-in-from-left-1/2 data-[state=open]:slide-in-from-top-[48%] sm:rounded-lg",
                    className
                )}
                {...props}
            >
                {children}
                <button
                    className="absolute right-4 top-4 rounded-sm opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:pointer-events-none data-[state=open]:bg-accent data-[state=open]:text-muted-foreground"
                    onClick={() => context?.onOpenChange(false)}
                >
                    <X className="h-4 w-4" />
                    <span className="sr-only">Close</span>
                </button>
            </div>
        </DialogPortal>
    )
})
DialogContent.displayName = "DialogContent"

const DialogHeader = ({
    className,
    ...props
}: React.HTMLAttributes<HTMLDivElement>) => (
    <div
        className={cn(
            "flex flex-col space-y-1.5 text-center sm:text-left",
            className
        )}
        {...props}
    />
)
DialogHeader.displayName = "DialogHeader"

const DialogFooter = ({
    className,
    ...props
}: React.HTMLAttributes<HTMLDivElement>) => (
    <div
        className={cn(
            "flex flex-col-reverse sm:flex-row sm:justify-end sm:space-x-2",
            className
        )}
        {...props}
    />
)
DialogFooter.displayName = "DialogFooter"

const DialogTitle = React.forwardRef<
    HTMLHeadingElement,
    React.HTMLAttributes<HTMLHeadingElement>
>(({ className, ...props }, ref) => (
    <h3
        ref={ref}
        className={cn(
            "text-lg font-semibold leading-none tracking-tight",
            className
        )}
        {...props}
    />
))
DialogTitle.displayName = "DialogTitle"

const DialogDescription = React.forwardRef<
    HTMLParagraphElement,
    React.HTMLAttributes<HTMLParagraphElement>
>(({ className, ...props }, ref) => (
    <p
        ref={ref}
        className={cn("text-sm text-muted-foreground", className)}
        {...props}
    />
))
DialogDescription.displayName = "DialogDescription"

// Exports
const DialogClose = React.forwardRef<
    HTMLButtonElement,
    React.ButtonHTMLAttributes<HTMLButtonElement>
>(({ className, ...props }, ref) => {
    const context = React.useContext(DialogContext)
    return (
        <button
            ref={ref}
            className={className}
            onClick={() => context?.onOpenChange(false)}
            {...props}
        />
    )
})
DialogClose.displayName = "DialogClose"

export {
    Dialog,
    DialogPortal,
    DialogOverlay,
    DialogClose,
    DialogTrigger,
    DialogContent,
    DialogHeader,
    DialogFooter,
    DialogTitle,
    DialogDescription,
}
