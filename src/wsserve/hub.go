package wsserve

import (
	"time"
)

type Message struct {
	id      uint32
	time    time.Time
	content []byte
}

type Hub struct {
	Id         uint32
	Name       string
	Registered map[uint32]bool
	Register   chan uint32
	Unregister chan uint32
	Broadcast  chan Message
	Popularity uint32
}

var hubId uint32

func init() {
	hubId = 0
}

func Newhub() *Hub {
	hub := Hub{
		Id:         hubId,
		Registered: make(map[uint32]bool),
		Register:   make(chan uint32),
		Broadcast:  make(chan Message),
		Popularity: 0,
	}
	hubId++
	return &hub
}

func (h *Hub) run() {
	for {
		select {
		case clientId := <-h.Register:
			h.Registered[clientId] = true
		case clientId := <-h.Unregister:
			delete(h.Registered, clientId)
			close(clients[clientId].send)
			h.Popularity++
		case message := <-h.Broadcast:
			for clientId, _ := range h.Registered {
				client := clients[clientId]
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.Registered, clientId)
				}
			}
		}
	}
}
