package controller

import (
	"chatroom/config"
	"chatroom/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HubController struct {
	Err error
}

func (h *HubController) Create(context *gin.Context) {
	userId := context.PostForm("userId")
	hubName := context.PostForm("hubName")

	lastInsertId, err := model.Hub.Create(hubName, userId)
	if err != nil {
		Error.Println(err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Some error happend"})
		return
	}

	err = model.Register.Insert(lastInsertId, userId, config.MEMBER_MODERATOR)
	if err != nil {
		Error.Println(err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Some error happend"})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"hubId":   lastInsertId,
		"hubName": hubName,
		"msg":     "ok"})
}

func (h *HubController) List(context *gin.Context) {
	userId := context.Param("userId")

	result, err := model.Register.GetHubList(userId)

	if err != nil {
		Error.Println(err.Error())
		context.JSON(http.StatusInternalServerError, nil)
		return
	}
	fmt.Println(result)
	context.JSON(http.StatusOK, result)
}
