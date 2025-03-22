import { type RouteConfig, index, layout, route } from "@react-router/dev/routes";

export default [
    index("routes/home.tsx"),
    layout("routes/room/layout.tsx", [
        // route("room/@:UserId" ,"routes/room/room.tsx"),
        route("room/:serverId/:roomId" ,"routes/room/room.tsx"),
    ])
] satisfies RouteConfig;
