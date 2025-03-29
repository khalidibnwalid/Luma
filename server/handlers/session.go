package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/khalidibnwalid/Luma/core"
	"github.com/khalidibnwalid/Luma/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// LoginHandler
func (ctx *HandlerContext) PostSession(w http.ResponseWriter, r *http.Request) {
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

	user := models.NewUser().WithUsername(req.Username)
	// check if the user exists
	if err := user.FindByUsername(ctx.Db); err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// check if the password is correct
	if err := core.VerifyHashWithSalt(req.Password, user.HashedPassword); err != nil {
		if err == core.ErrHashVerificationFailed {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	token, err := core.GenerateJwtToken(ctx.JwtSecret, user.ID.Hex())
	if err != nil {
		fmt.Println(err)
	}

	json, _ := json.Marshal(user)
	w.Header().Add("Set-Cookie", core.SerializeCookieWithToken(token))
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
