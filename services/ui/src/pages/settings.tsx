import { useState } from "react"
import {
    User,
    Bell,
    Shield,
    Key,
    Save,
    Mail,
    Slack,
    Smartphone
} from "lucide-react"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import { Switch } from "@/components/ui/switch"
import { cn } from "@/lib/utils"

const sidebarNavItems = [
    { title: "General", icon: User, id: "general" },
    { title: "Notifications", icon: Bell, id: "notifications" },
    { title: "Security", icon: Shield, id: "security" },
    { title: "API Access", icon: Key, id: "api" },
]

export default function SettingsPage() {
    const [activeTab, setActiveTab] = useState("general")

    return (
        <div className="flex flex-col space-y-8 lg:flex-row lg:space-x-12 lg:space-y-0 p-8 pt-6">
            <aside className="-mx-4 lg:w-1/5">
                <nav className="flex space-x-2 lg:flex-col lg:space-x-0 lg:space-y-1">
                    {sidebarNavItems.map((item) => (
                        <Button
                            key={item.id}
                            variant="ghost"
                            className={cn(
                                "justify-start hover:bg-muted",
                                activeTab === item.id ? "bg-muted" : ""
                            )}
                            onClick={() => setActiveTab(item.id)}
                        >
                            <item.icon className="mr-2 h-4 w-4" />
                            {item.title}
                        </Button>
                    ))}
                </nav>
            </aside>
            <div className="flex-1 lg:max-w-2xl">
                <div className="space-y-6">
                    <div>
                        <h3 className="text-lg font-medium">Settings</h3>
                        <p className="text-sm text-muted-foreground">
                            Manage your organization preferences and security policies.
                        </p>
                    </div>
                    <Separator />

                    {/* GENERAL SETTINGS */}
                    {activeTab === "general" && (
                        <div className="space-y-6">
                            <Card>
                                <CardHeader>
                                    <CardTitle>Organization Profile</CardTitle>
                                    <CardDescription>
                                        This is how others will see you on the site.
                                    </CardDescription>
                                </CardHeader>
                                <CardContent className="space-y-4">
                                    <div className="space-y-2">
                                        <label className="text-sm font-medium">Organization Name</label>
                                        <Input defaultValue="Acme Corp" />
                                    </div>
                                    <div className="space-y-2">
                                        <label className="text-sm font-medium">Support Email</label>
                                        <Input defaultValue="support@acme.com" />
                                    </div>
                                    <Button>
                                        <Save className="mr-2 h-4 w-4" /> Save Changes
                                    </Button>
                                </CardContent>
                            </Card>
                        </div>
                    )}

                    {/* NOTIFICATIONS SETTINGS */}
                    {activeTab === "notifications" && (
                        <div className="space-y-6">
                            <Card>
                                <CardHeader>
                                    <CardTitle>Alert Channels</CardTitle>
                                    <CardDescription>
                                        Configure how you receive critical system alerts.
                                    </CardDescription>
                                </CardHeader>
                                <CardContent className="space-y-4">
                                    <div className="flex items-center justify-between rounded-lg border p-4">
                                        <div className="space-y-0.5">
                                            <div className="flex items-center font-medium">
                                                <Mail className="mr-2 h-4 w-4" /> Email Notifications
                                            </div>
                                            <div className="text-sm text-muted-foreground">Receive daily summaries and critical alerts via email.</div>
                                        </div>
                                        <Switch defaultChecked />
                                    </div>
                                    <div className="flex items-center justify-between rounded-lg border p-4">
                                        <div className="space-y-0.5">
                                            <div className="flex items-center font-medium">
                                                <Slack className="mr-2 h-4 w-4" /> Slack Integration
                                            </div>
                                            <div className="text-sm text-muted-foreground">Post critical incidents to a dedicated Slack channel.</div>
                                        </div>
                                        <Switch />
                                    </div>
                                </CardContent>
                            </Card>
                        </div>
                    )}

                    {/* SECURITY SETTINGS */}
                    {activeTab === "security" && (
                        <div className="space-y-6">
                            <Card>
                                <CardHeader>
                                    <CardTitle>Access Policies</CardTitle>
                                    <CardDescription>
                                        Enforce security standards for all users.
                                    </CardDescription>
                                </CardHeader>
                                <CardContent className="space-y-4">
                                    <div className="flex items-center justify-between rounded-lg border p-4">
                                        <div className="space-y-0.5">
                                            <div className="flex items-center font-medium">
                                                <Smartphone className="mr-2 h-4 w-4" /> Enforce MFA
                                            </div>
                                            <div className="text-sm text-muted-foreground">Require Multi-Factor Authentication for all admins.</div>
                                        </div>
                                        <Switch defaultChecked />
                                    </div>
                                    <div className="space-y-2">
                                        <label className="text-sm font-medium">Session Timeout (minutes)</label>
                                        <Input type="number" defaultValue="60" />
                                    </div>
                                    <Button variant="destructive">Revoke All Sessions</Button>
                                </CardContent>
                            </Card>
                        </div>
                    )}

                    {/* API SETTINGS */}
                    {activeTab === "api" && (
                        <div className="space-y-6">
                            <Card>
                                <CardHeader>
                                    <CardTitle>API Keys</CardTitle>
                                    <CardDescription>
                                        Manage API keys for external integrations and CI/CD pipelines.
                                    </CardDescription>
                                </CardHeader>
                                <CardContent className="space-y-4">
                                    <div className="rounded-md bg-muted p-4">
                                        <div className="flex items-center justify-between">
                                            <div className="space-y-1">
                                                <p className="font-medium text-sm">Production Key (Primary)</p>
                                                <p className="text-xs text-muted-foreground font-mono">sk_live_...4f8a</p>
                                            </div>
                                            <Button variant="outline" size="sm">Roll Key</Button>
                                        </div>
                                    </div>
                                    <Button>
                                        <Key className="mr-2 h-4 w-4" /> Generate New Key
                                    </Button>
                                </CardContent>
                            </Card>
                        </div>
                    )}
                </div>
            </div>
        </div>
    )
}
