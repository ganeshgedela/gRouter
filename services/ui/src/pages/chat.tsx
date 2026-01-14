import { useState, useRef, useEffect } from "react"
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { atomDark } from 'react-syntax-highlighter/dist/esm/styles/prism'
import {
    Send,
    Bot,
    User,
    Sparkles,
    Zap,
    Activity,
    AlertTriangle,
    StopCircle,
    Copy,
    Check
} from "lucide-react"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Avatar } from "@/components/ui/avatar"
import { cn } from "@/lib/utils"

interface Message {
    id: string
    role: "user" | "assistant"
    content: string
    timestamp: Date
}

const SUGGESTED_QUERIES = [
    { icon: Zap, label: "Optimize BGP Config", prompt: "Analyze my current BGP configuration for optimization opportunities." },
    { icon: Activity, label: "Check Network Health", prompt: "Generate a summary of network latency and packet loss for the last hour." },
    { icon: AlertTriangle, label: "Explain Recent Alerts", prompt: "Explain the root cause of the recent 'High CPU' alerts in US-West." },
]

// Smart Mock Response Engine
const generateAIResponse = (query: string): string => {
    const lowerQuery = query.toLowerCase()

    if (lowerQuery.includes("bgp") || lowerQuery.includes("router")) {
        return `I've analyzed the BGP configuration for **router-us-west-1**.

It appears that neighbor \`10.0.5.2\` is experiencing route flaps. Here is a recommended configuration change to dampen these routes:

\`\`\`cisco
router bgp 65000
 neighbor 10.0.5.2 remote-as 65001
 ! Configure route dampening
 bgp dampening 15 750 2000 60
 !
 address-family ipv4
  neighbor 10.0.5.2 activate
  neighbor 10.0.5.2 soft-reconfiguration inbound
 exit-address-family
\`\`\`

**Breakdown of changes:**
1. **bgp dampening**: Penalizes flapping routes.
2. **soft-reconfiguration**: Allows policy changes without hard resets.

Would you like me to apply this configuration to the candidate datastore?`
    }

    if (lowerQuery.includes("health") || lowerQuery.includes("latency") || lowerQuery.includes("loss")) {
        return `### Network Health Summary (Last 1 Hour)

| Metric | Region | Value | Status |
| :--- | :--- | :--- | :--- |
| **Latency P95** | US-West | 24ms | ✅ Healthy |
| **Latency P95** | EU-Central | 45ms | ✅ Healthy |
| **Packet Loss** | Global | 0.001% | ✅ Healthy |
| **CPU Load** | router-us-east-2 | **88%** | ⚠️ Warning |

The network is generally healthy, but **router-us-east-2** is showing elevated CPU usage. This correlates with the scheduled backup job running at 13:00 UTC.

**Recommendation:**
Monitor the CPU load for the next 15 minutes. If it exceeds 95%, the auto-scaler will provision additional capacity.`
    }

    if (lowerQuery.includes("alert") || lowerQuery.includes("cpu")) {
        return `### Alert Analysis: High CPU in US-West

The alert **CPU_HIGH_WARN** was triggered at **14:32 UTC**.

**Root Cause Analysis:**
*   **Trigger**: BIRD routing daemon consumed 85% CPU.
*   **Correlation**: Simultaneous intake of 50k Full Table routes from peer \`AS2914\`.
*   **Impact**: Control plane latency increased by 150ms. Data plane forwarding remained unaffected.

**Suggested Remediation:**
Enable BGP generic receive offload (GRO) on the network interface to reduce CPU interrupts.

\`\`\`bash
# Enable GRO on eth0
ethtool -K eth0 gro on
\`\`\`
`
    }

    return `I can help you with that. I am a specialized Network Copilot trained on your infrastructure telemetry and configuration.

I can assist with:
*   **Traffic Analysis**: breakdown of ingress/egress flows
*   **Configuration Generation**: Cisco/Juniper/Linux syntax
*   **Troubleshooting**: Automated root cause analysis of alerts

Please provide more specific details about what you'd like to analyze.`
}

