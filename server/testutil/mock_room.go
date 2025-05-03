package testutil

import (
	"testing"

	"github.com/google/uuid"
	"github.com/khalidibnwalid/Luma/core"
	"github.com/khalidibnwalid/Luma/models"
	"gorm.io/gorm"
)

func MockRoom(t *testing.T, db *gorm.DB, userId uuid.UUID, server ...*models.RoomsServer) *models.RoomWithStatus {
	t.Helper()
	var _server *models.RoomsServer
	if len(server) > 0 {
		_server = server[0]
	} else {
		_server, _, _ = MockRoomsServer(t, db)
	}
	roomName, _ := core.GenerateRandomString(10)
	groupname, _ := core.GenerateRandomString(4)

	room := &models.Room{
		ServerID:  _server.ID,
		GroupName: groupname,
		Name:      roomName,
		Type:      "direct",
	}

	room.Create(db)

	status := &models.RoomUserStatus{
		UserID:   userId,
		RoomID:   room.ID,
		ServerID: _server.ID,
	}

	status.Create(db)

	t.Cleanup(func() {
		room.Delete(db)
		status.Delete(db)
	})

	return &models.RoomWithStatus{
		Room:   room,
		Status: status,
	}
}
