package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ID        uuid.UUID `gorm:"primarykey;type:uuid;default:gen_random_uuid()" json:"id"`
	AuthorID  uuid.UUID `gorm:"column:author_id;type:uuid;index" json:"authorId"`
	ServerID  uuid.UUID `gorm:"column:server_id;type:uuid;index" json:"serverId"`
	RoomID    uuid.UUID `gorm:"column:room_id;type:uuid;index" json:"roomId"`
	Content   string    `gorm:"column:content" json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	// Author shouldn't be stored in the database, only AuthorID
	Author User `gorm:"-" json:"author"`
}

// TableName specifies the table name for Message
func (Message) TableName() string {
	return "messages"
}

func NewMessage() *Message {
	return &Message{}
}

func (msg *Message) WithContent(content string) *Message {
	msg.Content = content
	return msg
}

func (msg *Message) WithID(id uuid.UUID) *Message {
	msg.ID = id
	return msg
}

func (msg *Message) WithAuthorID(authorID uuid.UUID) *Message {
	msg.AuthorID = authorID
	return msg
}

func (msg *Message) WithServerID(serverID uuid.UUID) *Message {
	msg.ServerID = serverID
	return msg
}

func (msg *Message) WithRoomID(roomID uuid.UUID) *Message {
	msg.RoomID = roomID
	return msg
}

func (msg *Message) Create(db *gorm.DB) error {
	result := db.Create(msg)
	return result.Error
}

func (msg *Message) Update(db *gorm.DB) error {
	result := db.Save(msg)
	return result.Error
}

func (msg *Message) Delete(db *gorm.DB) error {
	result := db.Unscoped().Delete(msg)
	return result.Error
}

func (msg *Message) FindByID(db *gorm.DB, id ...uuid.UUID) error {
	var _id uuid.UUID

	if len(id) > 0 {
		_id = id[0]
	} else {
		_id = msg.ID
	}

	result := db.First(msg, _id)
	return result.Error
}

// GetMessagesByRoom retrieves messages for a specific room with optional pagination
func (msg *Message) GetMessagesByRoom(db *gorm.DB, roomID uuid.UUID, limit int) ([]Message, error) {
	var messages []Message

	if limit <= 0 {
		limit = 50 // Default limit
	}

	result := db.Preload("Author").
		Where("room_id = ?", roomID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages)

	return messages, result.Error
}
