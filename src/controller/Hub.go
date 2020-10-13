package controller

import (
	"local/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HubController struct {
	Err error
}

func (h *HubController) Create(context *gin.Context) {
	userId := context.PostForm("userId")
	hubName := context.PostForm("hubName")

	err := model.Hub.Create(userId, hubName)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Some error happend"})

		return
	}

	context.JSON(http.StatusOK, gin.H{
		"msg": "ok"})
}
