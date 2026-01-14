"use client"

import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet'
import L from 'leaflet'
import { useEffect, useState } from 'react'

// Fix for default marker icon in React Leaflet
delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
    iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png',
    iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png',
    shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png',
});

// Custom marker styles using DivIcon for custom CSS pulsing
const createCustomIcon = (status: 'active' | 'warning' | 'error') => {
    const colorClass = status === 'active' ? 'bg-emerald-500' : status === 'warning' ? 'bg-amber-500' : 'bg-red-500';
    return L.divIcon({
        className: 'custom-div-icon',
        html: `<div class="relative flex items-center justify-center w-4 h-4">
                 <span class="animate-ping absolute inline-flex h-full w-full rounded-full opacity-75 ${colorClass}"></span>
                 <span class="relative inline-flex rounded-full h-3 w-3 ${colorClass} border-2 border-white"></span>
               </div>`,
        iconSize: [16, 16],
        iconAnchor: [8, 8],
        popupAnchor: [0, -10]
    })
}

const regions = [
    { id: "us-east", lat: 37.7749, lng: -77.4128, label: "US-East (N. Virginia)", status: "active", latency: "24ms" },
    { id: "us-west", lat: 45.5152, lng: -122.6784, label: "US-West (Oregon)", status: "active", latency: "45ms" },
    { id: "eu-central", lat: 50.1109, lng: 8.6821, label: "EU-Central (Frankfurt)", status: "active", latency: "82ms" },
    { id: "ap-southeast", lat: 1.3521, lng: 103.8198, label: "AP-Southeast (Singapore)", status: "warning", latency: "145ms" },
    { id: "uk-london", lat: 51.5074, lng: -0.1278, label: "UK-South (London)", status: "active", latency: "78ms" },
    { id: "sa-east", lat: -23.5505, lng: -46.6333, label: "SA-East (São Paulo)", status: "active", latency: "112ms" },
    { id: "ap-northeast", lat: 35.6762, lng: 139.6503, label: "AP-Northeast (Tokyo)", status: "active", latency: "130ms" },
]

export function MapView() {
    const [isMounted, setIsMounted] = useState(false)

    useEffect(() => {
        setIsMounted(true)
    }, [])

    if (!isMounted) {
        return <div className="h-[400px] w-full bg-slate-950 rounded-lg flex items-center justify-center text-slate-500">Loading Map...</div>
    }

    return (
        <div className="relative h-[400px] w-full bg-muted rounded-lg overflow-hidden border border-border shadow-2xl z-0">
            <MapContainer
                center={[20, 0]}
                zoom={2}
                scrollWheelZoom={false}
                style={{ height: "100%", width: "100%", zIndex: 0 }}
                attributionControl={true}
            >
                {/* Light tiles - full visibility */}
                <TileLayer
                    attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                    url="https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png"
                />

                {regions.map((region) => (
                    <Marker
                        key={region.id}
                        position={[region.lat, region.lng]}
                        icon={createCustomIcon(region.status as any)}
                    >
                        <Popup className="custom-popup">
                            <div className="text-sm font-semibold">{region.label}</div>
                            <div className="text-xs text-muted-foreground">Latency: {region.latency}</div>
                            <div className="text-xs capitalize mt-1">Status: <span className={region.status === 'active' ? 'text-emerald-600' : 'text-amber-600'}>{region.status}</span></div>
                        </Popup>
                    </Marker>
                ))}
            </MapContainer>

            <div className="absolute bottom-4 left-4 z-[400] bg-slate-900/80 backdrop-blur-md border border-slate-800 p-3 rounded-md pointer-events-none">
                <div className="text-xs font-mono text-slate-400 mb-2">NETWORK STATUS</div>
                <div className="flex items-center gap-2 text-xs font-mono text-emerald-400">
                    <div className="h-1.5 w-1.5 rounded-full bg-emerald-500" />
                    Global Mesh: STABLE
                </div>
                <div className="flex items-center gap-2 text-xs font-mono text-slate-300 mt-1">
                    <div className="h-1.5 w-1.5 rounded-full bg-blue-500" />
                    Avg Latency: 42ms
                </div>
            </div>
        </div>
    )
}
