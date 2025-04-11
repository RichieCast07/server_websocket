package application

import (
	"net/http"
	"websocket/WebSocketNew/domain"

	"github.com/gorilla/websocket"
)

type WebsocketService struct {
	upgrader websocket.Upgrader
}

func NewWebsocketService() *WebsocketService {
	return &WebsocketService{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (ws *WebsocketService) HandleConnection(
	w http.ResponseWriter, r *http.Request, userID string,
) error {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	session := domain.NewSession(conn, userID)
	session.StartHandling()  // This should now work with the capitalized method name

	return nil
}
