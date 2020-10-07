package wsservice

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Action   uint32 `json:"action"`
	UserId   uint32 `json:"userId"`
	UserName string `json:"userName"`
	HubId    uint32 `json:"hubId"`
	HubName  string `json:"hubName"`
	Content  string `json:"content"`
}

type Client struct {
	Id   uint32
	Name string
	Hubs map[uint32]bool
	conn *websocket.Conn
	send chan Message
}

func (c *Client) ReadPump() {
	defer func() {
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		var m Message
		json.Unmarshal(message, &m)
		fmt.Printf("%s get message %v\n", c.Name, m)
		go c.HandleAction(&m)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			var marshalMsg []byte

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			marshalMsg, err := json.Marshal(msg)

			if err == nil {
				c.conn.WriteMessage(websocket.TextMessage, marshalMsg)
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) HandleAction(message *Message) {
	switch message.Action {
	case SEND: // 傳送訊息到聊天室
		if _, isExist := c.Hubs[message.HubId]; isExist {
			hubs[message.HubId].Broadcast <- *message
		}
	case INVITE:
		clientId := message.UserId
		message.UserId = c.Id // 將被邀請人改成邀請人
		message.UserName = c.Name
		message.HubName = hubs[message.HubId].Name

		hubs[message.HubId].Inviting[clientId] = true // 被邀請人邀請中
		clients[clientId].send <- *message
	case ANSWER:
		if hubs[message.HubId].Inviting[message.UserId] { // 答覆的人的確在聊天室邀請中
			answer, err := strconv.ParseUint(message.Content, 10, 32)
			delete(hubs[message.HubId].Inviting, message.UserId)
			if err == nil && answer == 1 {
				hubs[message.HubId].Register <- message.UserId // 加入到聊天室中
			}
		}
	}
}
