import { NavLink } from "react-router-dom";
import {
    LayoutDashboard,
    Building2,
    Users,
    Network,
    Activity,
    Bell,
    Bot,
    Settings,
    Shield,
    FileText
} from "lucide-react";
import { cn } from "@/lib/utils";

import { PolarisLogo } from "@/components/polaris-logo";

const sidebarItems = [
    { icon: LayoutDashboard, label: "Dashboard", href: "/" },
    { icon: Building2, label: "Tenants", href: "/tenants" },
    { icon: Users, label: "Users", href: "/users" },
    { icon: Network, label: "Network", href: "/routers" },
    { icon: FileText, label: "Policies", href: "/policies" },
    { icon: Shield, label: "Security", href: "/security" },
    { icon: Activity, label: "Metrics", href: "/metrics" },
    { icon: Bell, label: "Alerts", href: "/alerts" },
    { icon: Bot, label: "Assistant", href: "/chat" },
    { icon: Settings, label: "Settings", href: "/settings" },
];

export function Sidebar() {
    return (
        <div className="flex h-full w-64 flex-col border-r bg-background/60 backdrop-blur-md text-card-foreground transition-colors duration-500">
            <div className="flex h-14 items-center border-b px-4 gap-3">
                <PolarisLogo className="w-8 h-8" />
                <span className="text-xl font-bold tracking-tight bg-gradient-to-r from-primary to-purple-400 bg-clip-text text-transparent">
                    Polaris
                </span>
            </div>
            <nav className="flex-1 overflow-y-auto py-4">
                <ul className="space-y-1 px-2">
                    {sidebarItems.map((item) => (
                        <li key={item.href}>
                            <NavLink
                                to={item.href}
                                className={({ isActive }) =>
                                    cn(
                                        "flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-all duration-200",
                                        "hover:bg-primary/5 hover:text-primary",
                                        isActive
                                            ? "bg-primary/10 text-primary shadow-sm"
                                            : "text-muted-foreground"
                                    )
                                }
                            >
                                {({ isActive }) => (
                                    <>
                                        <item.icon className={cn("h-4 w-4", isActive && "fill-primary/20")} />
                                        {item.label}
                                    </>
                                )}
                            </NavLink>
                        </li>
                    ))}
                </ul>
            </nav>
            <div className="border-t p-4">
                <div className="text-xs text-muted-foreground">
                    v1.0.0
                </div>
            </div>
        </div>
    );
}
