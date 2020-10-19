package controller

import (
	"chatroom/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InviteController struct {
	Err error
}

func (h *InviteController) List(context *gin.Context) {
	userId := context.Param("userId")
	inviteList, _ := model.Invite.GetInviteList(userId)
	context.JSON(http.StatusOK, inviteList)
}
