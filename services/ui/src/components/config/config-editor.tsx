import { useState, useEffect } from "react"
import {
    Save,
    RotateCcw,
    FileDiff,
} from "lucide-react"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { cn } from "@/lib/utils"

// --- Diff Utility (LCS Simplification) ---
type DiffLine = {
    type: 'same' | 'added' | 'removed'
    content: string
    oldLineNo?: number
    newLineNo?: number
}

function computeDiff(oldText: string, newText: string): DiffLine[] {
    const oldLines = oldText.split('\n')
    const newLines = newText.split('\n')

    // Naive LCS implementation for demo
    const m = oldLines.length
    const n = newLines.length
    const dp = Array(m + 1).fill(0).map(() => Array(n + 1).fill(0))

    for (let i = 1; i <= m; i++) {
        for (let j = 1; j <= n; j++) {
            if (oldLines[i - 1] === newLines[j - 1]) {
                dp[i][j] = dp[i - 1][j - 1] + 1
            } else {
                dp[i][j] = Math.max(dp[i - 1][j], dp[i][j - 1])
            }
        }
    }

    let result: DiffLine[] = []
    let i = m, j = n
    while (i > 0 || j > 0) {
        if (i > 0 && j > 0 && oldLines[i - 1] === newLines[j - 1]) {
            result.push({ type: 'same', content: oldLines[i - 1], oldLineNo: i, newLineNo: j })
            i--; j--;
        } else if (j > 0 && (i === 0 || dp[i][j - 1] >= dp[i - 1][j])) {
            result.push({ type: 'added', content: newLines[j - 1], newLineNo: j })
            j--
        } else if (i > 0 && (j === 0 || dp[i][j - 1] < dp[i - 1][j])) {
            result.push({ type: 'removed', content: oldLines[i - 1], oldLineNo: i })
            i--
        }
    }
    return result.reverse()
}

// --- Components ---

interface ConfigEditorProps {
    initialConfig: string
    routerId: string
    onSave?: (newConfig: string) => void
}

const CONFIG_VERSIONS = [
    { id: 'v3-current', date: '2024-01-10 14:30', author: 'admin', label: 'Current' },
    { id: 'v2', date: '2024-01-09 09:15', author: 'system', label: 'Auto-Backup' },
    { id: 'v1', date: '2023-12-25 10:00', author: 'admin', label: 'Initial Deploy' },
]

const MOCK_CONFIGS: Record<string, string> = {
    'v3-current': `interface eth0
 description Uplink to Core
 ip address 10.0.1.1/30
 no shutdown
!
interface eth1
 description Local LAN
 ip address 192.168.1.1/24
 no shutdown
!
router bgp 65001
 neighbor 10.0.1.2 remote-as 65002
 neighbor 10.0.1.2 description Core-Router-01
!
ip route 0.0.0.0 0.0.0.0 10.0.1.2`,

    'v2': `interface eth0
 description Uplink to Core
 ip address 10.0.1.1/30
 no shutdown
!
interface eth1
 description Local LAN
 ip address 192.168.1.1/24
 no shutdown
!
router bgp 65001
 neighbor 10.0.1.2 remote-as 65002
!
ip route 0.0.0.0 0.0.0.0 10.0.1.2`
}

