import { PlusIcon } from "lucide-react"
import { Outlet } from "react-router"
import { twJoin } from "tailwind-merge"
import { ThemeProvider } from "~/components/providers/theme-provider"
import { Button } from "~/components/ui/button"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "~/components/ui/tooltip"
import type { Server } from "~/types/server"
import type { Route } from "../../+types/root"

export default function RoomLayout(params: Route.LoaderArgs) {
    console.log(params.params.serverId)
    const active = "general"
    const servers: Server[] = [
        { name: "general", rooms: [], id: "1" },
        { name: "music", rooms: [], id: "1" },
        { name: "gaming", rooms: [], id: "1" },
        { name: "programming", rooms: [], id: "1" },
        { name: "design", rooms: [], id: "1" },
    ]

    return (
        <div className="w-full h-screen flex">
            <ThemeProvider defaultTheme="dark">
                <section className="w-16 flex flex-col items-center py-3 gap-2 overflow-y-auto">

                    {servers.map((server) => (
                        <TooltipProvider key={server.name}>
                            <Tooltip>
                                <TooltipTrigger asChild>
                                    <Button
                                        variant="contrast"
                                        className={twJoin(`size-12`, active === server.name && "bg-primary text-background")}
                                    >
                                        <span className="capitalize font-bold text-lg">{server.name.charAt(0).toUpperCase()}</span>
                                    </Button>
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
                <Outlet />
            </ThemeProvider>
        </div>
    )
}