package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/khalidibnwalid/Luma/core"
	"github.com/khalidibnwalid/Luma/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Signup Handler / Create a new user
func (s *HandlerContext) PostUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// TODO: needs a validator for username and password
	user := models.NewUser().WithUsername(req.Username)

	var err error
	// user exists? if so it won't return an error
	// if it does return an error of mongo.ErrNoDocuments, we assume that the user doesn't exist
	if err = user.FindByUsername(s.Db); err == nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// if not a 'Not Found' error, then it is a real error
	if err != mongo.ErrNoDocuments {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	user.WithPassword(req.Password)
	if err := user.Create(s.Db); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// create a token for the user
	token, err := core.GenerateJwtToken(s.JwtSecret, user.ID.Hex())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	json, _ := json.Marshal(user)

	w.Header().Add("Set-Cookie", core.SerializeCookieWithToken(token))
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
