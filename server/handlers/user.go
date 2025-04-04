package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/khalidibnwalid/Luma/core"
	"github.com/khalidibnwalid/Luma/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// TODO add a validator
// Signup Handler / Create a new user
func (s *HandlerContext) PostUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		newErrorResponse(w, http.StatusBadRequest, enumBadRequest)
		return
	}

	if req.Username == "" {
		newErrorResponse(w, http.StatusBadRequest, enumUsernameRequired)
		return
	}

	if req.Password == "" {
		newErrorResponse(w, http.StatusBadRequest, enumPasswordRequired)
		return
	}

	if req.Email == "" {
		newErrorResponse(w, http.StatusBadRequest, enumEmailRequired)
		return
	}
	if !validateEmail(req.Email) {
		newErrorResponse(w, http.StatusBadRequest, enumEmailInvalid)
		return
	}
	// TODO
	// if len(req.Password) < 6 {
	// 	http.Error(w, enumPasswordTooShort, http.StatusBadRequest)
	// 	return
	// }
	user := models.NewUser().WithUsername(req.Username).WithEmail(req.Email).WithPassword(req.Password)

	var err error
	// user exists? if so it won't return an error
	// if it does return an error of mongo.ErrNoDocuments, we assume that the user doesn't exist
	if err = user.FindByUsername(s.Db); err == nil {
		newErrorResponse(w, http.StatusBadRequest, enumUsernameExists)
		return
	} else if err != mongo.ErrNoDocuments {
		// if not a 'Not Found' error, then it is a real error
		newErrorResponse(w, http.StatusInternalServerError, enumInternalServerError)
		return
	}

	if err := user.Create(s.Db); err != nil {
		newErrorResponse(w, http.StatusInternalServerError, enumInternalServerError)
		log.Printf("Error creating user: %v", err)
		return
	}

	// create a token for the user
	token, err := core.GenerateJwtToken(s.JwtSecret, user.ID.Hex())
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, enumInternalServerError)
		return
	}

	json, _ := json.Marshal(user)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", core.SerializeCookieWithToken(token))
	w.WriteHeader(http.StatusCreated)
	w.Write(json)
}
