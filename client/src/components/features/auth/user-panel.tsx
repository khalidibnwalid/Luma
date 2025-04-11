import { useAuth } from "@/components/providers/auth-provider"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import { DropdownMenu, DropdownMenuContent, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"
import { HeadphonesIcon, MicIcon, SettingsIcon } from "lucide-react"
import { twJoin } from "tailwind-merge"

export default function UserPanel({
    className
}: {
    className?: string
}) {
    const { user } = useAuth()
    if (!user) return null

    return (
        <DropdownMenu>
            <DropdownMenuTrigger>
                <Avatar className={twJoin(className, " rounded-2xl size-12 cursor-pointer border-2")}>
                    <AvatarImage src="" alt="User" />
                    <AvatarFallback className="rounded-none capitalize font-bold text-lg">
                        {user.username.charAt(0).toUpperCase() || "_"}
                    </AvatarFallback>
                </Avatar>
            </DropdownMenuTrigger>
            <DropdownMenuContent
                className=" w-56 bg-background/60 backdrop-blur-md border-2 rounded-lg shadow-lg"
                side="right"
                align="end"
            >
                <DropdownMenuLabel className="flex items-center gap-2">
                    <Avatar className=" rounded-lg size-8 cursor-pointer border-2">
                        <AvatarImage src="" alt="User" />
                        <AvatarFallback className="rounded-none">{user.username.charAt(0).toUpperCase() || "_"}</AvatarFallback>
                    </Avatar>
                    <span className="text-lg">
                        {user.username}
                    </span>
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                <div className="flex gap-2">
                    <Button variant="ghost" size="icon" className="size-8 rounded-full">
                        <MicIcon className="size-4" />
                    </Button>
                    <Button variant="ghost" size="icon" className="size-8 rounded-full">
                        <HeadphonesIcon className="size-4" />
                    </Button>
                    <Button variant="ghost" size="icon" className="size-8 rounded-full">
                        <SettingsIcon className="size-4" />
                    </Button>
                </div>
            </DropdownMenuContent>
        </DropdownMenu>
    )
}