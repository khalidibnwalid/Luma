import { useServerContext } from "@/components/layouts/server-layout";
import { Button } from "@/components/ui/button";
import { ContextMenuContent, ContextMenuItem } from "@/components/ui/context-menu";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import http from "@/lib/http";
import { mutateRoomsCache } from "@/lib/queries/rooms";
import { Room, RoomType } from "@/types/room";
import { ContextMenu, ContextMenuTrigger } from "@radix-ui/react-context-menu";
import { useMutation } from "@tanstack/react-query";
import { AudioLinesIcon, MessageCircleIcon } from "lucide-react";
import { useRef, useState } from "react";

const SERVER_URL = "http://localhost:8080/v1/servers";

export default function ServerSidebarContextMenu({
    className,
    children
}: {
    className?: string;
    children?: React.ReactNode;
}) {
    const [openDialog, setOpenDialog] = useState<null | 'chat' | 'voice'>(null);
    const closeDialog = () => setOpenDialog(null);

    return (
        <>
            <ContextMenu>
                <ContextMenuTrigger className={className}>
                    {children}
                </ContextMenuTrigger>

                <ContextMenuContent className="w-56 bg-background/50 backdrop-blur-lg">
                    <ContextMenuItem onSelect={() => setOpenDialog('chat')}>
                        <MessageCircleIcon className="size-4" />
                        Create Chat Room
                    </ContextMenuItem>

                    <ContextMenuItem onSelect={() => setOpenDialog('voice')}>
                        <AudioLinesIcon className="size-4" />
                        Create Voice Room
                    </ContextMenuItem>
                </ContextMenuContent>
            </ContextMenu>

            <CreateChatRoomDialog open={openDialog === 'chat'} close={closeDialog} />
            <CreateVoiceRoomDialog open={openDialog === 'voice'} close={closeDialog} />
        </>
    );
}

type DialogProps = {
    open: boolean
    close: () => void
}

function CreateChatRoomDialog({ open, close }: DialogProps) {
    const nameRef = useRef<HTMLInputElement>(null);
    const { activeServer } = useServerContext();

    const mutation = useMutation<Room, Error, string, unknown>({
        mutationFn: async (name) =>
            await http(`${SERVER_URL}/${activeServer.id}/rooms`).post({
                name,
                type: RoomType.ServerRoom
            }),
        onSuccess: (room) => {
            mutateRoomsCache(activeServer.id).add(room)
            close();
        }
    });

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (nameRef.current) mutation.mutate(nameRef.current.value.trim());
    };

    return (
        <Dialog open={open} onOpenChange={close} modal={false}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Create Chat Room</DialogTitle>
                </DialogHeader>

                <form onSubmit={handleSubmit}>
                    <div className="grid gap-4 py-4">
                        <div className="grid gap-2">
                            <Label htmlFor="chat-name">Name</Label>
                            <Input
                                id="chat-name"
                                ref={nameRef}
                                placeholder="Enter room name"
                                disabled={mutation.isPending}
                                required
                            />
                        </div>
                    </div>
                    <div className="flex justify-end gap-2">
                        <Button type="button" variant="outline" onClick={close}>Cancel</Button>
                        <Button
                            type="submit"
                            disabled={mutation.isPending}
                        >
                            {mutation.isPending ? 'Creating...' : 'Create'}
                        </Button>
                    </div>
                </form>
            </DialogContent>
        </Dialog>
    );
}

function CreateVoiceRoomDialog({ open, close }: DialogProps) {
    const nameRef = useRef<HTMLInputElement>(null);
    const { activeServer } = useServerContext();

    const mutation = useMutation<Room, Error, string, unknown>({
        mutationFn: async (name) =>
            await http(`${SERVER_URL}/${activeServer.id}/rooms`).post({
                name,
                type: RoomType.ServerRoom
            }),
        onSuccess: (room) => {
            mutateRoomsCache(activeServer.id).add(room)
            close();
        }
    });

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (nameRef.current) mutation.mutate(nameRef.current.value.trim());
    };

    return (
        <Dialog open={open} onOpenChange={close} modal={false}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Create Voice Room</DialogTitle>
                </DialogHeader>

                <form onSubmit={handleSubmit}>
                    <div className="grid gap-4 py-4">
                        <div className="grid gap-2">
                            <Label htmlFor="voice-name">Name</Label>
                            <Input
                                id="voice-name"
                                ref={nameRef}
                                placeholder="Enter room name"
                                disabled={mutation.isPending}
                                required
                            />
                        </div>
                    </div>
                    <div className="flex justify-end gap-2">
                        <Button type="button" variant="outline" onClick={close}>Cancel</Button>
                        <Button
                            type="submit"
                            disabled={mutation.isPending}
                        >
                            {mutation.isPending ? 'Creating...' : 'Create'}
                        </Button>
                    </div>
                </form>
            </DialogContent>
        </Dialog>
    );
}