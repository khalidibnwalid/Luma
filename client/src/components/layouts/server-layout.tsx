import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area"
import { useRoomsQuery } from "@/lib/queries/rooms"
import type { Room } from "@/types/room"
import { RoomsServer } from "@/types/rooms-server"
import { Hash } from "lucide-react"
import Link from "next/link"
import { useRouter } from "next/router"
import { createContext, useContext } from "react"
import { twJoin } from "tailwind-merge"
import ServerSidebarContextMenu from "../features/server/sidebar-context-menu"
import { useAppLayoutContext } from "./app-layout"

interface ServerLayoutContext {
    rooms: Room[]
    activeServer: RoomsServer
}

const ServerContext = createContext({} as ServerLayoutContext)
export const useServerContext = () => useContext(ServerContext)

export default function ServerLayout({
    children
}: {
    children: React.ReactNode
}) {
    const router = useRouter()
    const roomId = router.query.roomId as string

    const context = useAppLayoutContext()
    const { activeServer } = context

    const { data: rooms, initCache } = useRoomsQuery(activeServer?.id)
    initCache()

    const wrappedContext: ServerLayoutContext = {
        rooms: rooms ?? [],
        activeServer: activeServer!
    }

    if (!activeServer) return <div>Server not found</div>

    return (
        <ServerContext.Provider value={wrappedContext}>
            <main className="flex h-screen w-full">
                <ServerSidebarContextMenu className="w-60 bg-foreground/5 flex flex-col">
                    <div className="p-4 shadow-md">
                        <h2 className="font-bold text-lg capitalize">{activeServer?.name}</h2>
                    </div>

                    <ScrollArea className="flex-1 overflow-y-auto">
                        <ChatRooms activeRoomId={roomId} rooms={rooms} />
                        <ScrollBar />
                    </ScrollArea>
                </ServerSidebarContextMenu>
                {children}
            </main >
        </ServerContext.Provider>
    )
}

function ChatRooms({
    activeRoomId,
    rooms = []
}: {
    activeRoomId?: string
    rooms?: Room[]
}) {
    const groups = rooms?.reduce((acc, room) => {
        if (!acc[room.groupName])
            acc[room.groupName] = []
        acc[room.groupName].push(room)
        return acc
    }, {} as { [key: string]: Room[] }) || []

    return (
        <div className="p-2">
            {Object.entries(groups).map(([groupName, rooms]) => (
                <div key={groupName} className="my-2">
                    <h3 className="my-2 text-sm text-foreground/50 font-medium">{groupName}</h3>
                    <div className="grid gap-1">
                        {rooms.map(room => (
                            <RoomPanel key={room.id} room={room} isActive={room.id === activeRoomId} />
                        ))}
                    </div>
                </div>
            ))}
        </div>
    )
}

function RoomPanel({
    isActive,
    room
}: {
    isActive: boolean,
    room: Room
}) {
    const { activeServer } = useAppLayoutContext()

    return (
        <Link href={`/chat/${activeServer!.id}/${room.id}`}         >
            <button
                className={twJoin(' w-full px-2 py-1 text-sm flex items-center gap-1.5 rounded duration-200 text-foreground/70',
                    isActive ? "bg-accent text-foreground/100" : "hover:bg-accent"
                )}
            >
                <Hash className="size-4" />
                <span className="capitalize">{room.name}</span>
            </button>
        </Link>

    )
}