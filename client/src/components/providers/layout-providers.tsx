import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import React from "react";
import AuthProvider from './auth-provider';

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
        // @ts-expect-error - it exists
        window.toggleDevtools = () => setShowDevtools((old) => !old)
    }, [])

    return (
        <AuthProvider>
            <QueryClientProvider client={queryClient}>
                {children}
                <ReactQueryDevtools />
                {showDevtools && (
                    <React.Suspense fallback={null}>
                        <ReactQueryDevtoolsProduction />
                    </React.Suspense>
                )}
            </QueryClientProvider>
        </AuthProvider>
    )
}
