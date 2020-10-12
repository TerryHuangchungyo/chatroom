package main

import (
	"local/controller"
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

	// 登入成功
	router.GET("/loginSuccess", func(context *gin.Context) {
		context.HTML(http.StatusOK, "loginSuccess.html", nil)
	})

	// 註冊
	router.GET("/signup", func(context *gin.Context) {
		context.HTML(http.StatusOK, "signup.html", nil)
	})

	router.POST("/signup", controller.User.Signup)

	// websocket服務
	router.GET("/chat/:id", controller.WebSocket.Serve)

	// 新增聊天室
	router.POST("/hub", controller.Hub.Create)

	router.Run(":8080")
}
