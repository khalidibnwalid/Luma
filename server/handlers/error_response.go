package handlers

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// newErrorResponse creates a new error response with the given code and error message.
// It sets the Content-Type header to application/json and writes the error message to the response body in a standard format.
func newErrorResponse(w http.ResponseWriter, code int, err string, msg ...string) {
	w.Header().Set("Content-Type", "application/json")

	res := errorResponse{
		Error: err,
	}

	if len(msg) > 0 {
		res.Message = msg[0]
	}

	j, _ := json.Marshal(res)
	w.WriteHeader(code)
	w.Write(j)
}

// user registration error codes
const (
	EnumUsernameExists   = "USERNAME_EXISTS"
	EnumUsernameRequired = "USERNAME_REQUIRED"
	EnumUsernameInvalid  = "USERNAME_INVALID"
	EnumUserDoesNotExist = "USER_DOES_NOT_EXIST"
	EnumPasswordRequired = "PASSWORD_REQUIRED"
	EnumPasswordInvalid  = "PASSWORD_INVALID"
	EnumEmailRequired    = "EMAIL_REQUIRED"
	EnumEmailInvalid     = "EMAIL_INVALID"
)

const (
	EnumInternalServerError = "INTERNAL_SERVER_ERROR"
	EnumNotFound            = "NOT_FOUND"
	EnumBadRequest          = "BAD_REQUEST"
	EnumUnauthorized        = "UNAUTHORIZED"
)

const (
	EnumServerIdRequired = "SERVER_ID_REQUIRED"
	EnumServerIdInvalid  = "SERVER_ID_INVALID"
	EnumServerNotFound   = "SERVER_NOT_FOUND"
	EnumServerTypeRequired = "SERVER_TYPE_REQUIRED"
	EnumServerNameRequired = "SERVER_NAME_REQUIRED"
	
)
