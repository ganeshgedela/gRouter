// CIDR Utilities for network calculations and validations

import type { CIDRInfo } from "@/types/network"

/**
 * Validate CIDR notation format
 */
export function validateCIDR(cidr: string): boolean {
    const cidrRegex = /^(\d{1,3}\.){3}\d{1,3}\/\d{1,2}$/
    if (!cidrRegex.test(cidr)) return false

    const [ip, prefix] = cidr.split('/')
    const prefixNum = parseInt(prefix, 10)

    // Validate prefix length
    if (prefixNum < 0 || prefixNum > 32) return false

    // Validate each octet
    const octets = ip.split('.').map(Number)
    return octets.every(octet => octet >= 0 && octet <= 255)
}

/**
 * Convert IP address string to 32-bit integer
 */
function ipToInt(ip: string): number {
    return ip.split('.').reduce((acc, octet) => (acc << 8) + parseInt(octet, 10), 0) >>> 0
}

/**
 * Convert 32-bit integer to IP address string
 */
function intToIp(int: number): string {
    return [
        (int >>> 24) & 255,
        (int >>> 16) & 255,
        (int >>> 8) & 255,
        int & 255
    ].join('.')
}

/**
 * Calculate detailed CIDR information
 */
export function calculateCIDRInfo(cidr: string): CIDRInfo | null {
    if (!validateCIDR(cidr)) return null

    const [ip, prefix] = cidr.split('/')
    const prefixNum = parseInt(prefix, 10)
    const ipInt = ipToInt(ip)

    // Calculate subnet mask
    const mask = (0xFFFFFFFF << (32 - prefixNum)) >>> 0
    const wildcardMask = ~mask >>> 0

    // Calculate network address
    const networkInt = (ipInt & mask) >>> 0

    // Calculate broadcast address
    const broadcastInt = (networkInt | wildcardMask) >>> 0

    // Calculate total and usable IPs
    const totalIPs = Math.pow(2, 32 - prefixNum)
    const usableIPs = totalIPs > 2 ? totalIPs - 2 : totalIPs // Exclude network and broadcast

    // Calculate first and last usable IPs
    const firstUsableInt = totalIPs > 2 ? networkInt + 1 : networkInt
    const lastUsableInt = totalIPs > 2 ? broadcastInt - 1 : broadcastInt

    return {
        cidr,
        networkAddress: intToIp(networkInt),
        broadcastAddress: intToIp(broadcastInt),
        firstUsableIP: intToIp(firstUsableInt),
        lastUsableIP: intToIp(lastUsableInt),
        totalIPs,
        usableIPs,
        subnetMask: intToIp(mask),
        wildcardMask: intToIp(wildcardMask)
    }
}

/**
 * Check if a subnet CIDR is contained within a VPC CIDR
 */
export function isSubnetOf(subnetCIDR: string, vpcCIDR: string): boolean {
    if (!validateCIDR(subnetCIDR) || !validateCIDR(vpcCIDR)) return false

    const [subnetIP, subnetPrefix] = subnetCIDR.split('/')
    const [vpcIP, vpcPrefix] = vpcCIDR.split('/')

    const subnetPrefixNum = parseInt(subnetPrefix, 10)
    const vpcPrefixNum = parseInt(vpcPrefix, 10)

    // Subnet must have equal or larger prefix (smaller network)
    if (subnetPrefixNum < vpcPrefixNum) return false

    // Check if subnet's network address falls within VPC range
    const vpcMask = (0xFFFFFFFF << (32 - vpcPrefixNum)) >>> 0
    const subnetInt = ipToInt(subnetIP)
    const vpcInt = ipToInt(vpcIP)

    return ((subnetInt & vpcMask) >>> 0) === ((vpcInt & vpcMask) >>> 0)
}

/**
 * Suggest subnet divisions for a given VPC CIDR
 */
export function suggestSubnets(vpcCIDR: string, count: number): string[] {
    const vpcInfo = calculateCIDRInfo(vpcCIDR)
    if (!vpcInfo) return []

    const [vpcIP, vpcPrefix] = vpcCIDR.split('/')
    const vpcPrefixNum = parseInt(vpcPrefix, 10)

    // Calculate new prefix length for subnets
    const bitsNeeded = Math.ceil(Math.log2(count))
    const newPrefix = vpcPrefixNum + bitsNeeded

    if (newPrefix > 30) return [] // Not enough space for subnets

    const subnets: string[] = []
    const subnetSize = Math.pow(2, 32 - newPrefix)
    const vpcInt = ipToInt(vpcInfo.networkAddress)

    for (let i = 0; i < count && i < Math.pow(2, bitsNeeded); i++) {
        const subnetInt = vpcInt + (i * subnetSize)
        subnets.push(`${intToIp(subnetInt)}/${newPrefix}`)
    }

    return subnets
}

/**
 * Calculate IP utilization percentage
 */
export function calculateUtilization(usedIPs: number, totalIPs: number): number {
    if (totalIPs === 0) return 0
    return Math.round((usedIPs / totalIPs) * 100 * 10) / 10 // Round to 1 decimal
}

/**
 * Format CIDR for display with color coding suggestion
 */
export function getUtilizationColor(percentage: number): string {
    if (percentage < 70) return 'green'
    if (percentage < 90) return 'yellow'
    return 'red'
}
