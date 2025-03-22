import type { Room } from "./room"

export interface Server {
    id: string
    name: string
    rooms: Room[]
}