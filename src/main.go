package main

import (
	"local/wsservice"
	"net/http"
	"strconv"

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
	router.GET("/chat/:id", func(context *gin.Context) {
		id, _ := strconv.ParseUint(context.Param("id"), 10, 32)
		wsservice.ServeWs(context.Writer, context.Request, uint32(id))
	})

	// 新增使用者
	router.POST("/user", func(context *gin.Context) {
		name := context.PostForm("username")
		client, _ := wsservice.CreateClient(name)
		context.JSON(http.StatusOK, gin.H{"id": client.GetId(), "username": client.GetName()})
	})

	// 新增聊天室
	router.POST("/hub", func(context *gin.Context) {
		c, _ := strconv.ParseUint(context.PostForm("userId"), 10, 32)
		creater := uint32(c)
		hubname := context.PostForm("hubname")
		hub, _ := wsservice.CreateHub(hubname, creater)
		context.JSON(http.StatusOK, gin.H{"id": hub.GetId(), "hubname": hub.GetName()})
	})

	router.Run(":8080")
}
