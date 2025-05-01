package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ServerUserStatus tracks the status of a user in a server
type ServerUserStatus struct {
	gorm.Model
	UserID   uuid.UUID `gorm:"column:user_id;type:uuid;index" json:"userId"`
	ServerID uuid.UUID `gorm:"column:server_id;type:uuid;index" json:"serverId"`
	Nickname string    `gorm:"column:nickname" json:"nickname"`
	Roles    []string  `gorm:"type:text[];column:roles" json:"roles"`
}

func (ServerUserStatus) TableName() string {
	return "server_user_status"
}

type RoomsServerWithStatus struct {
	*RoomsServer `gorm:"embedded"`
	Status       ServerUserStatus `gorm:"embedded" json:"status"`
}

func NewServerUserStatus() *ServerUserStatus {
	return &ServerUserStatus{}
}

// The ID should be in UUID format
func (s *ServerUserStatus) WithUserID(userID uuid.UUID) *ServerUserStatus {
	s.UserID = userID
	return s
}

// The ID should be in UUID format
func (s *ServerUserStatus) WithServerID(serverID uuid.UUID) *ServerUserStatus {
	s.ServerID = serverID
	return s
}

func (s *ServerUserStatus) WithNickname(nickname string) *ServerUserStatus {
	s.Nickname = nickname
	return s
}

func (s *ServerUserStatus) WithRoles(roles []string) *ServerUserStatus {
	s.Roles = roles
	return s
}

func (s *ServerUserStatus) Create(db *gorm.DB) error {
	result := db.Create(s)
	return result.Error
}

func (s *ServerUserStatus) Update(db *gorm.DB) error {
	result := db.Save(s)
	return result.Error
}

// use user ID and server ID to delete the status
func (s *ServerUserStatus) Delete(db *gorm.DB) error {
	result := db.Unscoped().Where("user_id = ? AND server_id = ?", s.UserID, s.ServerID).Delete(s)
	return result.Error
}

// by user and server ID
func (s *ServerUserStatus) Find(db *gorm.DB) error {
	result := db.Model(&ServerUserStatus{}).
		Where("user_id = ? AND server_id = ?", s.UserID, s.ServerID).
		First(s)
	return result.Error
}

func (s *ServerUserStatus) GetServers(db *gorm.DB) ([]RoomsServerWithStatus, error) {
	var serversWithStatus []RoomsServerWithStatus

	// Join the RoomsServer table with ServerUserStatus
	result := db.Model(&RoomsServer{}).
		Joins("JOIN server_user_status ON rooms_servers.id = server_user_status.server_id").
		Where("server_user_status.user_id = ?", s.UserID).
		Select("rooms_servers.*, server_user_status.*").
		Find(&serversWithStatus)

	return serversWithStatus, result.Error
}
