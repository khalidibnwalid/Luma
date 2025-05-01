package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const messagesLimit = 50
const RoomsCollection = "rooms"

type Room struct {
	gorm.Model
	ID        uuid.UUID `gorm:"primarykey;type:uuid;default:gen_random_uuid()" json:"id"`
	ServerID  uuid.UUID `gorm:"column:server_id;type:uuid;index" json:"serverId"`
	Name      string    `gorm:"column:name" json:"name"`
	GroupName string    `gorm:"column:group_name" json:"groupName"`
	Type      string    `gorm:"column:type" json:"type"` // direct, server room, server voice room, or users group
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Room) TableName() string {
	return "rooms"
}

func NewRoom() *Room {
	return &Room{}
}

func (r *Room) WithID(id uuid.UUID) *Room {
	r.ID = id
	return r
}

func (r *Room) WithServerID(serverID uuid.UUID) *Room {
	r.ServerID = serverID
	return r
}

func (r *Room) WithName(name string) *Room {
	r.Name = name
	return r
}

func (r *Room) WithGroupName(groupName string) *Room {
	r.GroupName = groupName
	return r
}

func (r *Room) WithType(roomType string) *Room {
	r.Type = roomType
	return r
}

// FindByID finds a room by its ID
func (r *Room) FindByID(db *gorm.DB, id ...uuid.UUID) error {
	var _id uuid.UUID

	if len(id) > 0 {
		_id = id[0]
	} else {
		_id = r.ID
	}

	result := db.First(r, _id)
	return result.Error
}

// GetMessages retrieves messages for a room
func (r *Room) GetMessages(db *gorm.DB, limit ...int) ([]Message, error) {
	var messages []Message
	var _limit int
	if len(limit) > 0 {
		_limit = limit[0]
	} else {
		_limit = messagesLimit
	}

	// Use joins to fetch messages with author information
	result := db.Model(&Message{}).
		Joins("JOIN users ON messages.author_id = users.id").
		Where("messages.room_id = ?", r.ID).
		Order("messages.created_at DESC").
		Limit(_limit).
		Find(&messages)

	return messages, result.Error
}

func (r *Room) Create(db *gorm.DB) error {
	result := db.Create(r)
	return result.Error
}

func (r *Room) Delete(db *gorm.DB) error {
	result := db.Unscoped().Delete(r)
	return result.Error
}

func (r *Room) Update(db *gorm.DB) error {
	result := db.Save(r)
	return result.Error
}
