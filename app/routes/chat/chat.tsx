import { Bell, Hash, Headphones, Mic, PencilIcon, Plus, PlusIcon, Search, Settings, TrashIcon, Users } from "lucide-react"
import { twJoin } from "tailwind-merge"
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar"
import { Button } from "~/components/ui/button"
import { Input } from "~/components/ui/input"
import { ScrollArea } from "~/components/ui/scroll-area"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "~/components/ui/tooltip"

export default function DiscordClone() {
    const active = "general"
    const activeChannel = "welcome"

    const onlineUsers = ["User 1", "User 2", "User 3", "User 4", "User 5"]
    const msgs: Message[] = [
        { user: "User 1", message: "Hello everyone! How are you all doing today? I hope everyone is having a great day!", date: new Date() },
        { user: "User 2", message: "Hi User 1! I'm doing well, thank you. How about you?", date: new Date() },
        { user: "User 3", message: "How's it going? I just finished working on a new project and I'm really excited about it!", date: new Date() },
        { user: "User 4", message: "Good morning! I just had a great breakfast and I'm ready to start the day!", date: new Date() },
        { user: "User 5", message: "What's up? I'm just chilling and enjoying some free time.", date: new Date() }
    ]

    return (
        <main className="flex h-screen">
            <section className="w-16 flex flex-col items-center py-3 gap-2 overflow-y-auto">

                {["general", "gaming", "dev", "dev", "art"].map((server) => (
                    <TooltipProvider key={server}>
                        <Tooltip>
                            <TooltipTrigger asChild>
                                <button
                                    className={twJoin(`w-12 h-12 flex items-center justify-center rounded-2xl border-2 duration-200 aspect-square`,
                                        active === server
                                            ? "bg-primary text-background"
                                            : "hover:bg-primary hover:text-background"
                                    )}
                                >
                                    <span className="capitalize font-bold text-lg">{server.charAt(0).toUpperCase()}</span>
                                </button>
                            </TooltipTrigger>
                            <TooltipContent side="right">
                                <p className="capitalize">{server}</p>
                            </TooltipContent>
                        </Tooltip>
                    </TooltipProvider>
                ))}

                <TooltipProvider>
                    <Tooltip>
                        <TooltipTrigger asChild>
                            <Button
                                variant="ghost"
                                className="rounded-2xl border-2  duration-200 p-0 w-12 h-12 flex items-center justify-center"
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

            <section className="w-60 bg-foreground/5 flex flex-col">
                <div className="p-4 shadow-md">
                    <h2 className="font-bold text-lg capitalize">{active}</h2>
                </div>

                <ScrollArea className="flex-1">
                    <div className="p-2">
                        <ChatChannel activeChannel="welcome" channelName="Text Channels" channels={["welcome", "general", "rules", "announcements"]} />
                        <ChatChannel channelName="Voice Channels" channels={["General", "Music", "Gaming"]} />
                    </div>
                </ScrollArea>

                <UserChannel user="User" subText="#000" />
            </section>

            <section className="flex-1 flex flex-col">
                <div className="h-12 flex items-center px-4">
                    <div className="flex items-center">
                        <Hash className="h-5 w-5 text-zinc-400 mr-1" />
                        <h3 className="font-bold text-white capitalize">{activeChannel}</h3>
                    </div>
                    <div className="ml-auto flex items-center gap-4">
                        <Button variant="ghost" size="icon" className="h-8 w-8 rounded-full">
                            <Bell className="h-5 w-5 text-zinc-400" />
                        </Button>
                        <Button variant="ghost" size="icon" className="h-8 w-8 rounded-full">
                            <Users className="h-5 w-5 text-zinc-400" />
                        </Button>
                        <div className="relative">
                            <Search className="h-5 w-5 text-zinc-400 absolute left-2 top-1/2 transform -translate-y-1/2" />
                            <Input
                                placeholder="Search"
                                className="pl-9 h-8 bg-zinc-900 border-zinc-700  w-40 focus:w-60 transition-all duration-300"
                            />
                        </div>
                    </div>
                </div>

                <ScrollArea className="flex-1 p-4">
                    <div className="grid gap-y-3">
                        <div className="flex flex-col items-center justify-center text-center p-8">
                            <div className="w-16 h-16 bg-foreground text-background rounded-full flex items-center justify-center mb-4">
                                <Hash className="h-8 w-8" />
                            </div>
                            <h2 className="text-2xl font-bold mb-1">Welcome to #{activeChannel}!</h2>
                            <p className="text-zinc-400 max-w-md">
                                This is the start of the #{activeChannel} channel in the {active} server.
                            </p>
                        </div>

                        {msgs.map((msg) => (
                            <ChatMesaage key={msg.user} {...msg} />
                        ))}
                    </div>
                </ScrollArea>

                <div className="py-5 px-4">
                    <div className="flex items-center px-3 gap-x-3">
                        <Button variant="ghost" size="icon" className="rounded-full">
                            <Plus />
                        </Button>
                        <Input
                            placeholder={`Message #${activeChannel}`}
                            className="border-0 bg-transparent focus-visible:ring-0 focus-visible:ring-offset-0"
                        />
                    </div>
                </div>

            </section>

            <section className="w-60 p-3 hidden md:block">
                <h3 className="text-xs font-semibold text-zinc-400 uppercase tracking-wide mb-2 py-2">Online — {onlineUsers.length}</h3>
                <div className="space-y-2">
                    {onlineUsers.map((user, i) => (
                        <div key={user} className="flex items-center gap-2">
                            <div className="relative">
                                <Avatar className="h-8 w-8">
                                    <AvatarImage src="" alt={`User ${user}`} />
                                    <AvatarFallback>{user.charAt(0)}</AvatarFallback>
                                </Avatar>
                                <span className="absolute bottom-0 right-0 w-2 h-2 bg-green-500 rounded-full"></span>
                            </div>
                            <div>
                                <p className="text-sm font-medium ">{user}</p>
                                <p className="text-xs text-zinc-400">{i % 2 === 0 ? "Playing a game" : "Online"}</p>
                            </div>
                        </div>
                    ))}
                </div>
            </section>
        </main >
    )
}

