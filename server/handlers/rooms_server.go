package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/khalidibnwalid/Luma/middlewares"
	"github.com/khalidibnwalid/Luma/models"
)

func (ctx *ServerContext) validateRoomsServerID(w http.ResponseWriter, r *http.Request) (*models.RoomsServer, error) {
	rCtx := r.Context()
	serverID := r.PathValue("id")
	w.Header().Set("Content-Type", "application/json")

	if serverID == "" {
		newErrorResponse(w, http.StatusBadRequest, EnumServerIdRequired)
		return &models.RoomsServer{}, errors.New(EnumServerIdRequired)
	}

	uuidServerID, err := uuid.Parse(serverID)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, EnumServerIdInvalid)
		return &models.RoomsServer{}, errors.New(EnumServerIdInvalid)
	}

	serverData := models.RoomsServer{}
	if err := serverData.FindByID(ctx.Database.Client.WithContext(rCtx), uuidServerID); err != nil {
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
	userID := rCtx.Value(middlewares.CtxUserIDKey).(uuid.UUID)

	servers, err := models.NewServerUserStatus().WithUserID(userID).GetServers(ctx.Database.Client.WithContext(rCtx))
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
		newErrorResponse(w, http.StatusBadRequest, EnumBadRequest)
		return
	}
	if t.Name == "" {
		newErrorResponse(w, http.StatusBadRequest, EnumServerNameRequired)
		return
	}

	userID := rCtx.Value(middlewares.CtxUserIDKey).(uuid.UUID)

	server := models.NewRoomsServer().WithOwnerID(userID)
	server.Name = t.Name

	if err := server.Create(ctx.Database.Client.WithContext(rCtx)); err != nil {
		newErrorResponse(w, http.StatusInternalServerError, EnumInternalServerError)
		return
	}

	userStatus := models.NewServerUserStatus().WithUserID(userID).WithServerID(server.ID)
	if err := userStatus.Create(ctx.Database.Client.WithContext(rCtx)); err != nil {
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
	userId := rCtx.Value(middlewares.CtxUserIDKey).(uuid.UUID)
	server, err := ctx.validateRoomsServerID(w, r)
	if err != nil {
		return
	}

	rooms, err := server.GetRooms(ctx.Database.Client.WithContext(rCtx), userId)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, EnumInternalServerError)
		return
	}
	json, _ := json.Marshal(rooms)
	log.Println("rooms", rooms)
	log.Println("json", string(json))
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

	userId := rCtx.Value(middlewares.CtxUserIDKey).(uuid.UUID)

	room := models.Room{
		ServerID:  server.ID,
		Type:      t.Type,
		Name:      t.Name,
		GroupName: t.GroupName,
	}

	if err := room.Create(ctx.Database.Client.WithContext(rCtx)); err != nil {
		newErrorResponse(w, http.StatusInternalServerError, EnumInternalServerError)
		return
	}

	status := models.RoomUserStatus{
		UserID:   userId,
		RoomID:   room.ID,
		ServerID: server.ID,
	}

	status.Create(ctx.Database.Client.WithContext(rCtx))

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

	userID := rCtx.Value(middlewares.CtxUserIDKey).(uuid.UUID)

	userStatus := models.NewServerUserStatus().WithUserID(userID).WithServerID(server.ID)
	if err := userStatus.Create(ctx.Database.Client.WithContext(rCtx)); err != nil {
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
