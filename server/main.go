package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/khalidibnwalid/Luma/core"
	"github.com/khalidibnwalid/Luma/handlers"
	"github.com/khalidibnwalid/Luma/middlewares"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {
	var (
		client *mongo.Client
		err    error
	)

	// Environment Variables
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("Loaded .env file")

	mongoUri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")
	port := ":" + os.Getenv("PORT")
	jwtSecret := os.Getenv("JWT_SECRET")

	// MongoDB
	if client, err = core.CreateMongoClient(mongoUri); err != nil {
		panic(err)
	}

	if err = core.PingDB(client, "Luma"); err != nil {
		panic(err)
	}

	defer func() {
		err = client.Disconnect(context.Background())
	}()
	log.Printf("Connected to MongoDB\n")

	ctx := &handlers.HandlerContext{
		Db:        client.Database(dbName),
		Client:    client,
		JwtSecret: jwtSecret,
	}

	rooms := handlers.NewRooms()

	// Server
	authedRoutes := core.NewApp()
	authedRoutes.Use(middlewares.Logging)
	authedRoutes.Use(middlewares.JwtAuthBuilder(jwtSecret))

	authedRoutes.HandleFunc("GET /user/{username}", ctx.UserGET)
	authedRoutes.HandleFunc("/room/{id}", ctx.RoomWS(rooms))
	// app.HandleFunc("POST /user", ) // signup
	// app.HandleFunc("PATCH /user/{id}", )
	// app.HandleFunc("DELETE /user/{id}", )

	unAuthedRoutes := core.NewApp()
	unAuthedRoutes.Use(middlewares.Logging)
	unAuthedRoutes.HandleFunc("POST /login", ctx.AuthLogin)
	// unAuthedRoutes.HandleFunc("GET /login", ctx.AuthGET) //test

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", authedRoutes.Mux))
	v1.Handle("/v1/auth/", http.StripPrefix("/v1/auth", unAuthedRoutes.Mux))

	server := http.Server{
		Addr:    port,
		Handler: v1,
	}

	log.Printf("Server Listening on port %s\n", strings.Replace(port, ":", "", 1))
	server.ListenAndServe()
}