function ChatChannel({
    activeChannel,
    channelName,
    channels
}: {
    activeChannel?: string,
    channelName: string
    channels: string[]
}) {

    return (
        <div className="my-2">
            <h3 className="p-2 text-xs font-semibold text-foreground uppercase tracking-wide">{channelName}</h3>
            <div className="grid gap-1">
                {channels.map((channel) => (
                    <button
                        key={channel}
                        className={twJoin(' w-full px-2 py-1 text-sm flex items-center gap-1.5 rounded duration-200 text-foreground/70',
                            activeChannel === channel ? "bg-accent text-foreground/100" : "hover:bg-accent"
                        )}
                    >
                        <Hash className="h-4 w-4" />
                        <span className="capitalize">{channel}</span>
                    </button>
                ))}
            </div>
        </div>
    )
}

function UserChannel({
    user,
    subText
}: {
    user: string
    subText: string
}) {
    return (
        <div className="p-2 flex items-center gap-2">
            <Avatar className="h-8 w-8">
                <AvatarImage src="" alt="User" />
                <AvatarFallback>U</AvatarFallback>
            </Avatar>

            <div className="flex-1 min-w-0">
                <p className="text-sm font-medium">{user}</p>
                <p className="text-xs text-zinc-400 truncate">{subText}</p>
            </div>

            <div className="flex gap-1">
                <Button variant="ghost" size="icon" className="h-8 w-8 rounded-full">
                    <Mic className="h-4 w-4" />
                </Button>
                <Button variant="ghost" size="icon" className="h-8 w-8 rounded-full">
                    <Headphones className="h-4 w-4" />
                </Button>
                <Button variant="ghost" size="icon" className="h-8 w-8 rounded-full">
                    <Settings className="h-4 w-4" />
                </Button>
            </div>
        </div>
    )
}

interface Message {
    user: string,
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
            <Avatar className="h-10 w-10 mt-0.5">
                <AvatarImage src="" alt={user} />
                <AvatarFallback>{user.charAt(0)}</AvatarFallback>
            </Avatar>
            <div className="flex-1">
                <div className="flex items-baseline gap-x-1">
                    <h4 className="font-medium">{user}</h4>
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