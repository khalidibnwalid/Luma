
export enum AuthError {
    USERNAME_EXISTS = "USERNAME_EXISTS",
    USERNAME_REQUIRED = "USERNAME_REQUIRED",
    USERNAME_INVALID = "USERNAME_INVALID",
    PASSWORD_REQUIRED = "PASSWORD_REQUIRED",
    PASSWORD_INVALID = "PASSWORD_INVALID",
    EMAIL_REQUIRED = "EMAIL_REQUIRED",
    EMAIL_INVALID = "EMAIL_INVALID",
}

export const AuthErrorMessages: Record<AuthError, string> = {
    [AuthError.USERNAME_EXISTS]: "Username already exists",
    [AuthError.USERNAME_REQUIRED]: "Username is required",
    [AuthError.USERNAME_INVALID]: "Username is invalid",
    [AuthError.PASSWORD_REQUIRED]: "Password is required",
    [AuthError.PASSWORD_INVALID]: "Password is invalid",
    [AuthError.EMAIL_REQUIRED]: "Email is required",
    [AuthError.EMAIL_INVALID]: "Email is invalid",
};

export interface ErrorResponse<T> {
    error: T;
    message?: string;
}

