import { useState, useRef, useEffect } from "react"
import { Bot, Send, X, Minimize2, Maximize2, Sparkles, User } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card } from "@/components/ui/card"
import { cn } from "@/lib/utils"

interface Message {
    id: string
    role: "user" | "assistant"
    content: string
    timestamp: Date
}

export function AIAssistantWidget() {
    const [isOpen, setIsOpen] = useState(false)
    const [isMinimized, setIsMinimized] = useState(false)
    const [inputValue, setInputValue] = useState("")
    const [isTyping, setIsTyping] = useState(false)
    const scrollRef = useRef<HTMLDivElement>(null)

    const [messages, setMessages] = useState<Message[]>([
        {
            id: "1",
            role: "assistant",
            content: "Hello! I'm Polaris AI. I can help you orchestrate your network, analyze traffic patterns, or troubleshoot router issues. How can I assist you today?",
            timestamp: new Date()
        }
    ])

    const scrollToBottom = () => {
        if (scrollRef.current) {
            scrollRef.current.scrollTop = scrollRef.current.scrollHeight
        }
    }

    useEffect(() => {
        scrollToBottom()
    }, [messages, isOpen, isTyping])

    const handleSendMessage = () => {
        if (!inputValue.trim()) return

        const userMsg: Message = {
            id: Date.now().toString(),
            role: "user",
            content: inputValue,
            timestamp: new Date()
        }

        setMessages(prev => [...prev, userMsg])
        setInputValue("")
        setIsTyping(true)

        // Mock AI Response
        setTimeout(() => {
            const aiMsg: Message = {
                id: (Date.now() + 1).toString(),
                role: "assistant",
                content: "I've analyzed the request. Based on current telemetry from US-EAST-1, I recommend scaling up the bandwidth allocation. Would you like me to execute this action?",
                timestamp: new Date()
            }
            setMessages(prev => [...prev, aiMsg])
            setIsTyping(false)
        }, 1500)
    }

    const handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.key === "Enter" && !e.shiftKey) {
            e.preventDefault()
            handleSendMessage()
        }
    }

    if (!isOpen) {
        return (
            <Button
                onClick={() => setIsOpen(true)}
                className="fixed bottom-6 right-6 h-14 w-14 rounded-full shadow-2xl bg-primary hover:bg-primary/90 transition-all duration-300 hover:scale-105 z-50 group"
            >
                <Bot className="h-8 w-8 text-primary-foreground group-hover:rotate-12 transition-transform" />
                <span className="absolute -top-1 -right-1 flex h-4 w-4">
                    <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-sky-400 opacity-75"></span>
                    <span className="relative inline-flex rounded-full h-4 w-4 bg-sky-500 border-2 border-background"></span>
                </span>
            </Button>
        )
    }

    return (
        <div className={cn(
            "fixed bottom-6 right-6 z-50 transition-all duration-300 ease-in-out",
            isMinimized ? "w-80 h-14" : "w-[400px] h-[600px]"
        )}>
            <Card className="h-full border-primary/20 shadow-2xl overflow-hidden flex flex-col bg-background/95 backdrop-blur-supports-[backdrop-filter]:bg-background/60">
                {/* Header */}
                <div
                    className="flex items-center justify-between p-4 bg-primary/5 border-b cursor-pointer"
                    onClick={() => setIsMinimized(!isMinimized)}
                >
                    <div className="flex items-center gap-3">
                        <div className="h-8 w-8 rounded-lg bg-primary/10 flex items-center justify-center border border-primary/20">
                            <Bot className="h-5 w-5 text-primary" />
                        </div>
                        <div>
                            <h3 className="font-semibold text-sm flex items-center gap-2">
                                Polaris Assistant
                                <span className="flex h-2 w-2 rounded-full bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.5)]"></span>
                            </h3>
                            {!isMinimized && (
                                <p className="text-[10px] text-muted-foreground font-mono uppercase tracking-wider">
                                    AI Orchestrator Online
                                </p>
                            )}
                        </div>
                    </div>
                    <div className="flex items-center gap-1">
                        <Button variant="ghost" size="icon" className="h-8 w-8" onClick={(e) => { e.stopPropagation(); setIsMinimized(!isMinimized); }}>
                            {isMinimized ? <Maximize2 className="h-4 w-4" /> : <Minimize2 className="h-4 w-4" />}
                        </Button>
                        <Button variant="ghost" size="icon" className="h-8 w-8 hover:text-destructive" onClick={(e) => { e.stopPropagation(); setIsOpen(false); }}>
                            <X className="h-4 w-4" />
                        </Button>
                    </div>
                </div>

                {/* Chat Area */}
                {!isMinimized && (
                    <>
                        <div className="flex-1 p-4 overflow-y-auto" ref={scrollRef}>
                            <div className="space-y-4">
                                {messages.map((msg) => (
                                    <div
                                        key={msg.id}
                                        className={cn(
                                            "flex gap-3",
                                            msg.role === "user" ? "flex-row-reverse" : "flex-row"
                                        )}
                                    >
                                        <div className={cn(
                                            "h-8 w-8 rounded-full flex items-center justify-center shrink-0 border",
                                            msg.role === "user" ? "bg-muted border-border" : "bg-primary/10 border-primary/20"
                                        )}>
                                            {msg.role === "user" ? <User className="h-4 w-4 text-muted-foreground" /> : <Bot className="h-4 w-4 text-primary" />}
                                        </div>
                                        <div className={cn(
                                            "rounded-lg p-3 text-sm max-w-[80%]",
                                            msg.role === "user"
                                                ? "bg-primary text-primary-foreground"
                                                : "bg-muted/50 border border-border/50 text-foreground"
                                        )}>
                                            <p>{msg.content}</p>
                                            <span className="text-[10px] opacity-50 mt-1 block">
                                                {msg.timestamp.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                                            </span>
                                        </div>
                                    </div>
                                ))}
                                {isTyping && (
                                    <div className="flex gap-3">
                                        <div className="h-8 w-8 rounded-full bg-primary/10 flex items-center justify-center border border-primary/20">
                                            <Bot className="h-4 w-4 text-primary" />
                                        </div>
                                        <div className="bg-muted/50 border border-border/50 rounded-lg p-3 flex gap-1 items-center">
                                            <span className="h-2 w-2 rounded-full bg-primary/40 animate-bounce [animation-delay:-0.3s]"></span>
                                            <span className="h-2 w-2 rounded-full bg-primary/40 animate-bounce [animation-delay:-0.15s]"></span>
                                            <span className="h-2 w-2 rounded-full bg-primary/40 animate-bounce"></span>
                                        </div>
                                    </div>
                                )}
                            </div>
                        </div>

                        {/* Input Area */}
                        <div className="p-4 bg-background/50 border-t">
                            <div className="flex gap-2">
                                <Input
                                    placeholder="Ask Polaris to configure routers..."
                                    className="flex-1 bg-background/50 focus-visible:ring-primary/20"
                                    value={inputValue}
                                    onChange={(e) => setInputValue(e.target.value)}
                                    onKeyDown={handleKeyDown}
                                />
                                <Button size="icon" onClick={handleSendMessage} disabled={!inputValue.trim() || isTyping}>
                                    <Send className="h-4 w-4" />
                                </Button>
                            </div>
                            <div className="mt-2 flex items-center justify-center gap-2">
                                <span className="text-[10px] text-muted-foreground flex items-center gap-1">
                                    <Sparkles className="h-3 w-3 text-sky-500" />
                                    Powered by Cortex LLM
                                </span>
                            </div>
                        </div>
                    </>
                )}
            </Card>
        </div>
    )
}
