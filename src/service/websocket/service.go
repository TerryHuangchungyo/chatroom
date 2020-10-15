package websocket

import (
	"chatroom/config"
	"chatroom/model"
	"context"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

const (
	MESSAGE   = iota // 傳送訊息到聊天室
	INVITE           // 邀請加入聊天室
	ANSWER           // 答覆聊天室邀請
	BROADCAST        // 系統廣播
)

const (
	// 寫訊息容許等待的時間，取決於網路狀況
	writeWait = 10 * time.Second

	// 讀取pong訊息的等待時間
	pongWait = 60 * time.Second

	// 向客戶端撰寫ping訊息，容許等待時間
	pingPeriod = (pongWait * 9) / 10

	// 設定從客戶端最大可讀的訊息大小，以byte為基數
	maxMessageSize = 2048
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

var redisOpt = redis.Options{
	Addr:     config.REDIS.Host + ":" + strconv.FormatInt(int64(config.REDIS.Port), 10),
	Password: config.REDIS.Password,
	DB:       config.REDIS.Db,
}

var ctx = context.Background()

var clients sync.Map
var hubs sync.Map

func init() {
	clients = sync.Map{}
	hubs = sync.Map{}
}

/***
 * 提供websocket服務
 */
func Serve(w http.ResponseWriter, r *http.Request, userId string) {
	conn, err := upgrader.Upgrade(w, r, nil) // 將HTTP協議升級成Websocket協議
	if err != nil {
		log.Println(err)
		return
	}

	var client *Client
	if item, isExist := clients.Load(userId); isExist {
		client = item.(*Client)
		client.wsConn.WriteMessage(websocket.CloseMessage, []byte{})
		client.wsConn.Close()
		client.wsConn = conn
	} else {
		userName, _ := model.User.GetUserName(userId)
		client = &Client{id: userId,
			name:   userName,
			wsConn: conn,
			hubs:   make(map[int64]bool),
			sub:    redis.NewClient(&redisOpt).PSubscribe(ctx),
			mail:   make(chan *Message)}
		clients.Store(userId, client)
	}

	go client.ReadPump()
	go client.WritePump()
}
