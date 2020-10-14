package controller

import (
	"chatroom/model"
	"log"
	"os"
)

var logFile *os.File
var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

var User UserController
var Hub HubController
var WebSocket WebSocketController

func init() {
	// 初始化logger 紀錄錯誤資訊
	logFile, err := os.OpenFile("./log/controller.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	Info = log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stderr, "Warning ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(logFile, "Error ", log.Ldate|log.Ltime|log.Lshortfile)

	User = UserController{}
	Hub = HubController{}
	WebSocket = WebSocketController{}
}

func Destroy() {
	logFile.Close()
	model.Destroy()
}
