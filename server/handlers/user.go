package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/khalidibnwalid/Luma/core"
	"github.com/khalidibnwalid/Luma/middlewares"
	"github.com/khalidibnwalid/Luma/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (s *ServerContext) GetUser(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()
	user := models.NewUser()
	userID := rCtx.Value(middlewares.CtxUserIDKey).(string)
	user.WithHexID(userID)

	if err := user.FindByID(s.Db, rCtx); err != nil {
		if err == mongo.ErrNoDocuments {
			// unlikely to happen, but just in case a user delete their account and use their token again
			newErrorResponse(w, http.StatusNotFound, enumUserDoesNotExist)
			return
		}
		newErrorResponse(w, http.StatusInternalServerError, enumInternalServerError)
		return
	}

	json, _ := json.Marshal(user)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// TODO add a validator
// Signup Handler / Create a new user
func (s *ServerContext) PostUser(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()

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
	if err = user.FindByUsername(s.Db, rCtx); err == nil {
		newErrorResponse(w, http.StatusBadRequest, enumUsernameExists)
		return
	} else if err != mongo.ErrNoDocuments {
		// if not a 'Not Found' error, then it is a real error
		newErrorResponse(w, http.StatusInternalServerError, enumInternalServerError)
		return
	}

	if err := user.Create(s.Db, rCtx); err != nil {
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
