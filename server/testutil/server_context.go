package testutil

import (
	"context"
	"testing"

	"github.com/khalidibnwalid/Luma/core"
	"github.com/khalidibnwalid/Luma/handlers"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Connect to the database and return a new ServerContext for testing purposes.
func NewTestingContext(t *testing.T) handlers.ServerContext {
	t.Helper()

	var (
		client *mongo.Client
		err    error
	)

	if client, err = core.CreateMongoClient("mongodb://root:example@localhost:27017/"); err != nil {
		t.Fatalf("MongoDB connection error: %v", err)
	}

	if err = core.PingDB(client, "Testing"); err != nil {
		t.Fatalf("MongoDB ping error: %v", err)
	}

	ctx := &handlers.ServerContext{
		Db:        client.Database("Testing"),
		Client:    client,
		JwtSecret: "SECRET",
	}

	t.Cleanup(func() {
		client.Disconnect(context.Background())
	})

	return *ctx
}
