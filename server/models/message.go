package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Message struct {
	ID        bson.ObjectID `bson:"_id" json:"id"`
	AuthorID  string        `bson:"author_id" json:"authorId"`
	RoomID    string        `bson:"room_id" json:"roomId"`
	Message   string        `bson:"message" json:"message"`
	CreatedAt int64         `bson:"created_at" json:"createdAt"`
	UpdatedAt int64         `bson:"updated_at" json:"updatedAt"`
}

func (msg *Message) Create(db *mongo.Database) error {
	msg.ID = bson.NewObjectID()
	msg.CreatedAt = time.Now().Unix()
	msg.UpdatedAt = time.Now().Unix()

	coll := db.Collection("messages")
	if _, err := coll.InsertOne(context.TODO(), msg); err != nil {
		return err
	}

	return nil
}

// default limit is 50
func (msg *Message) GetAllMessages(db *mongo.Database, room_id string, limit ...int64) ([]Message, error) {
	coll := db.Collection("messages")
	limitVal := int64(50)
	if len(limit) > 0 {
		limitVal = limit[0]
	}
	opts := options.Find().SetLimit(limitVal)

	cursor, err := coll.Find(context.TODO(), bson.M{"room_id": room_id}, opts)
	if err != nil {
		return nil, err
	}

	var messages []Message
	if err := cursor.All(context.Background(), &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func (msg *Message) Update(db *mongo.Database) error {
	msg.UpdatedAt = time.Now().Unix()
	coll := db.Collection("messages")
	if _, err := coll.UpdateOne(context.TODO(), bson.M{"_id": msg.ID}, bson.M{"$set": msg}); err != nil {
		return err
	}
	return nil
}

func (msg *Message) Delete(db *mongo.Database) error {
	coll := db.Collection("messages")
	if _, err := coll.DeleteOne(context.TODO(), bson.M{"_id": msg.ID}); err != nil {
		return err
	}
	return nil
}
