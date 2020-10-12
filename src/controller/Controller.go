package controller

import (
	"local/model"
)

var User UserController
var Hub HubController
var WebSocket WebSocketController

func init() {
	User = UserController{}
	Hub = HubController{}
	WebSocket = WebSocketController{}
}

func Destroy() {
	model.Destroy()
}
