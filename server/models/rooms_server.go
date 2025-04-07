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

func (rs *RoomsServer) Create(db *mongo.Database) error {
	rs.ID = bson.NewObjectID()
	rs.CreatedAt = time.Now().Unix()
	rs.UpdatedAt = time.Now().Unix()

	coll := db.Collection(RoomsServerCollection)
	if _, err := coll.InsertOne(context.TODO(), rs); err != nil {
		return err
	}

	return nil
}

// You can provide the ID as a parameter or use the ID from the struct
func (rs *RoomsServer) FindById(db *mongo.Database, id ...string) error {
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

	if err := coll.FindOne(context.TODO(), bson.M{"_id": objId}).Decode(&rs); err != nil {
		return err
	}
	return nil
}

// You can provide the UserId as a parameter or use the ID from the struct
// func (rs *RoomsServer) GetAllServersOfUser(db *mongo.Database, ownerID ...string) ([]RoomsServer, error) {
// 	coll := db.Collection(ServerCollection)

// 	var (
// 		OwenrID string
// 		err     error
// 		cursor  *mongo.Cursor
// 	)

// 	// we need the ID as a hex string
// 	if len(ownerID) > 0 {
// 		OwenrID = ownerID[0]
// 	} else {
// 		OwenrID = rs.ID.Hex()
// 	}

// 	if cursor, err = coll.Find(context.TODO(), bson.M{"owner_id": OwenrID}); err != nil {
// 		return nil, err
// 	}

// 	var servers []RoomsServer
// 	if err := cursor.All(context.Background(), &servers); err != nil {
// 		return nil, err
// 	}

// 	return servers, nil
// }

// You can provide the ServerId as a parameter or use the ID from the struct
func (rs *RoomsServer) GetRooms(db *mongo.Database, serverID ...string) ([]Room, error) {
	coll := db.Collection("rooms")

	var (
		objId bson.ObjectID
		err   error
	)
	if len(serverID) > 0 {
		if objId, err = bson.ObjectIDFromHex(serverID[0]); err != nil {
			return nil, err
		}
	} else {
		objId = rs.ID
	}

	cursor, err := coll.Find(context.TODO(), bson.M{"server_id": objId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var rooms []Room
	if err := cursor.All(context.Background(), &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (rs *RoomsServer) Update(db *mongo.Database) error {
	rs.UpdatedAt = time.Now().Unix()
	coll := db.Collection(RoomsServerCollection)
	if _, err := coll.UpdateOne(context.TODO(), bson.M{"_id": rs.ID}, bson.M{"$set": rs}); err != nil {
		return err
	}
	return nil
}

func (rs *RoomsServer) Delete(db *mongo.Database) error {
	coll := db.Collection(RoomsServerCollection)
	if _, err := coll.DeleteOne(context.TODO(), bson.M{"_id": rs.ID}); err != nil {
		return err
	}
	return nil
}
