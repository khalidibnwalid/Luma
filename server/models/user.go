package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type User struct {
	ID             bson.ObjectID `bson:"_id" json:"id"`
	Username       string        `bson:"username"`
	HashedPassword string        `bson:"hashed_password" json:"-"`
	CreatedAt      int64         `bson:"createdAt"`
	UpdatedAt      int64         `bson:"updatedAt"`
}

func (u *User) Create(db *mongo.Database) error {
	u.ID = bson.NewObjectID()
	u.CreatedAt = time.Now().Unix()
	u.UpdatedAt = time.Now().Unix()

	coll := db.Collection("users")
	if _, err := coll.InsertOne(context.TODO(), u); err != nil {
		return err
	}

	return nil
}

func (u *User) FindByUsername(db *mongo.Database, username string) error {
	coll := db.Collection("users")
	if err := coll.FindOne(context.TODO(), bson.M{"username": username}).Decode(&u); err != nil {
		return err
	}
	return nil
}

func (u *User) FindByID(db *mongo.Database, id string) error {
	coll := db.Collection("users")
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	if err := coll.FindOne(context.TODO(), bson.M{"_id": objId}).Decode(&u); err != nil {
		return err
	}
	return nil
}

func (u *User) Update(db *mongo.Database) error {
	u.UpdatedAt = time.Now().Unix()
	coll := db.Collection("users")
	if _, err := coll.UpdateOne(context.TODO(), bson.M{"_id": u.ID}, bson.M{"$set": u}); err != nil {
		return err
	}
	return nil
}

func (u *User) Delete(db *mongo.Database) error {
	coll := db.Collection("users")
	if _, err := coll.DeleteOne(context.TODO(), bson.M{"_id": u.ID}); err != nil {
		return err
	}
	return nil
}
