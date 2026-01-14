// Network Types for VPCs and Subnets Management

export type CloudProvider = 'AWS' | 'GCP' | 'Azure' | 'On-Premise'
export type VPCStatus = 'active' | 'pending' | 'error'
export type SubnetType = 'public' | 'private'
export type SubnetStatus = 'active' | 'pending' | 'error'

export interface VPC {
    id: string
    name: string
    description: string
    cidr: string // e.g., "10.0.0.0/16"
    region: string
    cloudProvider: CloudProvider
    routerIds: string[] // Associated routers
    status: VPCStatus
    createdAt: string
    tags: Record<string, string>
}

export interface Subnet {
    id: string
    name: string
    vpcId: string
    cidr: string // e.g., "10.0.1.0/24"
    availabilityZone: string
    routerId: string // Primary router for this subnet
    type: SubnetType
    availableIPs: number
    usedIPs: number
    status: SubnetStatus
    createdAt: string
    tags: Record<string, string>
}

export interface CIDRInfo {
    cidr: string
    networkAddress: string
    broadcastAddress: string
    firstUsableIP: string
    lastUsableIP: string
    totalIPs: number
    usableIPs: number
    subnetMask: string
    wildcardMask: string
}
