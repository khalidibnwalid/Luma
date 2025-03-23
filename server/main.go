package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/khalidibnwalid/Luma/core"
	"github.com/khalidibnwalid/Luma/handlers"
	"github.com/khalidibnwalid/Luma/middlewares"
	"github.com/khalidibnwalid/Luma/models"
)

func main() {
	env := models.GetEnv()
	ctx := handlers.NewHandlerContext(env.MongoUri, env.DbName, env.JwtSecret)
	defer func() {
		_ = ctx.Client.Disconnect(context.Background())
	}()
	topicStore := core.NewTopicStore()

	// Server
	authedRoutes := core.NewApp()
	authedRoutes.Use(middlewares.Logging)
	authedRoutes.Use(middlewares.JwtAuthBuilder(env.JwtSecret))

	// user data routes
	authedRoutes.HandleFunc("GET /user/{username}", ctx.UserGET)
	// server routes
	// authedRoutes.HandleFunc("GET /server/{id}", ctx.RoomsServerGET)
	// authedRoutes.HandleFunc("POST /server", ctx.RoomsServerPOST)
	// room routes
	authedRoutes.HandleFunc("GET /room/{id}/messages", ctx.RoomMessagesGET)
	authedRoutes.HandleFunc("/room/{id}", ctx.RoomWS(topicStore))

	unAuthedRoutes := core.NewApp()
	unAuthedRoutes.Use(middlewares.Logging)
	unAuthedRoutes.HandleFunc("POST /login", ctx.AuthLogin)

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", authedRoutes.Mux))
	v1.Handle("/v1/auth/", http.StripPrefix("/v1/auth", unAuthedRoutes.Mux))

	server := http.Server{
		Addr:    env.Port,
		Handler: v1,
	}

	log.Printf("Server Listening on port %s\n", strings.Replace(env.Port, ":", "", 1))
	server.ListenAndServe()
}
