package models

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const RoomUserStatusCollection = "room_user_status"

// RoomUserStatus tracks the read status of a room for a user
type RoomUserStatus struct {
	ID            bson.ObjectID `bson:"_id" json:"id"`
	UserID        string        `bson:"user_id" json:"userId"`
	ServerID      string        `bson:"server_id" json:"serverId"`
	RoomID        string        `bson:"room_id" json:"roomId"`
	LastReadMsgID string        `bson:"last_read_msg_id" json:"lastReadMsgId"`
}

type RoomWithStatus struct {
	*Room
	Status *RoomUserStatus `json:"status"`
}

func NewRoomUserStatus() *RoomUserStatus {
	return &RoomUserStatus{}
}

func (r *RoomUserStatus) WithHexID(id string) *RoomUserStatus {
	objID, _ := bson.ObjectIDFromHex(id)
	r.ID = objID
	return r
}

func (r *RoomUserStatus) WithObjID(id bson.ObjectID) *RoomUserStatus {
	r.ID = id
	return r
}

// The ID should in HEX format like "xxxxxxxxxxxxxxxxxxxxxxxx" not ObjectID("xxxxxxxxxxxxxxxxxxxxxxxx")
func (r *RoomUserStatus) WithUserID(userID string) *RoomUserStatus {
	r.UserID = userID
	return r
}

// The ID should in HEX format like "xxxxxxxxxxxxxxxxxxxxxxxx" not ObjectID("xxxxxxxxxxxxxxxxxxxxxxxx")
func (r *RoomUserStatus) WithRoomID(roomID string) *RoomUserStatus {
	r.RoomID = roomID
	return r
}

func (r *RoomUserStatus) Create(db *mongo.Database, ctx context.Context) error {
	r.ID = bson.NewObjectID()

	coll := db.Collection(RoomUserStatusCollection)
	if _, err := coll.InsertOne(ctx, r); err != nil {
		return err
	}

	return nil
}

// only updates the LastReadMsgID fields
func (r *RoomUserStatus) Update(db *mongo.Database, ctx context.Context) error {
	coll := db.Collection(RoomUserStatusCollection)

	filter := bson.M{"_id": r.ID}
	update := bson.M{
		"$set": bson.M{
			"last_read_msg_id": r.LastReadMsgID,
		},
	}

	if _, err := coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

// needs room_id and last_read_msg_id fields to be set in the struct
func (r *RoomUserStatus) UpdateAllUsersStatus(db *mongo.Database, ctx context.Context, users []string) error {
	coll := db.Collection(ServerUserStatusCollection)

	update := bson.M{
		"$set": bson.M{
			"last_read_msg_id": r.LastReadMsgID,
		},
	}

	filter := bson.M{
		"user_id": bson.M{"$in": users},
		"room_id": r.RoomID,
	}

	if _, err := coll.UpdateMany(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (r *RoomUserStatus) Delete(db *mongo.Database, ctx context.Context) error {
	coll := db.Collection(RoomUserStatusCollection)

	filter := bson.M{"_id": r.ID}
	if _, err := coll.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}

func (r *RoomUserStatus) FindById(db *mongo.Database, ctx context.Context, id ...string) error {
	coll := db.Collection(RoomUserStatusCollection)

	var (
		objId bson.ObjectID
		err   error
	)

	if len(id) > 0 {
		if objId, err = bson.ObjectIDFromHex(id[0]); err != nil {
			return err
		}
	} else {
		objId = r.ID
	}

	filter := bson.M{"_id": objId}
	err = coll.FindOne(ctx, filter).Decode(r)
	if err != nil {
		return err
	}

	return nil
}

func (r *RoomUserStatus) FindByUserIdAndRoomId(db *mongo.Database, ctx context.Context) error {
	coll := db.Collection(RoomUserStatusCollection)

	filter := bson.M{
		"user_id": r.UserID,
		"room_id": r.RoomID,
	}
	err := coll.FindOne(ctx, filter).Decode(&r)
	if err != nil {
		return err
	}

	return nil
}
