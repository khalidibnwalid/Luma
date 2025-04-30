package handlers

import (
	"log"

	"github.com/khalidibnwalid/Luma/core"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// TODO remove mongoDB Db and Client
type ServerContext struct {
	Db        *mongo.Database
	Client    *mongo.Client
	Database  *core.DBClient
	JwtSecret string
}

func NewServerContext(mongoUri, postgresUri, dbName, jwtSecret string) *ServerContext {
	var (
		client *mongo.Client
		err    error
	)

	log.Printf("Connecting to MongoDB...")

	// MongoDB 	// TODO delete this
	if client, err = core.CreateMongoClient(mongoUri); err != nil {
		log.Panicf("MongoDB connection error: %v", err)
	}

	log.Printf("Pinging MongoDB\n")

	if err = core.PingDB(client, "Luma"); err != nil {
		log.Panicf("MongoDB ping error: %v", err)
	}

	log.Printf("MongoDB Connected!\n")

	// Postgres
	db, err := core.CreateClient(postgresUri)
	if err != nil {
		log.Panicf("Postgres connection error: %v", err)
	}

	log.Printf("Pinging Postgres\n")
	if err = db.Client.Exec("SELECT 1").Error; err != nil {
		log.Panicf("Postgres ping error: %v", err)
	}

	log.Printf("Postgres Connected!\n")

	ctx := &ServerContext{
		Db:        client.Database(dbName),
		Client:    client,
		Database:  db,
		JwtSecret: jwtSecret,
	}
	return ctx
}
