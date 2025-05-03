package testutil

import (
	"testing"

	"github.com/khalidibnwalid/Luma/core"
	"github.com/khalidibnwalid/Luma/models"
	"gorm.io/gorm"
)

// if user is not nil, it will be used to create the server
func MockRoomsServer(t *testing.T, db *gorm.DB, user ...*models.User) (*models.RoomsServer, *models.ServerUserStatus, *models.User) {
	t.Helper()

	var _user *models.User
	if len(user) > 0 {
		_user = user[0]
	} else {
		_user, _ = MockUser(t, db)
	}

	server := models.NewRoomsServer().WithOwnerID(_user.ID)
	server.Name, _ = core.GenerateRandomString(10)
	server.Create(db)
	status := models.NewServerUserStatus().WithUserID(_user.ID).WithServerID(server.ID)

	t.Cleanup(func() {
		server.Delete(db)
		status.Delete(db)
	})

	return server, status, _user
}
