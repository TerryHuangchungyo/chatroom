package wsserve

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var guestId uint32 = 0
var hubs map[uint32]*Hub
var clients map[uint32]*Client

func init() {
	hubs = make(map[uint32]*Hub)
	clients = make(map[uint32]*Client)
}

func Serve(hubId uint32, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{Id: guestId, Hubs: make([]Hub, 5), conn: conn, send: make(chan Message, 256)}
	guestId++
	hubs[hubId].Register <- client.Id
	clients[client.Id] = client
}
