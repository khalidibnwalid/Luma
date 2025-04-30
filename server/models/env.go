package models

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	// TODO delete this
	MongoUri  string
	DbName    string
	//
	Port      string
	JwtSecret string
	PostgresUri string
}

func GetEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Env{
		// TODO delete this
		MongoUri:  os.Getenv("MONGO_URI"),
		DbName:    os.Getenv("DB_NAME"),
		//
		Port:      ":" + os.Getenv("PORT"),
		JwtSecret: os.Getenv("JWT_SECRET"),
		PostgresUri: os.Getenv("POSTGRES_URI"),
	}
}
