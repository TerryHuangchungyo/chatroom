package controller

import (
	"chatroom/service/websocket"

	"github.com/gin-gonic/gin"
)

type WebSocketController struct {
	Err error
}

func (w *WebSocketController) Serve(context *gin.Context) {

	userId := context.Param("userId")
	websocket.Serve(context.Writer, context.Request, userId)
}
