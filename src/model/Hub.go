package model

import (
	"database/sql"
)

/*HubModel ...
針對聊天室資料表操作的 Model
*/
type HubModel struct {
	tableName string
	db        *sql.DB
}

/*
描述:
新增新的聊天室

輸入:
* userId:string  使用者Id
* hubName:string 聊天室名稱

輸出:
* lastInsertId:int64 聊天室Id
* err:error 		   錯誤類別，讓外部的程式看需不需要處理
*/
func (model *HubModel) Create(hubName string, userId string) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO " + model.tableName +
		" ( hubName, ownerId, createTime) VALUE( ?, ?, CURRENT_TIMESTAMP())")

	if err != nil {
		Error.Println(err.Error())
		return -1, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(hubName, userId)
	if err != nil {
		Error.Println(err.Error())
		return -1, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		Error.Println(err.Error())
		return -1, err
	}

	return lastInsertId, err
}
