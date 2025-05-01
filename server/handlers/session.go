package handlers

import (
	"encoding/json"
	"net/http"
	"net/mail"

	"github.com/khalidibnwalid/Luma/core"
	"github.com/khalidibnwalid/Luma/models"
	"gorm.io/gorm"
)

// TODO add a validator
// TODO forget password
// LoginHandler
func (s *ServerContext) PostSession(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()
	var req struct {
		UsernameOrEmail string `json:"username"`
		Password        string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		newErrorResponse(w, http.StatusBadRequest, EnumBadRequest)
		return
	}

	if req.UsernameOrEmail == "" {
		newErrorResponse(w, http.StatusBadRequest, EnumUsernameInvalid)
		return
	}

	if req.Password == "" {
		newErrorResponse(w, http.StatusBadRequest, EnumPasswordInvalid)
		return
	}
	var user *models.User
	// check if the user exists
	if validateEmail(req.UsernameOrEmail) {
		user = models.NewUser().WithEmail(req.UsernameOrEmail)
		if err := user.FindByEmail(s.Database.Client.WithContext(rCtx)); err != nil {
			if err == gorm.ErrRecordNotFound {
				newErrorResponse(w, http.StatusUnauthorized, EnumUserDoesNotExist)
				return
			}
			newErrorResponse(w, http.StatusInternalServerError, EnumInternalServerError)
			return
		}
	} else {
		user = models.NewUser().WithUsername(req.UsernameOrEmail)
		if err := user.FindByUsername(s.Database.Client.WithContext(rCtx)); err != nil {
			if err == gorm.ErrRecordNotFound {
				newErrorResponse(w, http.StatusUnauthorized, EnumUserDoesNotExist)
				return
			}
			newErrorResponse(w, http.StatusInternalServerError, EnumInternalServerError)
			return
		}
	}

	// check if the password is correct
	if err := core.VerifyHashWithSalt(req.Password, user.HashedPassword); err != nil {
		if err == core.ErrHashVerificationFailed {
			newErrorResponse(w, http.StatusUnauthorized, EnumPasswordInvalid)
			return
		}
		newErrorResponse(w, http.StatusInternalServerError, EnumInternalServerError)
		return
	}

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

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Just replace the cookie with an empty token
func (s *ServerContext) DeleteSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", core.SerializeCookieWithToken(""))

	newOkResponse(w, http.StatusOK, EnumLoggedOut)
}
