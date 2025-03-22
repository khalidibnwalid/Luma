package handlers

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/khalidibnwalid/Luma/models"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (ctx *HandlerContext) RoomWS(rooms *Rooms) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := r.PathValue("id")
		if roomID == "" {
			http.Error(w, "Room ID is required", http.StatusBadRequest)
			return
		}

		roomData := models.Room{}
		if err := roomData.FindById(ctx.Db, roomID); err != nil {
			http.Error(w, "Room not found", http.StatusNotFound)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}
		defer conn.Close()

		log.Printf("Room [%s] Connected\n", roomID)
		room := rooms.GetOrCreateRoom(roomData.ID.String())
		room.subscribe(conn)
		defer room.unsubscribe(conn)

		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}
			room.publish(messageType, p)
			log.Printf("Room [%s] Received: %s\n", roomID, p)
		}
	}
}

type Room struct {
	ID      string
	Clients map[*websocket.Conn]bool
	mu      sync.Mutex // Mutex to protect the Clients map
}

func (r *Room) subscribe(conn *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Clients[conn] = true
}

func (r *Room) unsubscribe(conn *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Clients, conn)
	conn.Close()
}

func (r *Room) publish(messageType int, message []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for ws := range r.Clients {
		if err := ws.WriteMessage(messageType, message); err != nil {
			log.Println("Broadcast error:", err)
			r.unsubscribe(ws)
		}
	}
}

type Rooms struct {
	Rooms map[string]*Room
}

func NewRooms() *Rooms {
	return &Rooms{
		Rooms: make(map[string]*Room),
	}
}

func (rs *Rooms) GetOrCreateRoom(id string) *Room {
	if room, exists := rs.Rooms[id]; exists {
		return room
	}
	room := &Room{
		ID:      id,
		Clients: make(map[*websocket.Conn]bool),
	}
	rs.Rooms[id] = room
	return room
}
