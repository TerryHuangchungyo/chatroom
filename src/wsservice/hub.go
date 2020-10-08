package wsservice

import (
	"fmt"
)

type Hub struct {
	id        uint32
	name      string
	clients   map[uint32]bool
	inviting  map[uint32]bool
	register  chan uint32
	broadcast chan Message
}

func (h *Hub) GetId() uint32 {
	return h.id
}

func (h *Hub) GetName() string {
	return h.name
}

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
				clients[client].send <- message
			}
		}
	}
}
