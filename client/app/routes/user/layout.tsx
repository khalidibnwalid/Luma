import { PlusIcon } from "lucide-react"
import { useMemo } from "react"
import { NavLink, Outlet } from "react-router"
import { twJoin } from "tailwind-merge"
import LayoutProviders from "~/components/providers/layout-providers"
import { Button } from "~/components/ui/button"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "~/components/ui/tooltip"
import { useServersQuery } from "~/lib/queries/rooms-server"
import type { RoomsServer } from "~/types/rooms-server"
import type { Route } from "../../+types/root"

export interface ServerLayoutContext {
    activeServer?: RoomsServer,
    servers: RoomsServer[]
}

export default function Layout({ params }: Route.LoaderArgs) {
    const { data: servers, initCache, isSuccess } = useServersQuery()
    initCache()

    const activeServer = useMemo(() => servers?.find((server) => server.id === params.serverId),
        [params.serverId, JSON.stringify(servers)])

    if (!isSuccess) return <div>Loading...</div> // temporary

    const context: ServerLayoutContext = {
        activeServer,
        servers: servers ?? []
    }

    return (
        <LayoutProviders>
            <div className="w-full h-screen flex">
                <section className="w-16 flex flex-col items-center py-3 gap-2 overflow-y-auto">
                    {servers?.map((server) => (
                        <TooltipProvider key={server.name}>
                            <Tooltip>
                                <TooltipTrigger asChild>
                                    <NavLink
                                        to={`/server/${server.id}`}
                                    >
                                        <Button
                                            variant="contrast"
                                            className={twJoin(`size-12`, activeServer?.name === server.name && "bg-primary text-background")}
                                        >
                                            <span className="capitalize font-bold text-lg">{server.name.charAt(0).toUpperCase()}</span>
                                        </Button>
                                    </NavLink>
                                </TooltipTrigger>
                                <TooltipContent side="right">
                                    <p className="capitalize">{server.name}</p>
                                </TooltipContent>
                            </Tooltip>
                        </TooltipProvider>
                    ))}

                    <TooltipProvider>
                        <Tooltip>
                            <TooltipTrigger asChild>
                                <Button
                                    variant="ghost"
                                    className="rounded-2xl border-2  duration-200 p-0 size-12 flex items-center justify-center"
                                >
                                    <PlusIcon />
                                </Button>
                            </TooltipTrigger>
                            <TooltipContent side="right">
                                <p>Add a Server</p>
                            </TooltipContent>
                        </Tooltip>
                    </TooltipProvider>
                </section>
                <Outlet context={context} />
            </div>
        </LayoutProviders>
    )
}