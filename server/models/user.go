package models

import (
	"context"
	"time"

	"github.com/khalidibnwalid/Luma/core"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const UserCollection = "users"

type User struct {
	ID             bson.ObjectID `bson:"_id" json:"id"`
	Username       string        `bson:"username" json:"username"`
	Email          string        `bson:"email" json:"-"`
	HashedPassword string        `bson:"hashed_password" json:"-"`
	CreatedAt      int64         `bson:"created_at" json:"createdAt"`
	UpdatedAt      int64         `bson:"updated_at" json:"updatedAt"`
}

func NewUser(username ...string) *User {
	if len(username) == 0 {
		return &User{}
	}

	return &User{
		Username: username[0],
	}
}

func (u *User) WithHexID(id string) *User {
	objID, _ := bson.ObjectIDFromHex(id)
	u.ID = objID
	return u
}

func (u *User) WithObjID(id bson.ObjectID) *User {
	u.ID = id
	return u
}

func (u *User) WithUsername(username string) *User {
	u.Username = username
	return u
}

func (u *User) WithEmail(email string) *User {
	u.Email = email
	return u
}

// Automatically hashes the password
func (u *User) WithPassword(unhashedPassword string) *User {
	hashedPassword, salt, _ := core.CreateHashWithSalt(unhashedPassword)
	u.HashedPassword = core.SerializeHashWithSalt(hashedPassword, salt)
	return u
}

func (u *User) Create(db *mongo.Database, ctx context.Context) error {
	u.ID = bson.NewObjectID()
	u.CreatedAt = time.Now().Unix()
	u.UpdatedAt = time.Now().Unix()

	coll := db.Collection(UserCollection)
	if _, err := coll.InsertOne(ctx, u); err != nil {
		return err
	}

	return nil
}

// You can provide ID as a parameter or in the struct
func (u *User) FindByUsername(db *mongo.Database, ctx context.Context, username ...string) error {
	coll := db.Collection(UserCollection)
	var _username string

	if len(username) > 0 {
		_username = username[0]
	} else {
		_username = u.Username
	}

	if err := coll.FindOne(ctx, bson.M{"username": _username}).Decode(&u); err != nil {
		return err
	}
	return nil
}

// You can provide Email as a parameter or in the struct
func (u *User) FindByEmail(db *mongo.Database, ctx context.Context, email ...string) error {
	coll := db.Collection(UserCollection)
	var _email string

	if len(email) > 0 {
		_email = email[0]
	} else {
		_email = u.Email
	}

	if err := coll.FindOne(ctx, bson.M{"email": _email}).Decode(&u); err != nil {
		return err
	}
	return nil
}

func (u *User) FindByID(db *mongo.Database, ctx context.Context, id ...string) error {
	coll := db.Collection(UserCollection)

	var (
		objId bson.ObjectID
		err   error
	)

	if len(id) > 0 {
		if objId, err = bson.ObjectIDFromHex(id[0]); err != nil {
			return err
		}
	} else {
		objId = u.ID
	}

	if err := coll.FindOne(ctx, bson.M{"_id": objId}).Decode(&u); err != nil {
		return err
	}

	return nil
}

func (u *User) Update(db *mongo.Database, ctx context.Context) error {
	u.UpdatedAt = time.Now().Unix()
	coll := db.Collection(UserCollection)
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": u.ID}, bson.M{"$set": u}); err != nil {
		return err
	}
	return nil
}

func (u *User) Delete(db *mongo.Database, ctx context.Context) error {
	coll := db.Collection(UserCollection)
	if _, err := coll.DeleteOne(ctx, bson.M{"_id": u.ID}); err != nil {
		return err
	}
	return nil
}
