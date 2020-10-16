package websocket

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

type Message struct {
	Action   uint32 `json:"action"`   // 動作
	UserId   string `json:"userId"`   // 使用者帳號
	UserName string `json:"userName"` // 使用者名稱
	HubId    int64  `json:"hubId"`    // 聊天室ID
	HubName  string `json:"hubName"`  // 聊天室名稱
	Content  string `json:"content"`  // 訊息內容
}

type Client struct {
	id     string          // 使用者ID
	name   string          // 使用者名稱
	hubs   map[int64]bool  // 使用者擁有的聊天室
	wsConn *websocket.Conn // 使用者所使用的websocket連線
	sub    *redis.PubSub   // 訂閱redis，連線的客戶端
	mail   chan *Message   // 要送給使用者操作的訊息，像是邀請加入聊天室
}

func (c *Client) GetId() string {
	return c.id
}

func (c *Client) GetName() string {
	return c.name
}

/***
 * 從Websocket客戶端讀取訊息，並將原始JSON格式的資料解析後，
 * 交由HandleAction函式來處理後續動作
 */
func (c *Client) ReadPump() {
	defer func() {
		c.wsConn.Close()
	}()

	c.wsConn.SetReadLimit(maxMessageSize)
	c.wsConn.SetReadDeadline(time.Now().Add(pongWait))
	c.wsConn.SetPongHandler(func(string) error { c.wsConn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.wsConn.ReadMessage()

		if err != nil {
			log.Printf("Client %s %s %v", c.id, c.name, err)
			break
		}
		c.HandleAction(message)
	}
}

/***
 * 從send中拿取message，使用將資料封裝成JSON格式後，傳輸給websocket客戶端
 */
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.wsConn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.mail:
			if !ok {
				c.wsConn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			var marshalMsg []byte

			c.wsConn.SetWriteDeadline(time.Now().Add(writeWait))
			marshalMsg, err := json.Marshal(msg)

			if err == nil {
				c.wsConn.WriteMessage(websocket.TextMessage, marshalMsg)
			}
		case msg, ok := <-c.sub.Channel():
			if !ok {
				c.wsConn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// fmt.Printf("Client %s Get %s\n", c.name, msg.Payload)
			c.wsConn.SetWriteDeadline(time.Now().Add(writeWait))
			c.wsConn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
		case <-ticker.C:
			c.wsConn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.wsConn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

/***
 * 根據message的action code來處理封包
 */
func (c *Client) HandleAction(message []byte) {
	var unmarshalMessage = &Message{}
	json.Unmarshal(message, unmarshalMessage)

	if _, ok := c.hubs[unmarshalMessage.HubId]; !ok {
		return
	}

	var hub *Hub
	if item, isExist := hubs.Load(unmarshalMessage.HubId); isExist {
		hub = item.(*Hub)
	} else {
		// check mysql databases whether exists such hub

		hub = &Hub{unmarshalMessage.HubId,
			unmarshalMessage.HubName,
			make(map[string]bool),
			redis.NewClient(&redisOpt),
			make(chan []byte, 64)}

		go hub.run() // 運行Hub
		hubs.Store(hub.id, hub)
		c.sub.PSubscribe(ctx, "hub:"+strconv.FormatInt(hub.id, 10))
	}

	// 如果訊息的使用者名稱沒有的話，補上
	if unmarshalMessage.UserName == "" {
		unmarshalMessage.UserName = c.name
	}

	// 如果訊息的聊天室名稱沒有的話，補上
	if unmarshalMessage.HubName == "" {
		unmarshalMessage.HubName = hub.name
	}

	message, _ = json.Marshal(*unmarshalMessage)

	switch unmarshalMessage.Action {
	case MESSAGE: // 傳送訊息到Hub，讓Hub備份訊息到Mysql並publish到redis
		hub.docker <- message
	case INVITE:
		clientId := unmarshalMessage.UserId
		unmarshalMessage.UserId = c.id // 將被邀請人改成邀請人

		hub.inviting[clientId] = true // 被邀請人邀請中

		// 如果使用者在線中的話，直接將訊息寄給他
		if item, isExist := clients.Load(clientId); isExist {
			item.(*Client).mail <- unmarshalMessage
		} else {
			// 否則存入資料庫，待使用者上線，一併寄出
		}
	case ANSWER:
		if hub.inviting[unmarshalMessage.UserId] { // 答覆的人的確在聊天室邀請中
			answer, err := strconv.ParseUint(unmarshalMessage.Content, 10, 32)
			delete(hub.inviting, unmarshalMessage.UserId)
			if err == nil && answer == 1 {
				c.sub.Subscribe(ctx, "hub:"+strconv.FormatInt(unmarshalMessage.HubId, 10)) // 訂閱此聊天室的redis頻道
				c.hubs[unmarshalMessage.HubId] = true                                      // 新增使用者擁有的聊天室
			}
		}
	}
}
