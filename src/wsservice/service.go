package wsservice

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	SEND   = iota // 傳送訊息到聊天室
	REPLY         // 聊天室訊息回覆
	INVITE        // 邀請加入聊天室
	ANSWER        // 答覆聊天室邀請
)

const (
	// 寫訊息容許等待的時間，取決於網路狀況
	writeWait = 10 * time.Second

	// 讀取pong訊息的等待時間
	pongWait = 60 * time.Second

	// 向客戶端撰寫ping訊息，容許等待時間
	pingPeriod = (pongWait * 9) / 10

	// 設定從客戶端最大可讀的訊息大小，以byte為基數
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var userId uint32
var hubId uint32
var clients []*Client
var hubs []*Hub

func init() {
	hubId = 0
	userId = 0
}

/***
 * 提供websocket服務
 */
func ServeWs(w http.ResponseWriter, r *http.Request, id uint32) {
	conn, err := upgrader.Upgrade(w, r, nil) // 將HTTP協議升級成Websocket協議
	if err != nil {
		log.Println(err)
		return
	}

	client := clients[id]
	client.conn = conn

	go client.ReadPump()
	go client.WritePump()
}

/***
 * 創造使用者
 */
func CreateClient(name string) (*Client, error) {
	client := &Client{id: userId, name: name, hubs: make(map[uint32]bool), send: make(chan Message, 256)}
	clients = append(clients, client)
	userId++
	fmt.Printf("New User %d %s Created", client.id, client.name)
	fmt.Println(clients)
	return client, nil
}

/***
 * 創造聊天室
 */
func CreateHub(hubname string, creater uint32) (*Hub, error) {
	hub := &Hub{id: hubId,
		name:      hubname,
		clients:   make(map[uint32]bool),
		inviting:  make(map[uint32]bool),
		register:  make(chan uint32),
		broadcast: make(chan Message),
	}

	hub.clients[creater] = true
	clients[creater].hubs[hub.id] = true
	go hub.run()

	hubs = append(hubs, hub)
	hubId++

	fmt.Printf("New Hub %d %s Created", hub.id, hub.name)
	fmt.Println(hubs)
	return hub, nil
}
