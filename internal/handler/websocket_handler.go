package handler

import (
	"log"
	"net/http"

	ws "github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/websocket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct {
	hub *ws.Hub
}

func NewWebSocketHandler(hub *ws.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

// HANDLE WEBSOCKET REQUEST FROM CLIENT
func (h *WebSocketHandler) ServeWs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &ws.Client{
		Hub:  h.hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}
	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
