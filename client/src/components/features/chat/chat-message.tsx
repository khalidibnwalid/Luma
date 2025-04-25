import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { User } from "@/types/user"
import { PencilIcon, TrashIcon } from "lucide-react"

export default function ChatMesaage({
    user,
    message,
    date: unixEpochInSec,
    id,
}: {
    user: User,
    message: string,
    date: number,
    id: string,
}) {
    const date = new Date(unixEpochInSec * 1000)
    return (
        <article id={`messageIdOf-${id}`} className="flex gap-3 group hover:bg-foreground/5 p-2 rounded-lg duration-200">
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
                <p className="text-foreground/85">
                    {message}
                </p>
            </div>
        </article>
    )
}