package handlers

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type HandlerContext struct {
	Db *mongo.Database
	Client *mongo.Client
	JwtSecret string
}
