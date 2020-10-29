package model

import (
	"chatroom/config"
	"database/sql"
	"strconv"
)

/*RegisterModel ...
針對訊息資料表操作的 Model
*/
type RegisterModel struct {
	tableName string
	db        *sql.DB
}

/*Insert ...
描述:
新的使用者加入聊天室

輸入:
* hubId:      int64   聊天室Id
* userId:     string  使用者Id
* memberType: int32   使用者型別
輸出:
資訊，由外面函式決定處理方法
*/
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

/*GetHubList ...
描述:
取得使用者所在的聊天室列表

輸入:
* userId:string 使用者Id

輸出:
* hubList:

	[] struct {
		HubId   int64  `json:"hubId"`
		HubName string `json:"hubName"`
	}

聊天室資訊列表，包含聊天室編號、聊天室名稱

* err: error 回傳錯誤資訊，由外面函式決定處理方法
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

func (model *RegisterModel) GetUserListByHubId(hubId int64) (
	[]struct {
		UserId   string `json:"userId"`
		UserName string `json:"userName"`
		Active   bool   `json:"active"`
	}, error) {

	// get alive userid in redis
	aliveUserIdSet := make(map[string]bool)

	redisStrSliceCmd := redisClient.SMembers(ctx, config.REDIS.HubUsersSetKeyPrefix+strconv.FormatInt(hubId, 10))
	aliveUserIdList, err := redisStrSliceCmd.Result()
	if err != nil {
		Error.Println(err.Error())
		return nil, err
	}

	for _, aliveUserId := range aliveUserIdList {
		aliveUserIdSet[aliveUserId] = true
	}

	// get user list by hubId
	stmt, err := db.Prepare("SELECT r.userId, u.userName FROM " + model.tableName +
		" r JOIN " + User.tableName + " u ON r.userId = u.userId WHERE r.hubId = ?")

	if err != nil {
		Error.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(hubId)
	if err != nil {
		Error.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	var userList []struct {
		UserId   string `json:"userId"`
		UserName string `json:"userName"`
		Active   bool   `json:"active"`
	}

	for rows.Next() {
		var userInfo struct {
			UserId   string `json:"userId"`
			UserName string `json:"userName"`
			Active   bool   `json:"active"`
		}

		if err = rows.Scan(&(userInfo.UserId), &(userInfo.UserName)); err != nil {
			Error.Println(err)
			return nil, err
		}

		if _, isExist := aliveUserIdSet[userInfo.UserId]; isExist {
			userInfo.Active = true
		} else {
			userInfo.Active = false
		}

		userList = append(userList, userInfo)
	}
	return userList, nil
}
