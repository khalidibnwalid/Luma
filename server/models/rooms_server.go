package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const RoomsServerCollection = "rooms_server"

type RoomsServer struct {
	ID        bson.ObjectID `bson:"_id" json:"id"`
	OwnerID   string        `bson:"owner_id" json:"ownerId"`
	Name      string        `bson:"name" json:"name"`
	CreatedAt int64         `bson:"created_at" json:"createdAt"`
	UpdatedAt int64         `bson:"updated_at" json:"updatedAt"`
}

func NewRoomsServer(rs ...RoomsServer) *RoomsServer {
	if len(rs) > 0 {
		return &rs[0]
	}
	return &RoomsServer{}
}

func (rs *RoomsServer) WithHexID(id string) *RoomsServer {
	objId, _ := bson.ObjectIDFromHex(id)
	rs.ID = objId
	return rs
}

func (rs *RoomsServer) WithObjID(id bson.ObjectID) *RoomsServer {
	rs.ID = id
	return rs
}

// The ID should in HEX format like "xxxxxxxxxxxxxxxxxxxxxxxx" not ObjectID("xxxxxxxxxxxxxxxxxxxxxxxx")
func (rs *RoomsServer) WithOwnerID(ownerID string) *RoomsServer {
	rs.OwnerID = ownerID
	return rs
}

func (rs *RoomsServer) Create(db *mongo.Database, ctx context.Context) error {
	rs.ID = bson.NewObjectID()
	rs.CreatedAt = time.Now().Unix()
	rs.UpdatedAt = time.Now().Unix()

	coll := db.Collection(RoomsServerCollection)
	if _, err := coll.InsertOne(ctx, rs); err != nil {
		return err
	}

	return nil
}

// You can provide the ID as a parameter or use the ID from the struct
func (rs *RoomsServer) FindById(db *mongo.Database, ctx context.Context, id ...string) error {
	coll := db.Collection(RoomsServerCollection)

	var (
		objId bson.ObjectID
		err   error
	)
	if len(id) > 0 {
		if objId, err = bson.ObjectIDFromHex(id[0]); err != nil {
			return err
		}
	} else {
		objId = rs.ID
	}

	if err := coll.FindOne(ctx, bson.M{"_id": objId}).Decode(&rs); err != nil {
		return err
	}
	return nil
}

func (rs *RoomsServer) Update(db *mongo.Database, ctx context.Context) error {
	rs.UpdatedAt = time.Now().Unix()
	coll := db.Collection(RoomsServerCollection)
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": rs.ID}, bson.M{"$set": rs}); err != nil {
		return err
	}
	return nil
}

func (rs *RoomsServer) Delete(db *mongo.Database, ctx context.Context) error {
	coll := db.Collection(RoomsServerCollection)
	if _, err := coll.DeleteOne(ctx, bson.M{"_id": rs.ID}); err != nil {
		return err
	}
	return nil
}


// You can provide the ServerId as a parameter or use the ID from the struct
// TODO: add user's role check
func (rs *RoomsServer) GetRooms(db *mongo.Database, ctx context.Context, userId string) ([]RoomWithStatus, error) {
	coll := db.Collection(RoomsCollection)

	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from": RoomUserStatusCollection,
				"let": bson.M{
					"roomId": "$_id",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": bson.A{
									bson.M{"$eq": bson.A{bson.M{"$toObjectId": "$room_id"}, "$$roomId"}},
									bson.M{"$eq": bson.A{"$server_id", rs.ID.Hex()}},
									bson.M{"$eq": bson.A{"$user_id", userId}},
								},
							},
						},
					},
				},
				"as": "status",
			},
		},
		{
			"$unwind": bson.M{
				"path": "$status",
			},
		},
		{
			"$match": bson.M{
				"status": bson.M{"$ne": nil},
			},
		},
		{ // for marshalling the RoomWithStatus struct, which has room as an embedded struct
			"$project": bson.M{
				"_id":    0,
				"status": 1,
				"room": bson.M{
					"$mergeObjects": bson.A{
						bson.M{"status": 0},
						"$$ROOT",
					},
				},
			},
		},
	}

	var (
		err    error
		cursor *mongo.Cursor
	)
	if cursor, err = coll.Aggregate(ctx, pipeline); err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rooms []RoomWithStatus
	if err = cursor.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}