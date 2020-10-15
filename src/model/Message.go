package model

import (
	"database/sql"
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
func (model *MessageModel) Store(hubId int64, userId string, content *string) error {
	stmt, err := db.Prepare("INSERT INTO " + model.tableName +
		"( hubId, userId, content, createTime ) VALUE( ?, ?, ?, CURRENT_TIMESTAMP() )")

	if err != nil {
		Error.Println(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(hubId, userId, content)
	if err != nil {
		Error.Println(err.Error())
		return err
	}

	return nil
}
