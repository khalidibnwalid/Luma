package testutil

import (
	"context"
	"testing"

	"github.com/khalidibnwalid/Luma/core"
	"github.com/khalidibnwalid/Luma/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func MockUser(t *testing.T, db *mongo.Database) (user *models.User, password string) {
	t.Helper()

	randomeName, _ := core.GenerateRandomString(10)
	user = models.NewUser().
		WithUsername(randomeName).
		WithPassword("123456789").
		WithEmail("" + randomeName + "@example.com")
	user.Create(db, context.Background())
	password = "123456789"

	t.Cleanup(func() {
		user.Delete(db, context.Background())
	})
	return user, password
}
