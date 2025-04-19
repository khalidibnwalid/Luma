package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/khalidibnwalid/Luma/core"
	"github.com/khalidibnwalid/Luma/handlers"
	"github.com/khalidibnwalid/Luma/middlewares"
	"github.com/khalidibnwalid/Luma/models"
	"github.com/khalidibnwalid/Luma/testutil"
)

func TestGetUser(t *testing.T) {
	ctx := testutil.NewTestingContext(t)
	t.Run("Should return user data", func(t *testing.T) {
		user, _ := testutil.MockUser(t, ctx.Db)

		r := httptest.NewRequest(http.MethodGet, "/users", nil)
		w := httptest.NewRecorder()

		// mock handler ctx
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, user.ID.Hex()))

		ctx.GetUser(w, r)

		var resBody map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
			t.Errorf("Wrong response format should be json: %v", err)
		}

		if resBody["username"] != user.Username {
			t.Errorf("Expected username %s, got %s", user.Username, resBody["username"])
		}

		if resBody["email"] != user.Email {
			t.Errorf("Expected email %s, got %s", user.Email, resBody["email"])
		}

		if resBody["id"] != user.ID.Hex() {
			t.Errorf("Expected id %s, got %s", user.ID.Hex(), resBody["id"])
		}

		if resBody["password"] != nil {
			t.Errorf("Expected password to be empty, got %s", resBody["password"])
		}
	})

	t.Run("Should return error if user not found", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/users", nil)
		w := httptest.NewRecorder()

		// mock handler ctx with invalid user ID
		r = r.WithContext(context.WithValue(r.Context(), middlewares.CtxUserIDKey, "invalidID"))

		ctx.GetUser(w, r)

		if w.Code != http.StatusNotFound && w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code 404 or 500, got %d", w.Code)
		}
	})
}

func TestPostUser(t *testing.T) {
	ctx := testutil.NewTestingContext(t)

	t.Run("Should create a new user and return user data and a session cookie", func(t *testing.T) {
		username, _ := core.GenerateRandomString(10)
		data := []byte(`{"username": "` + username + `", "password": "testpassword", "email": "` + username + `@example.com"}`)
		defer models.NewUser().WithUsername(username).Delete(ctx.Db, context.Background())

		r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/json")

		ctx.PostUser(w, r)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code 201, got %d", w.Code)
		}

		// Verify cookie is set
		cookies := w.Header().Values("Set-Cookie")
		if len(cookies) == 0 {
			t.Error("Expected session cookie to be set")
		}

		// Verify response body contains user data
		var resBody map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
			t.Errorf("Wrong response format should be json: %v", err)
		}

		t.Log("Response body: ", resBody)

		if resBody["username"] != username {
			t.Errorf("Expected username testuser, got %s", resBody["username"])
		}

		if resBody["email"] != username+"@example.com" {
			t.Errorf("Expected email got %s", resBody["email"])
		}

		if resBody["id"] == nil {
			t.Error("Expected id to be set")
		}

		if resBody["password"] != nil {
			t.Errorf("Expected password to be empty, got %s", resBody["password"])
		}
	})

	t.Run("Should return error with empty username", func(t *testing.T) {
		data := []byte(`{"username": "", "password": "testpassword", "email": "example@example.com"}`)
		r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/json")

		ctx.PostUser(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 400, got %d", w.Code)
		}

		reqBody := map[string]interface{}{}
		if err := json.Unmarshal(w.Body.Bytes(), &reqBody); err != nil {
			t.Errorf("Wrong response format should be json: %v", err)
		}

		if reqBody["error"] != handlers.EnumUsernameRequired {
			t.Errorf("Expected error %s, got %s", handlers.EnumUsernameRequired, reqBody["error"])
		}
	})

	t.Run("Should return error with empty password", func(t *testing.T) {
		data := []byte(`{"username": "testuser", "password": "", "email": "example@example.com"}`)
		r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/json")

		ctx.PostUser(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 400, got %d", w.Code)
		}

		reqBody := map[string]interface{}{}
		if err := json.Unmarshal(w.Body.Bytes(), &reqBody); err != nil {
			t.Errorf("Wrong response format should be json: %v", err)
		}

		if reqBody["error"] != handlers.EnumPasswordRequired {
			t.Errorf("Expected error %s, got %s", handlers.EnumPasswordRequired, reqBody["error"])
		}
	})

	t.Run("Should return error with empty email", func(t *testing.T) {
		data := []byte(`{"username": "testuser", "password": "testpassword", "email": ""}`)
		r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/json")

		ctx.PostUser(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 400, got %d", w.Code)
		}

		reqBody := map[string]interface{}{}
		if err := json.Unmarshal(w.Body.Bytes(), &reqBody); err != nil {
			t.Errorf("Wrong response format should be json: %v", err)
		}

		if reqBody["error"] != handlers.EnumEmailRequired {
			t.Errorf("Expected error %s, got %s", handlers.EnumEmailRequired, reqBody["error"])
		}
	})

	t.Run("Should return error with invalid email", func(t *testing.T) {
		data := []byte(`{"username": "testuser", "password": "testpassword", "email": "invalid-email"}`)
		r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/json")

		ctx.PostUser(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 400, got %d", w.Code)
		}

		reqBody := map[string]interface{}{}
		if err := json.Unmarshal(w.Body.Bytes(), &reqBody); err != nil {
			t.Errorf("Wrong response format should be json: %v", err)
		}

		if reqBody["error"] != handlers.EnumEmailInvalid {
			t.Errorf("Expected error %s, got %s", handlers.EnumEmailInvalid, reqBody["error"])
		}
	})

	t.Run("Should return error with duplicate username", func(t *testing.T) {
		// First create a user
		user, _ := testutil.MockUser(t, ctx.Db)

		// Try to create another user with the same username
		data := []byte(`{"username": "` + user.Username + `", "password": "testpassword", "email": "another@example.com"}`)
		r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/json")

		ctx.PostUser(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 400, got %d", w.Code)
		}

		reqBody := map[string]interface{}{}
		if err := json.Unmarshal(w.Body.Bytes(), &reqBody); err != nil {
			t.Errorf("Wrong response format should be json: %v", err)
		}

		if reqBody["error"] != handlers.EnumUsernameExists {
			t.Errorf("Expected error %s, got %s", handlers.EnumUsernameExists, reqBody["error"])
		}
	})

	t.Run("Should return error with invalid JSON", func(t *testing.T) {
		data := []byte(`{invalid json}`)
		r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/json")

		ctx.PostUser(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 400, got %d", w.Code)
		}
	})
}
