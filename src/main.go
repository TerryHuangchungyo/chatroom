package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: %s host:port", os.Args[0])
		os.Exit(1)
	}
	address := os.Args[1]

	// 初始化一個http服務
	router := gin.Default()

	// 載入要使用的HTML template
	router.LoadHTMLGlob("view/*")

	// 回傳首頁
	router.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", nil)
	})

	// socket服務
	router.GET("/chatroom/:id", func(context *gin.Context) {
		hubId := context.Param("id")
		wsserve
	})

	router.Run(address)
}
