package model

import (
	"chatroom/core"
	"database/sql"
	"time"
)

/*MessageModel ...
針對訊息資料表操作的 Model
*/
type MessageModel struct {
	tableName string
	db        *sql.DB
}

/*Store ...
儲存聊天室訊息

輸入:
* hubId:    int64   聊天室Id
* userId:   string  使用者Id
* context: *string  聊天室內容

輸出:
* err:error 在內部會紀錄log，回傳讓呼叫的程式看需不需要處理此種錯誤
*/
func (model *MessageModel) Store(hubId int64, userId string, content string, createTime time.Time) error {
	stmt, err := db.Prepare("INSERT INTO " + model.tableName +
		"( hubId, userId, content, createTime ) VALUE( ?, ?, ?, ? )")

	if err != nil {
		Error.Println(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(hubId, userId, content, createTime)
	if err != nil {
		Error.Println(err.Error())
		return err
	}

	return nil
}

func (model *MessageModel) GetHistoryMessages(hubId int64, limit int64) ([]core.Message, error) {
	stmt, err := db.Prepare("SELECT m.hubId, h.hubName, m.userId, u.userName, content, m.createTime FROM " + model.tableName +
		" m JOIN " + User.tableName + " u ON m.userId = u.userId JOIN " + Hub.tableName + " h ON m.hubId = h.hubId WHERE m.hubId = ? " +
		" ORDER BY m.createTime ASC LIMIT ?")

	if err != nil {
		Error.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(hubId, limit)
	if err != nil {
		Error.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	var message core.Message
	var messageList []core.Message

	for rows.Next() {
		if err = rows.Scan(&(message.HubId), &(message.HubName), &(message.UserId), &(message.UserName), &(message.Content), &(message.CreateTime)); err != nil {
			Error.Println(err)
			return nil, err
		}
		message.Action = 0
		messageList = append(messageList, message)
	}
	return messageList, nil
}
