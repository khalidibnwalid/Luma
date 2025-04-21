package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/khalidibnwalid/Luma/handlers"
	"github.com/khalidibnwalid/Luma/middlewares"
	"github.com/khalidibnwalid/Luma/models"
	"github.com/khalidibnwalid/Luma/testutil"
)

func TestPostRoomsServer(t *testing.T) {
	ctx := testutil.NewTestingContext(t)

	t.Run("Should create a new server and returns its data with user status", func(t *testing.T) {
		user, _ := testutil.MockUser(t, ctx.Db)

		serverName := "New Test Server"
		data := []byte(`{"name": "` + serverName + `"}`)

		r := httptest.NewRequest(http.MethodPost, "/servers", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/json")
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID.Hex()))

		ctx.PostRoomsServer(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}

		var resBody map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
			t.Errorf("Wrong response format should be json: %v", err)
		}

		testutil.AssertInterface(t, map[string]interface{}{
			"name":    serverName,
			"ownerId": user.ID.Hex(),
			"status": map[string]interface{}{
				"userId":   user.ID.Hex(),
				"serverId": resBody["id"],
				"nickname": "",
			},
		}, resBody)

		if resBody["id"] == nil {
			t.Error("Expected server ID to be set")
		}

	})

	t.Run("Should return error with empty name", func(t *testing.T) {
		user, _ := testutil.MockUser(t, ctx.Db)

		data := []byte(`{"name": ""}`)

		r := httptest.NewRequest(http.MethodPost, "/servers", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/json")

		// Add user ID to context
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID.Hex()))

		ctx.PostRoomsServer(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 400, got %d", w.Code)
		}
	})

	t.Run("Should return error with invalid JSON", func(t *testing.T) {
		user, _ := testutil.MockUser(t, ctx.Db)

		data := []byte(`{invalid json}`)

		r := httptest.NewRequest(http.MethodPost, "/servers", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID.Hex()))
		r.Header.Set("Content-Type", "application/json")

		ctx.PostRoomsServer(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 400, got %d", w.Code)
		}
	})
}

func TestGetRoomsOfServer(t *testing.T) {
	ctx := testutil.NewTestingContext(t)

	t.Run("Should return rooms of server with user status", func(t *testing.T) {
		server, _, user := testutil.MockRoomsServer(t, ctx.Db)
		var rooms [2]*models.RoomWithStatus
		rooms[0] = testutil.MockRoom(t, ctx.Db, user.ID.Hex(), server)
		rooms[1] = testutil.MockRoom(t, ctx.Db, user.ID.Hex(), server)

		r := httptest.NewRequest(http.MethodGet, "/servers/"+server.ID.Hex()+"/rooms", nil)
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID.Hex()))
		r.SetPathValue("id", server.ID.Hex())

		ctx.GetRoomsOfServer(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}

		var resBody []map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
			t.Errorf("Wrong response format should be json: %v", err)
		}

		if len(resBody) != 2 {
			t.Error("Expected to have 2 rooms in the response got:", len(resBody))
		} else {
			for i, room := range resBody {
				originalRoom := rooms[i]

				testutil.AssertInterface(t, map[string]interface{}{
					"id":        originalRoom.ID.Hex(),
					"name":      originalRoom.Name,
					"type":      originalRoom.Type,
					"groupName": originalRoom.GroupName,
				}, room)

				if room["status"] == nil {
					t.Error("Expected status to be included in response")
				} else {
					status, ok := room["status"].(map[string]interface{})
					if !ok || status == nil {
						t.Error("Expected status to be included in response")
					} else {
						testutil.AssertInterface(t, map[string]interface{}{
							"userId":    user.ID.Hex(),
							"serverId":  server.ID.Hex(),
							"roomId":    originalRoom.ID.Hex(),
							"isCleared": originalRoom.Status.IsCleared,
						}, status)
					}
				}
			}
		}

	})

	t.Run("Should return error if server not found", func(t *testing.T) {
		_, _, user := testutil.MockRoomsServer(t, ctx.Db)

		r := httptest.NewRequest(http.MethodGet, "/servers/InvalidID/rooms", nil)
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID.Hex()))
		r.SetPathValue("id", "InvalidID")

		ctx.GetRoomsOfServer(w, r)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code 404, got %d", w.Code)
		}

		var resBody map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
			t.Errorf("Wrong response format should be json: %v", err)
		}

		testutil.AssertInterface(t, map[string]interface{}{
			"error": handlers.EnumServerNotFound,
		}, resBody)
	})
}

