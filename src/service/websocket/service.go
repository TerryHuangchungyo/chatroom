package websocket

import (
	"chatroom/config"
	"chatroom/core"
	"chatroom/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

const HISTORY_SIZE = 50

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
var Log *log.Logger
var redisClient *redis.Client

func init() {
	// 初始化logger 紀錄錯誤資訊
	logFile, err := os.OpenFile("./log/service.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	Log = log.New(logFile, "ERROR ", log.Ldate|log.Ltime|log.Lshortfile)

	clients = sync.Map{}
	hubs = sync.Map{}
	redisClient = redis.NewClient(&redisOpt)
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
			mail:   make(chan *core.Message)}

		list, err := model.Register.GetHubList(userId)

		if err == nil {
			for _, hubInfo := range list {
				client.hubs[hubInfo.HubId] = true
				userCnt, err := redisClient.SCard(ctx, config.REDIS.HubUsersSetKeyPrefix+strconv.FormatInt(hubInfo.HubId, 10)).Result()
				redisClient.SAdd(ctx, config.REDIS.HubUsersSetKeyPrefix+strconv.FormatInt(hubInfo.HubId, 10), client.id)

				if err != nil {
					Log.Println(err.Error())
				}

				if userCnt == 0 {
					var hub *Hub
					if _, isExist := hubs.Load(hubInfo.HubId); !isExist {
						hub = &Hub{hubInfo.HubId,
							hubInfo.HubName,
							redis.NewClient(&redisOpt),
							make(chan *core.Message, 64),
						}

						go hub.run() // 運行Hub
						hubs.Store(hub.id, hub)
						client.sub.PSubscribe(ctx, config.REDIS.ChannelKeyPrefix+strconv.FormatInt(hub.id, 10))
					}
				}
			}
		}

		clients.Store(userId, client)
		redisClient.SAdd(ctx, config.REDIS.UserAliveSet, client.id)
	}

	go client.ReadPump()
	go client.WritePump()
}

/*Destroy ...
關閉客戶端Client，將使用者Id從Redis移除，並檢查是否有Hub(聊天室人數為0)需要被關閉
*/
func Destroy(userId string) {
	item, isExist := clients.Load(userId)

	if isExist {
		clients.Delete(userId)
		client := item.(*Client)
		client.Destroy()
	}
	redisClient.SRem(ctx, config.REDIS.UserAliveSet, userId)
}

/*OwnerRegist ...
描述:
使用者創建新的聊天室必須記錄到使用者擁有的hubs中
*/
func OwnerRegist(userId string, hubId int64) error {
	item, isExist := clients.Load(userId)

	if isExist {
		var client *Client
		client = item.(*Client)
		client.hubs[hubId] = true
		redisClient.SAdd(ctx, config.REDIS.HubUsersSetKeyPrefix+strconv.FormatInt(hubId, 10), client.id)
		return nil
	}
	return fmt.Errorf("No such client %s running", userId)
}

/*GetHubHistoryMessage ...
描述：
獲取聊天室歷史訊息
*/
func GetHubHistoryMessage(hubId int64) ([]core.Message, error) {
	// 從Redis獲取歷史訊息，如果沒有就從資料庫抓取,並寫入redis
	var redisClient = redis.NewClient(&redisOpt)

	hubId_str := strconv.FormatInt(hubId, 10)
	historyIsExist, _ := redisClient.Exists(ctx, config.REDIS.HubHistoryKeyPrefix+hubId_str).Result()
	var historyMessage []core.Message

	if historyIsExist == 1 {
		marshalMessageList, _ := redisClient.LRange(ctx, config.REDIS.HubHistoryKeyPrefix+hubId_str, 0, -1).Result()

		for _, marshalMessage := range marshalMessageList {
			var msg core.Message
			json.Unmarshal([]byte(marshalMessage), &msg)
			historyMessage = append(historyMessage, msg)
		}
	} else {
		var err error
		historyMessage, err = model.Message.GetHistoryMessages(hubId, HISTORY_SIZE)

		if err != nil {
			return nil, err
		}

		for _, msg := range historyMessage {
			marshalMsg, err := json.Marshal(msg)

			if err != nil {
				redisClient.Del(ctx, config.REDIS.HubHistoryKeyPrefix+hubId_str)
				return nil, err
			}

			redisClient.RPush(ctx, config.REDIS.HubHistoryKeyPrefix+hubId_str, marshalMsg)
		}
	}

	return historyMessage, nil
}
