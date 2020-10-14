package model

import (
	"chatroom/config"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)
var db *sql.DB
var User UserModel
var Hub HubModel
var Register RegisterModel

func init() {
	// 初始化logger 紀錄錯誤資訊
	logFile, err := os.OpenFile("./log/model.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	Info = log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stderr, "Warning ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(logFile, "Error ", log.Ldate|log.Ltime|log.Lshortfile)

	conn := config.DATABASE.Connection
	host := config.DATABASE.Host
	port := config.DATABASE.Port
	user := config.DATABASE.User
	pass := config.DATABASE.Password
	dbname := config.DATABASE.Dbname

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", user, pass, host, port, dbname)
	db, err = sql.Open(conn, dataSource)

	if err != nil {
		panic(err)
	}

	// Initial Model
	User = UserModel{"Users", db}
	Hub = HubModel{"Hubs", db}
	Register = RegisterModel{"Registers", db}
}

func Destroy() {
	db.Close()
}