func TestPostRoomToServer(t *testing.T) {
	ctx := testutil.NewTestingContext(t)

	t.Run("Should create a new room in server and return the room with user status", func(t *testing.T) {
		server, _, user := testutil.MockRoomsServer(t, ctx.Db)

		roomName := "new-text-channel"
		roomType := "text"
		groupName := "Text Channels"

		data := []byte(`{"type": "` + roomType + `", "name": "` + roomName + `", "groupName": "` + groupName + `"}`)

		r := httptest.NewRequest(http.MethodPost, "/servers/"+server.ID.Hex()+"/rooms", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/json")
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID.Hex()))
		r.SetPathValue("id", server.ID.Hex())

		ctx.PostRoomToServer(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}

		var resBody map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
			t.Errorf("Wrong response format should be json: %v", err)
		}

		testutil.AssertInterface(t, map[string]interface{}{
			"serverId":  server.ID.Hex(),
			"name":      roomName,
			"type":      roomType,
			"groupName": groupName,
			"status": map[string]interface{}{
				"userId":    user.ID.Hex(),
				"serverId":  server.ID.Hex(),
				"roomId":    resBody["id"],
				"isCleared": true,
			},
		}, resBody)

		// Check if status is included
		status, ok := resBody["status"].(map[string]interface{})
		if !ok || status == nil {
			t.Error("Expected status to be included in response")
		}
	})

	t.Run("Should return error with empty name", func(t *testing.T) {
		server, _, user := testutil.MockRoomsServer(t, ctx.Db)

		data := []byte(`{"type": "type", "name": ""}`)

		r := httptest.NewRequest(http.MethodPost, "/servers/"+server.ID.Hex()+"/rooms", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/json")
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID.Hex()))
		r.SetPathValue("id", server.ID.Hex())

		ctx.PostRoomToServer(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 400, got %d", w.Code)
		}

		var resBody map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
			t.Errorf("Wrong response format should be json: %v", err)
		}
		testutil.AssertInterface(t, map[string]interface{}{
			"error": handlers.EnumServerNameRequired,
		}, resBody)
	})

	t.Run("Should return error with empty type", func(t *testing.T) {
		server, _, user := testutil.MockRoomsServer(t, ctx.Db)

		data := []byte(`{"type": "", "name": "test-room", "groupName": "Text Channels"}`)

		r := httptest.NewRequest(http.MethodPost, "/servers/"+server.ID.Hex()+"/rooms", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/json")
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID.Hex()))
		r.SetPathValue("id", server.ID.Hex())

		ctx.PostRoomToServer(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 400, got %d", w.Code)
		}

		var resBody map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
			t.Errorf("Wrong response format should be json: %v", err)
		}
		testutil.AssertInterface(t, map[string]interface{}{
			"error": handlers.EnumServerTypeRequired,
		}, resBody)
	})

	t.Run("Should return error if server not found", func(t *testing.T) {
		_, _, user := testutil.MockRoomsServer(t, ctx.Db)

		r := httptest.NewRequest(http.MethodPost, "/servers/InvalidID/rooms", nil)
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID.Hex()))
		r.SetPathValue("id", "InvalidID")

		ctx.GetRoomsOfServer(w, r)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code 404, got %d", w.Code)
		}

		var resBody map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
			t.Errorf("Wrong response format should be json: %v", err)
		}

		testutil.AssertInterface(t, map[string]interface{}{
			"error": handlers.EnumServerNotFound,
		}, resBody)
	})
}

func TestJoinServer(t *testing.T) {
	ctx := testutil.NewTestingContext(t)

	t.Run("Should join a server and return the server with user status", func(t *testing.T) {
		server, _, user := testutil.MockRoomsServer(t, ctx.Db)

		r := httptest.NewRequest(http.MethodGet, "/servers/"+server.ID.Hex(), nil)
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID.Hex()))
		r.SetPathValue("id", server.ID.Hex())

		ctx.JoinServer(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}

		var resBody map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
			t.Errorf("Wrong response format should be json: %v", err)
		}

		testutil.AssertInterface(t, map[string]interface{}{
			"id":   server.ID.Hex(),
			"name": server.Name,
			"ownerId": server.OwnerID,
			"status": map[string]interface{}{
				"userId":   user.ID.Hex(),
				"serverId": server.ID.Hex(),
				"nickname": "",
			},
		}, resBody)

	})

	t.Run("Should return error if server not found", func(t *testing.T) {
		_, _, user := testutil.MockRoomsServer(t, ctx.Db)

		r := httptest.NewRequest(http.MethodPost, "/servers/InvalidID", nil)
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID.Hex()))
		r.SetPathValue("id", "InvalidID")

		ctx.JoinServer(w, r)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code 404, got %d", w.Code)
		}

		var resBody map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
			t.Errorf("Wrong response format should be json: %v", err)
		}

		testutil.AssertInterface(t, map[string]interface{}{
			"error": handlers.EnumServerNotFound,
		}, resBody)
	})
}
