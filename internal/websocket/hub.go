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
	Clients          map[int]*Client
	incomingMessages chan clientMessage
	editingSessions  map[string]*Client
	sessionCommands  chan sessionCommand
}

type EditingSession struct {
	UserID    int
	Entity    string
	ContextID int
}

type sessionCommand struct {
	action     string
	client     *Client
	sessionKey string
	payload    map[string]interface{}
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
		Clients:          make(map[int]*Client),
		incomingMessages: make(chan clientMessage),
		editingSessions:  make(map[string]*Client),
		sessionCommands:  make(chan sessionCommand),
	}
}

// RUN HUB AS GOROUTINE
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if oldClient, ok := h.Clients[client.UserID]; ok {
				close(oldClient.Send)
			}
			h.Clients[client.UserID] = client
			log.Printf("WebSocket client registered for UserID: %d", client.UserID)
		case client := <-h.Unregister:
			if registeredClient, ok := h.Clients[client.UserID]; ok && registeredClient == client {
				h.cleanupClientSessions(client)
				delete(h.Clients, client.UserID)
				close(client.Send)
				log.Printf("WebSocket client unregistered for UserID: %d", client.UserID)
			}
		case message := <-h.Broadcast:
			for _, client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client.UserID)
				}
			}
		case clientMsg := <-h.incomingMessages:
			h.handleIncomingMessage(clientMsg.client, clientMsg.message)
		case cmd := <-h.sessionCommands:
			switch cmd.action {
			case "start":
				if _, exists := h.editingSessions[cmd.sessionKey]; !exists {
					h.editingSessions[cmd.sessionKey] = cmd.client
					log.Printf("User %d started editing %s", cmd.client.UserID, cmd.sessionKey)
					broadcastMsg, _ := NewMessage("EDITING_STARTED", cmd.payload)
					h.broadcastToOthers(broadcastMsg, cmd.client)
				}
			case "finish":
				if holder, exists := h.editingSessions[cmd.sessionKey]; exists && holder == cmd.client {
					delete(h.editingSessions, cmd.sessionKey)
					log.Printf("User %d finished editing %s", cmd.client.UserID, cmd.sessionKey)
					broadcastMsg, _ := NewMessage("EDITING_FINISHED", cmd.payload)
					h.broadcastToOthers(broadcastMsg, cmd.client)
				}
			}
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

	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		log.Printf("Invalid payload format for event: %s", msg.Event)
		return
	}

	entity, ok1 := payload["entity"].(string)
	contextIDFloat, ok2 := payload["context_id"].(float64)
	if !ok1 || !ok2 {
		log.Printf("Payload missing 'entity' or 'context_id' for event: %s", msg.Event)
		return
	}
	contextID := int(contextIDFloat)
	sessionKey := fmt.Sprintf("%s:%d", entity, contextID)

	switch msg.Event {
	case "START_EDITING":

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
	for userID, client := range h.Clients {
		if userID != exclude.UserID {
			client.Send <- message
		}
	}
}

func (h *Hub) ReleaseLock(client *Client, entity string, contextID int) {
	sessionKey := fmt.Sprintf("%s:%d", entity, contextID)

	if holder, exists := h.editingSessions[sessionKey]; exists && holder == client {
		delete(h.editingSessions, sessionKey)
		log.Printf("User %d released lock for %s", client.UserID, sessionKey)

		payload := map[string]interface{}{
			"entity":     entity,
			"context_id": contextID,
		}
		broadcastMsg, err := NewMessage("EDITING_FINISHED", payload)

		if err != nil {
			log.Printf("CRITICAL: Failed to create websocket message for finish reorder: %v", err)
		} else {
			h.broadcastToOthers(broadcastMsg, client)
		}

	}
}

func (h *Hub) GetClientByUserID(userID int) *Client {
	if client, ok := h.Clients[userID]; ok {
		return client
	}
	return nil
}
