package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/khalidibnwalid/Luma/middlewares"
	"github.com/khalidibnwalid/Luma/models"
)

func (ctx *ServerContext) validateRoomsServerID(w http.ResponseWriter, r *http.Request) (*models.RoomsServer, error) {
	rCtx := r.Context()
	serverID := r.PathValue("id")
	if serverID == "" {
		w.Header().Set("Content-Type", "application/json")
		newErrorResponse(w, http.StatusBadRequest, EnumServerIdRequired)
		return &models.RoomsServer{}, errors.New(EnumServerIdRequired)
	}

	serverData := models.RoomsServer{}
	if err := serverData.FindById(ctx.Db, rCtx, serverID); err != nil {
		w.Header().Set("Content-Type", "application/json")
		newErrorResponse(w, http.StatusNotFound, EnumServerNotFound)
		return &models.RoomsServer{}, errors.New(EnumServerNotFound)
	}

	return &serverData, nil
}

// Deprecated: use GetUserRoomsServer instead, might keep it for refreshing data,
// but I might also remove it and push notification to the user
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
	rCtx := r.Context()
	userID := rCtx.Value(middlewares.CtxUserIDKey).(string)

	servers, err := models.NewServerUserStatus().WithUserID(userID).GetServers(ctx.Db, rCtx)
	if err != nil {
		http.Error(w, "Error getting servers", http.StatusInternalServerError)
		return
	}

	json, _ := json.Marshal(servers)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (ctx *ServerContext) PostRoomsServer(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()
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

	userID := rCtx.Value(middlewares.CtxUserIDKey).(string)

	server := models.NewRoomsServer().WithOwnerID(userID)
	server.Name = t.Name

	if err := server.Create(ctx.Db, rCtx); err != nil {
		newErrorResponse(w, http.StatusInternalServerError, EnumInternalServerError)
		return
	}

	userStatus := models.NewServerUserStatus().WithUserID(userID).WithServerID(server.ID.Hex())
	if err := userStatus.Create(ctx.Db, rCtx); err != nil {
		newErrorResponse(w, http.StatusInternalServerError, EnumInternalServerError)
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

func (ctx *ServerContext) GetRoomsOfServer(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()
	userId := rCtx.Value(middlewares.CtxUserIDKey).(string)
	server, err := ctx.validateRoomsServerID(w, r)
	if err != nil {
		return
	}

	rooms, err := server.GetRooms(ctx.Db, rCtx, userId)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, EnumInternalServerError)
		return
	}

	json, _ := json.Marshal(rooms)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (ctx *ServerContext) PostRoomToServer(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()
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
		newErrorResponse(w, http.StatusBadRequest, EnumBadRequest)
		return
	}
	if t.Type == "" {
		newErrorResponse(w, http.StatusBadRequest, EnumServerTypeRequired)
		return
	}

	if t.Name == "" {
		newErrorResponse(w, http.StatusBadRequest, EnumServerNameRequired)
		return
	}

	userId := rCtx.Value(middlewares.CtxUserIDKey).(string)

	room := models.Room{
		ServerID:  server.ID.Hex(),
		Type:      t.Type,
		Name:      t.Name,
		GroupName: t.GroupName,
	}

	if err := room.Create(ctx.Db, rCtx); err != nil {
		newErrorResponse(w, http.StatusInternalServerError, EnumInternalServerError)
		return
	}

	status := models.RoomUserStatus{
		UserID:    userId,
		RoomID:    room.ID.Hex(),
		ServerID:  server.ID.Hex(),
	}

	status.Create(ctx.Db, rCtx)

	RoomWithStatus := models.RoomWithStatus{
		Room:   &room,
		Status: &status,
	}

	json, _ := json.Marshal(RoomWithStatus)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (ctx *ServerContext) JoinServer(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()
	server, err := ctx.validateRoomsServerID(w, r)
	if err != nil {
		return
	}

	userID := rCtx.Value(middlewares.CtxUserIDKey).(string)

	userStatus := models.NewServerUserStatus().WithUserID(userID).WithServerID(server.ID.Hex())
	if err := userStatus.Create(ctx.Db, rCtx); err != nil {
		newErrorResponse(w, http.StatusInternalServerError, EnumInternalServerError)
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
