package server

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Create a separate file for text chat functionality and WebSocket management.

// Implement a map to store connected WebSocket clients for text chat.
var connectedClients = make(map[*websocket.Conn]bool)

// Use a mutex for synchronized access to connectedClients.
var connectedClientsMu sync.Mutex

// Handle text chat messages sent by clients.
func handleTextChat(roomID, text string) {
	message := map[string]interface{}{
		"text":      text,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}

	broadcastMsg := broadcastMsg{
		Message: message,
		RoomID:  roomID,
	}

	broadcast <- broadcastMsg
}

// JoinRoomRequestHandler joins the client in a particular room.
func JoinRoomRequestHandler2(w http.ResponseWriter, r *http.Request) {
	roomID, ok := r.URL.Query()["roomID"]

	if !ok {
		log.Println("roomID missing in URL Parameters")
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return
	}

	// Insert the WebSocket connection into the connectedClients map.
	connectedClientsMu.Lock()
	connectedClients[ws] = true
	connectedClientsMu.Unlock()

	AllRooms.InsertIntoRoom(roomID[0], false, ws)

	// Start a goroutine to handle WebSocket messages.
	go func() {
		defer func() {
			// Remove the WebSocket connection from connectedClients when it's closed.
			connectedClientsMu.Lock()
			defer connectedClientsMu.Unlock()
			delete(connectedClients, ws)
		}()

		for {
			var msg broadcastMsg

			err := ws.ReadJSON(&msg.Message)
			if err != nil {
				log.Println("Read Error:", err)
				return
			}

			msg.Client = ws
			msg.RoomID = roomID[0]

			// Handle text chat messages separately from WebRTC messages.
			if text, ok := msg.Message["text"].(string); ok {
				handleTextChat(msg.RoomID, text)
			}

			broadcast <- msg
		}
	}()
}
