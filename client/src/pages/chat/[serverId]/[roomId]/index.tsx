import ChatMesaage from "@/components/features/chat/chat-message"
import EmojiSelector from "@/components/features/chat/emoji-selector"
import ChatTopBar from "@/components/features/chat/top-bar"
import ChatUsersSidebar from "@/components/features/chat/users-bar"
import AppLayout from "@/components/layouts/app-layout"
import ServerLayout, { useServerContext } from "@/components/layouts/server-layout"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area"
import { useIntersectionObserver } from "@/lib/hooks/useIntersectionObserver"
import http from "@/lib/http"
import { useMessagesQuery } from "@/lib/queries/message"
import { mutateRoomsCache } from "@/lib/queries/rooms"
import type { MessageResponse } from "@/types/message"
import { RoomUserStatus } from "@/types/user-status"
import { useMutation } from "@tanstack/react-query"
import { Hash, Plus, SmilePlusIcon } from "lucide-react"
import { useRouter } from "next/router"
import { ReactElement, useEffect, useRef, useState } from "react"
import useWebSocket from 'react-use-websocket'

const SERVERS_URL = "http://localhost:8080/v1/rooms"


//TODO scroll to the lastReadMsgId on load
export default function Page() {
    const router = useRouter()
    const roomId = router.query.roomId as string
    const { activeServer, rooms } = useServerContext()
    const activeRoom = rooms?.find(room => room.id === roomId)

    const { data: messages, isSuccess } = useMessagesQuery(roomId)
    const socketUrl = 'ws://localhost:8080/v1/rooms/' + roomId
    const { sendMessage, lastMessage } = useWebSocket(socketUrl);

    const input = useRef<HTMLInputElement>(null);
    const [newMessageHistory, setMessageHistory] = useState<MessageResponse[]>([]);
    const allMessages = [...(messages ?? []), ...newMessageHistory.filter(msg => msg.roomId === roomId)]

    const [lastVisibleMessageId, setLastVisibleMessageId] = useState<string>(activeRoom?.status?.lastReadMsgId || "")

    const mutateRoom = useMutation({
        mutationFn: () =>
            http(SERVERS_URL + '/' + roomId + '/status').patch({ lastReadMsgId: lastVisibleMessageId }),
        onSuccess: () => mutateRoomsCache(activeServer.id, activeRoom?.id).partialUpdate({
            status: {
                lastReadMsgId: lastVisibleMessageId || "",
            } as RoomUserStatus // :DDDDD
        }),
    })

    const {
        isIntersecting: isBottomInView,
        ref: endOfScroll,
        setRef: bindEndOfScroll
    } = useIntersectionObserver({ threshold: 1 })

    const isSnappedToBottom = endOfScroll && isBottomInView

    useEffect(() => {
        if (!isSnappedToBottom) return
        endOfScroll.scrollIntoView({ behavior: "smooth" })
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [newMessageHistory, messages])

    useEffect(() => {
        if (lastMessage === null) return

        const data = JSON.parse(lastMessage.data) as MessageResponse
        setMessageHistory((prev) => prev.concat(data))
        markAsUnread(data.id)
    }, [lastMessage]);

    if (!activeRoom) return <div>Room not found</div>
    if (!isSuccess) return <div>Server not found</div>

    function markAsUnread(msgId: string) {
        setLastVisibleMessageId(msgId)
        mutateRoom.mutate()
    }

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
    const mockUsers = ["User 1", "User 2", "User 3", "User 4", "User 5"]

    return (
        <>
            <section className="flex-1 flex flex-col">
                <ChatTopBar room={activeRoom} />
                {/* {(lastVisibleMessageId === activeRoom?.status?.lastReadMsgId && !isSnappedToBottom) &&
                    <button
                        className="py-1 w-full text-center rounded-b-2xl bg-foreground/10 hover:bg-foreground/15 duration-150"
                        onClick={() => endOfScroll?.scrollIntoView({ behavior: "smooth" })}
                    >
                        Mark As Read
                    </button>
                } */}
                <ScrollArea className="flex-1 p-4 overflow-y-auto">
                    <div className="grid gap-y-3">
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
                            <>
                                {(msg.id === lastVisibleMessageId && i <= (allMessages.length - 5))
                                    && <span className="w-full text-center rounded-2xl border-foreground/10 border-t-4"></span>}
                                <ChatMesaage
                                    id={msg.id}
                                    key={i + msg.content}
                                    user={msg.author}
                                    message={msg.content}
                                    date={msg.createdAt}
                                    markAsUnreadFn={markAsUnread}
                                />
                            </>
                        ))}

                        <div ref={bindEndOfScroll}></div>
                    </div>
                    <ScrollBar />
                </ScrollArea>

                <div className="py-5 px-4">
                    <div className="flex items-center px-3 gap-x-2">
                        <EmojiSelector onEmojiSelect={e => input.current!.value += e}>
                            <Button variant="ghost" size="icon" className="rounded-full text">
                                <SmilePlusIcon />
                            </Button>
                        </EmojiSelector>
                        <Button variant="ghost" size="icon" className="rounded-full" onClick={sendFormattedMessage}>
                            <Plus />
                        </Button>
                        <Input
                            ref={input}
                            onKeyDown={inputOnKeyDown}
                            placeholder={`Message #${activeRoom.name}`}
                            className="border-0 md:text-md py-5 bg-transparent focus-visible:ring-0 focus-visible:ring-offset-0"
                        />
                    </div>
                </div>

            </section>

            <ChatUsersSidebar users={mockUsers} />
        </>
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