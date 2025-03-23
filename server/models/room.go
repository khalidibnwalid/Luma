package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Room struct {
	ID        bson.ObjectID `bson:"_id" json:"id"`
	ServerID  bson.ObjectID `bson:"server_id" json:"serverId"`
	Type      string        `bson:"type" json:"type"` // direct, server room, or group
	CreatedAt int64         `bson:"created_at" json:"createdAt"`
	UpdatedAt int64         `bson:"updated_at" json:"updatedAt"`
}

func (r *Room) Create(db *mongo.Database) error {
	r.ID = bson.NewObjectID()
	r.CreatedAt = time.Now().Unix()
	r.UpdatedAt = time.Now().Unix()

	coll := db.Collection("rooms")
	if _, err := coll.InsertOne(context.TODO(), r); err != nil {
		return err
	}

	return nil
}

func (r *Room) FindById(db *mongo.Database, id string) error {
	coll := db.Collection("rooms")

	// Convert the string ID to a bson.ObjectID, since the ID in the database is an ObjectID
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	if err := coll.FindOne(context.TODO(), bson.M{"_id": objId}).Decode(&r); err != nil {
		return err
	}
	return nil
}
