import { type RouteConfig, index, layout, route } from "@react-router/dev/routes";

export default [
    index("routes/home.tsx"),
    layout("routes/chat/layout.tsx", [
        route("chat" ,"routes/chat/chat.tsx")
    ])
] satisfies RouteConfig;
