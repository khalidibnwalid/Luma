import { Hash, Headphones, Mic, Settings } from "lucide-react"
import { Outlet, useOutletContext } from "react-router"
import { twJoin } from "tailwind-merge"
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar"
import { Button } from "~/components/ui/button"
import { ScrollArea, ScrollBar } from "~/components/ui/scroll-area"
import type { ServerLayoutContext } from "../layout"
import type { Route } from "./+types/layout"

export default function ServerLayout({ params: { serverId, roomId } }: Route.LoaderArgs) {
    const context = useOutletContext<ServerLayoutContext>()
    const { activeServer } = context

    return (
        <main className="flex h-screen w-full">
            <section className="w-60 bg-foreground/5 flex flex-col">
                <div className="p-4 shadow-md">
                    <h2 className="font-bold text-lg capitalize">{activeServer?.name}</h2>
                </div>

                <ScrollArea className="flex-1 overflow-y-auto">
                    <div className="p-2">
                        <ChatRoom activeRoom="welcome" roomName="Text Rooms" rooms={["welcome", "general", "rules", "announcements"]} />
                    </div>
                    <ScrollBar />
                </ScrollArea>
                {/* {readyState === ReadyState.OPEN &&
                    <div className="p-2 text-xs text-center text-foreground/50 bg-foreground/5">
                        <p>Connected</p>
                    </div>
                } */}
                <UserPanel user="User" subText="#000" />
            </section>

            <Outlet context={context} />
        </main >
    )
}

function ChatRoom({
    activeRoom,
    roomName,
    rooms
}: {
    activeRoom?: string,
    roomName: string
    rooms: string[]
}) {

    return (
        <div className="my-2">
            <h3 className="p-2 text-xs font-semibold text-foreground uppercase tracking-wide">{roomName}</h3>
            <div className="grid gap-1">
                {rooms.map((room) => (
                    <button
                        key={room}
                        className={twJoin(' w-full px-2 py-1 text-sm flex items-center gap-1.5 rounded duration-200 text-foreground/70',
                            activeRoom === room ? "bg-accent text-foreground/100" : "hover:bg-accent"
                        )}
                    >
                        <Hash className="size-4" />
                        <span className="capitalize">{room}</span>
                    </button>
                ))}
            </div>
        </div>
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