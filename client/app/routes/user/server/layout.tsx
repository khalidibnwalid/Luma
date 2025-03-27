import { Hash, Headphones, Mic, Settings } from "lucide-react"
import { NavLink, Outlet, useOutletContext } from "react-router"
import { twJoin } from "tailwind-merge"
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar"
import { Button } from "~/components/ui/button"
import { ScrollArea, ScrollBar } from "~/components/ui/scroll-area"
import { useRoomsQuery } from "~/lib/queries/rooms"
import type { Room } from "~/types/room"
import type { ServerLayoutContext } from "../layout"
import type { Route } from "./+types/layout"

export default function ServerLayout({ params: { roomId } }: Route.LoaderArgs) {
    const context = useOutletContext<ServerLayoutContext>()
    const { activeServer } = context
    if (!activeServer) return <div>Server not found</div>

    const { data: rooms } = useRoomsQuery(activeServer.id)
    const wrappedContext = { ...context, rooms }

    return (
        <main className="flex h-screen w-full">
            <section className="w-60 bg-foreground/5 flex flex-col">
                <div className="p-4 shadow-md">
                    <h2 className="font-bold text-lg capitalize">{activeServer?.name}</h2>
                </div>

                <ScrollArea className="flex-1 overflow-y-auto">
                    <ChatRoom activeRoomId={roomId} rooms={rooms} />
                    <ScrollBar />
                </ScrollArea>
                <UserPanel user="User" subText="#000" />
            </section>
            <Outlet context={wrappedContext} />
        </main >
    )
}

function ChatRoom({
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
    const {activeServer} = useOutletContext<ServerLayoutContext>()

    return (
        <NavLink
            to={`/server/${activeServer!.id}/${room.id}`}
        >
            <button
                className={twJoin(' w-full px-2 py-1 text-sm flex items-center gap-1.5 rounded duration-200 text-foreground/70',
                    isActive ? "bg-accent text-foreground/100" : "hover:bg-accent"
                )}
            >
                <Hash className="size-4" />
                <span className="capitalize">{room.name}</span>
            </button>
        </NavLink>

    )
}

function UserPanel({
    user,
    subText
}: {
    user: string
    subText: string
}) {
    return (
        <div className="p-2 flex items-center gap-2">
            <Avatar className="size-8">
                <AvatarImage src="" alt="User" />
                <AvatarFallback>U</AvatarFallback>
            </Avatar>

            <div className="flex-1 min-w-0">
                <p className="text-sm font-medium">{user}</p>
                <p className="text-xs text-foreground/50 truncate">{subText}</p>
            </div>

            <div className="flex gap-1">
                <Button variant="ghost" size="icon" className="size-8 rounded-full">
                    <Mic className="size-4" />
                </Button>
                <Button variant="ghost" size="icon" className="size-8 rounded-full">
                    <Headphones className="size-4" />
                </Button>
                <Button variant="ghost" size="icon" className="size-8 rounded-full">
                    <Settings className="size-4" />
                </Button>
            </div>
        </div>
    )
}