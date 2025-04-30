package core

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBClient struct {
	Client *gorm.DB
}

// postgres URI: postgres://<username>:<password>@<host>:<port>/<dbname>?sslmode=disable
func CreateClient(uri string) (*DBClient, error) {
	pg, err := ParsePostgresURI(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Postgres URI: %w", err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", pg.Host, pg.Username, pg.Password, pg.Database, pg.Port, pg.SSLMode, pg.TimeZone),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
	}), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return &DBClient{
		Client: db,
	}, nil
}

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

type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
	TimeZone string
}

// postgres URI: postgres://<username>:<password>@<host>:<port>/<dbname>?sslmode=disable
func ParsePostgresURI(uri string) (*PostgresConfig, error) {
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URI: %w", err)
	}

	if parsedURL.Scheme != "postgres" {
		return nil, fmt.Errorf("invalid scheme: %s, expected postgres", parsedURL.Scheme)
	}

	// Extract username and password from userinfo
	username := ""
	password := ""
	if parsedURL.User != nil {
		username = parsedURL.User.Username()
		password, _ = parsedURL.User.Password()
	}

	// Extract host and port
	host := parsedURL.Hostname()
	port := parsedURL.Port()
	if port == "" {
		port = "5432" // Default PostgreSQL port
	}

	// Extract database name from path
	database := strings.TrimPrefix(parsedURL.Path, "/")

	// Extract sslmode from query parameters
	sslMode := "disable"
	timeZone := "UTC"
	if parsedURL.Query().Get("sslmode") != "" {
		sslMode = parsedURL.Query().Get("sslmode")
	}
	if parsedURL.Query().Get("TimeZone") != "" {
		timeZone = parsedURL.Query().Get("TimeZone")
	}

	return &PostgresConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Database: database,
		SSLMode:  sslMode,
		TimeZone: timeZone,
	}, nil
}
