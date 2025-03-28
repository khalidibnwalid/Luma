import { ThemeProvider } from 'next-themes'


export default function RootProvider({
    children
}: {
    children: React.ReactNode
}) {
    return (
        <ThemeProvider defaultTheme="dark" attribute="class" storageKey='luma-theme'>
            {children}
        </ThemeProvider>
    )
}