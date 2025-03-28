import AppLayout from "@/components/layouts/app-layout"
import ServerLayout, { useServerContext } from "@/components/layouts/server-layout"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area"
import { useMessagesQuery } from "@/lib/queries/message"
import type { MessageResponse } from "@/types/message"
import type { User } from "@/types/user"
import { Bell, Hash, PencilIcon, Plus, Search, TrashIcon, Users } from "lucide-react"
import { useRouter } from "next/router"
import { ReactElement, useEffect, useRef, useState } from "react"
import useWebSocket from 'react-use-websocket'

export default function Page() {
    const router = useRouter()
    const roomId = router.query.roomId as string
    const { activeServer, rooms } = useServerContext()
    const activeRoom = rooms?.find(room => room.id === roomId)

    const { data: messages, isSuccess } = useMessagesQuery(roomId)

    const socketUrl = 'ws://localhost:8080/v1/rooms/' + roomId + "?jwt=" + localStorage.getItem("token")
    const { sendMessage, lastMessage } = useWebSocket(socketUrl);

    const [newMessageHistory, setMessageHistory] = useState<MessageResponse[]>([]);
    const allMessages = [...(messages ?? []), ...newMessageHistory.filter(msg => msg.roomId === roomId)]

    const input = useRef<HTMLInputElement>(null);

    useEffect(() => {
        if (lastMessage !== null) setMessageHistory((prev) => prev.concat(JSON.parse(lastMessage.data)));
    }, [lastMessage]);

    if (!activeRoom) return <div>Room not found</div>
    if (!isSuccess) return <div>Server not found</div>

    function sendFormattedMessage() {
        const value = input.current?.value
        if (!value) return
        sendMessage(JSON.stringify({ content: value } as MessageResponse))
        input.current!.value = ""
    }

    function inputOnKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
        if (e.key === "Enter") sendFormattedMessage()
    }

    /// mocks
    const onlineUsers = ["User 1", "User 2", "User 3", "User 4", "User 5"]

    return (
        <>

            <section className="flex-1 flex flex-col">
                <div className="h-12 flex items-center px-4">
                    <div className="flex items-center">
                        <Hash className="size-5 text-foreground/50 mr-1" />
                        <h3 className="font-bold text-white capitalize">{activeRoom?.name}</h3>
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
                    <div  className="grid gap-y-3">
                        <div className="flex flex-col items-center justify-center text-center p-8">
                            <div className="size-16 bg-foreground text-background rounded-full flex items-center justify-center mb-4">
                                <Hash className="size-8" />
                            </div>
                            <h2 className="text-2xl font-bold mb-1">Welcome to #{activeRoom?.name}!</h2>
                            <p className="text-foreground/50 max-w-md">
                                This is the start of the #{activeRoom?.name} Room in the {activeServer?.name} server.
                            </p>
                        </div>

                        {allMessages.map((msg, i) => (
                            <ChatMesaage
                                key={i + msg.content}
                                user={msg.author}
                                message={msg.content}
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
                            placeholder={`Message #${activeRoom.name}`}
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
        </>
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


Page.getLayout = function getLayout(page: ReactElement) {
    return (
        <AppLayout>
            <ServerLayout>
                {page}
            </ServerLayout>
        </AppLayout>
    )
}