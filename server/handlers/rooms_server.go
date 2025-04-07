package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/khalidibnwalid/Luma/middlewares"
	"github.com/khalidibnwalid/Luma/models"
)

func (ctx *ServerContext) validateRoomsServerID(w http.ResponseWriter, r *http.Request) (models.RoomsServer, error) {
	serverID := r.PathValue("id")
	if serverID == "" {
		http.Error(w, "Server ID is required", http.StatusBadRequest)
		return models.RoomsServer{}, nil
	}

	serverData := models.RoomsServer{}
	if err := serverData.FindById(ctx.Db, serverID); err != nil {
		http.Error(w, "Server not found", http.StatusNotFound)
		return models.RoomsServer{}, errors.New("Server not found")
	}

	return serverData, nil
}

func (ctx *ServerContext) GetRoomsServer(w http.ResponseWriter, r *http.Request) {
	server, err := ctx.validateRoomsServerID(w, r)
	if err != nil {
		return
	}

	json, _ := json.Marshal(server)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

}

// get all servers of a user
func (ctx *ServerContext) GetUserRoomsServer(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.CtxUserIDKey).(string)

	servers, err := models.NewServerUserStatus().WithUserID(userID).GetServers(ctx.Db)
	if err != nil {
		http.Error(w, "Error getting servers", http.StatusInternalServerError)
		return
	}

	json, _ := json.Marshal(servers)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (ctx *ServerContext) PostRoomsServer(w http.ResponseWriter, r *http.Request) {
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

func (ctx *ServerContext) GetRoomsOfServer(w http.ResponseWriter, r *http.Request) {
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

func (ctx *ServerContext) PostRoomToServer(w http.ResponseWriter, r *http.Request) {
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

func (ctx *ServerContext) JoinServer(w http.ResponseWriter, r *http.Request) {
	server, err := ctx.validateRoomsServerID(w, r)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, enumBadRequest)
		return
	}

	userID := r.Context().Value(middlewares.CtxUserIDKey).(string)

	userStatus := models.NewServerUserStatus().WithUserID(userID).WithServerID(server.ID.Hex())
	if err := userStatus.Create(ctx.Db); err != nil {
		newErrorResponse(w, http.StatusInternalServerError, enumInternalServerError)
		return
	}

	serverWithStatus := &models.RoomsServerWithStatus{
		RoomsServer: server,
		Status:      *userStatus,
	}

	json, _ := json.Marshal(serverWithStatus)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
