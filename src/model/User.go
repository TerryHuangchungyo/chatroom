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
		Error.Println(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, userName, password)
	if err != nil {
		Error.Println(err.Error())
		return err
	}

	return nil
}

func (model *UserModel) GetUserName(userId string) (string, error) {
	stmt, err := db.Prepare("SELECT userName FROM " + model.tableName +
		" WHERE userId = ?")

	if err != nil {
		Error.Println(err.Error())
		return "", err
	}
	defer stmt.Close()

	var name string
	err = stmt.QueryRow(userId).Scan(&name)

	if err != nil {
		Error.Println(err.Error())
		return "", err
	}

	return name, nil
}

func (model *UserModel) GetPassword(userId string) (string, error) {
	stmt, err := db.Prepare("SELECT password FROM " + model.tableName +
		" WHERE userId = ?")

	if err != nil {
		Error.Println(err.Error())
		return "", err
	}
	defer stmt.Close()

	var password string
	err = stmt.QueryRow(userId).Scan(&password)

	if err != nil {
		Error.Println(err.Error())
		return "", err
	}

	return password, nil
}
