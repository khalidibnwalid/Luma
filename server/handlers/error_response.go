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
func newErrorResponse(w http.ResponseWriter, code int ,err string ,msg ...string) {
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
	enumUsernameExists   = "USERNAME_EXISTS"
	enumUsernameRequired = "USERNAME_REQUIRED"
	enumUsernameInvalid  = "USERNAME_INVALID"
	enumUserDoesNotExist = "USER_DOES_NOT_EXIST"
	enumPasswordRequired = "PASSWORD_REQUIRED"
	enumPasswordInvalid  = "PASSWORD_INVALID"
	enumEmailRequired    = "EMAIL_REQUIRED"
	enumEmailInvalid     = "EMAIL_INVALID"
)

const (
	enumInternalServerError = "INTERNAL_SERVER_ERROR"
	enumNotFound            = "NOT_FOUND"
	enumBadRequest          = "BAD_REQUEST"
	enumUnauthorized        = "UNAUTHORIZED"
)