const ResponseMessage = ({ content }: { content: string }) => {
    const [copied, setCopied] = useState(false)

    const handleCopy = (text: string) => {
        navigator.clipboard.writeText(text)
        setCopied(true)
        setTimeout(() => setCopied(false), 2000)
    }

    return (
        <ReactMarkdown
            remarkPlugins={[remarkGfm]}
            components={{
                code({ node, inline, className, children, ...props }: any) {
                    const match = /language-(\w+)/.exec(className || '')
                    return !inline && match ? (
                        <div className="relative group my-4 rounded-lg overflow-hidden border border-border/50 bg-zinc-950">
                            <div className="absolute right-2 top-2 z-10 opacity-0 group-hover:opacity-100 transition-opacity">
                                <Button
                                    variant="ghost"
                                    size="icon"
                                    className="h-6 w-6 text-zinc-400 hover:text-white bg-zinc-800/50 hover:bg-zinc-700"
                                    onClick={() => handleCopy(String(children).replace(/\n$/, ''))}
                                >
                                    {copied ? <Check className="h-3 w-3" /> : <Copy className="h-3 w-3" />}
                                </Button>
                            </div>
                            <div className="bg-zinc-900/50 px-4 py-1.5 text-xs font-mono text-zinc-400 border-b border-border/50 uppercase select-none">
                                {match[1]}
                            </div>
                            <SyntaxHighlighter
                                {...props}
                                style={atomDark}
                                language={match[1]}
                                PreTag="div"
                                customStyle={{ margin: 0, padding: '1rem', background: 'transparent' }}
                            >
                                {String(children).replace(/\n$/, '')}
                            </SyntaxHighlighter>
                        </div>
                    ) : (
                        <code {...props} className={cn("bg-muted px-1.5 py-0.5 rounded font-mono text-xs", className)}>
                            {children}
                        </code>
                    )
                },
                table({ children }) {
                    return (
                        <div className="my-4 w-full overflow-y-auto rounded-lg border">
                            <table className="w-full text-sm">
                                {children}
                            </table>
                        </div>
                    )
                },
                thead({ children }) {
                    return <thead className="bg-muted/50 text-left font-medium">{children}</thead>
                },
                tbody({ children }) {
                    return <tbody className="divide-y">{children}</tbody>
                },
                tr({ children }) {
                    return <tr className="m-0 border-t p-0 even:bg-muted/10">{children}</tr>
                },
                th({ children }) {
                    return <th className="px-4 py-2 text-left font-bold [&[align=center]]:text-center [&[align=right]]:text-right">{children}</th>
                },
                td({ children }) {
                    return <td className="px-4 py-2 text-left [&[align=center]]:text-center [&[align=right]]:text-right">{children}</td>
                },
                ul({ children }) {
                    return <ul className="my-3 ml-6 list-disc [&>li]:mt-1">{children}</ul>
                },
                ol({ children }) {
                    return <ol className="my-3 ml-6 list-decimal [&>li]:mt-1">{children}</ol>
                },
                h3({ children }) {
                    return <h3 className="mt-6 mb-2 text-lg font-bold tracking-tight">{children}</h3>
                },
                p({ children }) {
                    return <p className="leading-7 [&:not(:first-child)]:mt-3">{children}</p>
                },
            }}
        >
            {content}
        </ReactMarkdown>
    )
}


