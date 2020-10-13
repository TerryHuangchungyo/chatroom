package controller

import (
	"github.com/gin-gonic/gin"
)

type HubController struct {
	Err error
}

func (h *HubController) Create(context *gin.Context) {
	// userId := context.PostForm("userId")
	// hubName := context.PostForm("hubName")
}
