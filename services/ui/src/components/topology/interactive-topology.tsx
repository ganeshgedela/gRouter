"use client"

import { useState, useEffect } from 'react'
import { MapContainer, TileLayer, Marker, Polyline } from 'react-leaflet'
import L from 'leaflet'
import { Server, Wifi, X } from 'lucide-react'
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { cn } from "@/lib/utils"

// Fix for default marker icon in React Leaflet
delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
    iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png',
    iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png',
    shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png',
});

type NodeType = 'core' | 'edge' | 'gateway'
type LinkStatus = 'active' | 'degraded' | 'down'

interface Node {
    id: string
    label: string
    type: NodeType
    status: 'online' | 'offline' | 'alert'
    lat: number
    lng: number
    region: string
    load: number
    asn: string
}

interface Link {
    id: string
    source: string
    target: string
    status: LinkStatus
    bandwidth: string
    latency: string
}

const NODES: Node[] = [
    { id: 'n1', label: 'US-West-Core', type: 'core', status: 'online', lat: 45.5152, lng: -122.6784, region: 'Oregon', load: 45, asn: 'AS65001' },
    { id: 'n2', label: 'US-East-Core', type: 'core', status: 'online', lat: 37.7749, lng: -77.4128, region: 'Virginia', load: 62, asn: 'AS65001' },
    { id: 'n3', label: 'EU-West-GW', type: 'gateway', status: 'alert', lat: 51.5074, lng: -0.1278, region: 'London', load: 88, asn: 'AS65002' },
    { id: 'n4', label: 'AP-South-Edge', type: 'edge', status: 'online', lat: 1.3521, lng: 103.8198, region: 'Singapore', load: 23, asn: 'AS65003' },
    { id: 'n5', label: 'SA-East-Edge', type: 'edge', status: 'offline', lat: -23.5505, lng: -46.6333, region: 'Sao Paulo', load: 0, asn: 'AS65004' },
    { id: 'n6', label: 'EU-Central', type: 'core', status: 'online', lat: 50.1109, lng: 8.6821, region: 'Frankfurt', load: 34, asn: 'AS65002' },
]

const LINKS: Link[] = [
    { id: 'l1', source: 'n1', target: 'n2', status: 'active', bandwidth: '10 Gbps', latency: '45ms' },
    { id: 'l2', source: 'n2', target: 'n3', status: 'degraded', bandwidth: '5 Gbps', latency: '78ms' },
    { id: 'l3', source: 'n2', target: 'n5', status: 'active', bandwidth: '2.5 Gbps', latency: '112ms' },
    { id: 'l4', source: 'n3', target: 'n6', status: 'active', bandwidth: '10 Gbps', latency: '12ms' },
    { id: 'l5', source: 'n4', target: 'n6', status: 'active', bandwidth: '5 Gbps', latency: '160ms' },
    { id: 'l6', source: 'n1', target: 'n5', status: 'down', bandwidth: '1 Gbps', latency: '-' },
]

const createCustomIcon = (status: 'online' | 'offline' | 'alert') => {
    const colorClass = status === 'online' ? 'bg-emerald-500' : status === 'alert' ? 'bg-amber-500' : 'bg-red-500';
    return L.divIcon({
        className: 'custom-div-icon',
        html: `<div class="relative flex items-center justify-center w-6 h-6">
                 <span class="animate-ping absolute inline-flex h-full w-full rounded-full opacity-75 ${colorClass}"></span>
                 <span class="relative inline-flex items-center justify-center rounded-full h-5 w-5 ${colorClass} text-white border-2 border-white shadow-sm">
                    <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><rect width="20" height="8" x="2" y="2" rx="2" ry="2"/><rect width="20" height="8" x="2" y="14" rx="2" ry="2"/><line x1="6" x2="6.01" y1="6" y2="6"/><line x1="6" x2="6.01" y1="18" y2="18"/></svg>
                 </span>
               </div>`,
        iconSize: [24, 24],
        iconAnchor: [12, 12],
    })
}

