package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/khalidibnwalid/Luma/handlers"
	"github.com/khalidibnwalid/Luma/middlewares"
	"github.com/khalidibnwalid/Luma/models"
	"github.com/khalidibnwalid/Luma/testutil"
)

func TestPostRoomsServer(t *testing.T) {
	ctx := testutil.NewTestingContext(t)

	t.Run("Should create a new server and returns its data with user status", func(t *testing.T) {
		user, _ := testutil.MockUser(t, ctx.Database.Client)

		serverName := "New Test Server"
		data := []byte(`{"name": "` + serverName + `"}`)

		r := httptest.NewRequest(http.MethodPost, "/servers", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/json")
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID))

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
			"ownerId": user.ID.String(),
			"status": map[string]interface{}{
				"userId":   user.ID.String(),
				"serverId": resBody["id"],
				"nickname": "",
			},
		}, resBody)

		if resBody["id"] == nil {
			t.Error("Expected server ID to be set")
		}

	})

	t.Run("Should return error with empty name", func(t *testing.T) {
		user, _ := testutil.MockUser(t, ctx.Database.Client)

		data := []byte(`{"name": ""}`)

		r := httptest.NewRequest(http.MethodPost, "/servers", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/json")

		// Add user ID to context
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID))

		ctx.PostRoomsServer(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 400, got %d", w.Code)
		}
	})

	t.Run("Should return error with invalid JSON", func(t *testing.T) {
		user, _ := testutil.MockUser(t, ctx.Database.Client)

		data := []byte(`{invalid json}`)

		r := httptest.NewRequest(http.MethodPost, "/servers", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID))
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
		server, _, user := testutil.MockRoomsServer(t, ctx.Database.Client)
		var rooms [2]*models.RoomWithStatus
		rooms[0] = testutil.MockRoom(t, ctx.Database.Client, user.ID, server)
		rooms[1] = testutil.MockRoom(t, ctx.Database.Client, user.ID, server)

		r := httptest.NewRequest(http.MethodGet, "/servers/"+server.ID.String()+"/rooms", nil)
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID))
		r.SetPathValue("id", server.ID.String())

		ctx.GetRoomsOfServer(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}
		t.Log(w.Body.String())

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
					"id":        originalRoom.ID.String(),
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
							"userId":   user.ID,
							"serverId": server.ID.String(),
							"roomId":   originalRoom.ID.String(),
						}, status)
					}
				}
			}
		}

	})

	t.Run("Should return error if server id is invalid", func(t *testing.T) {
		_, _, user := testutil.MockRoomsServer(t, ctx.Database.Client)

		r := httptest.NewRequest(http.MethodGet, "/servers/InvalidID/rooms", nil)
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID))
		r.SetPathValue("id", "InvalidID")

		ctx.GetRoomsOfServer(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 400, got %d", w.Code)
		}

		var resBody map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
			t.Errorf("Wrong response format should be json: %v", err)
		}

		testutil.AssertInterface(t, map[string]interface{}{
			"error": handlers.EnumServerIdInvalid,
		}, resBody)
	})
}

func TestGetUserRoomsServer(t *testing.T) {
	ctx := testutil.NewTestingContext(t)

	t.Run("Should return all servers user has joined", func(t *testing.T) {
		user, _ := testutil.MockUser(t, ctx.Database.Client)

		// Create multiple servers for the user
		server1, _, _ := testutil.MockRoomsServer(t, ctx.Database.Client, user)
		server2, _, _ := testutil.MockRoomsServer(t, ctx.Database.Client, user)

		r := httptest.NewRequest(http.MethodGet, "/user/servers", nil)
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID))

		ctx.GetUserRoomsServer(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}

		var resBody []map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
			t.Errorf("Wrong response format should be json: %v", err)
		}

		if len(resBody) < 2 {
			t.Errorf("Expected at least 2 servers in response, got %d", len(resBody))
		}

		// Check if both created servers are in the response
		foundServer1 := false
		foundServer2 := false

		for _, server := range resBody {
			if server["id"] == server1.ID.String() {
				foundServer1 = true
				testutil.AssertInterface(t, map[string]interface{}{
					"id":      server1.ID.String(),
					"name":    server1.Name,
					"ownerId": server1.OwnerID.String(),
				}, server)

				// Check if status is included
				status, ok := server["status"].(map[string]interface{})
				if !ok || status == nil {
					t.Error("Expected status to be included in server1 response")
				} else {
					testutil.AssertInterface(t, map[string]interface{}{
						"userId":   user.ID.String(),
						"serverId": server1.ID.String(),
					}, status)
				}
			}

			if server["id"] == server2.ID.String() {
				foundServer2 = true
				testutil.AssertInterface(t, map[string]interface{}{
					"id":      server2.ID.String(),
					"name":    server2.Name,
					"ownerId": server2.OwnerID.String(),
				}, server)

				// Check if status is included
				status, ok := server["status"].(map[string]interface{})
				if !ok || status == nil {
					t.Error("Expected status to be included in server2 response")
				} else {
					testutil.AssertInterface(t, map[string]interface{}{
						"userId":   user.ID.String(),
						"serverId": server2.ID.String(),
					}, status)
				}
			}
		}

		if !foundServer1 {
			t.Error("Server 1 not found in response")
		}

		if !foundServer2 {
			t.Error("Server 2 not found in response")
		}
	})

	t.Run("Should return empty array when user has no servers", func(t *testing.T) {
		// Create a new user with no servers
		user, _ := testutil.MockUser(t, ctx.Database.Client)

		r := httptest.NewRequest(http.MethodGet, "/user/servers", nil)
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID))

		ctx.GetUserRoomsServer(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}

		var resBody []map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
			t.Errorf("Wrong response format should be json: %v", err)
		}

		if len(resBody) != 0 {
			t.Errorf("Expected empty array, got array with %d items", len(resBody))
		}
	})
}

