package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/khalidibnwalid/Luma/core"
	"github.com/khalidibnwalid/Luma/middlewares"
	"github.com/khalidibnwalid/Luma/models"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type msgRes struct {
	models.Message
	AuthorID string      `json:"-"`
	Author   models.User `json:"author"`
}

func (ctx *HandlerContext) RoomWS(store *core.TopicStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := r.PathValue("id")
		if roomID == "" {
			http.Error(w, "Room ID is required", http.StatusBadRequest)
			return
		}

		roomData := models.Room{}
		if err := roomData.FindById(ctx.Db, roomID); err != nil {
			http.Error(w, "Room not found", http.StatusNotFound)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}
		defer conn.Close()

		log.Printf("Room [%s] Connected\n", roomID)
		room := store.GetOrCreateRoom(roomData.ID.String())
		room.Subscribe(conn)
		defer room.Unsubscribe(conn)

		userId := r.Context().Value(middlewares.CtxUserIDKey).(string)
		userData := models.User{}
		userData.FindByID(ctx.Db, userId)

		for {
			var msg msgRes // only includes the msg

			// needs a validator
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Println("Read error:", err)
			}

			msg.RoomID = roomID
			msg.AuthorID = userId
			msg.Create(ctx.Db)
			msg.Author = userData

			json, _ := json.Marshal(msg)

			room.Publish([]byte(json))
		}
	}
}
