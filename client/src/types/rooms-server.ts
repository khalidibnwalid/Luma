import { ServerUserStatus } from "./user-status"

export interface RoomsServer {
    id: string
    name: string
    ownerId: string
    createdAt: number
    updatedAt: number
    status: ServerUserStatus
}
