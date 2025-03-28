import { useQuery } from "@tanstack/react-query";
import { queryClient } from "@/components/providers/layout-providers";
import type { RoomsServer } from "@/types/rooms-server";
import http from "../../lib/http";

const SERVERS_URL = "http://localhost:8080/v1/servers"

export function useServersQuery() {
    const usequery = useQuery<RoomsServer[]>({
        queryKey: ["servers"],
        queryFn: async () => await http(SERVERS_URL).get(),
    }, queryClient);

    // populate cache with server data
    function initCache() {
        if (!usequery.data) return;

        usequery.data.forEach((server) => {
            queryClient.setQueryData(["server", server.id], server);
        });
    }

    return { initCache, ...usequery };
}

export function useSingleServerQuery(serverId: string) {
    return useQuery<RoomsServer>({
        queryKey: ["server", serverId],
        queryFn: async () => await http(`${SERVERS_URL}/${serverId}`).get(),
        enabled: !!serverId,
    }, queryClient);
}