
export enum RoomType {
    ServerRoom = "server_room",
    ServerVoiceRoom = "server_voice_room",
    UsersGroup = "users_group",
    Direct = "direct"
}

export interface Room {
    id: string
    serverId: string
    name: string
    groupName: string
    type: RoomType
    createdAt: number
    updatedAt: number
}