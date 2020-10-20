package main

import (
	"chatroom/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	defer controller.Destroy()

	// 初始化一個http服務
	router := gin.Default()

	// 載入要使用的HTML template
	router.LoadHTMLGlob("view/*")

	// 綁定靜態文件目錄
	router.Static("/asset", "./asset")

	// 回傳首頁
	router.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "home.html", nil)
	})

	// 登入
	router.GET("/login", func(context *gin.Context) {
		context.HTML(http.StatusOK, "login.html", nil)
	})

	router.POST("/login", controller.User.Login)

	// 登入後聊天室頁面
	router.GET("/chatroom", func(context *gin.Context) {
		userId := context.Query("userId")
		context.HTML(http.StatusOK, "chatroom.html", gin.H{"userId": userId})
	})

	// 登出
	router.GET("/logout/:userId", controller.User.Logout)

	// 註冊
	router.GET("/signup", func(context *gin.Context) {
		context.HTML(http.StatusOK, "signup.html", nil)
	})

	router.POST("/signup", controller.User.Signup)

	// websocket服務
	router.GET("/chat/:userId", controller.WebSocket.Serve)

	// 聊天室列表
	router.GET("/hub/:userId", controller.Hub.List)

	// 新增聊天室
	router.POST("/hub", controller.Hub.Create)

	// 聊天室歷史訊息
	router.GET("/history/:hubId", controller.Hub.GetHistoryMessage)

	// 使用者被邀請進聊天室訊息
	router.GET("/invite/:userId", controller.Invite.List)
	router.Run(":8080")
}
