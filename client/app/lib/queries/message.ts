import { useQuery } from "@tanstack/react-query";
import { queryClient } from "~/components/providers/layout-providers";
import type { Message, MessageResponse } from "~/types/message";
import type { Room } from "~/types/room";
import http from "../http";

const SERVERS_URL = "http://localhost:8080/v1/rooms"

export function useMessagesQuery(roomId: string) {
    const usequery = useQuery<MessageResponse[]>({
        queryKey: ["messages", roomId],
        queryFn: async () => await http(SERVERS_URL + '/' + roomId + '/messages').get(),
    }, queryClient);

    return usequery;
}

export function mutateMessagesCache(roomId: string, message: string) {
    queryClient.setQueryData(["messages", roomId], (old: Room[]) => {
        return [...old, message]
    });
}
