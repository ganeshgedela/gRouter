import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { SecurityDashboard } from "@/components/security/security-dashboard"
import { AuditLog } from "@/components/security/audit-log"
import { DDoSControls } from "@/components/security/ddos-controls"

export default function SecurityPage() {
    return (
        <div className="flex-1 space-y-4 p-8 pt-6">
            <div className="flex items-center justify-between space-y-2">
                <div>
                    <h2 className="text-3xl font-bold tracking-tight">Security Operations</h2>
                    <p className="text-muted-foreground">
                        Real-time threat monitoring and mitigation center.
                    </p>
                </div>
            </div>

            <Tabs defaultValue="overview" className="space-y-4">
                <TabsList>
                    <TabsTrigger value="overview">SOC Overview</TabsTrigger>
                    <TabsTrigger value="audit">Audit Logs</TabsTrigger>
                    <TabsTrigger value="mitigation">DDoS Mitigation</TabsTrigger>
                </TabsList>

                <TabsContent value="overview" className="space-y-4">
                    <SecurityDashboard />
                </TabsContent>

                <TabsContent value="audit" className="space-y-4">
                    <AuditLog />
                </TabsContent>

                <TabsContent value="mitigation" className="space-y-4">
                    <DDoSControls />
                </TabsContent>
            </Tabs>
        </div>
    )
}
