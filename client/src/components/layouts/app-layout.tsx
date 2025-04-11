import LayoutProviders from "@/components/providers/layout-providers"
import { Button } from "@/components/ui/button"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip"
import { useServersQuery } from "@/lib/queries/rooms-server"
import type { RoomsServer } from "@/types/rooms-server"
import { PlusIcon } from "lucide-react"
import Link from "next/link"
import { useRouter } from "next/router"
import { createContext, useContext } from "react"
import { twJoin } from "tailwind-merge"
import AddServerDialog from "../features/server/add-server-dialog"
import { ScrollArea, ScrollBar } from "../ui/scroll-area"
import UserPanel from "../features/auth/user-panel"
import { Separator } from "../ui/separator"

interface UserLayoutContext {
    activeServer?: RoomsServer,
    servers: RoomsServer[],
}

const UserLayoutContext = createContext({} as UserLayoutContext)
export const useAppLayoutContext = () => useContext(UserLayoutContext)

export default function AppLayout({
    children
}: {
    children: React.ReactNode
}) {
    const router = useRouter()
    const serverId = router.query.serverId as string

    const { data: servers, initCache } = useServersQuery()
    initCache()

    const activeServer = servers?.find((server) => server.id === serverId)

    const context: UserLayoutContext = {
        activeServer,
        servers: servers ?? []
    }

    return (
        <UserLayoutContext.Provider value={context}>
            <LayoutProviders>
                <div className="w-full h-screen flex">
                    <ScrollArea className="absolute">
                        <div className="w-16 py-3 pb-24 flex flex-col items-center gap-2 overflow-y-auto">
                            {servers?.map((server, i) => (
                                <TooltipProvider key={i + "server" + server.name}>
                                    <Tooltip>
                                        <TooltipTrigger asChild>
                                            <Link href={`/chat/${server.id}`}>
                                                <Button
                                                    variant="contrast"
                                                    className={twJoin(`size-12`, activeServer?.name === server.name && "bg-primary text-background")}
                                                >
                                                    <span className="capitalize font-bold text-lg">{server.name.charAt(0).toUpperCase()}</span>
                                                </Button>
                                            </Link>
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
                                        <AddServerDialog>
                                            <Button
                                                variant="ghost"
                                                className="rounded-2xl border-2  duration-200 p-0 size-12 flex items-center justify-center"
                                            >
                                                <PlusIcon />
                                            </Button>
                                        </AddServerDialog>
                                    </TooltipTrigger>
                                    <TooltipContent side="right">
                                        <p>Add a Server</p>
                                    </TooltipContent>
                                </Tooltip>
                            </TooltipProvider>
                        </div>
                        <ScrollBar />
                        <article className="absolute bottom-0 p-2 bg-background space-y-2" >
                            <Separator />
                            <UserPanel />
                        </article>
                    </ScrollArea>

                    {children}
                </div>
            </LayoutProviders>
        </UserLayoutContext.Provider>
    )
}