export function ConfigEditor({ initialConfig, routerId, onSave }: ConfigEditorProps) {
    const [mode, setMode] = useState<'edit' | 'diff'>('edit')
    const [currentConfig, setCurrentConfig] = useState(initialConfig || MOCK_CONFIGS['v3-current'])
    const [compareVersion, setCompareVersion] = useState('v2')
    const [diffLines, setDiffLines] = useState<DiffLine[]>([])

    useEffect(() => {
        if (mode === 'diff') {
            const oldConf = MOCK_CONFIGS[compareVersion] || ''
            setDiffLines(computeDiff(oldConf, currentConfig))
        }
    }, [mode, compareVersion, currentConfig])

    return (
        <Card className="h-[600px] flex flex-col overflow-hidden">
            <CardHeader className="flex flex-row items-center justify-between py-3 px-6 bg-muted/20 border-b">
                <div className="flex items-center gap-4">
                    <CardTitle className="text-base font-medium flex items-center gap-2">
                        <FileDiff className="h-4 w-4" />
                        Configuration
                    </CardTitle>
                    <div className="flex bg-muted rounded-lg p-0.5">
                        <button
                            onClick={() => setMode('edit')}
                            className={cn(
                                "px-3 py-1 text-xs font-medium rounded-md transition-all",
                                mode === 'edit' ? "bg-background shadow-sm text-foreground" : "text-muted-foreground hover:text-foreground"
                            )}
                        >
                            Edit
                        </button>
                        <button
                            onClick={() => setMode('diff')}
                            className={cn(
                                "px-3 py-1 text-xs font-medium rounded-md transition-all",
                                mode === 'diff' ? "bg-background shadow-sm text-foreground" : "text-muted-foreground hover:text-foreground"
                            )}
                        >
                            Compare
                        </button>
                    </div>
                </div>

                <div className="flex items-center gap-2">
                    {mode === 'diff' && (
                        <div className="flex items-center gap-2 text-sm mr-4">
                            <span className="text-muted-foreground">Compare with:</span>
                            <select
                                className="bg-transparent border rounded-md text-xs h-8 px-2"
                                value={compareVersion}
                                onChange={(e) => setCompareVersion(e.target.value)}
                            >
                                {CONFIG_VERSIONS.map(v => (
                                    <option key={v.id} value={v.id}>{v.id} ({v.label})</option>
                                ))}
                            </select>
                        </div>
                    )}
                    <Button variant="outline" size="sm" onClick={() => setCurrentConfig(MOCK_CONFIGS['v3-current'])}>
                        <RotateCcw className="mr-2 h-3 w-3" /> Reset
                    </Button>
                    <Button size="sm" onClick={() => onSave?.(currentConfig)}>
                        <Save className="mr-2 h-3 w-3" /> Save Changes
                    </Button>
                </div>
            </CardHeader>

            <CardContent className="flex-1 p-0 overflow-hidden relative">
                {mode === 'edit' ? (
                    <div className="h-full relative">
                        {/* Line Numbers */}
                        <div className="absolute left-0 top-0 bottom-0 w-12 bg-muted/10 border-r flex flex-col text-right text-xs font-mono text-muted-foreground pt-4 pr-2 select-none">
                            {currentConfig.split('\n').map((_, i) => (
                                <span key={i} className="leading-6">{i + 1}</span>
                            ))}
                        </div>
                        {/* Editor Area */}
                        <textarea
                            className="w-full h-full resize-none bg-transparent p-4 pl-14 font-mono text-sm leading-6 focus:outline-none"
                            value={currentConfig}
                            onChange={(e) => setCurrentConfig(e.target.value)}
                            spellCheck={false}
                        />
                    </div>
                ) : (
                    <div className="h-full overflow-auto font-mono text-sm">
                        {/* Diff View */}
                        <div className="divide-y divide-border/20">
                            {diffLines.map((line, i) => (
                                <div
                                    key={i}
                                    className={cn(
                                        "grid grid-cols-[3rem_3rem_1fr] min-w-full",
                                        line.type === 'added' ? "bg-emerald-500/10" :
                                            line.type === 'removed' ? "bg-red-500/10" : ""
                                    )}
                                >
                                    {/* Old Line No */}
                                    <div className="text-right px-2 py-0.5 text-xs text-muted-foreground select-none border-r border-border/10 bg-muted/5">
                                        {line.type !== 'added' && line.oldLineNo}
                                    </div>
                                    {/* New Line No */}
                                    <div className="text-right px-2 py-0.5 text-xs text-muted-foreground select-none border-r border-border/10">
                                        {line.type !== 'removed' && line.newLineNo}
                                    </div>
                                    {/* Content */}
                                    <div className="px-4 py-0.5 whitespace-pre break-all relative group">
                                        {line.type === 'added' && <span className="absolute left-1 text-emerald-500 font-bold">+</span>}
                                        {line.type === 'removed' && <span className="absolute left-1 text-red-500 font-bold">-</span>}
                                        <span className={cn(
                                            line.type === 'added' ? "text-emerald-700 dark:text-emerald-400" :
                                                line.type === 'removed' ? "text-red-700 dark:text-red-400 line-through opacity-70" :
                                                    "text-foreground"
                                        )}>
                                            {line.content || ' '}
                                        </span>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                )}
            </CardContent>
        </Card>
    )
}
