import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import React from "react";
import { ThemeProvider } from "./theme-provider";
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'

export const queryClient = new QueryClient()

const ReactQueryDevtoolsProduction = React.lazy(() =>
    import('@tanstack/react-query-devtools/build/modern/production.js').then(
        (d) => ({
            default: d.ReactQueryDevtools,
        }),
    ),
)

export default function LayoutProviders({ children }: { children: React.ReactNode }) {
    const [showDevtools, setShowDevtools] = React.useState(false)

    React.useEffect(() => {
        // @ts-expect-error
        window.toggleDevtools = () => setShowDevtools((old) => !old)
    }, [])

    return (
        <ThemeProvider defaultTheme="dark">
            <QueryClientProvider client={queryClient}>
                {children}
                <ReactQueryDevtools />
                {showDevtools && (
                    <React.Suspense fallback={null}>
                        <ReactQueryDevtoolsProduction />
                    </React.Suspense>
                )}
            </QueryClientProvider>
        </ThemeProvider>
    )
}
