package wsservice

import (
	"fmt"
)

var userId uint32
var clients []*Client

func init() {
	userId = 0
}

func CreateClient(name string) (error) {
	client := &Client{Id: userId, Name: name}
	clients = append(clients, client)
	userId++
	fmt.Printf("New User %d %s", client.Id, client.Name)
	fmt.Println(clients)
	return nil
}

func CreateHub( name string ) (error) {
	
}
