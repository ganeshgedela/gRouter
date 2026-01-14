import { ModeToggle } from "@/components/mode-toggle";
import { Button } from "@/components/ui/button";
import { Search } from "lucide-react";
import { Breadcrumbs } from "@/components/layout/breadcrumbs";

export function Header() {
    return (
        <header className="flex h-14 items-center justify-between border-b bg-background/60 px-4 lg:px-6 backdrop-blur-md supports-[backdrop-filter]:bg-background/60 sticky top-0 z-40 transition-colors duration-500">
            <div className="flex items-center gap-4">
                <div className="hidden md:flex">
                    <Breadcrumbs />
                </div>
            </div>
            <div className="flex items-center gap-2">
                <Button variant="outline" size="sm" className="hidden w-64 justify-between text-muted-foreground sm:inline-flex bg-muted/50">
                    <span className="flex items-center">
                        <Search className="mr-2 h-4 w-4" />
                        Search...
                    </span>
                    <kbd className="pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium text-muted-foreground opacity-100">
                        <span className="text-xs">⌘</span>K
                    </kbd>
                </Button>
                <div className="w-px h-6 bg-border mx-2" />
                <ModeToggle />
                <div className="h-8 w-8 rounded-full bg-primary/10 flex items-center justify-center">
                    <span className="text-xs font-medium">JD</span>
                </div>
            </div>
        </header>
    );
}
