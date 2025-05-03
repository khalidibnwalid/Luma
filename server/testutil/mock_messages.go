package testutil

import (
	"testing"

	"github.com/google/uuid"
	"github.com/khalidibnwalid/Luma/core"
	"github.com/khalidibnwalid/Luma/models"
	"gorm.io/gorm"
)

func MockMessages(t *testing.T, db *gorm.DB, num uint16, userID uuid.UUID, room ...*models.RoomWithStatus) ([]*models.Message, *models.RoomWithStatus) {
	t.Helper()

	var _room *models.RoomWithStatus
	if len(room) > 0 {
		_room = room[0]
	} else {
		_mockServer, _, user := MockRoomsServer(t, db)
		_room = MockRoom(t, db, user.ID, _mockServer)
	}

	mockMessages := make([]*models.Message, num)

	for i := 0; i < int(num); i++ {
		msg, _ := core.GenerateRandomString(20)
		mockMessages[i] = &models.Message{
			RoomID:   _room.ID,
			AuthorID: userID,
			Content:  msg,
		}

		mockMessages[i].Create(db)
	}

	t.Cleanup(func() {
		for _, msg := range mockMessages {
			msg.Delete(db)
		}
	})

	return mockMessages, _room

}
