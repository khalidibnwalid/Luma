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
	RoomID        string        `bson:"room_id" json:"roomId"`
	LastReadMsgID string        `bson:"last_read_msg_id" json:"lastReadMsgId"`
	// IsCleared indicates that the user has read all messages in the room,
	// thus each new messge will reset the LastReadMsgID to the new message ID.
	// if false, the LastReadMsgID won't be updated when a new message arrives
	IsCleared bool `bson:"is_cleared" json:"isCleared"`
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

func (r *RoomUserStatus) Create(db *mongo.Database) error {
	r.ID = bson.NewObjectID()

	coll := db.Collection(RoomUserStatusCollection)
	if _, err := coll.InsertOne(context.TODO(), r); err != nil {
		return err
	}

	return nil
}

func (r *RoomUserStatus) FindById(db *mongo.Database, id ...string) error {
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
	err = coll.FindOne(context.TODO(), filter).Decode(r)
	if err != nil {
		return err
	}

	return nil
}
