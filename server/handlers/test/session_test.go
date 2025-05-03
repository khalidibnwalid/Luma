package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/khalidibnwalid/Luma/core"
	"github.com/khalidibnwalid/Luma/handlers"
	"github.com/khalidibnwalid/Luma/testutil"
)

func TestDeleteSession(t *testing.T) {
	r := httptest.NewRequest(http.MethodDelete, "/session", nil)
	w := httptest.NewRecorder()

	ctx := &handlers.ServerContext{
		JwtSecret: "SECRET",
	}

	ctx.DeleteSession(w, r)

	// not that important
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	t.Run("Should clean auth cookies", func(t *testing.T) {
		cookies := w.Result().Cookies()
		some := false
		for _, cookie := range cookies {
			if cookie.Name == core.JwtSessionCookieName {
				if cookie.Value != "" {
					t.Errorf("Expected cookie %s to be empty, got %s", core.JwtSessionCookieName, cookie.Value)
				}
				some = true
				break
			}
		}

		if !some {
			t.Errorf("Expected cookie %s to be set", core.JwtSessionCookieName)
		}
	})

	t.Run("Should respond with the proper format", func(t *testing.T) {
		var resBody map[string]string

		if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
			t.Errorf("Error unmarshalling response: %v", err)
		}

		if resBody["message"] != handlers.EnumLoggedOut {
			t.Errorf("Expected response body %s, got %s", handlers.EnumLoggedOut, resBody["message"])
		}
	})

}

func TestPostSession(t *testing.T) {
	ctx := testutil.NewTestingContext(t)

	// treat "Should login" as "Should return a cookie session + and user data"
	t.Run("Should return a cookie and user data if username and password are correct", func(t *testing.T) {
		user, pass := testutil.MockUser(t, ctx.Database.Client)
		defer user.Delete(ctx.Database.Client)

		data := []byte(`{"username": "` + user.Username + `", "password": "` + pass + `"}`)

		r := httptest.NewRequest(http.MethodPost, "/session", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/json")

		t.Log("User: ", user.Username)
		ctx.PostSession(w, r)

		cookies := w.Result().Cookies()
		some := false
		for _, cookie := range cookies {
			if cookie.Name == core.JwtSessionCookieName {
				if cookie.Value == "" {
					t.Errorf("Expected cookie %s to be set", core.JwtSessionCookieName)
				}
				some = true
				break
			}
		}
		if !some {
			t.Errorf("Expected cookie %s to be set", core.JwtSessionCookieName)
		}

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code 201, got %d", w.Code)
		}

		var resBody map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
			t.Errorf("Wrong response format should be json %v", err)
		}

		testutil.AssertInterface(t, map[string]interface{}{
			"username": user.Username,
			"email":    user.Email,
			"id":       user.ID.String(),
			"password": nil,
		}, resBody)
	})

	t.Run("Should return a cookie and user data if email and password are correct", func(t *testing.T) {
		user, pass := testutil.MockUser(t, ctx.Database.Client)
		defer user.Delete(ctx.Database.Client)

		data := []byte(`{"username": "` + user.Username + `", "password": "` + pass + `"}`)

		r := httptest.NewRequest(http.MethodPost, "/session", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/json")

		ctx.PostSession(w, r)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code 201, got %d", w.Code)
		}

		cookies := w.Result().Cookies()
		some := false
		for _, cookie := range cookies {
			if cookie.Name == core.JwtSessionCookieName {
				if cookie.Value == "" {
					t.Errorf("Expected cookie %s to be set", core.JwtSessionCookieName)
				}
				some = true
				break
			}
		}
		if !some {
			t.Errorf("Expected cookie %s to be set", core.JwtSessionCookieName)
		}

		var resBody map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
			t.Errorf("Wrong response format: %v", err)
		}

		testutil.AssertInterface(t, map[string]interface{}{
			"username": user.Username,
			"email":    user.Email,
			"id":       user.ID.String(),
			"password": nil,
		}, resBody)
	})

	t.Run("Should not login if password is wrong", func(t *testing.T) {
		user, _ := testutil.MockUser(t, ctx.Database.Client)
		defer user.Delete(ctx.Database.Client)

		data := []byte(`{"username": "` + user.Username + `", "password": "wrongpassword"}`)

		r := httptest.NewRequest(http.MethodPost, "/session", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/json")

		ctx.PostSession(w, r)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code 401, got %d", w.Code)
		}

		var resBody map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
			t.Errorf("Wrong response format: %v", err)
		}

		testutil.AssertInterface(t, map[string]interface{}{
			"error": handlers.EnumPasswordInvalid,
		}, resBody)

	})
	t.Run("Should not login if username or email is empty", func(t *testing.T) {
		data := []byte(`{"username": "", "password": "wrongpassword"}`)
		r := httptest.NewRequest(http.MethodPost, "/session", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", "application/json")

		ctx.PostSession(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 400, got %d", w.Code)
		}

		var resBody map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
			t.Errorf("Wrong response format: %v", err)
		}
	})

	t.Run("Should return a 400 if the body is empty", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/session", nil)
		w := httptest.NewRecorder()

		ctx.PostSession(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 400, got %d", w.Code)
		}
	})

}
