package websocket

import (
	"chatroom/config"
	"chatroom/model"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type Hub struct {
	id     int64         // 聊天室Id
	name   string        // 聊天室名稱
	pub    *redis.Client // 連線redis的client
	docker chan *Message // 訊息接收的channel，之後會publish到redis跟mysql備份
}

func (h *Hub) GetId() int64 {
	return h.id
}

func (h *Hub) GetName() string {
	return h.name
}

/***
 * 運行聊天室，主要的工作有
 * 1. 加入使用者
 * 2. 訊息備份到mysql
 * 3. 訊息publish到redis料庫
 */
func (h *Hub) run() {
	fmt.Println("Hub " + h.name + " is running")
	for {
		select {
		case message := <-h.docker:
			message.CreateTime = time.Now()

			marshalMessage, err := json.Marshal(message)
			if err != nil {
				Log.Println(err.Error())
				return
			}

			// Mysql備份
			err = model.Message.Store(h.id, message.UserId, message.Conten, message.CreateTime)

			// Publish到Redis
			if err == nil {
				h.pub.Publish(ctx, config.REDIS.ChannelKeyPrefix+strconv.FormatInt(h.id, 10), marshalMessage)
			}
		}
	}
}
