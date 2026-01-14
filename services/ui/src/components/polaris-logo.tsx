import { cn } from "@/lib/utils"

interface LogoProps extends React.SVGProps<SVGSVGElement> {
    className?: string;
}

export function PolarisLogo({ className, ...props }: LogoProps) {
    return (
        <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 100 100"
            className={cn("w-10 h-10 text-primary", className)}
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
            {...props}
        >
            {/* Defs for gradients/effects */}
            <defs>
                <filter id="glow" x="-20%" y="-20%" width="140%" height="140%">
                    <feGaussianBlur stdDeviation="2" result="blur" />
                    <feComposite in="SourceGraphic" in2="blur" operator="over" />
                </filter>
            </defs>

            {/* Central North Star - The Core Orchestrator */}
            <path
                d="M50 20 L55 45 L80 50 L55 55 L50 80 L45 55 L20 50 L45 45 Z"
                className="fill-primary/20 animate-pulse"
                strokeWidth="2"
            />

            {/* Inner Core */}
            <circle cx="50" cy="50" r="4" className="fill-primary" />

            {/* Orbiting Nodes - Representing Clouds/Routers */}
            {/* Node 1 (Top Right) */}
            <circle cx="75" cy="25" r="3" className="fill-current opacity-80" />
            <path d="M75 25 L58 42" className="opacity-40" strokeDasharray="4 2" />

            {/* Node 2 (Bottom Left) */}
            <circle cx="25" cy="75" r="3" className="fill-current opacity-80" />
            <path d="M25 75 L42 58" className="opacity-40" strokeDasharray="4 2" />

            {/* Node 3 (Bottom Right - Far) */}
            <circle cx="80" cy="80" r="2" className="fill-current opacity-60" />
            <path d="M80 80 L55 55" className="opacity-20" />

            {/* Orbital Ring - Connectivity */}
            <path
                d="M85 50 A 35 35 0 0 1 15 50"
                className="opacity-20"
                strokeWidth="1"
            />
            <path
                d="M15 50 A 35 35 0 0 1 25 25"
                className="opacity-20"
                strokeWidth="1"
            />
        </svg>
    )
}
