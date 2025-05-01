package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomUserStatus struct {
	gorm.Model
	UserID        uuid.UUID `gorm:"column:user_id;type:uuid;index" json:"userId"`
	ServerID      uuid.UUID `gorm:"column:server_id;type:uuid;index" json:"serverId"`
	RoomID        uuid.UUID `gorm:"column:room_id;type:uuid;index" json:"roomId"`
	LastReadMsgID uuid.UUID `gorm:"column:last_read_msg_id;type:uuid" json:"lastReadMsgId"`
}

func (RoomUserStatus) TableName() string {
	return "room_user_status"
}

type RoomWithStatus struct {
	*Room  `gorm:"embedded"`
	Status *RoomUserStatus `gorm:"embedded" json:"status"`
}

func NewRoomUserStatus() *RoomUserStatus {
	return &RoomUserStatus{}
}

func (r *RoomUserStatus) WithUserID(userID uuid.UUID) *RoomUserStatus {
	r.UserID = userID
	return r
}

func (r *RoomUserStatus) WithServerID(serverID uuid.UUID) *RoomUserStatus {
	r.ServerID = serverID
	return r
}

func (r *RoomUserStatus) WithRoomID(roomID uuid.UUID) *RoomUserStatus {
	r.RoomID = roomID
	return r
}

func (r *RoomUserStatus) WithLastReadMsgID(msgID uuid.UUID) *RoomUserStatus {
	r.LastReadMsgID = msgID
	return r
}

func (r *RoomUserStatus) Create(db *gorm.DB) error {
	result := db.Create(r)
	return result.Error
}

// only updates the LastReadMsgID fields
func (r *RoomUserStatus) Update(db *gorm.DB) error {
	result := db.Model(r).Update("last_read_msg_id", r.LastReadMsgID)
	return result.Error
}

// needs room_id and last_read_msg_id fields to be set in the struct
func (r *RoomUserStatus) UpdateAllUsersStatus(db *gorm.DB, users []uuid.UUID) error {
	result := db.Model(&RoomUserStatus{}).
		Where("user_id IN ? AND room_id = ?", users, r.RoomID).
		Update("last_read_msg_id", r.LastReadMsgID)

	return result.Error
}

func (r *RoomUserStatus) Delete(db *gorm.DB) error {
	result := db.Unscoped().Delete(r)
	return result.Error
}

// Find by user and room ID
func (r *RoomUserStatus) Find(db *gorm.DB) error {
	result := db.Model(&RoomUserStatus{}).
		Where("user_id = ? AND room_id = ?", r.UserID, r.RoomID).
		First(r)
	return result.Error
}
