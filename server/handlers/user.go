package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/khalidibnwalid/Luma/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (s *HandlerContext) UserGET(w http.ResponseWriter, req *http.Request) {
	user := models.NewUser().WithUsername(req.PathValue("username"))
	if err := user.FindByUsername(s.Db); err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonUser, err := json.Marshal(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonUser)
}
