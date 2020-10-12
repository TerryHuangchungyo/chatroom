package controller

var User UserController
var Hub HubController
var WebSocket WebSocketController

func init() {
	User = UserController{}
	Hub = HubController{}
	WebSocket = WebSocketController{}
}
