import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";

export default function ChatUsersSidebar({
    // TOODO: replace with actual data from the server
    users,
}: {
    users: string[]
}) {

    return (
        <section className="w-60 p-3 hidden md:block">
            <h3 className="text-xs font-semibold text-foreground/50 uppercase tracking-wide mb-2 py-2">Online — {users.length}</h3>
            <div className="space-y-2">
                {users.map((user, i) => (
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
    )
}