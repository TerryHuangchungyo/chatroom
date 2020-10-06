package main

import (
	"net/http"
	
	"local/wsservice"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化一個http服務
	router := gin.Default()

	// 載入要使用的HTML template
	router.LoadHTMLGlob("view/*")

	// 回傳首頁
	router.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", nil)
	})

	// 新增聊天室
	router.POST( "/chatroom", func( context *gin.Context ) {
		name := context.PostForm( "hubname" )
		wsservice.CreateHub( name )
	})

	// 新增使用者
	router.POST("/user", func(context *gin.Context) {
		name := context.PostForm("username")
		wsservice.CreateClient(name)
	})

	router.Run(":8080")
}
