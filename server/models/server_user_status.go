package models

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const ServerUserStatusCollection = "server_user_status"

// ServerUserStatus tracks the status of a user in a server
type ServerUserStatus struct {
	ID       bson.ObjectID `bson:"_id" json:"id"`
	UserID   string        `bson:"user_id" json:"userId"`
	ServerID string        `bson:"server_id" json:"serverId"`
	Nickname string        `bson:"nickname" json:"nickname"`
	Roles    []string      `bson:"roles" json:"roles"`
}

func NewServerUserStatus() *ServerUserStatus {
	return &ServerUserStatus{}
}

func (s *ServerUserStatus) WithHexID(id string) *ServerUserStatus {
	objID, _ := bson.ObjectIDFromHex(id)
	s.ID = objID
	return s
}

func (s *ServerUserStatus) WithObjID(id bson.ObjectID) *ServerUserStatus {
	s.ID = id
	return s
}

// The ID should in HEX format like "xxxxxxxxxxxxxxxxxxxxxxxx" not ObjectID("xxxxxxxxxxxxxxxxxxxxxxxx")
func (s *ServerUserStatus) WithUserID(userID string) *ServerUserStatus {
	s.UserID = userID
	return s
}

// The ID should in HEX format like "xxxxxxxxxxxxxxxxxxxxxxxx" not ObjectID("xxxxxxxxxxxxxxxxxxxxxxxx")
func (s *ServerUserStatus) WithServerID(serverID string) *ServerUserStatus {
	s.ServerID = serverID
	return s
}

func (s *ServerUserStatus) Create(db *mongo.Database) error {
	s.ID = bson.NewObjectID()

	coll := db.Collection(ServerUserStatusCollection)
	if _, err := coll.InsertOne(context.TODO(), s); err != nil {
		return err
	}

	return nil
}

func (s *ServerUserStatus) Update(db *mongo.Database) error {
	coll := db.Collection(ServerUserStatusCollection)

	update := bson.M{
		"$set": bson.M{
			"nickname": s.Nickname,
			"roles":    s.Roles,
		},
	}

	filter := bson.M{"_id": s.ID}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServerUserStatus) Delete(db *mongo.Database) error {
	coll := db.Collection(ServerUserStatusCollection)
	filter := bson.M{"_id": s.ID}
	_, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServerUserStatus) FindById(db *mongo.Database, id ...string) error {
	coll := db.Collection(ServerUserStatusCollection)

	var (
		objId bson.ObjectID
		err   error
	)

	if len(id) > 0 {
		objId, err = bson.ObjectIDFromHex(id[0])
		if err != nil {
			return err
		}
	} else {
		objId = s.ID
	}

	filter := bson.M{"_id": objId}
	err = coll.FindOne(context.TODO(), filter).Decode(s)
	if err != nil {
		return err
	}

	return nil
}
