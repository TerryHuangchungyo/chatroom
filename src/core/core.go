package core

import "time"

type Message struct {
	Action     uint32    `json:"action"`   // 動作
	UserId     string    `json:"userId"`   // 使用者帳號
	UserName   string    `json:"userName"` // 使用者名稱
	HubId      int64     `json:"hubId"`    // 聊天室ID
	HubName    string    `json:"hubName"`  // 聊天室名稱
	Content    string    `json:"content"`  // 訊息內容
	CreateTime time.Time `json:"time"`
}
