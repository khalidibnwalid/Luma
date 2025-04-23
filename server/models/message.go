package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Message struct {
	ID        bson.ObjectID `bson:"_id" json:"id"`
	AuthorID  string        `bson:"author_id" json:"authorId"`
	ServerID  string        `bson:"server_id" json:"serverId"`
	RoomID    string        `bson:"room_id" json:"roomId"`
	Content   string        `bson:"content" json:"content"`
	CreatedAt int64         `bson:"created_at" json:"createdAt"`
	UpdatedAt int64         `bson:"updated_at" json:"updatedAt"`
	// Author shouldn't be stored in the database, only AuthorID
	Author User `bson:"author" json:"author"`
}

func NewMessage() *Message {
	return &Message{}
}

func (msg *Message) WithMessage(message string) *Message {
	msg.Content = message
	return msg
}

func (msg *Message) WithAuthorID(authorID string) *Message {
	msg.AuthorID = authorID
	return msg
}

func (msg *Message) WithRoomID(roomID string) *Message {
	msg.RoomID = roomID
	return msg
}

func (msg *Message) Create(db *mongo.Database, ctx context.Context) error {
	msg.ID = bson.NewObjectID()
	msg.CreatedAt = time.Now().Unix()
	msg.UpdatedAt = time.Now().Unix()

	coll := db.Collection("messages")
	if _, err := coll.InsertOne(ctx, msg); err != nil {
		return err
	}

	return nil
}

func (msg *Message) Update(db *mongo.Database, ctx context.Context) error {
	msg.UpdatedAt = time.Now().Unix()
	coll := db.Collection("messages")
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": msg.ID}, bson.M{"$set": msg}); err != nil {
		return err
	}
	return nil
}

func (msg *Message) Delete(db *mongo.Database, ctx context.Context) error {
	coll := db.Collection("messages")
	if _, err := coll.DeleteOne(ctx, bson.M{"_id": msg.ID}); err != nil {
		return err
	}
	return nil
}
