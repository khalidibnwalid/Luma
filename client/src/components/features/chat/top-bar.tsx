import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Room } from "@/types/room";
import { BellIcon, HashIcon, SearchIcon, UsersIcon } from "lucide-react";

export default function ChatTopBar({
    room,
}: {
    room: Room
}) {

    return (
        <div className="h-12 flex items-center px-4">
            <div className="flex items-center">
                <HashIcon className="size-5 text-foreground/50 mr-1" />
                <h3 className="font-bold text-foreground capitalize">{room.name}</h3>
            </div>
            <div className="ml-auto flex items-center gap-4">
                <Button variant="ghost" size="icon" className="size-8 rounded-full">
                    <BellIcon className="size-5 text-foreground/50" />
                </Button>
                <Button variant="ghost" size="icon" className="size-8 rounded-full">
                    <UsersIcon className="size-5 text-foreground/50" />
                </Button>
                <div className="relative">
                    <SearchIcon className="size-5 text-foreground/50 absolute left-2 top-1/2 transform -translate-y-1/2" />
                    <Input
                        placeholder="Search"
                        className="pl-9 h-8 bg-foreground/30 border-foreground/20 w-40 focus:w-60 transition-all duration-300"
                    />
                </div>
            </div>
        </div>
    )
}