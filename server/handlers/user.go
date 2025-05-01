package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/khalidibnwalid/Luma/core"
	"github.com/khalidibnwalid/Luma/middlewares"
	"github.com/khalidibnwalid/Luma/models"
	"gorm.io/gorm"
)

func (s *ServerContext) GetUser(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()
	user := models.NewUser()
	userID := rCtx.Value(middlewares.CtxUserIDKey).(uuid.UUID)
	user.WithID(userID)

	if err := user.FindByID(s.Database.Client.WithContext(rCtx)); err != nil {
		if err == gorm.ErrRecordNotFound {
			// unlikely to happen, but just in case a user delete their account and use their token again
			newErrorResponse(w, http.StatusNotFound, EnumUserDoesNotExist)
			return
		}
		newErrorResponse(w, http.StatusInternalServerError, EnumInternalServerError)
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
		newErrorResponse(w, http.StatusBadRequest, EnumBadRequest)
		return
	}

	if req.Username == "" {
		newErrorResponse(w, http.StatusBadRequest, EnumUsernameRequired)
		return
	}

	if req.Password == "" {
		newErrorResponse(w, http.StatusBadRequest, EnumPasswordRequired)
		return
	}

	if req.Email == "" {
		newErrorResponse(w, http.StatusBadRequest, EnumEmailRequired)
		return
	}
	if !validateEmail(req.Email) {
		newErrorResponse(w, http.StatusBadRequest, EnumEmailInvalid)
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
	// if it does return an error of gorm.ErrRecordNotFound, we assume that the user doesn't exist
	if err = user.FindByUsername(s.Database.Client.WithContext(rCtx)); err == nil {
		newErrorResponse(w, http.StatusBadRequest, EnumUsernameExists)
		return
	} else if err != gorm.ErrRecordNotFound {
		// if not a 'Not Found' error, then it is a real error
		newErrorResponse(w, http.StatusInternalServerError, EnumInternalServerError)
		return
	}

	if err := user.Create(s.Database.Client); err != nil {
		newErrorResponse(w, http.StatusInternalServerError, EnumInternalServerError)
		log.Printf("Error creating user: %v", err)
		return
	}

	// create a token for the user
	token, err := core.GenerateJwtToken(s.JwtSecret, user.ID.String())
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, EnumInternalServerError)
		return
	}

	json, _ := json.Marshal(user)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", core.SerializeCookieWithToken(token))
	w.WriteHeader(http.StatusCreated)
	w.Write(json)
}
