package testutil

import (
	"context"
	"testing"

	"github.com/khalidibnwalid/Luma/core"
	"github.com/khalidibnwalid/Luma/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func MockMessages(t *testing.T, db *mongo.Database, num uint16, userID string, room ...*models.RoomWithStatus) ([]*models.Message, *models.RoomWithStatus) {
	t.Helper()

	var _room *models.RoomWithStatus
	if len(room) > 0 {
		_room = room[0]
	} else {
		_mockServer, _, user := MockRoomsServer(t, db)
		_room = MockRoom(t, db, user.ID.Hex(), _mockServer)
	}

	mockMessages := make([]*models.Message, num)

	for i := 0; i < int(num); i++ {
		msg, _ := core.GenerateRandomString(20)
		mockMessages[i] = &models.Message{
			RoomID:   _room.ID.Hex(),
			AuthorID: userID,
			Content:  msg,
		}

		mockMessages[i].Create(db, context.Background())
	}

	t.Cleanup(func() {
		for _, msg := range mockMessages {
			msg.Delete(db, context.Background())
		}
	})

	return mockMessages, _room

}
