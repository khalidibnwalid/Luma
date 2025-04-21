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

func NewServerContext(mongoUri, dbName, jwtSecret string) *ServerContext {
	var (
		client *mongo.Client
		err    error
	)

	log.Printf("Connecting to MongoDB...")

	// MongoDB
	if client, err = core.CreateMongoClient(mongoUri); err != nil {
		log.Panicf("MongoDB connection error: %v", err)
	}

	log.Printf("Pinging MongoDB\n")

	if err = core.PingDB(client, "Luma"); err != nil {
		log.Panicf("MongoDB ping error: %v", err)
	}

	log.Printf("MongoDB Connected!\n")

	ctx := &ServerContext{
		Db:        client.Database(dbName),
		Client:    client,
		JwtSecret: jwtSecret,
	}
	return ctx
}
