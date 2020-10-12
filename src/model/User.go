package model

import (
	"database/sql"
)

type UserModel struct {
	tableName string
	db        *sql.DB
}

func (model *UserModel) Create(userId string, userName string, password string) error {
	stmt, err := db.Prepare("INSERT INTO " + model.tableName +
		"( userId, userName, password, createTime ) VALUE( ?, ?, ?, CURRENT_TIMESTAMP() )")

	if err != nil {
		return err
	}

	_, err = stmt.Exec(userId, userName, password)
	if err != nil {
		return err
	}

	return nil
}
