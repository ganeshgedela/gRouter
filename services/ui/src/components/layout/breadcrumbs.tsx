import { ChevronRight, Home } from "lucide-react"
import { useLocation, Link } from "react-router-dom"

export function Breadcrumbs() {
    const location = useLocation()
    const pathnames = location.pathname.split("/").filter((x) => x)

    return (
        <nav className="flex items-center space-x-1 text-sm text-muted-foreground">
            <Link
                to="/"
                className="flex items-center hover:text-foreground transition-colors"
                title="Dashboard"
            >
                <Home className="h-4 w-4" />
            </Link>
            {pathnames.length > 0 && <ChevronRight className="h-4 w-4" />}

            {pathnames.map((value, index) => {
                const to = `/${pathnames.slice(0, index + 1).join("/")}`;
                const isLast = index === pathnames.length - 1;
                const title = value.charAt(0).toUpperCase() + value.slice(1);

                return (
                    <div key={to} className="flex items-center space-x-1">
                        {isLast ? (
                            <span className="font-medium text-foreground">{title}</span>
                        ) : (
                            <Link
                                to={to}
                                className="hover:text-foreground transition-colors"
                            >
                                {title}
                            </Link>
                        )}
                        {!isLast && <ChevronRight className="h-4 w-4" />}
                    </div>
                );
            })}
        </nav>
    )
}
