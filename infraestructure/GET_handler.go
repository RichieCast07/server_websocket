package infrastructure

import (
	"net/http"
	"websocket/application"

	"github.com/gin-gonic/gin"
)

type WebsocketHandler struct {
	wsService application.WebsocketService
}

func NewWebsocketHandler(
	appService application.WebsocketService,
) *WebsocketHandler {
	return &WebsocketHandler{
		wsService: appService,
	}
}

func (wsH *WebsocketHandler) Upgrade(ctx *gin.Context) {
	userID := ctx.Query("user_id")

	err := wsH.wsService.HandleConnection(ctx.Writer, ctx.Request, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to upgrade"})
	}
}
