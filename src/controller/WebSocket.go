package controller

import (
	"github.com/gin-gonic/gin"
)

type WebSocketController struct {
	Err error
}

func (w *WebSocketController) Serve(context *gin.Context) {
	// id, _ := strconv.ParseUint(context.Param("id"), 10, 32)
	// websocket.ServeWs(context.Writer, context.Request, uint32(id))
}
