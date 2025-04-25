import { useQuery } from "@tanstack/react-query";
import { queryClient } from "@/components/providers/layout-providers";
import type { Room } from "@/types/room";
import http from "../../lib/http";

const SERVERS_URL = "http://localhost:8080/v1/servers"

export function useRoomsQuery(serverId?: string) {
    const usequery = useQuery<Room[]>({
        queryKey: ["rooms", serverId],
        queryFn: async () => await http(SERVERS_URL + '/' + serverId + '/rooms').get(),
    }, queryClient);

    // populate cache with rooms data
    function initCache() {
        if (!usequery.data) return;

        usequery.data.forEach((room) => {
            queryClient.setQueryData(["room", room.id], room);
        });
    }

    return { initCache, ...usequery };
}



export function mutateRoomsCache(serverId: string, roomId?: string) {
    const allRooms = ["rooms", serverId]
    const key = ["room", roomId]

    function add(room: Room) {
        queryClient.setQueryData(allRooms, (old: Room[]) => {
            return [...old, room]
        });

        queryClient.setQueryData(key, room);
    }

    function remove(roomID?: string) {
        queryClient.setQueryData(allRooms, (old: Room[]) => {
            return old.filter((room) => room.id !== (roomID || roomId))
        });
        queryClient.setQueryData(key, undefined);
    }

    function update(room: Partial<Room>) {
        queryClient.setQueryData(allRooms, (old: Room[]) => {
            return old.map((r) => r.id === room.id ? room : r)
        })
        queryClient.setQueryData(key, room);
    }

    function partialUpdate(room: Partial<Room>) {
        queryClient.setQueryData(allRooms, (old: Room[]) => {
            return old.map((r) => r.id === roomId ? { ...r, ...room } : r)
        })
        queryClient.setQueryData(key, room);
    }


    return {
        add,
        remove,
        update,
        partialUpdate,
    }
}