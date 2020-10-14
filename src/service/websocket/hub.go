package websocket

import (
	"fmt"
)

type Hub struct {
	id        uint32          // 聊天室Id
	name      string          // 聊天室名稱
	clients   map[uint32]bool // 聊天室有的使用者
	inviting  map[string]bool // 正在邀請的使用者
	register  chan string     // 等待註冊的使用者
	broadcast chan Message    // 廣播，訊息會發給所有使用者
}

func (h *Hub) GetId() uint32 {
	return h.id
}

func (h *Hub) GetName() string {
	return h.name
}

/***
 * 運行聊天室，主要的工作有
 * 1. 加入使用者
 * 2. 轉發廣播的訊息給所有使用者
 */
func (h *Hub) run() {
	fmt.Println("Hub " + h.name + " is running")
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			fmt.Printf("New User[%s] add in Hub[%s]\n", clients[client].name, h.name)
		case message := <-h.broadcast:
			message.Action = REPLY
			message.UserName = clients[message.UserId].name
			message.HubName = hubs[message.HubId].name
			for client, _ := range h.clients {
				// fmt.Printf("Message send to %s", clients[client].name)
				clients[client].send <- message
			}
		}
	}
}
