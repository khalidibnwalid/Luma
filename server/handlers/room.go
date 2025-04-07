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
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (ctx *ServerContext) validateRoomID(w http.ResponseWriter, r *http.Request) (models.Room, error) {
	roomID := r.PathValue("id")
	if roomID == "" {
		http.Error(w, "Room ID is required", http.StatusBadRequest)
		return models.Room{}, nil
	}

	roomData := models.Room{}
	if err := roomData.FindById(ctx.Db, roomID); err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return models.Room{}, nil
	}

	return roomData, nil
}

func (ctx *ServerContext) WSRoom(store *core.TopicStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		room, err := ctx.validateRoomID(w, r)
		if err != nil {
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}
		defer conn.Close()

		log.Printf("Room [%s] Connected\n", room.ID.Hex())
		roomTopic := store.GetOrCreateRoom(room.ID.Hex())
		roomTopic.Subscribe(conn)
		defer roomTopic.Unsubscribe(conn)

		userId := r.Context().Value(middlewares.CtxUserIDKey).(string)
		user := models.NewUser().WithHexID(userId)
		user.FindByID(ctx.Db) // not finding the user

		for {
			var body struct {
				Content string `json:"content"`
			}

			// needs a validator
			err := conn.ReadJSON(&body)
			if err != nil {
				log.Println("Read error:", err)
				roomTopic.Unsubscribe(conn)
				break
			}

			var msg models.Message
			msg.Content = body.Content
			msg.RoomID = room.ID.Hex()
			msg.AuthorID = user.ID.Hex()
			msg.Create(ctx.Db)

			msg.Author = *user

			json, _ := json.Marshal(msg)
			roomTopic.Publish([]byte(json))
		}
	}
}

func (ctx *ServerContext) GETRoomMessages(w http.ResponseWriter, r *http.Request) {
	room, err := ctx.validateRoomID(w, r)
	if err != nil {
		return
	}

	messages, err := room.GetMessages(ctx.Db)
	if err != nil {
		http.Error(w, "Error fetching messages", http.StatusInternalServerError)
		return
	}

	json, _ := json.Marshal(messages)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
