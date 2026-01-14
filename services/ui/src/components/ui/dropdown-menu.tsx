import * as React from "react"
import { Check, ChevronRight, Circle } from "lucide-react"
import { cn } from "@/lib/utils"

interface DropdownContextType {
    open: boolean
    setOpen: (open: boolean) => void
}

const DropdownContext = React.createContext<DropdownContextType | undefined>(undefined)

const DropdownMenu = ({ children }: { children: React.ReactNode }) => {
    const [open, setOpen] = React.useState(false)
    const dropdownRef = React.useRef<HTMLDivElement>(null)

    React.useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
                setOpen(false)
            }
        }
        document.addEventListener("mousedown", handleClickOutside)
        return () => document.removeEventListener("mousedown", handleClickOutside)
    }, [])

    return (
        <DropdownContext.Provider value={{ open, setOpen }}>
            <div className="relative inline-block text-left" ref={dropdownRef}>
                {children}
            </div>
        </DropdownContext.Provider>
    )
}

const DropdownMenuTrigger = React.forwardRef<
    HTMLButtonElement,
    React.ButtonHTMLAttributes<HTMLButtonElement> & { asChild?: boolean }
>(({ className, children, onClick, asChild, ...props }, ref) => {
    const context = React.useContext(DropdownContext)
    if (!context) throw new Error("DropdownMenuTrigger must be used within DropdownMenu")

    const Comp = asChild ? React.Slot : "button" // React.Slot is not standard, simplifying for no-dep

    // Handling asChild simply by rendering the child and cloning it with handlers if needed, 
    // but for simplicity in this no-dep version, we will assume children is a valid element if asChild is true.
    // Actually, standard usage in this project seems to handle 'asChild' somewhat loosely or properly via Radix.
    // Since we don't have Radix, we'll strip 'asChild' and just wrap. 
    // Wait, the usages I added in tenants.tsx use `asChild`.
    // <DropdownMenuTrigger asChild><Button .../></DropdownMenuTrigger>

    // Implementation for 'asChild' without Radix Slot:
    if (asChild && React.isValidElement(children)) {
        return React.cloneElement(children as React.ReactElement, {
            onClick: (e: React.MouseEvent) => {
                context.setOpen(!context.open)
                children.props.onClick?.(e)
            },
            // @ts-ignore
            ref: ref
        })
    }

    return (
        <button
            ref={ref}
            onClick={(e) => {
                context.setOpen(!context.open)
                onClick?.(e)
            }}
            className={cn("inline-flex items-center justify-center", className)}
            {...props}
        >
            {children}
        </button>
    )
})
DropdownMenuTrigger.displayName = "DropdownMenuTrigger"

const DropdownMenuContent = React.forwardRef<
    HTMLDivElement,
    React.HTMLAttributes<HTMLDivElement> & { align?: "start" | "end" | "center" }
>(({ className, align = "center", ...props }, ref) => {
    const context = React.useContext(DropdownContext)
    if (!context) throw new Error("DropdownMenuContent must be used within DropdownMenu")

    if (!context.open) return null

    const alignmentClass = {
        start: "left-0",
        end: "right-0",
        center: "left-1/2 -translate-x-1/2"
    }[align]

    return (
        <div
            ref={ref}
            className={cn(
                "absolute z-50 mt-2 min-w-[8rem] overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-md animate-in fade-in-0 zoom-in-95",
                alignmentClass,
                className
            )}
            {...props}
        />
    )
})
DropdownMenuContent.displayName = "DropdownMenuContent"

const DropdownMenuItem = React.forwardRef<
    HTMLDivElement,
    React.HTMLAttributes<HTMLDivElement> & {
        inset?: boolean
    }
>(({ className, inset, ...props }, ref) => (
    <div
        ref={ref}
        className={cn(
            "relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none transition-colors hover:bg-accent hover:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50",
            inset && "pl-8",
            className
        )}
        {...props}
    />
))
DropdownMenuItem.displayName = "DropdownMenuItem"

const DropdownMenuCheckboxItem = React.forwardRef<
    HTMLDivElement,
    React.HTMLAttributes<HTMLDivElement> & {
        checked?: boolean
        onCheckedChange?: (checked: boolean) => void
    }
>(({ className, children, checked, onCheckedChange, ...props }, ref) => (
    <div
        ref={ref}
        className={cn(
            "relative flex cursor-default select-none items-center rounded-sm py-1.5 pl-8 pr-2 text-sm outline-none transition-colors hover:bg-accent hover:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50",
            className
        )}
        onClick={() => onCheckedChange?.(!checked)}
        {...props}
    >
        <span className="absolute left-2 flex h-3.5 w-3.5 items-center justify-center">
            {checked && <Check className="h-4 w-4" />}
        </span>
        {children}
    </div>
))
DropdownMenuCheckboxItem.displayName = "DropdownMenuCheckboxItem"

const DropdownMenuLabel = React.forwardRef<
    HTMLDivElement,
    React.HTMLAttributes<HTMLDivElement> & {
        inset?: boolean
    }
>(({ className, inset, ...props }, ref) => (
    <div
        ref={ref}
        className={cn(
            "px-2 py-1.5 text-sm font-semibold",
            inset && "pl-8",
            className
        )}
        {...props}
    />
))
DropdownMenuLabel.displayName = "DropdownMenuLabel"

const DropdownMenuSeparator = React.forwardRef<
    HTMLDivElement,
    React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
    <div
        ref={ref}
        className={cn("-mx-1 my-1 h-px bg-muted", className)}
        {...props}
    />
))
DropdownMenuSeparator.displayName = "DropdownMenuSeparator"

export {
    DropdownMenu,
    DropdownMenuTrigger,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuCheckboxItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
}
