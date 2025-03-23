import type { User } from "./user";

export interface Message {
    id: string
    authorId: string
    roomId: string
    message: string
    createdAt: number
    updatedAt: number
}

export interface MessageResponse extends Message {
    author: User
}

export interface MessageCreate {
    message: string
}