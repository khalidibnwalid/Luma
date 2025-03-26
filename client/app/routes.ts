import { type RouteConfig, index, layout, route } from "@react-router/dev/routes";

export default [
    index("routes/home.tsx"),
    layout("routes/user/layout.tsx", [
        layout("routes/user/server/layout.tsx", [
            // route("room/@:UserId" ,"routes/room/room.tsx"),
            route("server/:serverId/:roomId?", "routes/user/server/room.tsx"),
        ])
    ])
] satisfies RouteConfig;
