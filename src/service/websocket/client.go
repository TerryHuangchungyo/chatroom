package websocket

import (
	"chatroom/config"
	"chatroom/core"
	"chatroom/model"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

type Client struct {
	id     string             // 使用者ID
	name   string             // 使用者名稱
	hubs   map[int64]bool     // 使用者擁有的聊天室
	wsConn *websocket.Conn    // 使用者所使用的websocket連線
	sub    *redis.PubSub      // 訂閱redis，連線的客戶端
	mail   chan *core.Message // 要送給使用者操作的訊息，像是邀請加入聊天室
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
	var unmarshalMessage = &core.Message{}
	json.Unmarshal(message, unmarshalMessage)

	if _, ok := c.hubs[unmarshalMessage.HubId]; !ok {
		return
	}

	// 如果訊息的使用者名稱沒有的話，補上
	if unmarshalMessage.UserName == "" {
		unmarshalMessage.UserName = c.name
	}

	// 如果訊息的聊天室名稱沒有的話，補上
	if unmarshalMessage.HubName == "" {
		name, _ := model.Hub.GetHubName(unmarshalMessage.HubId)
		unmarshalMessage.HubName = name
	}

	message, _ = json.Marshal(*unmarshalMessage)

	switch unmarshalMessage.Action {
	case MESSAGE: // 傳送訊息到Hub，讓Hub備份訊息到Mysql並publish到redis
		var hub *Hub
		if item, isExist := hubs.Load(unmarshalMessage.HubId); isExist {
			hub = item.(*Hub)
		} else {
			// check mysql databases whether exists such hub

			hub = &Hub{unmarshalMessage.HubId,
				unmarshalMessage.HubName,
				redis.NewClient(&redisOpt),
				make(chan *core.Message, 64),
			}

			go hub.run() // 運行Hub
			hubs.Store(hub.id, hub)
			c.sub.PSubscribe(ctx, config.REDIS.ChannelKeyPrefix+strconv.FormatInt(hub.id, 10))
		}

		hub.docker <- unmarshalMessage
	case INVITE:
		clientId := unmarshalMessage.UserId

		_, err := model.User.GetUserName(clientId)
		if err != nil { // 沒有該使用者
			return
		}

		// 將邀請的訊息紀錄在Mysql資料庫中
		err = model.Invite.CreateOrUpdate(unmarshalMessage.HubId, unmarshalMessage.UserId, c.id)
		if err != nil { // 邀請失敗
			return
		}

		if item, isExist := clients.Load(clientId); isExist {
			// 如果使用者在線中的話，直接將訊息寄給他
			unmarshalMessage.UserId = c.id // 將被邀請人改成邀請人
			item.(*Client).mail <- unmarshalMessage
		}
		// case ANSWER:
		// 	if hub.inviting[unmarshalMessage.UserId] { // 答覆的人的確在聊天室邀請中
		// 		answer, err := strconv.ParseUint(unmarshalMessage.Content, 10, 32)
		// 		delete(hub.inviting, unmarshalMessage.UserId)
		// 		if err == nil && answer == 1 {
		// 			c.sub.Subscribe(ctx, "hub:"+strconv.FormatInt(unmarshalMessage.HubId, 10)) // 訂閱此聊天室的redis頻道
		// 			c.hubs[unmarshalMessage.HubId] = true                                      // 新增使用者擁有的聊天室
		// 		}
		// 	}
	}
}

func (client *Client) Destroy() {
	close(client.mail)
	client.sub.Close()
	client.wsConn.Close()

	// 從redis聊天室擁有的使用者移除，如果聊天室人數為零就從map移出，並釋放聊天室所擁有的資源
	for hubId, _ := range client.hubs {
		redisClient.SRem(ctx, config.REDIS.HubUsersSetKeyPrefix+strconv.FormatInt(hubId, 10), client.id)
		userCnt, err := redisClient.SCard(ctx, config.REDIS.HubUsersSetKeyPrefix+strconv.FormatInt(hubId, 10)).Result()
		if err != nil {
			Log.Println(err.Error())
			continue
		}

		if userCnt != 0 {
			continue
		}

		item, isExist := hubs.Load(hubId)
		if isExist {
			hubs.Delete(hubId)
			hub := item.(*Hub)
			hub.Destroy()
		}
	}
}
