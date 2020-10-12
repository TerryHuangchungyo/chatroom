package main

import (
	"local/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化一個http服務
	router := gin.Default()

	// 載入要使用的HTML template
	router.LoadHTMLGlob("view/*")

	// 綁定靜態文件目錄
	router.Static("/asset", "./asset")

	router.GET("/test", func(context *gin.Context) {
		context.HTML(http.StatusOK, "test.html", nil)
	})

	// 回傳首頁
	router.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", nil)
	})

	// websocket服務
	router.GET("/chat/:id", controller.WebSocket.Serve)

	// 新增使用者
	router.POST("/user", controller.User.Create)

	// 新增聊天室
	router.POST("/hub", controller.Hub.Create)

	router.Run(":8080")
}
