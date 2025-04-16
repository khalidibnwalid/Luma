import { Button } from "@/components/ui/button"
import { DropdownMenu, DropdownMenuContent, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area"
import { Separator } from "@/components/ui/separator"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { AppleIcon, CatIcon, HeartIcon, HistoryIcon, ImageIcon, PlaneIcon, Smile, SmileIcon, Sticker, VolleyballIcon, WatchIcon } from "lucide-react"
import { useState } from "react"

const EMOJIS = [
    {
        id: "recent",
        name: "Recent",
        emojis: ["😀", "😂", "❤️", "👍", "🔥", "✨", "🎉"],
        icon: HistoryIcon,
    }, // TODO add real recent
    {
        id: "smileys",
        name: "Smileys",
        emojis: ["😀", "😃", "😄", "😁", "😆", "😅", "😂", "🤣", "😊", "😇", "🙂", "🙃", "😉", "😌", "😍", "🥰", "😘", "👍", "👎", "👊", "✊", "🤛", "🤜", "🤞", "✌️", "🤟", "🤘", "👌", "👈", "👉", "👆", "👇", "☝️", "✋", "🤚", "🖐️", "🖖", "👋", "🤙"],
        icon: SmileIcon,
    },
    {
        id: "nature",
        name: "Nature",
        emojis: ["🐶", "🐱", "🐭", "🐹", "🐰", "🦊", "🐻", "🐼", "🐨", "🐯", "🦁", "🐮", "🐷", "🐸", "🐵", "🐔", "🐧"],
        icon: CatIcon,
    },
    {
        id: "food",
        name: "Food",
        emojis: ["🍏", "🍎", "🍐", "🍊", "🍋", "🍌", "🍉", "🍇", "🍓", "🍈", "🍒", "🍑", "🥭", "🍍", "🥥", "🥝", "🍅"],
        icon: AppleIcon,
    },
    {
        id: "activities",
        name: "Activities",
        emojis: ["⚽", "🏀", "🏈", "⚾", "🥎", "🎾", "🏐", "🏉", "🥏", "🎱", "🪀", "🏓", "🏸", "🏒", "🏑", "🥍", "🏏"],
        icon: VolleyballIcon,
    },
    {
        id: "travel",
        name: "Travel",
        emojis: ["🚗", "🚕", "🚙", "🚌", "🚎", "🏎️", "🚓", "🚑", "🚒", "🚐", "🚚", "🚛", "🚜", "🛴", "🚲", "🛵", "🏍️"],
        icon: PlaneIcon,
    },
    {
        id: "objects",
        name: "Objects",
        emojis: ["⌚", "📱", "💻", "⌨️", "🖥️", "🖨️", "🖱️", "🖲️", "🕹️", "🗜️", "💽", "💾", "💿", "📀", "📼", "📷", "📸"],
        icon: WatchIcon,
    },
    {
        id: "symbols",
        name: "Symbols",
        emojis: ["❤️", "🧡", "💛", "💚", "💙", "💜", "🖤", "🤍", "🤎", "💔", "❣️", "💕", "💞", "💓", "💗", "💖", "💘"],
        icon: HeartIcon,
    },
]

export default function EmojiSelector({
    onEmojiSelect,
    children,
}: {
    onEmojiSelect: (emoji: string) => void
    children: React.ReactNode
}) {
    const [activeCategory, setActiveCategory] = useState("recent")

    return (
        <DropdownMenu>
            <DropdownMenuTrigger>
                {children}
            </DropdownMenuTrigger>
            <DropdownMenuContent side="top" className="border rounded-lg shadow-lg w-md bg-background/50 backdrop-blur-lg">
                <Tabs defaultValue="emojis" className="w-full">
                    <TabsList className="w-full grid grid-cols-3 sticky top-0 z-10 bg-background/50 backdrop-blur-lg">
                        <TabsTrigger value="emojis" className="flex items-center gap-2">
                            <Smile className="h-4 w-4" />
                            <span>Emojis</span>
                        </TabsTrigger>
                        <TabsTrigger value="stickers" className="flex items-center gap-2">
                            <Sticker className="h-4 w-4" />
                            <span>Stickers</span>
                        </TabsTrigger>
                        <TabsTrigger value="gifs" className="flex items-center gap-2">
                            <ImageIcon className="h-4 w-4" />
                            <span>GIFs</span>
                        </TabsTrigger>
                    </TabsList>

                    <ScrollArea className="h-56">
                        <div className="max-w-md h-full">
                            <TabsContent value="emojis" className="px-2">
                                <div className="grid gap-1 h-full">
                                    <ScrollArea >
                                        <div className="px-2 flex gap-x-1 overflow-x-auto">
                                            {EMOJIS.map((category) => (
                                                <Button
                                                    key={category.id}
                                                    variant={activeCategory === category.id ? "contrast" : "ghost"}
                                                    className="text-lg px-2 hover:bg-foreground/10 hover:text-foreground"
                                                    onClick={() => setActiveCategory(category.id)}
                                                >

                                                    <category.icon className="size-4" />
                                                </Button>
                                            ))}
                                        </div>
                                        <ScrollBar orientation="horizontal" />
                                    </ScrollArea>
                                    <Separator className="my-2" />
                                    <div className="grid grid-cols-8 gap-1 w-full h-full">
                                        {EMOJIS
                                            .find((category) => category.id === activeCategory)
                                            ?.emojis.map((emoji, index) => (
                                                <Button
                                                    key={index}
                                                    variant="ghost"
                                                    size="icon"
                                                    className="text-3xl p-0"
                                                    onClick={() => onEmojiSelect(emoji)}
                                                >
                                                    {emoji}
                                                </Button>
                                            ))}
                                    </div>
                                </div>
                            </TabsContent>

                            <TabsContent value="stickers" className="p-4 h-48 flex items-center justify-center text-muted-foreground">
                                <p>Stickers not implemented yet</p>
                            </TabsContent>

                            <TabsContent value="gifs" className="p-4 h-48 flex items-center justify-center text-muted-foreground">
                                <p>GIFs not implemented yet</p>
                            </TabsContent>
                            <ScrollBar />
                        </div>
                    </ScrollArea>


                </Tabs>

            </DropdownMenuContent>
        </DropdownMenu >
    )
}
