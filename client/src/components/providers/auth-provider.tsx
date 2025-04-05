import http from "@/lib/http";
import { User } from "@/types/user";
import { useQuery } from "@tanstack/react-query";
import { useRouter } from "next/router";
import { createContext, useContext, useEffect } from "react";
import { queryClient } from "./layout-providers";

const SERVER_URL = "http://localhost:8080/v1/users";

interface AuthContext {
    user: User;
}

const UserContext = createContext<AuthContext>({} as AuthContext);

export default function AuthProvider({
    children
}: {
    children: React.ReactNode
}) {
    const router = useRouter()

    const { data, isSuccess, isFetched } = useQuery<User>({
        queryKey: ['user'],
        queryFn: async () => await http(SERVER_URL).get(),

    }, queryClient)

    const isAuthed = !!data?.username && isSuccess

    useEffect(() => {
        if (!isAuthed && isFetched)
            router.push('/login')
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [isAuthed])

    return (
        <UserContext.Provider value={{ user: data } as AuthContext}>
            {children}
        </UserContext.Provider>
    )
};

export const useAuth = () => {
    return useContext(UserContext);
}
