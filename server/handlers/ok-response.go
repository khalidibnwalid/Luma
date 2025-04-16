package handlers

import (
	"encoding/json"
	"net/http"
)

const (
	enumLoggedOut = "LOGGED_OUT"
)

func newOkResponse(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	res := map[string]string{
		"message": msg,
	}

	j, _ := json.Marshal(res)
	w.Write(j)
}