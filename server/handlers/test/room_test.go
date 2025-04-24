package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/khalidibnwalid/Luma/core"
	"github.com/khalidibnwalid/Luma/handlers"
	"github.com/khalidibnwalid/Luma/middlewares"
	"github.com/khalidibnwalid/Luma/models"
	"github.com/khalidibnwalid/Luma/testutil"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestGETRoomMessages(t *testing.T) {
	ctx := testutil.NewTestingContext(t)

	// TODO : check for maximum messages returned
	t.Run("Should get messages of a room with each msg author data", func(t *testing.T) {
		user1, _ := testutil.MockUser(t, ctx.Db)
		room := testutil.MockRoom(t, ctx.Db, user1.ID.Hex())
		msgsOfUser1, _ := testutil.MockMessages(t, ctx.Db, 12, user1.ID.Hex(), room)

		user2, _ := testutil.MockUser(t, ctx.Db)
		msgsOfUser2, _ := testutil.MockMessages(t, ctx.Db, 10, user2.ID.Hex(), room)

		r := httptest.NewRequest(http.MethodGet, "/rooms/"+room.ID.Hex(), nil)
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, room.Status.UserID))
		r.Header.Set("Content-Type", "application/json")
		r.SetPathValue("id", room.ID.Hex())

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

		for i, msg := range resBody {

			if i < len(msgsOfUser1) {
				testutil.AssertInterface(t, map[string]interface{}{
					"content": msgsOfUser1[i].Content,
					"room_id": msgsOfUser1[i].RoomID,
					"author": map[string]interface{}{
						"id":       user1.ID.Hex(),
						"username": user1.Username,
						"email":    nil,
						"password": nil,
					},
					"author_id": msgsOfUser1[i].AuthorID,
					"server_id": msgsOfUser1[i].ServerID,
				}, msg)
			} else {
				j := i - len(msgsOfUser1)
				testutil.AssertInterface(t, map[string]interface{}{
					"content": msgsOfUser2[j].Content,
					"room_id": msgsOfUser2[j].RoomID,
					"author": map[string]interface{}{
						"id":       user2.ID.Hex(),
						"username": user2.Username,
						"email":    nil,
						"password": nil,
					},
					"author_id": msgsOfUser2[j].AuthorID,
					"server_id": msgsOfUser2[j].ServerID,
				}, msg)
			}
		}
	})

	// TODO: add enums for the point
	// t.Run("Should return error if room not found", func(t *testing.T) {
	// 	_, _, user := testutil.MockRoomsServer(t, ctx.Db)

	// 	r := httptest.NewRequest(http.MethodGet, "/rooms/InvalidID", nil)
	// 	w := httptest.NewRecorder()
	// 	r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID.Hex()))
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
func mockWSRoomHandler(t *testing.T, ctx handlers.ServerContext, userID, roomid string) http.HandlerFunc {
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
		user1, _ := testutil.MockUser(t, ctx.Db)
		room := testutil.MockRoom(t, ctx.Db, user1.ID.Hex())
		handler := mockWSRoomHandler(t, ctx, user1.ID.Hex(), room.ID.Hex())
		s := httptest.NewServer(handler)
		defer s.Close()

		u := "ws" + strings.TrimPrefix(s.URL, "http")

		conn, _, err := websocket.DefaultDialer.Dial(u, nil)
		if err != nil {
			t.Fatalf("err: %v || %v", err, conn)
		}
		defer conn.Close()

		ids := make([]string, 0)

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
				"authorId": user1.ID.Hex(),
				"roomId":   room.ID.Hex(),
				"serverId": room.ServerID,
				"content":  "hello",
				"author": map[string]interface{}{
					"id":       user1.ID.Hex(),
					"username": user1.Username,
					"email":    nil,
					"password": nil,
				},
			}, resBody)

			ids = append(ids, resBody["id"].(string))
		}
		t.Cleanup(func() {
			for _, id := range ids {
				objId, _ := bson.ObjectIDFromHex(id)
				msg := &models.Message{
					ID: objId,
				}
				msg.Delete(ctx.Db, context.Background())
			}
		})
	})
}
