import {
    Avatar,
    AvatarFallback,
    AvatarImage,
} from "@/components/ui/avatar"

export function RecentActivity() {
    return (
        <div className="space-y-8">
            <div className="flex items-center">
                <Avatar className="h-9 w-9">
                    {/* With the new Avatar component, Image will hide if it fails, showing Fallback */}
                    <AvatarImage src="/avatars/01.png" alt="Avatar" />
                    <AvatarFallback>JD</AvatarFallback>
                </Avatar>
                <div className="ml-4 space-y-1">
                    <p className="text-sm font-medium leading-none">Jackson Lee</p>
                    <p className="text-sm text-muted-foreground">
                        modified bgp-config-12
                    </p>
                </div>
                <div className="ml-auto font-medium text-xs text-muted-foreground">2 min ago</div>
            </div>
            <div className="flex items-center">
                <Avatar className="h-9 w-9">
                    <AvatarImage src="/avatars/02.png" alt="Avatar" />
                    <AvatarFallback>OD</AvatarFallback>
                </Avatar>
                <div className="ml-4 space-y-1">
                    <p className="text-sm font-medium leading-none">Olivia Davis</p>
                    <p className="text-sm text-muted-foreground">
                        deployed new firmware v1.2
                    </p>
                </div>
                <div className="ml-auto font-medium text-xs text-muted-foreground">15 min ago</div>
            </div>
            <div className="flex items-center">
                <Avatar className="h-9 w-9">
                    <AvatarImage src="/avatars/03.png" alt="Avatar" />
                    <AvatarFallback>IA</AvatarFallback>
                </Avatar>
                <div className="ml-4 space-y-1">
                    <p className="text-sm font-medium leading-none">Isabella Nguyen</p>
                    <p className="text-sm text-muted-foreground">
                        restarted router-us-west-2
                    </p>
                </div>
                <div className="ml-auto font-medium text-xs text-muted-foreground">1h ago</div>
            </div>
        </div>
    )
}
