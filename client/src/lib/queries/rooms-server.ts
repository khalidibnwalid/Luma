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


export function mutateServersCache(serverId: string) {
    const allServers = ["servers"] as const;
    const key = ["server", serverId]

    function add(server: RoomsServer) {
        queryClient.setQueryData(allServers, (old: RoomsServer[]) => {
            return [...old, server]
        });

        queryClient.setQueryData(key, server);
    }

    function remove(serverId: string) {
        queryClient.setQueryData(allServers, (old: RoomsServer[]) => {
            return old.filter((server) => server.id !== serverId)
        });
        queryClient.setQueryData(key, undefined);
    }

    function update(server: RoomsServer) {
        queryClient.setQueryData(allServers, (old: RoomsServer[]) => {
            return old.map((s) => s.id === server.id ? server : s)
        })
        queryClient.setQueryData(key, server);
    }

    return {
        add,
        remove,
        update,
    }

}