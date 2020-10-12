package controller

import (
	"local/wsservice"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HubController struct {
	Err error
}

func (h *HubController) Create(context *gin.Context) {
	c, _ := strconv.ParseUint(context.PostForm("userId"), 10, 32)
	creater := uint32(c)
	hubname := context.PostForm("hubname")
	hub, _ := wsservice.CreateHub(hubname, creater)
	context.JSON(http.StatusOK, gin.H{"id": hub.GetId(), "hubname": hub.GetName()})
}
