package wsservice

import (
	"fmt"
)

var userId uint32
var hubId uint32
var clients []*Client
var hubs []*Hub

func init() {
	hubId = 0
	userId = 0
}

/***
 * 創造使用者
 */
func CreateClient(name string) ( *Client, error) {
	client := &Client{Id: userId, Name: name, Hubs: make( map[uint32]bool)}
	clients = append(clients, client)
	userId++
	fmt.Printf("New User %d %s Created", client.Id, client.Name)
	fmt.Println(clients)
	return client, nil
}

/***
 * 創造聊天室
 */
func CreateHub( hubname string, creater uint32 ) ( *Hub, error){
	hub := &Hub{Id: hubId, Name: hubname, Clients: make( map[uint32]bool)}

	hub.Clients[creater] = true
	clients[creater].Hubs[hub.Id] = true
	
	hubs = append( hubs, hub )
	hubId++
	fmt.Printf("New Hub %d %s Created", hub.Id, hub.Name)
	fmt.Println(hubs)
	return hub, nil
}
