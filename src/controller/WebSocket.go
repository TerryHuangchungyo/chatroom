package controller

import (
	"local/wsservice"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WebSocketController struct {
	Err error
}

func (w *WebSocketController) Serve(context *gin.Context) {
	id, _ := strconv.ParseUint(context.Param("id"), 10, 32)
	wsservice.ServeWs(context.Writer, context.Request, uint32(id))
}
