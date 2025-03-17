import { ThemeProvider } from "~/components/providers/theme-provider"
import { Outlet } from "react-router"

export default function ChatLayout() {
    return (
        <div className="w-full h-screen">
            <ThemeProvider defaultTheme="dark">
                <Outlet />
            </ThemeProvider>
        </div>
    )
}