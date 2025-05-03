package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `gorm:"primarykey;type:uuid;default:gen_random_uuid()" json:"id"`
	AuthorID   uuid.UUID `gorm:"column:author_id;type:uuid;index" json:"-"` // will always be called with author joined
	ServerID   uuid.UUID `gorm:"column:server_id;type:uuid;index" json:"serverId"`
	RoomID     uuid.UUID `gorm:"column:room_id;type:uuid;index" json:"roomId"`
	Content    string    `gorm:"column:content" json:"content"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	// Relationships
	Author User        `gorm:"foreignKey:AuthorID;references:ID;constraint:OnDelete:CASCADE;" json:"author"`
	Room   Room        `gorm:"foreignKey:RoomID;references:ID;constraint:OnDelete:CASCADE;" json:"room"`
	Server RoomsServer `gorm:"foreignKey:ServerID;references:ID;constraint:OnDelete:CASCADE;" json:"server"`
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
