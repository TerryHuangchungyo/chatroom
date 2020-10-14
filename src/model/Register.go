package model

import (
	"database/sql"
)

type RegisterModel struct {
	tableName string
	db        *sql.DB
}

func (model *RegisterModel) Insert(hubId int64, userId string, memberType int32) error {
	stmt, err := db.Prepare("INSERT INTO " + model.tableName +
		"( hubId, userId, type, registerTime ) VALUE( ?, ?, ?, CURRENT_TIMESTAMP() )")

	if err != nil {
		Error.Println(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(hubId, userId, memberType)
	if err != nil {
		Error.Println(err.Error())
		return err
	}

	return nil
}

/*
取得使用者所在的聊天室列表

輸入:
* userId: string 使用者Id

輸出:

*/
func (model *RegisterModel) GetHubList(userId string) (
	[]struct {
		HubId   int64  `json:"hubId"`
		HubName string `json:"hubName"`
	}, error) {
	stmt, err := db.Prepare("SELECT r.hubId, h.hubName FROM " + model.tableName +
		" r JOIN " + Hub.tableName + " h ON r.hubId = h.hubId WHERE r.userId = ?")

	if err != nil {
		Error.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId)
	if err != nil {
		Error.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	var hubList []struct {
		HubId   int64  `json:"hubId"`
		HubName string `json:"hubName"`
	}

	for rows.Next() {
		var hubInfo struct {
			HubId   int64  `json:"hubId"`
			HubName string `json:"hubName"`
		}

		if err = rows.Scan(&(hubInfo.HubId), &(hubInfo.HubName)); err != nil {
			Error.Println(err)
			return nil, err
		}
		hubList = append(hubList, hubInfo)
	}
	return hubList, nil
}
