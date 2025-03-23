package models

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	MongoUri  string
	DbName    string
	Port      string
	JwtSecret string
}

func GetEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoUri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")
	port := ":" + os.Getenv("PORT")
	jwtSecret := os.Getenv("JWT_SECRET")

	return &Env{
		MongoUri:  mongoUri,
		DbName:    dbName,
		Port:      port,
		JwtSecret: jwtSecret,
	}
}
