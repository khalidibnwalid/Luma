import { useQuery } from "@tanstack/react-query";
import { queryClient } from "~/components/providers/layout-providers";
import type { Room } from "~/types/room";
import http from "../http";

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
