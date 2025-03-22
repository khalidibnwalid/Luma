package core

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Topic struct {
	ID      string
	Clients map[*websocket.Conn]bool
	mu      sync.Mutex // Mutex to protect the Clients map
}

type TopicStore struct {
	Topics map[string]*Topic
}

func NewTopicStore() *TopicStore {
	return &TopicStore{
		Topics: make(map[string]*Topic),
	}
}

func (s *TopicStore) GetOrCreateRoom(id string) *Topic {
	if topic, exists := s.Topics[id]; exists {
		return topic
	}
	topic := &Topic{
		ID:      id,
		Clients: make(map[*websocket.Conn]bool),
	}
	s.Topics[id] = topic
	return topic
}

func (t *Topic) Subscribe(conn *websocket.Conn) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Clients[conn] = true
}

func (t *Topic) Unsubscribe(conn *websocket.Conn) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.Clients, conn)
	conn.Close()
}

func (t *Topic) Publish(message []byte) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for ws := range t.Clients {
		if err := ws.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Broadcast error:", err)
			t.Unsubscribe(ws)
		}
	}
}
