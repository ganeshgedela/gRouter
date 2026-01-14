import { useState, useEffect } from "react"
import { useNavigate } from "react-router-dom"
import { useTheme } from "@/components/theme-provider"
import {
    Settings,
    Smile,
    LayoutDashboard,
    Server,
    Network,
    Activity,
    Bell,
    Bot,
    ArrowRight,
    Moon,
    Sun,
    Laptop,
    CreditCard,
    Users,
    Search
} from "lucide-react"

import { cn } from "@/lib/utils"

export function CommandMenu() {
    const [open, setOpen] = useState(false)
    const [search, setSearch] = useState("")
    const navigate = useNavigate()
    const { setTheme } = useTheme()

    useEffect(() => {
        const down = (e: KeyboardEvent) => {
            if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
                e.preventDefault()
                setOpen((open) => !open)
            }
        }
        document.addEventListener("keydown", down)
        return () => document.removeEventListener("keydown", down)
    }, [])

    const runCommand = (command: () => void) => {
        setOpen(false)
        command()
    }

    if (!open) return null

    const groups = [
        {
            heading: "Pages",
            items: [
                { icon: LayoutDashboard, name: "Dashboard", action: () => navigate("/") },
                { icon: Server, name: "Tenants", action: () => navigate("/tenants") },
                { icon: Users, name: "Users & IAM", action: () => navigate("/users") },
                { icon: Network, name: "Network Topology", action: () => navigate("/routers") },
                { icon: Activity, name: "Metrics", action: () => navigate("/metrics") },
                { icon: Bell, name: "Alerts", action: () => navigate("/alerts") },
                { icon: Bot, name: "AI Assistant", action: () => navigate("/chat") },
                { icon: Settings, name: "Settings", action: () => navigate("/settings") },
            ]
        },
        {
            heading: "Theme",
            items: [
                { icon: Moon, name: "Dark Mode", action: () => setTheme("dark") },
                { icon: Sun, name: "Light Mode", action: () => setTheme("light") },
                { icon: Laptop, name: "System", action: () => setTheme("system") },
            ]
        },
        {
            heading: "Actions",
            items: [
                { icon: Smile, name: "Send Feedback", action: () => console.log("Feedback") },
                { icon: CreditCard, name: "Billing", action: () => console.log("Billing") },
            ]
        }
    ]

    const filteredGroups = groups.map(group => ({
        ...group,
        items: group.items.filter(item =>
            item.name.toLowerCase().includes(search.toLowerCase())
        )
    })).filter(group => group.items.length > 0)

    return (
        <div className="fixed inset-0 z-50 flex items-start justify-center sm:items-center">
            <div
                className="fixed inset-0 bg-background/80 backdrop-blur-sm transition-all"
                onClick={() => setOpen(false)}
            />
            <div className="animate-in fade-in-0 zoom-in-95 relative z-50 w-full max-w-lg rounded-xl border bg-card shadow-2xl overflow-hidden mt-16 sm:mt-0">
                <div className="flex items-center border-b px-3">
                    <Search className="mr-2 h-4 w-4 shrink-0 opacity-50" />
                    <input
                        className="flex h-12 w-full rounded-md bg-transparent py-3 text-sm outline-none placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50"
                        placeholder="Type a command or search..."
                        value={search}
                        onChange={(e) => setSearch(e.target.value)}
                        autoFocus
                    />
                    <div className="text-xs text-muted-foreground bg-muted px-1.5 py-0.5 rounded border">
                        ESC
                    </div>
                </div>
                <div className="max-h-[300px] overflow-y-auto p-2">
                    {filteredGroups.length === 0 ? (
                        <div className="py-6 text-center text-sm text-muted-foreground">
                            No results found.
                        </div>
                    ) : (
                        filteredGroups.map((group, i) => (
                            <div key={i} className="mb-2">
                                <div className="px-2 py-1.5 text-xs font-medium text-muted-foreground">
                                    {group.heading}
                                </div>
                                {group.items.map((item, j) => (
                                    <div
                                        key={j}
                                        className={cn(
                                            "relative flex cursor-pointer select-none items-center rounded-sm px-2 py-2 text-sm outline-none hover:bg-accent hover:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50",
                                        )}
                                        onClick={() => runCommand(item.action)}
                                    >
                                        <item.icon className="mr-2 h-4 w-4" />
                                        <span>{item.name}</span>
                                        <ArrowRight className="ml-auto h-4 w-4 opacity-0 hover:opacity-100" />
                                    </div>
                                ))}
                            </div>
                        ))
                    )}
                </div>
                <div className="border-t bg-muted/50 px-4 py-2 text-xs text-muted-foreground">
                    <div className="flex justify-between">
                        <div>
                            Press <kbd className="pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium opacity-100">⌘J</kbd> or <kbd className="pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium opacity-100">↓</kbd> to navigate
                        </div>
                        <div>Cloud Native Orchestrator</div>
                    </div>
                </div>
            </div>
        </div>
    )
}
