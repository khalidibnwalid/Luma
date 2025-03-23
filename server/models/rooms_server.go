package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const ServerCollection = "rooms_server"

type RoomsServer struct {
	ID        bson.ObjectID `bson:"_id" json:"id"`
	OwnerID   bson.ObjectID `bson:"owner_id" json:"ownerId"`
	Name      string        `bson:"name" json:"name"`
	CreatedAt int64         `bson:"created_at" json:"createdAt"`
	UpdatedAt int64         `bson:"updated_at" json:"updatedAt"`
}

func (r *RoomsServer) Create(db *mongo.Database) error {
	r.ID = bson.NewObjectID()
	r.CreatedAt = time.Now().Unix()
	r.UpdatedAt = time.Now().Unix()

	coll := db.Collection(ServerCollection)
	if _, err := coll.InsertOne(context.TODO(), r); err != nil {
		return err
	}

	return nil
}

func (r *RoomsServer) FindById(db *mongo.Database, id string) error {
	coll := db.Collection(ServerCollection)

	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	if err := coll.FindOne(context.TODO(), bson.M{"_id": objId}).Decode(&r); err != nil {
		return err
	}
	return nil
}

func GetRoomsServersByOwner(db *mongo.Database, ownerID bson.ObjectID) ([]RoomsServer, error) {
	coll := db.Collection(ServerCollection)

	cursor, err := coll.Find(context.TODO(), bson.M{"owner_id": ownerID})
	if err != nil {
		return nil, err
	}

	var servers []RoomsServer
	if err := cursor.All(context.Background(), &servers); err != nil {
		return nil, err
	}

	return servers, nil
}

func (r *RoomsServer) Update(db *mongo.Database) error {
	r.UpdatedAt = time.Now().Unix()
	coll := db.Collection(ServerCollection)
	if _, err := coll.UpdateOne(context.TODO(), bson.M{"_id": r.ID}, bson.M{"$set": r}); err != nil {
		return err
	}
	return nil
}

func (r *RoomsServer) Delete(db *mongo.Database) error {
	coll := db.Collection(ServerCollection)
	if _, err := coll.DeleteOne(context.TODO(), bson.M{"_id": r.ID}); err != nil {
		return err
	}
	return nil
}