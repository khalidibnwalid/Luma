import { queryClient } from "@/components/providers/layout-providers";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import Spinner from "@/components/ui/spinner";
import http from "@/lib/http";
import { AuthError, AuthErrorMessages } from "@/types/errors";
import { useMutation } from "@tanstack/react-query";
import { Eye, EyeOff } from "lucide-react";
import Head from "next/head";
import Link from "next/link";
import { useRouter } from "next/router";
import { useState } from "react";
import { useForm } from "react-hook-form";

interface Form {
    username: string;
    password: string;
}

const LOGIN_URL = "http://localhost:8080/v1/auth/sessions"

export default function LoginPage() {
    const router = useRouter()

    const { register, handleSubmit } = useForm<Form>()
    const [error, setError] = useState<string | null>(null)
    const [showPassword, setShowPassword] = useState(false);

    const mutation = useMutation({
        mutationKey: ["user"],
        mutationFn: async ({ username, password }: Form) =>
            await http(LOGIN_URL).post({
                username,
                password,
            }),
        onSuccess: () => router.push("/chat"),
        onError: (error: Error) =>
            setError(AuthErrorMessages[error.message as AuthError]),
    }, queryClient)

    async function onSubmit({ username, password }: Form) {
        mutation.mutate({ username: username.trim(), password })
    }

    return (
        <main className="min-h-screen flex items-center justify-center ">
            <Head>
                <title>Luma - Login</title>
            </Head>
            <div className="flex flex-col items-center gap-6 min-w-4xl">
                {/* <h1 className="text-9xlxl font-bold text-center text-primary">
                    Luma
                </h1> */}
                <Card className="mx-auto w-full max-w-md overflow-hidden">
                    <CardHeader className="items-center grid gap-2 text-center">
                        <CardTitle className="text-2xl font-bold">
                            Login to your account
                        </CardTitle>
                        <CardDescription>
                            Welcome back! Please enter your details.
                        </CardDescription>
                    </CardHeader>
                    <CardContent>
                        <form className="grid gap-4" onSubmit={handleSubmit(onSubmit)}>
                            <Spinner />
                            <Button
                                type="button"
                                variant="outline"
                            >
                                <svg
                                    stroke="currentColor"
                                    fill="currentColor"
                                    strokeWidth="0"
                                    version="1.1"
                                    x="0px"
                                    y="0px"
                                    viewBox="0 0 48 48"
                                    enableBackground="new 0 0 48 48"
                                    className="mr-2 size-5"
                                    height="1em"
                                    width="1em"
                                    xmlns="http://www.w3.org/2000/svg"
                                >
                                    <path
                                        fill="#FFC107"
                                        d="M43.611,20.083H42V20H24v8h11.303c-1.649,4.657-6.08,8-11.303,8c-6.627,0-12-5.373-12-12	c0-6.627,5.373-12,12-12c3.059,0,5.842,1.154,7.961,3.039l5.657-5.657C34.046,6.053,29.268,4,24,4C12.955,4,4,12.955,4,24	c0,11.045,8.955,20,20,20c11.045,0,20-8.955,20-20C44,22.659,43.862,21.35,43.611,20.083z"
                                    ></path>
                                    <path
                                        fill="#FF3D00"
                                        d="M6.306,14.691l6.571,4.819C14.655,15.108,18.961,12,24,12c3.059,0,5.842,1.154,7.961,3.039l5.657-5.657	C34.046,6.053,29.268,4,24,4C16.318,4,9.656,8.337,6.306,14.691z"
                                    ></path>
                                    <path
                                        fill="#4CAF50"
                                        d="M24,44c5.166,0,9.86-1.977,13.409-5.192l-6.19-5.238C29.211,35.091,26.715,36,24,36	c-5.202,0-9.619-3.317-11.283-7.946l-6.522,5.025C9.505,39.556,16.227,44,24,44z"
                                    ></path>
                                    <path
                                        fill="#1976D2"
                                        d="M43.611,20.083H42V20H24v8h11.303c-0.792,2.237-2.231,4.166-4.087,5.571	c0.001-0.001,0.002-0.001,0.003-0.002l6.19,5.238C36.971,39.205,44,34,44,24C44,22.659,43.862,21.35,43.611,20.083z"
                                    ></path>
                                </svg>
                                Continue with Google
                                <span className="absolute inset-x-0 -bottom-px h-px bg-gradient-to-r from-transparent via-primary/50 to-transparent"></span>
                            </Button>
                            <Separator orientation="horizontal" />
                            <div className="grid gap-3 py-2">
                                <Label htmlFor="username">Username or Email</Label>
                                <Input
                                    id="username"
                                    type="text"
                                    autoComplete="username"
                                    placeholder="Username or Email"
                                    required
                                    {...register("username", {
                                        required: true,
                                        minLength: {
                                            value: 3,
                                            message: "Must be at least 3 characters",
                                        },
                                    })}
                                />
                                <Label htmlFor="password">Password</Label>

                                <div className="relative">
                                    <Input
                                        id="password"
                                        type={showPassword ? "text" : "password"}
                                        placeholder="••••••••"
                                        className="pr-10"
                                        autoComplete="new-password"
                                        minLength={8}
                                        required
                                        {...register("password", {
                                            required: true,
                                            minLength: {
                                                value: 8,
                                                message: "Password must be at least 8 characters",
                                            },
                                        })}
                                    />
                                    <Button
                                        type="button"
                                        variant="ghost"
                                        size="sm"
                                        className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                                        onClick={() => setShowPassword(!showPassword)}
                                        aria-label={showPassword ? "Hide password" : "Show password"}
                                    >
                                        {showPassword
                                            ? <EyeOff className="size-4 text-muted-foreground" />
                                            : <Eye className="size-4 text-muted-foreground" />
                                        }
                                    </Button>
                                </div>
                            </div>
                            {error && <p className="text-destructive text-sm text-center">{error}</p>}
                            <Button
                                type="submit"
                                className="group w-full"
                            >
                                Continue
                            </Button>

                            <span className="text-sm text-muted-foreground text-center">
                                Don&apos;t have an account? <Link href="/signup" className="text-primary">Signup Here!</Link>
                            </span>
                        </form>
                    </CardContent>
                </Card>
            </div>
        </main>
    );
}
