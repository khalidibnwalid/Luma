package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/khalidibnwalid/Luma/middlewares"
	"github.com/khalidibnwalid/Luma/models"
)

func (ctx *HandlerContext) validateRoomsServerID(w http.ResponseWriter, r *http.Request) (models.RoomsServer, error) {
	serverID := r.PathValue("id")
	if serverID == "" {
		http.Error(w, "Server ID is required", http.StatusBadRequest)
		return models.RoomsServer{}, nil
	}

	serverData := models.RoomsServer{}
	if err := serverData.FindById(ctx.Db, serverID); err != nil {
		http.Error(w, "Server not found", http.StatusNotFound)
		return models.RoomsServer{}, nil
	}

	return serverData, nil
}

func (ctx *HandlerContext) GetRoomsServer(w http.ResponseWriter, r *http.Request) {
	server, err := ctx.validateRoomsServerID(w, r)
	if err != nil {
		return
	}

	json, _ := json.Marshal(server)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

}

// get all servers of a user
func (ctx *HandlerContext) GetUserRoomsServer(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.CtxUserIDKey).(string)

	// temp to get all owned servers, will be changed to get all servers of a user
	servers, err := models.NewRoomsServer().WithOwnerID(userID).GetAllServersOfOwner(ctx.Db)
	if err != nil {
		log.Print(err)
		http.Error(w, "Error getting servers", http.StatusInternalServerError)
		return
	}

	json, _ := json.Marshal(servers)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (ctx *HandlerContext) PostRoomsServer(w http.ResponseWriter, r *http.Request) {
	var t struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if t.Name == "" {
		http.Error(w, "'name' is required", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(middlewares.CtxUserIDKey).(string)

	server := models.NewRoomsServer().WithOwnerID(userID)
	server.Name = t.Name

	if err := server.Create(ctx.Db); err != nil {
		http.Error(w, "Error creating server", http.StatusInternalServerError)
		return
	}

	json, _ := json.Marshal(server)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (ctx *HandlerContext) GetRoomsOfServer(w http.ResponseWriter, r *http.Request) {
	server, err := ctx.validateRoomsServerID(w, r)
	if err != nil {
		return
	}

	rooms, err := server.GetRooms(ctx.Db)
	if err != nil {
		http.Error(w, "Error getting rooms", http.StatusInternalServerError)
		return
	}

	json, _ := json.Marshal(rooms)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (ctx *HandlerContext) PostRoomToServer(w http.ResponseWriter, r *http.Request) {
	server, err := ctx.validateRoomsServerID(w, r)
	if err != nil {
		return
	}

	var t struct {
		Type      string `json:"type"`
		Name      string `json:"name"`
		GroupName string `json:"groupName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if t.Type == "" {
		http.Error(w, "'type' is required", http.StatusBadRequest)
		return
	}

	if t.Name == "" {
		http.Error(w, "'name' is required", http.StatusBadRequest)
		return
	}

	room := models.Room{
		ServerID:  server.ID,
		Type:      t.Type,
		Name:      t.Name,
		GroupName: t.GroupName,
	}

	if err := room.Create(ctx.Db); err != nil {
		http.Error(w, "Error creating room", http.StatusInternalServerError)
		return
	}

	json, _ := json.Marshal(room)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
