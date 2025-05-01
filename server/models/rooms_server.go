package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomsServer struct {
	gorm.Model
	ID        uuid.UUID `gorm:"primarykey;type:uuid;default:gen_random_uuid()" json:"id"`
	OwnerID   uuid.UUID `gorm:"column:owner_id;type:uuid" json:"ownerId"`
	Name      string    `gorm:"column:name" json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewRoomsServer() *RoomsServer {
	return &RoomsServer{}
}

func (rs *RoomsServer) WithID(id uuid.UUID) *RoomsServer {
	rs.ID = id
	return rs
}

func (rs *RoomsServer) WithOwnerID(ownerID uuid.UUID) *RoomsServer {
	rs.OwnerID = ownerID
	return rs
}

func (rs *RoomsServer) WithName(name string) *RoomsServer {
	rs.Name = name
	return rs
}

func (rs *RoomsServer) Create(db *gorm.DB) error {
	result := db.Create(rs)
	return result.Error
}

func (rs *RoomsServer) FindByID(db *gorm.DB, id ...uuid.UUID) error {
	var _id uuid.UUID

	if len(id) > 0 {
		_id = id[0]
	} else {
		_id = rs.ID
	}

	result := db.First(rs, _id)
	return result.Error
}

func (rs *RoomsServer) Update(db *gorm.DB) error {
	result := db.Save(rs)
	return result.Error
}

func (rs *RoomsServer) Delete(db *gorm.DB) error {
	result := db.Unscoped().Delete(rs)
	return result.Error
}

// GetRooms gets all rooms for a user in this server
func (rs *RoomsServer) GetRooms(db *gorm.DB, userID uuid.UUID) ([]RoomWithStatus, error) {
	var rooms []RoomWithStatus

	err := db.Model(&Room{}).
		Joins("JOIN room_user_status ON rooms.id = room_user_status.room_id").
		Where("room_user_status.server_id = ? AND room_user_status.user_id = ?", rs.ID, userID).
		Select("rooms.*, room_user_status.*").
		Find(&rooms).Error

	return rooms, err
}
