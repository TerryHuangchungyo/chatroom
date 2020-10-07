package wsservice

import (
	"fmt"
)

type Hub struct {
	Id        uint32
	Name      string
	Clients   map[uint32]bool
	Inviting  map[uint32]bool
	Register  chan uint32
	Broadcast chan Message
}

func (h *Hub) run() {
	fmt.Println("Hub is running")
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case message := <-h.Broadcast:
			message.Action = REPLY
			message.UserName = clients[message.UserId].Name
			message.HubName = hubs[message.HubId].Name
			for client, _ := range h.Clients {
				clients[client].send <- message
			}
		}
	}
}
