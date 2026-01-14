import { PolicyManager } from "@/components/policy/policy-manager"

export default function PoliciesPage() {
    return (
        <div className="flex-1 space-y-4 p-8 pt-6">
            <div className="flex items-center justify-between space-y-2">
                <h2 className="text-3xl font-bold tracking-tight">Access Control Policies</h2>
            </div>
            <PolicyManager />
        </div>
    )
}
