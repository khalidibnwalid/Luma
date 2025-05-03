package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
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

func (ctx *ServerContext) validateRoomID(w http.ResponseWriter, r *http.Request) (*models.Room, error) {
	rCtx := r.Context()
	roomID := r.PathValue("id")
	w.Header().Set("Content-Type", "application/json")

	if roomID == "" {
		newErrorResponse(w, http.StatusBadRequest, EnumBadRequest, "Room ID is required")
		return &models.Room{}, nil
	}

	uuidRoomID, err := uuid.Parse(roomID)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, EnumBadRequest, "Invalid Room ID format")
		return &models.Room{}, nil
	}

	roomData := models.Room{}
	if err := roomData.FindByID(ctx.Database.Client.WithContext(rCtx), uuidRoomID); err != nil {
		newErrorResponse(w, http.StatusNotFound, EnumNotFound, "Room not found")
		return &models.Room{}, nil
	}

	return &roomData, nil
}

func (ctx *ServerContext) WSRoom(store *core.TopicStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rCtx := r.Context()
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

		log.Printf("Room [%s] Connected\n", room.ID)
		roomTopic := store.GetOrCreateRoom(room.ID.String())
		roomTopic.Subscribe(conn)
		defer roomTopic.Unsubscribe(conn)

		userId := rCtx.Value(middlewares.CtxUserIDKey).(uuid.UUID)
		user := models.NewUser().WithID(userId)
		user.FindByID(ctx.Database.Client) // not finding the user

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
			msg.RoomID = room.ID
			msg.ServerID = room.ServerID
			msg.AuthorID = user.ID
			msg.Create(ctx.Database.Client.WithContext(rCtx))

			msg.Author = *user

			json, _ := json.Marshal(msg)
			roomTopic.Publish([]byte(json))
		}
	}
}

func (ctx *ServerContext) GETRoomMessages(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()
	room, err := ctx.validateRoomID(w, r)
	if err != nil {
		return
	}

	messages, err := room.GetMessages(ctx.Database.Client.WithContext(rCtx))
	if err != nil {
		http.Error(w, "Error fetching messages", http.StatusInternalServerError)
		return
	}

	json, _ := json.Marshal(messages)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (ctx *ServerContext) PatchRoomStatus(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()
	room, err := ctx.validateRoomID(w, r)
	if err != nil {
		return
	}

	userId := rCtx.Value(middlewares.CtxUserIDKey).(uuid.UUID)

	var body struct {
		LastReadMsgID string `json:"lastReadMsgId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		newErrorResponse(w, http.StatusBadRequest, EnumBadRequest, "LastReadMsgId is required")
		return
	}

	uuidLastMsgID, err := uuid.Parse(body.LastReadMsgID)
	if err != nil {
		log.Println("Error parsing LastReadMsgId:", err)
		newErrorResponse(w, http.StatusBadRequest, EnumBadRequest, "Invalid LastReadMsgId format")
		return
	}

	status := models.NewRoomUserStatus().WithUserID(userId).WithRoomID(room.ID)
	if err := status.Find(ctx.Database.Client.WithContext(rCtx)); err != nil {
		newErrorResponse(w, http.StatusNotFound, EnumNotFound, "Room status not found")
		return
	}

	if uuidLastMsgID == status.LastReadMsgID {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	status.LastReadMsgID = uuidLastMsgID

	if err := status.Update(ctx.Database.Client.WithContext(rCtx)); err != nil {
		newErrorResponse(w, http.StatusInternalServerError, EnumInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
