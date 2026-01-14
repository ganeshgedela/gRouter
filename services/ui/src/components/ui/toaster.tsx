import { useToast } from "@/lib/use-toast"
import { Toast } from "@/components/ui/toast"

export function Toaster() {
    const { toasts, dismiss } = useToast()

    return (
        <div className="fixed top-0 right-0 z-[100] flex max-h-screen w-full flex-col-reverse p-4 sm:bottom-0 sm:right-0 sm:top-auto sm:flex-col md:max-w-[420px] gap-2">
            {toasts.map(function ({ id, ...props }) {
                return (
                    <Toast key={id} {...props} onOpenChange={(open) => !open && dismiss(id)} />
                )
            })}
        </div>
    )
}
