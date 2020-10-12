package controller

import (
	"local/wsservice"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Err error
}

func (u *UserController) Create(context *gin.Context) {
	name := context.PostForm("username")
	client, _ := wsservice.CreateClient(name)
	context.JSON(http.StatusOK, gin.H{"id": client.GetId(), "username": client.GetName()})
}
