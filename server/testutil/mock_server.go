package testutil

import (
	"context"
	"testing"

	"github.com/khalidibnwalid/Luma/core"
	"github.com/khalidibnwalid/Luma/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// if user is not nil, it will be used to create the server
func MockRoomsServer(t *testing.T, db *mongo.Database, user ...models.User) (*models.RoomsServer, *models.ServerUserStatus, *models.User) {
	t.Helper()

	var _user *models.User
	if len(user) > 0 {
		_user = &user[0]
	} else {
		_user, _ = MockUser(t, db)
	}

	server := models.NewRoomsServer().WithOwnerID(_user.ID.Hex())
	server.Name, _ = core.GenerateRandomString(10)
	server.Create(db, context.Background())
	userState := models.NewServerUserStatus().WithUserID(_user.ID.Hex()).WithServerID(server.ID.Hex())

	t.Cleanup(func() {
		server.Delete(db, context.Background())
		userState.Delete(db, context.Background())
	})

	return server, userState, _user
}

func MockRoom(t *testing.T, db *mongo.Database, userId string, server ...*models.RoomsServer) *models.RoomWithStatus {
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
		ServerID:  _server.ID.Hex(),
		GroupName: groupname,
		Name:      roomName,
		Type:      "direct",
	}

	room.Create(db, context.Background())

	status := &models.RoomUserStatus{
		UserID:    userId,
		RoomID:    room.ID.Hex(),
		ServerID:  _server.ID.Hex(),
		IsCleared: true,
	}

	status.Create(db, context.Background())

	t.Cleanup(func() {
		room.Delete(db, context.Background())
		status.Delete(db, context.Background())
	})

	return &models.RoomWithStatus{
		Room:   room,
		Status: status,
	}
}
