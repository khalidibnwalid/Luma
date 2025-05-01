package testutil

import (
	"context"
	"testing"

	"github.com/khalidibnwalid/Luma/core"
	"github.com/khalidibnwalid/Luma/handlers"
	"github.com/khalidibnwalid/Luma/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Connect to the database and return a new ServerContext for testing purposes.
func NewTestingContext(t *testing.T) handlers.ServerContext {
	t.Helper()

	var (
		client *mongo.Client
		err    error
	)

	// TODO delete this
	if client, err = core.CreateMongoClient("mongodb://root:example@localhost:27017/"); err != nil {
		t.Fatalf("MongoDB connection error: %v", err)
	}

	if err = core.PingDB(client, "Testing"); err != nil {
		t.Fatalf("MongoDB ping error: %v", err)
	}

	db, err := core.CreateClient("postgres://admin:123qweasd@localhost:5432/testingLuma?sslmode=disable&TimeZone=UTC")
	if err != nil {
		t.Fatalf("Postgres connection error: %v", err)
	}

	db.Client.AutoMigrate(&models.User{}, &models.RoomsServer{}, &models.ServerUserStatus{}, &models.Room{}, &models.RoomsServerWithStatus{} ,&models.Message{})

	if err = db.Client.Exec("SELECT 1").Error; err != nil {
		t.Fatalf("Postgres ping error: %v", err)
	}

	ctx := &handlers.ServerContext{
		Db:        client.Database("Testing"),
		Client:    client,
		JwtSecret: "SECRET",
		Database:  db,
	}

	t.Cleanup(func() {
		client.Disconnect(context.Background())
	})

	return *ctx
}
