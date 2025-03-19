package core

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func CreateMongoClient(mongodbUrl string) (*mongo.Client, error) {
	var (
		client *mongo.Client
		err    error
	)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongodbUrl).SetServerAPIOptions(serverAPI)

	if client, err = mongo.Connect(opts); err != nil {
		return nil, err
	}

	return client, nil
}

func PingDB(client *mongo.Client, databaseName string) error {
	var result bson.M
	if err := client.Database(databaseName).RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		return err
	}

	return nil
}
