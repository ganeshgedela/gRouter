import * as React from "react"
import { createPortal } from "react-dom"
import { buttonVariants } from "@/components/ui/button"
import { cn } from "@/lib/utils"

// Reuse simpler contexts for AlertDialog
const AlertDialogContext = React.createContext<{
    open: boolean
    onOpenChange: (open: boolean) => void
} | undefined>(undefined)

const AlertDialog = ({ open, onOpenChange, children }: { open?: boolean; onOpenChange?: (open: boolean) => void; children: React.ReactNode }) => {
    const [isOpen, setIsOpen] = React.useState(open || false)

    React.useEffect(() => {
        if (open !== undefined) setIsOpen(open)
    }, [open])

    const handleOpenChange = (newOpen: boolean) => {
        setIsOpen(newOpen)
        onOpenChange?.(newOpen)
    }

    return (
        <AlertDialogContext.Provider value={{ open: isOpen, onOpenChange: handleOpenChange }}>
            {children}
        </AlertDialogContext.Provider>
    )
}

const AlertDialogTrigger = React.forwardRef<
    HTMLButtonElement,
    React.ButtonHTMLAttributes<HTMLButtonElement> & { asChild?: boolean }
>(({ className, onClick, children, ...props }, ref) => {
    const context = React.useContext(AlertDialogContext)
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
AlertDialogTrigger.displayName = "AlertDialogTrigger"

const AlertDialogPortal = ({ children }: { children: React.ReactNode }) => {
    const context = React.useContext(AlertDialogContext)
    if (!context?.open) return null
    return createPortal(
        <div className="fixed inset-0 z-50 flex items-end justify-center sm:items-center">
            {children}
        </div>,
        document.body
    )
}

const AlertDialogOverlay = React.forwardRef<
    HTMLDivElement,
    React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => {
    const context = React.useContext(AlertDialogContext)
    return (
        <div
            ref={ref}
            className={cn(
                "fixed inset-0 z-50 bg-background/80 backdrop-blur-sm transition-opacity animate-in fade-in",
                className
            )}
            onClick={() => context?.onOpenChange(false)}
            {...props}
        />
    )
})
AlertDialogOverlay.displayName = "AlertDialogOverlay"

const AlertDialogContent = React.forwardRef<
    HTMLDivElement,
    React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
    <div
        ref={ref}
        className={cn(
            "fixed z-50 grid w-full max-w-lg scale-100 gap-4 border bg-background p-6 opacity-100 shadow-lg animate-in fade-in-90 slide-in-from-bottom-10 sm:rounded-lg sm:zoom-in-90 sm:slide-in-from-bottom-0",
            className
        )}
        {...props}
    />
))
AlertDialogContent.displayName = "AlertDialogContent"

const AlertDialogHeader = ({
    className,
    ...props
}: React.HTMLAttributes<HTMLDivElement>) => (
    <div
        className={cn(
            "flex flex-col space-y-2 text-center sm:text-left",
            className
        )}
        {...props}
    />
)
AlertDialogHeader.displayName = "AlertDialogHeader"

const AlertDialogFooter = ({
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
AlertDialogFooter.displayName = "AlertDialogFooter"

const AlertDialogTitle = React.forwardRef<
    HTMLHeadingElement,
    React.HTMLAttributes<HTMLHeadingElement>
>(({ className, ...props }, ref) => (
    <h2
        ref={ref}
        className={cn("text-lg font-semibold", className)}
        {...props}
    />
))
AlertDialogTitle.displayName = "AlertDialogTitle"

const AlertDialogDescription = React.forwardRef<
    HTMLParagraphElement,
    React.HTMLAttributes<HTMLParagraphElement>
>(({ className, ...props }, ref) => (
    <p
        ref={ref}
        className={cn("text-sm text-muted-foreground", className)}
        {...props}
    />
))
AlertDialogDescription.displayName = "AlertDialogDescription"

const AlertDialogAction = React.forwardRef<
    HTMLButtonElement,
    React.ButtonHTMLAttributes<HTMLButtonElement>
>(({ className, onClick, ...props }, ref) => {
    const context = React.useContext(AlertDialogContext)
    return (
        <button
            ref={ref}
            className={cn(buttonVariants(), className)}
            onClick={(e) => {
                onClick?.(e)
                context?.onOpenChange(false) // Auto close on action
            }}
            {...props}
        />
    )
})
AlertDialogAction.displayName = "AlertDialogAction"

const AlertDialogCancel = React.forwardRef<
    HTMLButtonElement,
    React.ButtonHTMLAttributes<HTMLButtonElement>
>(({ className, onClick, ...props }, ref) => {
    const context = React.useContext(AlertDialogContext)
    return (
        <button
            ref={ref}
            className={cn(
                buttonVariants({ variant: "outline" }),
                "mt-2 sm:mt-0",
                className
            )}
            onClick={(e) => {
                onClick?.(e)
                context?.onOpenChange(false)
            }}
            {...props}
        />
    )
})
AlertDialogCancel.displayName = "AlertDialogCancel"

export {
    AlertDialog,
    AlertDialogPortal,
    AlertDialogOverlay,
    AlertDialogTrigger,
    AlertDialogContent,
    AlertDialogHeader,
    AlertDialogFooter,
    AlertDialogTitle,
    AlertDialogDescription,
    AlertDialogAction,
    AlertDialogCancel,
}
