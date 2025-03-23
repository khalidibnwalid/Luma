import { Bell, Hash, Headphones, Mic, PencilIcon, Plus, Search, Settings, TrashIcon, Users } from "lucide-react"
import { useEffect, useRef, useState } from "react"
import useWebSocket, { ReadyState } from 'react-use-websocket'
import { twJoin } from "tailwind-merge"
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar"
import { Button } from "~/components/ui/button"
import { Input } from "~/components/ui/input"
import { ScrollArea, ScrollBar } from "~/components/ui/scroll-area"
import type { MessageResponse } from "~/types/message"
import type { User } from "~/types/user"
import type { Route } from "./+types/room"

export default function RoomPage({ params: { serverId, roomId } }: Route.LoaderArgs) {
    const socketUrl = 'ws://localhost:8080/v1/room/' + roomId + "?jwt=" + localStorage.getItem("token")

    const input = useRef<HTMLInputElement>(null);

    const [messageHistory, setMessageHistory] = useState<MessageResponse[]>([]);
    const { sendMessage, lastMessage, readyState } = useWebSocket(socketUrl);

    useEffect(() => {
        if (lastMessage !== null) setMessageHistory((prev) => prev.concat(JSON.parse(lastMessage.data)));
    }, [lastMessage]);

    function sendFormattedMessage() {
        const value = input.current?.value
        if (!value) return
        sendMessage(JSON.stringify({ message: value } as MessageResponse))
        input.current!.value = ""
    }

    function inputOnKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
        if (e.key === "Enter") sendFormattedMessage()
    }

    /// mocks
    const active = "general"
    const activeRoom = "welcome"

    const onlineUsers = ["User 1", "User 2", "User 3", "User 4", "User 5"]

    return (
        <main className="flex h-screen w-full">
            <section className="w-60 bg-foreground/5 flex flex-col">
                <div className="p-4 shadow-md">
                    <h2 className="font-bold text-lg capitalize">{active}</h2>
                </div>

                <ScrollArea className="flex-1 overflow-y-auto">
                    <div className="p-2">
                        <ChatRoom activeRoom="welcome" roomName="Text Rooms" rooms={["welcome", "general", "rules", "announcements"]} />
                        <ChatRoom roomName="Voice Rooms" rooms={["General", "Music", "Gaming"]} />
                    </div>
                    <ScrollBar />
                </ScrollArea>
                {readyState === ReadyState.OPEN &&
                    <div className="p-2 text-xs text-center text-foreground/50 bg-foreground/5">
                        <p>Connected</p>
                    </div>
                }
                <UserPanel user="User" subText="#000" />
            </section>

            <section className="flex-1 flex flex-col">
                <div className="h-12 flex items-center px-4">
                    <div className="flex items-center">
                        <Hash className="size-5 text-foreground/50 mr-1" />
                        <h3 className="font-bold text-white capitalize">{activeRoom}</h3>
                    </div>
                    <div className="ml-auto flex items-center gap-4">
                        <Button variant="ghost" size="icon" className="size-8 rounded-full">
                            <Bell className="size-5 text-foreground/50" />
                        </Button>
                        <Button variant="ghost" size="icon" className="size-8 rounded-full">
                            <Users className="size-5 text-foreground/50" />
                        </Button>
                        <div className="relative">
                            <Search className="size-5 text-foreground/50 absolute left-2 top-1/2 transform -translate-y-1/2" />
                            <Input
                                placeholder="Search"
                                className="pl-9 h-8 bg-zinc-900 border-zinc-700 w-40 focus:w-60 transition-all duration-300"
                            />
                        </div>
                    </div>
                </div>

                <ScrollArea className="flex-1 p-4 overflow-y-auto">
                    <div className="grid gap-y-3 h-full">
                        <div className="flex flex-col items-center justify-center text-center p-8">
                            <div className="size-16 bg-foreground text-background rounded-full flex items-center justify-center mb-4">
                                <Hash className="size-8" />
                            </div>
                            <h2 className="text-2xl font-bold mb-1">Welcome to #{activeRoom}!</h2>
                            <p className="text-foreground/50 max-w-md">
                                This is the start of the #{activeRoom} Room in the {active} server.
                            </p>
                        </div>

                        {messageHistory.map((msg, i) => (
                            <ChatMesaage
                                key={i + msg.message}
                                user={msg.author}
                                message={msg.message}
                                date={new Date(msg.createdAt)}
                            />
                        ))}
                    </div>
                    <ScrollBar />
                </ScrollArea>

                <div className="py-5 px-4">
                    <div className="flex items-center px-3 gap-x-3">
                        <Button variant="ghost" size="icon" className="rounded-full" onClick={sendFormattedMessage}>
                            <Plus />
                        </Button>
                        <Input
                            ref={input}
                            onKeyDown={inputOnKeyDown}
                            placeholder={`Message #${activeRoom}`}
                            className="border-0 bg-transparent focus-visible:ring-0 focus-visible:ring-offset-0"
                        />
                    </div>
                </div>

            </section>

            <section className="w-60 p-3 hidden md:block">
                <h3 className="text-xs font-semibold text-foreground/50 uppercase tracking-wide mb-2 py-2">Online — {onlineUsers.length}</h3>
                <div className="space-y-2">
                    {onlineUsers.map((user, i) => (
                        <div key={user} className="flex items-center gap-2">
                            <div className="relative">
                                <Avatar className="size-8">
                                    <AvatarImage src="" alt={`User ${user}`} />
                                    <AvatarFallback>{user.charAt(0)}</AvatarFallback>
                                </Avatar>
                                <span className="absolute bottom-0 right-0 w-2 h-2 bg-green-500 rounded-full"></span>
                            </div>
                            <div>
                                <p className="text-sm font-medium ">{user}</p>
                                <p className="text-xs text-foreground/50">{i % 2 === 0 ? "Playing a game" : "Online"}</p>
                            </div>
                        </div>
                    ))}
                </div>
            </section>
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

interface Message {
    user: User,
    message: string,
    date: Date
}

function ChatMesaage({
    user,
    message,
    date
}: Message) {
    return (
        <article className="flex gap-3 group hover:bg-foreground/5 p-2 rounded-lg duration-200">
            <Avatar className="size-10 mt-0.5">
                {/* <AvatarImage src="" alt={user} /> */}
                <AvatarFallback>{user.username.charAt(0)}</AvatarFallback>
            </Avatar>
            <div className="flex-1">
                <div className="flex items-baseline gap-x-1">
                    <h4 className="font-medium">{user.username}</h4>
                    <span className="ml-2 text-xs text-foreground/50">
                        {new Intl.DateTimeFormat("en-US", {
                            hour: "numeric",
                            minute: "numeric"
                        }).format(date)}
                    </span>
                    <button>
                        <PencilIcon className=" text-foreground/50 " size={12} />
                    </button>
                    <button>
                        <TrashIcon className=" text-foreground/50 " size={12} />
                    </button>
                </div>
                <p className="text-zinc-300">
                    {message}
                </p>
            </div>
        </article>
    )
}