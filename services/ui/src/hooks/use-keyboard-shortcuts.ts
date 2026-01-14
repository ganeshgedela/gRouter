import { useEffect, useCallback } from 'react'

interface KeyboardShortcut {
    key: string
    ctrlKey?: boolean
    shiftKey?: boolean
    callback: () => void
    description: string
}

export function useKeyboardShortcuts(shortcuts: KeyboardShortcut[]) {
    const handleKeyDown = useCallback(
        (event: KeyboardEvent) => {
            const shortcut = shortcuts.find(
                (s) =>
                    s.key.toLowerCase() === event.key.toLowerCase() &&
                    (s.ctrlKey === undefined || s.ctrlKey === event.ctrlKey) &&
                    (s.shiftKey === undefined || s.shiftKey === event.shiftKey)
            )

            if (shortcut) {
                // Don't trigger if user is typing in an input
                const target = event.target as HTMLElement
                if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA') {
                    return
                }

                event.preventDefault()
                shortcut.callback()
            }
        },
        [shortcuts]
    )

    useEffect(() => {
        window.addEventListener('keydown', handleKeyDown)
        return () => window.removeEventListener('keydown', handleKeyDown)
    }, [handleKeyDown])
}

// Dashboard-specific shortcuts hook
export function useDashboardShortcuts({
    onRefresh,
    onToggleFilters,
    onSwitchTab,
    onShowHelp,
}: {
    onRefresh: () => void
    onToggleFilters: () => void
    onSwitchTab: (tab: number) => void
    onShowHelp: () => void
}) {
    useKeyboardShortcuts([
        {
            key: 'r',
            callback: onRefresh,
            description: 'Refresh data',
        },
        {
            key: 'f',
            callback: onToggleFilters,
            description: 'Toggle filters',
        },
        {
            key: '1',
            callback: () => onSwitchTab(0),
            description: 'Switch to Overview tab',
        },
        {
            key: '2',
            callback: () => onSwitchTab(1),
            description: 'Switch to Analytics tab',
        },
        {
            key: '3',
            callback: () => onSwitchTab(2),
            description: 'Switch to Reports tab',
        },
        {
            key: '?',
            shiftKey: true,
            callback: onShowHelp,
            description: 'Show keyboard shortcuts',
        },
    ])
}
