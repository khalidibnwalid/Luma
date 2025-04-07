package handlers

import (
	"log"

	"github.com/khalidibnwalid/Luma/core"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ServerContext struct {
	Db        *mongo.Database
	Client    *mongo.Client
	JwtSecret string
}

func NewServerContext(mongoUri, dbName, jwtSecret string) ServerContext {
	var (
		client *mongo.Client
		err    error
	)

	// MongoDB
	if client, err = core.CreateMongoClient(mongoUri); err != nil {
		panic(err)
	}

	if err = core.PingDB(client, "Luma"); err != nil {
		panic(err)
	}

	log.Printf("Connected to MongoDB\n")

	ctx := &ServerContext{
		Db:        client.Database(dbName),
		Client:    client,
		JwtSecret: jwtSecret,
	}
	return *ctx
}
