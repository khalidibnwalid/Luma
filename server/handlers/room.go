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

func (ctx *HandlerContext) validateRoomID(w http.ResponseWriter, r *http.Request) (models.Room, error) {
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

func (ctx *HandlerContext) RoomWS(store *core.TopicStore) http.HandlerFunc {
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

		log.Printf("Room [%s] Connected\n", room.ID.String())
		roomTopic := store.GetOrCreateRoom(room.ID.String())
		roomTopic.Subscribe(conn)
		defer roomTopic.Unsubscribe(conn)

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

			msg.RoomID = room.ID.String()
			msg.AuthorID = userId
			msg.Create(ctx.Db)
			msg.Author = userData

			json, _ := json.Marshal(msg)

			roomTopic.Publish([]byte(json))
		}
	}
}

func (ctx *HandlerContext) RoomMessagesGET(w http.ResponseWriter, r *http.Request) {
	room, err := ctx.validateRoomID(w, r)
	if err != nil {
		return
	}

	msg := models.Message{}
	messages, err := msg.GetAllMessages(ctx.Db, room.ID.String(), 50)
	if err != nil {
		http.Error(w, "Error fetching messages", http.StatusInternalServerError)
		return
	}

	json, _ := json.Marshal(messages)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