export function InteractiveTopology() {
    const [selectedNode, setSelectedNode] = useState<Node | null>(null)
    const [isMounted, setIsMounted] = useState(false)

    useEffect(() => {
        setIsMounted(true)
    }, [])

    if (!isMounted) {
        return <div className="h-[600px] flex items-center justify-center text-muted-foreground">Loading Topology Map...</div>
    }

    return (
        <div className="relative h-full w-full bg-muted rounded-lg overflow-hidden border border-border shadow-2xl z-0">
            <MapContainer
                center={[20, 0]}
                zoom={2}
                minZoom={2}
                scrollWheelZoom={false}
                style={{ height: "100%", width: "100%", zIndex: 0 }}
                attributionControl={true}
            >
                {/* Light tiles - full visibility */}
                <TileLayer
                    attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                    url="https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png"
                />

                {/* Render BGP Links as Polylines */}
                {LINKS.map(link => {
                    const start = NODES.find(n => n.id === link.source)
                    const end = NODES.find(n => n.id === link.target)

                    if (!start || !end) return null

                    const color = link.status === 'active' ? '#10b981' :
                        link.status === 'degraded' ? '#f59e0b' :
                            '#ef4444'

                    return (
                        <Polyline
                            key={link.id}
                            positions={[[start.lat, start.lng], [end.lat, end.lng]]}
                            pathOptions={{
                                color: color,
                                weight: 2,
                                opacity: 0.6,
                                dashArray: link.status === 'degraded' ? '5, 10' : undefined
                            }}
                        />
                    )
                })}

                {/* Render Router Nodes as Markers */}
                {NODES.map(node => (
                    <Marker
                        key={node.id}
                        position={[node.lat, node.lng]}
                        icon={createCustomIcon(node.status)}
                        eventHandlers={{
                            click: () => setSelectedNode(node)
                        }}
                    />
                ))}
            </MapContainer>

            {/* Status Overlay - Bottom Left (Dashboard Style) */}
            <div className="absolute bottom-4 left-4 z-[400] bg-slate-900/80 backdrop-blur-md border border-slate-800 p-3 rounded-md pointer-events-none">
                <div className="text-xs font-mono text-slate-400 mb-2">TOPOLOGY STATUS</div>
                <div className="flex items-center gap-2 text-xs font-mono text-emerald-400">
                    <div className="h-1.5 w-1.5 rounded-full bg-emerald-500" />
                    {NODES.filter(n => n.status === 'online').length} Nodes Online
                </div>
                <div className="flex items-center gap-2 text-xs font-mono text-slate-300 mt-1">
                    <div className="h-1.5 w-1.5 rounded-full bg-blue-500" />
                    {LINKS.filter(l => l.status === 'active').length}/{LINKS.length} BGP Links Active
                </div>
            </div>

            {/* Inspector Panel (Slide-over) */}
            <div className={cn(
                "absolute top-0 right-0 h-full w-80 bg-background/95 backdrop-blur border-l shadow-2xl transition-transform duration-300 ease-in-out p-6 overflow-y-auto z-[500]",
                selectedNode ? "translate-x-0" : "translate-x-full"
            )}>
                {selectedNode && (
                    <div className="space-y-6">
                        <div className="flex items-center justify-between">
                            <h3 className="font-bold text-lg">Node Details</h3>
                            <Button variant="ghost" size="icon" onClick={() => setSelectedNode(null)} className="h-8 w-8">
                                <X className="h-4 w-4" />
                            </Button>
                        </div>

                        <div className="space-y-4">
                            <div className="flex items-center gap-4">
                                <div className="h-12 w-12 rounded-lg bg-muted flex items-center justify-center">
                                    <Server className="h-6 w-6 text-foreground" />
                                </div>
                                <div>
                                    <div className="font-bold">{selectedNode.label}</div>
                                    <div className="text-sm text-muted-foreground">{selectedNode.asn}</div>
                                </div>
                            </div>

                            <div className="grid grid-cols-2 gap-2">
                                <div className="p-3 rounded-md bg-muted/50 border">
                                    <div className="text-xs text-muted-foreground mb-1">Status</div>
                                    <Badge variant={selectedNode.status === 'online' ? 'default' : 'destructive'} className="uppercase text-[10px]">
                                        {selectedNode.status}
                                    </Badge>
                                </div>
                                <div className="p-3 rounded-md bg-muted/50 border">
                                    <div className="text-xs text-muted-foreground mb-1">Region</div>
                                    <div className="font-medium text-sm">{selectedNode.region}</div>
                                </div>
                            </div>

                            <div className="space-y-2">
                                <div className="flex justify-between text-sm">
                                    <span className="text-muted-foreground">CPU Load</span>
                                    <span className="font-mono">{selectedNode.load}%</span>
                                </div>
                                <div className="h-2 w-full bg-secondary rounded-full overflow-hidden">
                                    <div
                                        className={cn("h-full transition-all", selectedNode.load > 80 ? "bg-red-500" : "bg-primary")}
                                        style={{ width: `${selectedNode.load}%` }}
                                    />
                                </div>
                            </div>

                            <div className="pt-4 border-t space-y-2">
                                <h4 className="font-medium text-sm">BGP Peers</h4>
                                <div className="space-y-2">
                                    {LINKS.filter(l => l.source === selectedNode.id || l.target === selectedNode.id).map(link => {
                                        const peerId = link.source === selectedNode.id ? link.target : link.source
                                        const peer = NODES.find(n => n.id === peerId)
                                        return (
                                            <div key={link.id} className="flex items-center justify-between p-2 rounded border bg-background/50 text-sm">
                                                <div className="flex items-center gap-2">
                                                    <Wifi className={cn(
                                                        "h-3 w-3",
                                                        link.status === 'active' ? "text-emerald-500" :
                                                            link.status === 'degraded' ? "text-amber-500" :
                                                                "text-red-500"
                                                    )} />
                                                    <span>{peer?.label}</span>
                                                </div>
                                                <div className="font-mono text-xs text-muted-foreground">{link.latency}</div>
                                            </div>
                                        )
                                    })}
                                </div>
                            </div>

                            <Button className="w-full">
                                Manage Router
                            </Button>
                        </div>
                    </div>
                )}
            </div>
        </div>
    )
}
