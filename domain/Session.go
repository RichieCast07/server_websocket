package domain

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Session struct {
	conn      *websocket.Conn
	SessionID string
}

// Add at the top with other imports and types
var (
	activeSessions = make(map[string]*Session)
	sessionsLock   sync.RWMutex
)

func NewSession(conn *websocket.Conn, userID string) *Session {
	session := &Session{
		conn:      conn,
		SessionID: userID,
	}

	// Add session to active sessions
	sessionsLock.Lock()
	activeSessions[userID] = session
	sessionsLock.Unlock()

	return session
}

// Add cleanup on disconnect
func (s *Session) cleanup() {
	sessionsLock.Lock()
	delete(activeSessions, s.SessionID)
	sessionsLock.Unlock()
}

// Modify readPump to include cleanup
func (s *Session) readPump() {
	defer func() {
		s.cleanup()
		s.conn.Close()
	}()

	for {
		messageType, p, err := s.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("Error %v", err)
				break
			}
		}

		if messageType != -1 {
			var recipe Recipe
			err := json.Unmarshal(p, &recipe)
			if err != nil {
				log.Println("Error parsing JSON:", err)
				continue
			}

			message := fmt.Sprintf("¡Nueva receta disponible! - %s", recipe.Title)
			log.Printf("%s", message)

			// Broadcast to all connected clients
			broadcastToAll(message)
		}

		time.Sleep(17 * time.Millisecond)
	}
}

// Add new function to broadcast to all sessions
func broadcastToAll(message string) {
	sessionsLock.RLock()
	defer sessionsLock.RUnlock()

	for _, session := range activeSessions {
		err := session.conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Printf("Error broadcasting to %s: %v", session.SessionID, err)
		}
	}
}

/* func (s *Session) writePump() {
	defer func() {
		s.conn.Close()
	}()

	for {

		messageType := websocket.TextMessage
		message := []byte("Message from server")

		err := s.conn.WriteMessage(messageType, message)

		if err != nil {
			log.Println("Write error: ", err)
			break
		}

		time.Sleep(10 * time.Second)
	}

	select {}
}
*/
// Add this struct at the top of the file
type Recipe struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Ingredients string `json:"ingredients"`
	Description string `json:"description"`
	UserID      int    `json:"user_id"`
}

// Update the broadcast function
func (s *Session) broadcast(messageType int, payloadbyte []byte) {
	var recipe Recipe
	err := json.Unmarshal(payloadbyte, &recipe)
	if err != nil {
		log.Println("Error parsing JSON:", err)
		return
	}

	// Format the message
	message := fmt.Sprintf("¡Nueva receta disponible! - %s", recipe.Title)

	err = s.conn.WriteMessage(messageType, []byte(message))
	if err != nil {
		log.Println("Broadcast error: ", err)
	}
}

// Change from startHandling to StartHandling (capitalize first letter)
func (s *Session) StartHandling() {
    log.Println(s.SessionID)
    s.readPump()
}
