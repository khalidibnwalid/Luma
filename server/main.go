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
	// authedRoutes.HandleFunc("GET /user/{username}", ctx.UserGET)

	// servers routes
	authedRoutes.HandleFunc("GET /servers", ctx.GetUserRoomsServer)
	authedRoutes.HandleFunc("POST /servers", ctx.PostRoomsServer)
	// server rooms routes
	authedRoutes.HandleFunc("GET /servers/{id}", ctx.GetRoomsServer)
	authedRoutes.HandleFunc("GET /servers/{id}/rooms", ctx.GetRoomsOfServer)
	authedRoutes.HandleFunc("POST /servers/{id}/rooms", ctx.PostRoomToServer)

	// room routes
	authedRoutes.HandleFunc("GET /rooms/{id}/messages", ctx.GETRoomMessages)
	authedRoutes.HandleFunc("/rooms/{id}", ctx.WSRoom(topicStore))

	// user routes
	authedRoutes.HandleFunc("GET /users", ctx.GetUser)

	// auth routes
	unAuthedRoutes := core.NewApp()
	unAuthedRoutes.Use(middlewares.Logging)
	unAuthedRoutes.HandleFunc("POST /sessions", ctx.PostSession)
	unAuthedRoutes.HandleFunc("POST /users", ctx.PostUser)

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", authedRoutes.Mux))
	v1.Handle("/v1/auth/", http.StripPrefix("/v1/auth", unAuthedRoutes.Mux))

	server := http.Server{
		Addr:    env.Port,
		Handler: middlewares.CORS(v1),
	}

	log.Printf("Server Listening on port %s\n", strings.Replace(env.Port, ":", "", 1))
	server.ListenAndServe()
}
