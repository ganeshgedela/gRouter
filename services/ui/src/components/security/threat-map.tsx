"use client"

import { useEffect, useState } from 'react'
import { MapContainer, TileLayer, Marker, Polyline, CircleMarker } from 'react-leaflet'
import L from 'leaflet'
import { Shield, AlertTriangle } from 'lucide-react'

// Fix for default marker icon
delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
    iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png',
    iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png',
    shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png',
});

interface ThreatSource {
    id: string
    country: string
    lat: number
    lng: number
    severity: 'high' | 'medium' | 'low'
    attackType: string
}

const THREAT_SOURCES: ThreatSource[] = [
    { id: 't1', country: 'Russia', lat: 55.7558, lng: 37.6173, severity: 'high', attackType: 'DDoS Flood' },
    { id: 't2', country: 'China', lat: 39.9042, lng: 116.4074, severity: 'high', attackType: 'Botnet C2' },
    { id: 't3', country: 'Brazil', lat: -23.5505, lng: -46.6333, severity: 'medium', attackType: 'SQL Injection' },
    { id: 't4', country: 'India', lat: 28.6139, lng: 77.2090, severity: 'low', attackType: 'Port Scan' },
    { id: 't5', country: 'Iran', lat: 35.6892, lng: 51.3890, severity: 'medium', attackType: 'SSH Brute' },
]

// Protected infrastructure (US location as example)
const PROTECTED_NODE = { lat: 37.7749, lng: -97.0, name: 'US Infrastructure' }

const createThreatIcon = (severity: 'high' | 'medium' | 'low') => {
    const color = severity === 'high' ? 'bg-red-500' : severity === 'medium' ? 'bg-amber-500' : 'bg-yellow-500'
    const size = severity === 'high' ? 'h-4 w-4' : 'h-3 w-3'

    return L.divIcon({
        className: 'threat-marker',
        html: `<div class="relative flex items-center justify-center ${size}">
                 <span class="animate-ping absolute inline-flex h-full w-full rounded-full opacity-75 ${color}"></span>
                 <span class="relative inline-flex rounded-full h-full w-full ${color} border border-white shadow-lg"></span>
               </div>`,
        iconSize: severity === 'high' ? [16, 16] : [12, 12],
        iconAnchor: severity === 'high' ? [8, 8] : [6, 6],
    })
}

const createProtectedIcon = () => {
    return L.divIcon({
        className: 'protected-marker',
        html: `<div class="relative flex items-center justify-center w-6 h-6">
                 <span class="relative inline-flex rounded-full h-5 w-5 bg-emerald-500 border-2 border-white shadow-lg items-center justify-center">
                    <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
                 </span>
               </div>`,
        iconSize: [24, 24],
        iconAnchor: [12, 12],
    })
}

export function ThreatMap() {
    const [isMounted, setIsMounted] = useState(false)

    useEffect(() => {
        setIsMounted(true)
    }, [])

    if (!isMounted) {
        return (
            <div className="h-full flex items-center justify-center bg-zinc-950/50 rounded text-muted-foreground">
                Loading Threat Map...
            </div>
        )
    }

    return (
        <div className="relative h-full w-full bg-muted rounded overflow-hidden">
            <MapContainer
                center={[30, 0]}
                zoom={2}
                minZoom={2}
                maxZoom={5}
                scrollWheelZoom={false}
                style={{ height: "100%", width: "100%", zIndex: 0 }}
                attributionControl={false}
            >
                {/* Dark map tiles */}
                <TileLayer
                    url="https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png"
                />

                {/* Attack arcs from threats to protected node */}
                {THREAT_SOURCES.map((threat) => {
                    const color = threat.severity === 'high' ? '#ef4444' :
                        threat.severity === 'medium' ? '#f59e0b' :
                            '#eab308'

                    return (
                        <Polyline
                            key={`arc-${threat.id}`}
                            positions={[
                                [threat.lat, threat.lng],
                                [PROTECTED_NODE.lat, PROTECTED_NODE.lng]
                            ]}
                            pathOptions={{
                                color: color,
                                weight: threat.severity === 'high' ? 2 : 1,
                                opacity: 0.4,
                                dashArray: '5, 10',
                            }}
                        />
                    )
                })}

                {/* Threat source markers */}
                {THREAT_SOURCES.map((threat) => (
                    <Marker
                        key={threat.id}
                        position={[threat.lat, threat.lng]}
                        icon={createThreatIcon(threat.severity)}
                    />
                ))}

                {/* Protected infrastructure marker */}
                <Marker
                    position={[PROTECTED_NODE.lat, PROTECTED_NODE.lng]}
                    icon={createProtectedIcon()}
                />

                {/* Pulsing circles around protected node */}
                <CircleMarker
                    center={[PROTECTED_NODE.lat, PROTECTED_NODE.lng]}
                    radius={30}
                    pathOptions={{
                        color: '#10b981',
                        fillColor: '#10b981',
                        fillOpacity: 0.1,
                        weight: 1,
                        opacity: 0.3,
                    }}
                />
            </MapContainer>

            {/* Stats Overlay - Bottom Left */}
            <div className="absolute bottom-4 left-4 z-[400] space-y-2">
                <div className="bg-red-500/90 backdrop-blur-sm border border-red-500 px-3 py-1.5 rounded-md flex items-center gap-2 text-white font-mono text-xs">
                    <AlertTriangle className="h-3 w-3" />
                    <span className="font-bold">{THREAT_SOURCES.filter(t => t.severity === 'high').length}</span>
                    <span>High Severity</span>
                </div>
                <div className="bg-emerald-500/90 backdrop-blur-sm border border-emerald-500 px-3 py-1.5 rounded-md flex items-center gap-2 text-white font-mono text-xs">
                    <Shield className="h-3 w-3" />
                    <span className="font-bold">Protected</span>
                </div>
            </div>

            {/* Legend */}
            <div className="absolute bottom-4 right-4 z-[400] bg-background/90 backdrop-blur-md border p-3 rounded-md text-xs space-y-1.5">
                <div className="font-semibold mb-2 text-muted-foreground">THREAT LEGEND</div>
                <div className="flex items-center gap-2">
                    <div className="w-3 h-3 rounded-full bg-red-500 border border-white animate-pulse"></div>
                    <span>Critical Threat</span>
                </div>
                <div className="flex items-center gap-2">
                    <div className="w-2.5 h-2.5 rounded-full bg-amber-500 border border-white"></div>
                    <span>Medium Threat</span>
                </div>
                <div className="flex items-center gap-2">
                    <div className="w-2 h-2 rounded-full bg-yellow-500 border border-white"></div>
                    <span>Low Threat</span>
                </div>
                <div className="border-t pt-1.5 mt-1.5 flex items-center gap-2">
                    <div className="w-4 h-4 rounded-full bg-emerald-500 border-2 border-white flex items-center justify-center">
                        <Shield className="h-2 w-2 text-white" />
                    </div>
                    <span>Protected Node</span>
                </div>
            </div>

            {/* Live indicator */}
            <div className="absolute top-4 right-4 z-[400] bg-background/90 backdrop-blur-md border px-2 py-1 rounded-md flex items-center gap-1.5">
                <div className="w-1.5 h-1.5 rounded-full bg-red-500 animate-pulse"></div>
                <span className="text-[10px] font-mono text-muted-foreground uppercase">Live</span>
            </div>
        </div>
    )
}
