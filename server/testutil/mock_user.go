package testutil

import (
	"testing"

	"github.com/khalidibnwalid/Luma/core"
	"github.com/khalidibnwalid/Luma/models"
	"gorm.io/gorm"
)

func MockUser(t *testing.T, db *gorm.DB) (user *models.User, password string) {
	t.Helper()

	randomeName, _ := core.GenerateRandomString(10)
	user = models.NewUser().
		WithUsername(randomeName).
		WithPassword("123456789").
		WithEmail("" + randomeName + "@example.com")
	user.Create(db)
	password = "123456789"

	t.Cleanup(func() {
		user.Delete(db)
	})
	return user, password
}