import { useState, useRef, useEffect } from "react"
import { Terminal as TerminalIcon, X, Maximize2, Minimize2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"

interface TerminalLine {
    type: "input" | "output" | "system" | "error"
    content: string
    timestamp: Date
}

interface WebTerminalProps {
    routerId: string
    onClose?: () => void
}

export function WebTerminal({ routerId, onClose }: WebTerminalProps) {
    const [lines, setLines] = useState<TerminalLine[]>([
        { type: "system", content: `Connecting to ${routerId}...`, timestamp: new Date() },
        { type: "system", content: "Connection established. Type 'help' for available commands.", timestamp: new Date() }
    ])
    const [input, setInput] = useState("")
    const [isMaximized, setIsMaximized] = useState(false)
    const inputRef = useRef<HTMLInputElement>(null)
    const scrollRef = useRef<HTMLDivElement>(null)

    useEffect(() => {
        if (scrollRef.current) {
            scrollRef.current.scrollTop = scrollRef.current.scrollHeight
        }
    }, [lines])

    const handleCommand = (cmd: string) => {
        const timestamp = new Date()
        const newLines = [...lines, { type: "input", content: `> ${cmd}`, timestamp } as TerminalLine]

        // Mock Command Processing
        const command = cmd.trim().toLowerCase()
        let response: TerminalLine

        switch (command) {
            case "help":
                response = {
                    type: "output",
                    content: "Available commands:\n  show ip int brief   Show interface summary\n  show bgp summary    Show BGP neighbor status\n  ping <ip>           Ping a destination\n  clear counters      Reset interface counters",
                    timestamp
                }
                break
            case "show ip int brief":
                response = {
                    type: "output",
                    content: "Interface    IP-Address      OK? Method Status                Protocol\neth0         10.0.0.1        YES DHCP   up                    up      \neth1         192.168.1.1     YES NVRAM  up                    up      \nlo           127.0.0.1       YES NVRAM  up                    up",
                    timestamp
                }
                break
            case "show bgp summary":
                response = {
                    type: "output",
                    content: "BGP router identifier 10.0.0.1, local AS number 65000\nNeighbor        V    AS MsgRcvd MsgSent   TblVer  InQ OutQ Up/Down  State/PfxRcd\n192.168.1.2     4 65001   12453   12344        0    0    0 2d14h           4",
                    timestamp
                }
                break
            case "clear counters":
                response = { type: "system", content: "Clear 'show interface' counters on all interfaces [confirm]", timestamp }
                break
            case "":
                response = { type: "output", content: "", timestamp }
                break
            default:
                if (command.startsWith("ping")) {
                    response = { type: "output", content: "PING 8.8.8.8 (8.8.8.8) 56(84) bytes of data.\n64 bytes from 8.8.8.8: icmp_seq=1 ttl=118 time=14.2 ms\n64 bytes from 8.8.8.8: icmp_seq=2 ttl=118 time=15.1 ms\n64 bytes from 8.8.8.8: icmp_seq=3 ttl=118 time=13.9 ms\n\n--- 8.8.8.8 ping statistics ---\n3 packets transmitted, 3 received, 0% packet loss, time 2003ms", timestamp }
                } else {
                    response = { type: "error", content: `Command not found: ${command}`, timestamp }
                }
        }

        if (response.content) {
            // Simulate network delay
            setTimeout(() => {
                setLines(prev => [...prev, response])
            }, 300)
        }
        setLines(newLines)
        setInput("")
    }

    const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === "Enter") {
            handleCommand(input)
        }
    }

    return (
        <div className={cn(
            "flex flex-col bg-slate-950 text-slate-50 font-mono text-sm rounded-lg overflow-hidden border border-slate-800 shadow-2xl transition-all duration-300",
            isMaximized ? "fixed inset-4 z-50 h-auto" : "h-[500px] w-full"
        )}>
            {/* Terminal Header */}
            <div className="flex items-center justify-between px-4 py-2 bg-slate-900 border-b border-slate-800">
                <div className="flex items-center gap-2">
                    <TerminalIcon className="h-4 w-4 text-emerald-500" />
                    <span className="font-semibold text-slate-300">root@{routerId}:~</span>
                </div>
                <div className="flex items-center gap-2">
                    <Button
                        variant="ghost"
                        size="icon"
                        className="h-6 w-6 text-slate-400 hover:text-slate-100 hover:bg-slate-800"
                        onClick={() => setIsMaximized(!isMaximized)}
                    >
                        {isMaximized ? <Minimize2 className="h-4 w-4" /> : <Maximize2 className="h-4 w-4" />}
                    </Button>
                    {onClose && (
                        <Button
                            variant="ghost"
                            size="icon"
                            className="h-6 w-6 text-slate-400 hover:text-red-400 hover:bg-slate-800"
                            onClick={onClose}
                        >
                            <X className="h-4 w-4" />
                        </Button>
                    )}
                </div>
            </div>

            {/* Terminal Output */}
            <div
                ref={scrollRef}
                className="flex-1 overflow-y-auto p-4 space-y-1 scrollbar-thin scrollbar-thumb-slate-700 scrollbar-track-transparent"
                onClick={() => inputRef.current?.focus()}
            >
                {lines.map((line, i) => (
                    <div key={i} className={cn("break-words whitespace-pre-wrap", {
                        "text-slate-400": line.type === "input",
                        "text-emerald-400": line.type === "output",
                        "text-blue-400": line.type === "system",
                        "text-red-400": line.type === "error"
                    })}>
                        {line.content}
                    </div>
                ))}
                <div className="flex items-center pt-2">
                    <span className="text-emerald-500 mr-2">➜</span>
                    <input
                        ref={inputRef}
                        type="text"
                        value={input}
                        onChange={(e) => setInput(e.target.value)}
                        onKeyDown={handleKeyDown}
                        className="flex-1 bg-transparent border-none outline-none text-slate-50 placeholder-slate-600 focus:ring-0 p-0"
                        autoFocus
                    />
                </div>
            </div>
        </div>
    )
}
