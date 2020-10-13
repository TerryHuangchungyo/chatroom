package model

import (
	"database/sql"
)

/**
 * 針對Hubs table 操作的 Model
 */
type HubModel struct {
	tableName string
	db        *sql.DB
}

func (model *HubModel) Create(userId string, hubName string) error {
	stmt, err := db.Prepare("INSERT INTO " + model.tableName +
		"( 'hubName', 'userId', 'createTime') VALUE( ?, ?, CURRENT_TIMESTAMP())")

	if err != nil {
		Error.Println(err.Error())
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(userId, hubName)
	if err != nil {
		Error.Println(err.Error())
		return err
	}

	return nil
}
