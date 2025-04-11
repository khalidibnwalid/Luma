import { queryClient } from "@/components/providers/layout-providers";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import http from "@/lib/http";
import { mutateServersCache } from "@/lib/queries/rooms-server";
import { RoomsServer } from "@/types/rooms-server";
import { useMutation } from "@tanstack/react-query";
import { ArrowRightIcon, PlusIcon } from "lucide-react";
import { FormEvent, useRef } from "react";

const SERVER_URL = "http://localhost:8080/v1/servers";

type Form = NewServerForm | JoinServerForm;

interface NewServerForm {
    name: string;
}

interface JoinServerForm {
    id: string;
}

export default function AddServerDialog({
    children
}: {
    children: React.ReactNode
}) {
    const nameRef = useRef<HTMLInputElement>(null);
    const idRef = useRef<HTMLInputElement>(null);

    const mutation = useMutation<RoomsServer, Error, Form | string, unknown>({
        mutationFn: async (data) => {
            const url = typeof data === 'string' ? `${SERVER_URL}/${data}` : SERVER_URL;
            return await http(url).post(typeof data === 'string' ? undefined : data);
        },
        onSuccess: (data) => {
            mutateServersCache(data.id).add(data);
            if (nameRef.current) nameRef.current.value = "";
            if (idRef.current) idRef.current.value = "";
        },
        onError: (error) => {
            console.error("Error creating server:", error);
        },
    }, queryClient);

    const onCreateServer = (e?: FormEvent<HTMLFormElement>) => {
        if (e) e.preventDefault();
        const name = nameRef.current?.value;
        if (!name) return;
        mutation.mutate({ name });
    };

    const onJoinServer = (e?: FormEvent<HTMLFormElement>) => {
        if (e) e.preventDefault();
        const id = idRef.current?.value;
        if (!id) return;
        mutation.mutate(id);
    };

    return (
        <Dialog>
            <DialogTrigger asChild>
                {children}
            </DialogTrigger>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle className="font-bold">Add Server</DialogTitle>
                </DialogHeader>

                <form onSubmit={onCreateServer} className="grid gap-y-2">
                    <span>Create</span>
                    <div className="flex items-center gap-2">
                        <Input id="name" ref={nameRef} placeholder="Enter server name" />
                        <Button type="submit" onClick={() => onCreateServer()}><PlusIcon /></Button>
                    </div>
                </form>

                <Separator className="mt-3" />

                <form onSubmit={onJoinServer} className="grid gap-y-2">
                    <span>Join</span>
                    <div className="flex items-center gap-2">
                        <Input id="id" ref={idRef} placeholder="Enter server ID" />
                        <Button type="submit" onClick={() => onJoinServer()}><ArrowRightIcon /></Button>
                    </div>
                </form>
            </DialogContent>
        </Dialog>
    );
}
