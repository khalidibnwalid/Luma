package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/khalidibnwalid/Luma/core"
	"github.com/khalidibnwalid/Luma/handlers"
	"github.com/khalidibnwalid/Luma/server/middlewares"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {
	mongodbUrl := "mongodb://root:example@localhost:27017/"

	var (
		client *mongo.Client
		err    error
	)
	// MongoDB
	if client, err = core.CreateMongoClient(mongodbUrl); err != nil {
		panic(err)
	}

	if err = core.PingDB(client, "Luma"); err != nil {
		panic(err)
	}

	defer func() {
		err = client.Disconnect(context.Background())
	}()
	fmt.Printf("Connected to MongoDB\n")

	ctx := &handlers.HandlerContext{
		Db:     client.Database("Luma"),
		Client: client,
	}

	// Server
	app := core.NewApp()
	app.Use(middlewares.Logging)

	app.HandleFunc("GET /user/{username}", ctx.UserGET)
	// app.HandleFunc("POST /user", ) // signup
	// app.HandleFunc("PATCH /user/{id}", )
	// app.HandleFunc("DELETE /user/{id}", )

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", app.Mux))

	server := http.Server{
		Addr:    ":8080",
		Handler: v1,
	}

	fmt.Printf("Server Listening on port 8080\n")
	server.ListenAndServe()
}
