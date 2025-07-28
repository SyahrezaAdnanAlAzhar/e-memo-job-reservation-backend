package websocket

import (
	"encoding/json"
	"fmt"
	"log"
)

type Hub struct {
	Broadcast        chan []byte
	Register         chan *Client
	Unregister       chan *Client
	Clients          map[*Client]bool
	incomingMessages chan clientMessage
	editingSessions  map[string]*Client
}

type EditingSession struct {
	UserID    int
	Entity    string
	ContextID int
}

type clientMessage struct {
	client  *Client
	message []byte
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:        make(chan []byte),
		Register:         make(chan *Client),
		Unregister:       make(chan *Client),
		Clients:          make(map[*Client]bool),
		incomingMessages: make(chan clientMessage),
		editingSessions:  make(map[string]*Client),
	}
}

// RUN HUB AS GOROUTINE
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			log.Println("WebSocket client registered")
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				h.cleanupClientSessions(client)
				delete(h.Clients, client)
				close(client.Send)
				log.Println("WebSocket client unregistered")
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		case clientMsg := <-h.incomingMessages:
			h.handleIncomingMessage(clientMsg.client, clientMsg.message)
		}
	}
}

func (h *Hub) BroadcastMessage(message []byte) {
	h.Broadcast <- message
}

func (h *Hub) handleIncomingMessage(client *Client, rawMessage []byte) {
	var msg Message
	if err := json.Unmarshal(rawMessage, &msg); err != nil {
		log.Printf("Error unmarshalling incoming message: %v", err)
		return
	}

	switch msg.Event {
	case "START_EDITING":
		payload, ok := msg.Payload.(map[string]interface{})
		if !ok {
			return
		}

		entity := payload["entity"].(string)
		contextID := int(payload["context_id"].(float64))
		sessionKey := fmt.Sprintf("%s:%d", entity, contextID)

		if _, exists := h.editingSessions[sessionKey]; !exists {
			h.editingSessions[sessionKey] = client
			log.Printf("User %d started editing %s", client.UserID, sessionKey)

			broadcastMsg, err := NewMessage("EDITING_STARTED", payload)

			if err != nil {
				log.Printf("CRITICAL: Failed to create broadcast message for start reorder: %v", err)
			} else {
				h.broadcastToOthers(broadcastMsg, client)
			}
		}

	case "FINISH_EDITING":
		payload, ok := msg.Payload.(map[string]interface{})
		if !ok {
			return
		}

		entity := payload["entity"].(string)
		contextID := int(payload["context_id"].(float64))
		sessionKey := fmt.Sprintf("%s:%d", entity, contextID)

		if holder, exists := h.editingSessions[sessionKey]; exists && holder == client {
			delete(h.editingSessions, sessionKey)
			log.Printf("User %d finished editing %s", client.UserID, sessionKey)

			broadcastMsg, err := NewMessage("EDITING_FINISHED", payload)

			if err != nil {
				log.Printf("CRITICAL: Failed to create broadcast message for finish reorder: %v", err)
			} else {
				h.broadcastToOthers(broadcastMsg, client)
			}

		}
	}
}

func (h *Hub) cleanupClientSessions(client *Client) {
	for key, holder := range h.editingSessions {
		if holder == client {
			delete(h.editingSessions, key)
			log.Printf("Cleaned up editing session %s for disconnected user %d", key, client.UserID)
		}
	}
}

func (h *Hub) broadcastToOthers(message []byte, exclude *Client) {
	for client := range h.Clients {
		if client != exclude {
			client.Send <- message
		}
	}
}