func TestPostRoomToServer(t *testing.T) {
	ctx := testutil.NewTestingContext(t)

	t.Run("Should create a new room in server and return the room with user status", func(t *testing.T) {
		server, _, user := testutil.MockRoomsServer(t, ctx.Database.Client)

		roomName := "new-text-channel"
		roomType := "text"
		groupName := "Text Channels"

		data := []byte(`{"type": "` + roomType + `", "name": "` + roomName + `", "groupName": "` + groupName + `"}`)

		r := httptest.NewRequest(http.MethodPost, "/servers/"+server.ID.String()+"/rooms", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/json")
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID))
		r.SetPathValue("id", server.ID.String())

		ctx.PostRoomToServer(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}

		var resBody map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
			t.Errorf("Wrong response format should be json: %v", err)
		}

		testutil.AssertInterface(t, map[string]interface{}{
			"serverId":  server.ID.String(),
			"name":      roomName,
			"type":      roomType,
			"groupName": groupName,
			"status": map[string]interface{}{
				"userId":   user.ID.String(),
				"serverId": server.ID.String(),
				"roomId":   resBody["id"],
			},
		}, resBody)

		// Check if status is included
		status, ok := resBody["status"].(map[string]interface{})
		if !ok || status == nil {
			t.Error("Expected status to be included in response")
		}
	})

	t.Run("Should return error with empty name", func(t *testing.T) {
		server, _, user := testutil.MockRoomsServer(t, ctx.Database.Client)

		data := []byte(`{"type": "type", "name": ""}`)

		r := httptest.NewRequest(http.MethodPost, "/servers/"+server.ID.String()+"/rooms", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/json")
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID))
		r.SetPathValue("id", server.ID.String())

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
		server, _, user := testutil.MockRoomsServer(t, ctx.Database.Client)

		data := []byte(`{"type": "", "name": "test-room", "groupName": "Text Channels"}`)

		r := httptest.NewRequest(http.MethodPost, "/servers/"+server.ID.String()+"/rooms", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/json")
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID))
		r.SetPathValue("id", server.ID.String())

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

	t.Run("Should return error if server id is invalid", func(t *testing.T) {
		_, _, user := testutil.MockRoomsServer(t, ctx.Database.Client)

		r := httptest.NewRequest(http.MethodPost, "/servers/InvalidID/rooms", nil)
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID))
		r.SetPathValue("id", "InvalidID")

		ctx.GetRoomsOfServer(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 400, got %d", w.Code)
		}

		var resBody map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
			t.Errorf("Wrong response format should be json: %v", err)
		}

		testutil.AssertInterface(t, map[string]interface{}{
			"error": handlers.EnumServerIdInvalid,
		}, resBody)
	})
}

func TestJoinServer(t *testing.T) {
	ctx := testutil.NewTestingContext(t)

	t.Run("Should join a server and return the server with user status", func(t *testing.T) {
		server, _, user := testutil.MockRoomsServer(t, ctx.Database.Client)

		r := httptest.NewRequest(http.MethodGet, "/servers/"+server.ID.String(), nil)
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID))
		r.SetPathValue("id", server.ID.String())

		ctx.JoinServer(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}

		var resBody map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
			t.Errorf("Wrong response format should be json: %v", err)
		}

		testutil.AssertInterface(t, map[string]interface{}{
			"id":      server.ID.String(),
			"name":    server.Name,
			"ownerId": server.OwnerID.String(),
			"status": map[string]interface{}{
				"userId":   user.ID.String(),
				"serverId": server.ID.String(),
				"nickname": "",
			},
		}, resBody)

	})

	t.Run("Should return error if server id is invalid", func(t *testing.T) {
		_, _, user := testutil.MockRoomsServer(t, ctx.Database.Client)

		r := httptest.NewRequest(http.MethodPost, "/servers/invalidID", nil)
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID))
		r.SetPathValue("id", "invalidID")

		ctx.JoinServer(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 400, got %d", w.Code)
		}

		var resBody map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
			t.Errorf("Wrong response format should be json: %v", err)
		}

		testutil.AssertInterface(t, map[string]interface{}{
			"error": handlers.EnumServerIdInvalid,
		}, resBody)
	})

	t.Run("Should return error if server not found", func(t *testing.T) {
		_, _, user := testutil.MockRoomsServer(t, ctx.Database.Client)

		mockUUID, _ := uuid.NewRandom()
		r := httptest.NewRequest(http.MethodPost, "/servers/"+mockUUID.String(), nil)
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID))
		r.SetPathValue("id", mockUUID.String())

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
