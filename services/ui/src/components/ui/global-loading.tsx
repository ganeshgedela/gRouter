import { useEffect, useState } from "react"
import { useLocation } from "react-router-dom"

export function GlobalLoading() {
    const location = useLocation()
    const [loading, setLoading] = useState(false)
    const [progress, setProgress] = useState(0)

    useEffect(() => {
        // Trigger loading on route change
        setLoading(true)
        setProgress(10)

        // Simulate fast load
        const timer1 = setTimeout(() => setProgress(60), 100)
        const timer2 = setTimeout(() => setProgress(90), 300)
        const timer3 = setTimeout(() => {
            setProgress(100)
            setTimeout(() => setLoading(false), 200)
        }, 500)

        return () => {
            clearTimeout(timer1)
            clearTimeout(timer2)
            clearTimeout(timer3)
        }
    }, [location.pathname])

    if (!loading) return null

    return (
        <div className="fixed top-0 left-0 right-0 z-[100] h-1 w-full bg-transparent">
            <div
                className="h-full bg-primary transition-all duration-300 ease-out shadow-[0_0_10px_#3b82f6]"
                style={{ width: `${progress}%` }}
            />
        </div>
    )
}
