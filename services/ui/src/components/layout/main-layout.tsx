import { useState, useEffect } from "react";
import { Outlet } from "react-router-dom";
import { Sidebar } from "./sidebar";
import { Header } from "./header";
import { CommandMenu } from "@/components/command-menu"
import { AIAssistantWidget } from "@/components/ai-assistant-widget"
import { ShortcutHelp } from "@/components/dashboard/shortcut-help"

export function MainLayout() {
    const [showShortcuts, setShowShortcuts] = useState(false);

    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            // Check for '?' (Shift + /)
            if (e.key === '?' && e.target instanceof HTMLElement && !e.target.matches('input, textarea')) {
                e.preventDefault();
                setShowShortcuts(true);
            }
        };

        window.addEventListener('keydown', handleKeyDown);
        return () => window.removeEventListener('keydown', handleKeyDown);
    }, []);

    return (
        <div className="flex h-screen overflow-hidden bg-background">
            <Sidebar />
            <div className="flex flex-1 flex-col overflow-hidden">
                <Header />
                <main className="flex-1 overflow-y-auto">
                    <Outlet />
                </main>
            </div>
            <CommandMenu />
            <AIAssistantWidget />
            <ShortcutHelp open={showShortcuts} onOpenChange={setShowShortcuts} />
        </div>
    );
}
