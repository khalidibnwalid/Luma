package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const messagesLimit = 50
const RoomsCollection = "rooms"

type Room struct {
	ID        bson.ObjectID `bson:"_id" json:"id"`
	ServerID  bson.ObjectID `bson:"server_id" json:"serverId"`
	Name      string        `bson:"name" json:"name"`
	GroupName string        `bson:"group_name" json:"groupName"`
	Type      string        `bson:"type" json:"type"` // direct, server room, server voice room, or users group,
	CreatedAt int64         `bson:"created_at" json:"createdAt"`
	UpdatedAt int64         `bson:"updated_at" json:"updatedAt"`
}

func NewRoom() *Room {
	return &Room{}
}

func (r *Room) WithHexID(id string) *Room {
	objID, _ := bson.ObjectIDFromHex(id)
	r.ID = objID
	return r
}

func (r *Room) WithObjID(id bson.ObjectID) *Room {
	r.ID = id
	return r
}

func (r *Room) Create(db *mongo.Database) error {
	r.ID = bson.NewObjectID()
	r.CreatedAt = time.Now().Unix()
	r.UpdatedAt = time.Now().Unix()

	coll := db.Collection(RoomsCollection)
	if _, err := coll.InsertOne(context.TODO(), r); err != nil {
		return err
	}

	return nil
}

// You can provide ID as a parameter or in the struct
func (r *Room) FindById(db *mongo.Database, id ...string) error {
	coll := db.Collection(RoomsCollection)

	// Convert the string ID to a bson.ObjectID, since the ID in the database is an ObjectID
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

	if err := coll.FindOne(context.TODO(), bson.M{"_id": objId}).Decode(&r); err != nil {
		return err
	}
	return nil
}

// You can provide ID as a parameter or in the struct
func (r *Room) GetMessages(db *mongo.Database, room_id ...string) ([]Message, error) {
	coll := db.Collection("messages")

	// limitVal := int64(messagesLimit)

	var (
		roomID string
		err    error
	)

	// we need the hexa ID
	if len(room_id) > 0 {
		roomID = room_id[0]
	} else {
		roomID = r.ID.Hex()
	}

	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from": "users",
				"let": bson.M{
					"author_id_str": "$author_id",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$eq": bson.A{
									"$_id",
									bson.M{"$toObjectId": "$$author_id_str"},
								},
							},
						},
					},
				},
				"as": "author",
			},
		},
		{
			"$unwind": bson.M{
				"path": "$author",
			},
		},
		{
			"$match": bson.M{
				"room_id": roomID,
			},
		},
	}

	cursor, err := coll.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// not providing author data
	var MessagesWithAuthors []Message
	if err := cursor.All(context.TODO(), &MessagesWithAuthors); err != nil {
		return nil, err
	}

	return MessagesWithAuthors, nil
}
