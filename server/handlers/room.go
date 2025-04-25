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

func (ctx *ServerContext) validateRoomID(w http.ResponseWriter, r *http.Request) (*models.Room, error) {
	rCtx := r.Context()
	roomID := r.PathValue("id")
	if roomID == "" {
		w.Header().Set("Content-Type", "application/json")
		newErrorResponse(w, http.StatusBadRequest, EnumBadRequest, "Room ID is required")
		return &models.Room{}, nil
	}

	roomData := models.Room{}
	if err := roomData.FindById(ctx.Db, rCtx, roomID); err != nil {
		w.Header().Set("Content-Type", "application/json")
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

		log.Printf("Room [%s] Connected\n", room.ID.Hex())
		roomTopic := store.GetOrCreateRoom(room.ID.Hex())
		roomTopic.Subscribe(conn)
		defer roomTopic.Unsubscribe(conn)

		userId := rCtx.Value(middlewares.CtxUserIDKey).(string)
		user := models.NewUser().WithHexID(userId)
		user.FindByID(ctx.Db, rCtx) // not finding the user

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
			msg.ServerID = room.ServerID
			msg.AuthorID = user.ID.Hex()
			msg.Create(ctx.Db, rCtx)

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

	messages, err := room.GetMessages(ctx.Db, rCtx)
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

	userId := rCtx.Value(middlewares.CtxUserIDKey).(string)

	var body struct {
		LastReadMsgID string `json:"lastReadMsgId"`
		IsCleared     bool   `json:"isCleared"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		newErrorResponse(w, http.StatusBadRequest, EnumBadRequest, "LastReadMsgId and IsCleared are required")
		return
	}

	status := models.NewRoomUserStatus().WithUserID(userId).WithRoomID(room.ID.Hex())
	if err := status.FindByUserIdAndRoomId(ctx.Db, rCtx); err != nil {
		newErrorResponse(w, http.StatusNotFound, EnumNotFound, "Room status not found")
		return
	}

	if body.IsCleared == status.IsCleared ||body.LastReadMsgID == status.LastReadMsgID {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if len(body.LastReadMsgID) > 0 {
		status.IsCleared = false
		status.LastReadMsgID = body.LastReadMsgID
	} else {
		status.IsCleared = true
		status.LastReadMsgID = ""
	}

	if err := status.Update(ctx.Db, rCtx); err != nil {
		newErrorResponse(w, http.StatusInternalServerError, EnumInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
