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

type RoomsServerWithStatus struct {
	*RoomsServer
	Status ServerUserStatus `json:"status"`
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

func (s *ServerUserStatus) Create(db *mongo.Database, ctx context.Context) error {
	s.ID = bson.NewObjectID()

	coll := db.Collection(ServerUserStatusCollection)
	if _, err := coll.InsertOne(ctx, s); err != nil {
		return err
	}

	return nil
}

func (s *ServerUserStatus) Update(db *mongo.Database, ctx context.Context) error {
	coll := db.Collection(ServerUserStatusCollection)

	update := bson.M{
		"$set": bson.M{
			"nickname": s.Nickname,
			"roles":    s.Roles,
		},
	}

	filter := bson.M{"_id": s.ID}
	_, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServerUserStatus) Delete(db *mongo.Database, ctx context.Context) error {
	coll := db.Collection(ServerUserStatusCollection)
	filter := bson.M{"_id": s.ID}
	_, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServerUserStatus) FindById(db *mongo.Database, ctx context.Context, id ...string) error {
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
	err = coll.FindOne(ctx, filter).Decode(s)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServerUserStatus) GetServers(db *mongo.Database, ctx context.Context) ([]RoomsServerWithStatus, error) {
	coll := db.Collection(RoomsServerCollection)

	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from": ServerUserStatusCollection,
				"let": bson.M{
					"serverId": "$_id",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": bson.A{
									bson.M{"$eq": bson.A{bson.M{"$toObjectId": "$server_id"}, "$$serverId"}},
									bson.M{"$eq": bson.A{"$user_id", s.UserID}},
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
		{ // for marshalling the RoomsServerWithStatus struct, which has servers as an embedded struct
			"$project": bson.M{
				"_id":    0,
				"status": 1,
				"roomsserver": bson.M{
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

	var servers []RoomsServerWithStatus
	if err = cursor.All(ctx, &servers); err != nil {
		return nil, err
	}

	return servers, nil
}
