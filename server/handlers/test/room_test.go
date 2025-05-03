package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/khalidibnwalid/Luma/core"
	"github.com/khalidibnwalid/Luma/handlers"
	"github.com/khalidibnwalid/Luma/middlewares"
	"github.com/khalidibnwalid/Luma/models"
	"github.com/khalidibnwalid/Luma/testutil"
)

func TestGETRoomMessages(t *testing.T) {
	ctx := testutil.NewTestingContext(t)

	// TODO : check for maximum messages returned
	t.Run("Should get messages of a room with each msg author data", func(t *testing.T) {
		user1, _ := testutil.MockUser(t, ctx.Database.Client)
		room := testutil.MockRoom(t, ctx.Database.Client, user1.ID)
		msgsOfUser1, _ := testutil.MockMessages(t, ctx.Database.Client, 10, user1.ID, room)

		user2, _ := testutil.MockUser(t, ctx.Database.Client)
		msgsOfUser2, _ := testutil.MockMessages(t, ctx.Database.Client, 10, user2.ID, room)

		r := httptest.NewRequest(http.MethodGet, "/rooms/"+room.ID.String(), nil)
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, room.Status.UserID))
		r.Header.Set("Content-Type", "application/json")
		r.SetPathValue("id", room.ID.String())

		ctx.GETRoomMessages(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		var resBody []map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
			t.Errorf("Wrong response format should be json: %v", err)
		}

		// TODO : check for maximum messages returned
		if len(resBody) != len(msgsOfUser1)+len(msgsOfUser2) {
			t.Errorf("Expected %d messages, got %d", len(msgsOfUser1), len(resBody))
		}

		// messages are in reversed order (newest first)
		for i, msg := range resBody {
			if i < len(msgsOfUser2) {
				msgIndex := len(msgsOfUser2) - 1 - i
				testutil.AssertInterface(t, map[string]interface{}{
					"id":      msgsOfUser2[msgIndex].ID.String(),
					"content": msgsOfUser2[msgIndex].Content,
					"room_id": msgsOfUser2[msgIndex].RoomID.String(),
					"author": map[string]interface{}{
						"id":       user2.ID.String(),
						"username": user2.Username,
						"email":    nil,
						"password": nil,
					},
					"author_id": msgsOfUser2[msgIndex].AuthorID.String(),
					"server_id": msgsOfUser2[msgIndex].ServerID.String(),
				}, msg)
			} else {
				j := i - len(msgsOfUser2)
				msgIndex := len(msgsOfUser1) - 1 - j
				testutil.AssertInterface(t, map[string]interface{}{
					"id":      msgsOfUser1[msgIndex].ID.String(),
					"content": msgsOfUser1[msgIndex].Content,
					"room_id": msgsOfUser1[msgIndex].RoomID.String(),
					"author": map[string]interface{}{
						"id":       user1.ID.String(),
						"username": user1.Username,
						"email":    nil,
						"password": nil,
					},
					"author_id": msgsOfUser1[msgIndex].AuthorID.String(),
					"server_id": msgsOfUser1[msgIndex].ServerID.String(),
				}, msg)
			}
		}
	})

	// TODO: add enums for the point
	// t.Run("Should return error if room not found", func(t *testing.T) {
	// 	_, _, user := testutil.MockRoomsServer(t, ctx.Database.Client)

	// 	r := httptest.NewRequest(http.MethodGet, "/rooms/InvalidID", nil)
	// 	w := httptest.NewRecorder()
	// 	r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID.String()))
	// 	r.SetPathValue("id", "InvalidID")

	// 	ctx.GETRoomMessages(w, r)

	// 	if w.Code != http.StatusNotFound {
	// 		t.Errorf("Expected status code 404, got %d", w.Code)
	// 	}

	// 	var resBody map[string]interface{}
	// 	if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
	// 		t.Errorf("Wrong response format should be json: %v", err)
	// 	}

	// 	testutil.AssertInterface(t, map[string]interface{}{
	// 		"error": "Room not found",
	// 	}, resBody)
	// })
}

// for websocket testing
func mockWSRoomHandler(t *testing.T, ctx handlers.ServerContext, userID uuid.UUID, roomid string) http.HandlerFunc {
	t.Helper()
	topicStore := core.NewTopicStore()
	handler := ctx.WSRoom(topicStore)

	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, userID))
		r.SetPathValue("id", roomid)
		handler(w, r)
	}
}

func TestGETRoom(t *testing.T) {
	ctx := testutil.NewTestingContext(t)

	t.Run("Should join a room", func(t *testing.T) {
		user, _ := testutil.MockUser(t, ctx.Database.Client)
		room := testutil.MockRoom(t, ctx.Database.Client, user.ID)
		handler := mockWSRoomHandler(t, ctx, user.ID, room.ID.String())
		s := httptest.NewServer(handler)
		defer s.Close()

		u := "ws" + strings.TrimPrefix(s.URL, "http")

		conn, _, err := websocket.DefaultDialer.Dial(u, nil)
		if err != nil {
			t.Fatalf("err: %v || %v", err, conn)
		}
		defer conn.Close()

		ids := make([]uuid.UUID, 0)

		for i := 0; i < 4; i++ {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(`{"content":"hello"}`)); err != nil {
				t.Fatalf("%v", err)
			}
			var resBody map[string]interface{}
			_, p, err := conn.ReadMessage()
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			json.Unmarshal(p, &resBody)

			testutil.AssertInterface(t, map[string]interface{}{
				"authorId": user.ID.String(),
				"roomId":   room.ID.String(),
				"serverId": room.ServerID.String(),
				"content":  "hello",
				"author": map[string]interface{}{
					"id":       user.ID.String(),
					"username": user.Username,
					"email":    nil,
					"password": nil,
				},
			}, resBody)

			uuidId, err := uuid.Parse(resBody["id"].(string))
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			ids = append(ids, uuidId)
		}

		t.Cleanup(func() {
			for _, id := range ids {
				models.NewMessage().
					WithID(id).
					Delete(ctx.Database.Client)
			}
		})
	})
}

func TestPatchRoomStatus(t *testing.T) {
	ctx := testutil.NewTestingContext(t)

	t.Run("Should update room status with lastReadMsgId", func(t *testing.T) {
		user, _ := testutil.MockUser(t, ctx.Database.Client)
		room := testutil.MockRoom(t, ctx.Database.Client, user.ID)
		msgs, _ := testutil.MockMessages(t, ctx.Database.Client, 1, user.ID, room)
		msg := msgs[0]

		data := []byte(`{"lastReadMsgId":"` + msg.ID.String() + `"}`)

		r := httptest.NewRequest(http.MethodPatch, "/rooms/"+room.ID.String()+"/status", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID))
		r.SetPathValue("id", room.ID.String())

		ctx.PatchRoomStatus(w, r)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status code %d, got %d", http.StatusNoContent, w.Code)
		}

		status := models.NewRoomUserStatus().WithRoomID(room.ID).WithUserID(user.ID)
		err := status.Find(ctx.Database.Client)
		if err != nil {
			t.Fatalf("Error refreshing room status: %v", err)
		}

		statusMap := map[string]interface{}{
			"roomId":        status.RoomID.String(),
			"userId":        status.UserID.String(),
			"lastReadMsgId": status.LastReadMsgID.String(),
		}

		testutil.AssertInterface(t, map[string]interface{}{
			"roomId":        room.ID.String(),
			"userId":        user.ID.String(),
			"lastReadMsgId": msg.ID.String(),
		}, statusMap)
	})
}
