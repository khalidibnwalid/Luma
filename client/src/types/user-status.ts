
export interface RoomUserStatus {
    id: string
    userId: string
    roomId: string
    lastReadMsgId: string
}

export interface ServerUserStatus {
    id: string
    userId: string
    serverId: string
    nickname?: string
    roles?: string[]
}