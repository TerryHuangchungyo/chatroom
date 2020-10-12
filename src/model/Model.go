package model

import (
	"database/sql"
	"fmt"
	"local/config"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var User UserModel

func init() {
	conn := config.DATABASE.Connection
	host := config.DATABASE.Host
	port := config.DATABASE.Port
	user := config.DATABASE.User
	pass := config.DATABASE.Password
	dbname := config.DATABASE.Dbname

	var err error
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", user, pass, host, port, dbname)
	db, err = sql.Open(conn, dataSource)

	if err != nil {
		panic(err)
	}

	// Initial Model
	User = UserModel{"Users", db}
}

func Destroy() {
	db.Close()
}