export default function ChatPage() {
    const [messages, setMessages] = useState<Message[]>([])
    const [inputValue, setInputValue] = useState("")
    const [isTyping, setIsTyping] = useState(false)
    const messagesEndRef = useRef<HTMLDivElement>(null)

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: "smooth" })
    }

    useEffect(() => {
        scrollToBottom()
    }, [messages, isTyping])

    const handleSendMessage = (content: string) => {
        if (!content.trim()) return

        const newUserMessage: Message = {
            id: Date.now().toString(),
            role: "user",
            content: content,
            timestamp: new Date()
        }

        setMessages(prev => [...prev, newUserMessage])
        setInputValue("")
        setIsTyping(true)

        // Generate Context-Aware Response
        const responseText = generateAIResponse(content)

        // Mock AI Delay
        setTimeout(() => {
            const aiResponse: Message = {
                id: (Date.now() + 1).toString(),
                role: "assistant",
                content: responseText,
                timestamp: new Date()
            }
            setMessages(prev => [...prev, aiResponse])
            setIsTyping(false)
        }, 1500)
    }

    const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === "Enter" && !e.shiftKey) {
            e.preventDefault()
            handleSendMessage(inputValue)
        }
    }

    return (
        <div className="flex h-full flex-col">
            {/* Header */}
            <div className="border-b bg-card px-6 py-4">
                <div className="flex items-center gap-3">
                    <div className="flex h-10 w-10 items-center justify-center rounded-full bg-primary/10">
                        <Bot className="h-6 w-6 text-primary" />
                    </div>
                    <div>
                        <h2 className="text-lg font-semibold">Network Copilot</h2>
                        <p className="text-xs text-muted-foreground flex items-center">
                            <span className="mr-1.5 h-2 w-2 rounded-full bg-emerald-500 animate-pulse" />
                            Online
                        </p>
                    </div>
                </div>
            </div>

            {/* Chat Area */}
            <div className="flex-1 overflow-y-auto p-4 md:p-6 space-y-6">
                {messages.length === 0 ? (
                    <div className="flex h-full flex-col items-center justify-center space-y-8 opacity-90">
                        <div className="text-center space-y-2">
                            <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-2xl bg-muted">
                                <Sparkles className="h-8 w-8 text-muted-foreground" />
                            </div>
                            <h3 className="text-2xl font-bold tracking-tight">How can I help you today?</h3>
                            <p className="text-muted-foreground max-w-sm mx-auto">
                                I can analyze telemetry, troubleshoot alerts, and generate configuration snippets.
                            </p>
                        </div>

                        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 w-full max-w-3xl">
                            {SUGGESTED_QUERIES.map((query, i) => (
                                <button
                                    key={i}
                                    onClick={() => handleSendMessage(query.prompt)}
                                    className="flex flex-col items-start gap-2 rounded-lg border bg-card p-4 text-left transition-colors hover:bg-accent/50 hover:border-primary/50"
                                >
                                    <query.icon className="h-5 w-5 text-primary" />
                                    <span className="font-medium text-sm">{query.label}</span>
                                </button>
                            ))}
                        </div>
                    </div>
                ) : (
                    messages.map((msg) => (
                        <div
                            key={msg.id}
                            className={cn(
                                "flex w-full items-start gap-4 max-w-4xl mx-auto",
                                msg.role === "user" ? "flex-row-reverse" : "flex-row"
                            )}
                        >
                            <Avatar className={cn("h-8 w-8 mt-1", msg.role === "assistant" ? "bg-primary/10 border-primary/20" : "bg-muted")}>
                                {msg.role === "assistant" ? (
                                    <Bot className="h-5 w-5 text-primary" />
                                ) : (
                                    <User className="h-5 w-5 text-muted-foreground" />
                                )}
                            </Avatar>

                            <div
                                className={cn(
                                    "relative flex flex-col gap-2 rounded-lg px-4 py-3 text-sm shadow-sm max-w-[85%]",
                                    msg.role === "user"
                                        ? "bg-primary text-primary-foreground"
                                        : "bg-muted/40 border"
                                )}
                            >
                                <div className="leading-relaxed">
                                    {msg.role === "user" ? (
                                        <div className="whitespace-pre-wrap">{msg.content}</div>
                                    ) : (
                                        <ResponseMessage content={msg.content} />
                                    )}
                                </div>
                                <span className={cn(
                                    "text-[10px] opacity-70 mt-1",
                                    msg.role === "user" ? "text-primary-foreground/70" : "text-muted-foreground"
                                )}>
                                    {msg.timestamp.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                                </span>
                            </div>
                        </div>
                    ))
                )}

                {isTyping && (
                    <div className="flex w-full items-start gap-4 max-w-4xl mx-auto">
                        <Avatar className="h-8 w-8 bg-primary/10 border-primary/20">
                            <Bot className="h-5 w-5 text-primary" />
                        </Avatar>
                        <div className="flex items-center gap-1 rounded-lg border bg-muted/40 px-4 py-3 h-[46px]">
                            <span className="h-2 w-2 animate-bounce rounded-full bg-muted-foreground/40 [animation-delay:-0.3s]" />
                            <span className="h-2 w-2 animate-bounce rounded-full bg-muted-foreground/40 [animation-delay:-0.15s]" />
                            <span className="h-2 w-2 animate-bounce rounded-full bg-muted-foreground/40" />
                        </div>
                    </div>
                )}
                <div ref={messagesEndRef} />
            </div>

            {/* Input Area */}
            <div className="border-t bg-background p-4 md:px-6 md:pb-6">
                <div className="mx-auto max-w-4xl relative flex items-center gap-2">
                    <Button variant="ghost" size="icon" className="shrink-0 text-muted-foreground hover:text-foreground">
                        <StopCircle className="h-5 w-5" />
                    </Button>
                    <Input
                        value={inputValue}
                        onChange={(e) => setInputValue(e.target.value)}
                        onKeyDown={handleKeyDown}
                        placeholder="Ask anything about your network infrastructure..."
                        className="flex-1 bg-muted/50 shadow-none border-0 focus-visible:ring-1 focus-visible:ring-primary/50 focus-visible:bg-background transition-all"
                        disabled={isTyping}
                        autoFocus
                    />
                    <Button
                        size="icon"
                        onClick={() => handleSendMessage(inputValue)}
                        disabled={!inputValue.trim() || isTyping}
                        className={cn("shrink-0 transition-all", inputValue.trim() ? "" : "opacity-50")}
                    >
                        <Send className="h-4 w-4" />
                    </Button>
                </div>
                <div className="mt-2 text-center text-[10px] text-muted-foreground">
                    AI can make mistakes. Verify important configuration changes.
                </div>
            </div>
        </div>
    )
}
