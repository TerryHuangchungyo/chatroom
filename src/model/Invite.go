package model

import (
	"chatroom/core"
	"database/sql"
)

/*InviteModel ...
針對邀請資料表操作的 Model
*/
type InviteModel struct {
	tableName string
	db        *sql.DB
}

/*
描述:
儲存邀請的紀錄，如果使用者上線後，一併發出

輸入:
* hubId:int64    聊天室Id
* userId:string  使用者Id
* invitor:string 邀請人Id

輸出:
* err:error      錯誤類別，讓外部的程式看需不需要處理
*/
func (model *InviteModel) CreateOrUpdate(hubId int64, userId string, invitor string) error {
	insertStmt, err := db.Prepare("INSERT INTO " + model.tableName +
		" ( hubId, userId, invitor, createTime) VALUE( ?, ?, ?,CURRENT_TIMESTAMP())")

	if err != nil {
		Error.Println(err.Error())
		return err
	}

	defer insertStmt.Close()

	updateStmt, err := db.Prepare("Update " + model.tableName + ` SET createTime = CURRENT_TIMESTAMP()
	 WHERE hubId = ? AND userId = ? AND invitor = ?;`)

	if err != nil {
		Error.Println(err.Error())
		return err
	}

	defer updateStmt.Close()

	_, err = insertStmt.Exec(hubId, userId, invitor)
	if err != nil {
		_, err = updateStmt.Exec(hubId, userId, invitor)

		if err != nil {
			Error.Println(err.Error())
			return err
		}
	}

	return nil
}

/*
描述:
獲得儲存邀請的紀錄，如果使用者上線後，一併發出

輸入:
* userId:string  使用者Id

輸出:
*
* err:error      錯誤類別，讓外部的程式看需不需要處理
*/
func (model *InviteModel) GetInviteList(userId string) ([]core.Message, error) {
	stmt, err := db.Prepare("SELECT h.hubId as hubId, h.hubName as hubName,i.invitor as invitorId, u.userName as invitorName, i.createTime FROM " +
		model.tableName + " i JOIN " + User.tableName + " u ON i.invitor = u.userId JOIN " +
		Hub.tableName + " h ON i.hubId = h.hubId WHERE i.userId = ?")

	if err != nil {
		Error.Println(err.Error())
		return nil, err
	}

	var inviteList []core.Message
	var inviteMessage core.Message

	rows, err := stmt.Query(userId)
	for rows.Next() {
		err = rows.Scan(&(inviteMessage.HubId), &(inviteMessage.HubName), &(inviteMessage.UserId),
			&(inviteMessage.UserName), &(inviteMessage.CreateTime))

		if err != nil {
			Error.Println(err.Error())
			continue
		}
		inviteList = append(inviteList, inviteMessage)
	}
	return inviteList, nil
}
