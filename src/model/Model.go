package model

import (
	"chatroom/config"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

var redisOpt = redis.Options{
	Addr:     config.REDIS.Host + ":" + strconv.FormatInt(int64(config.REDIS.Port), 10),
	Password: config.REDIS.Password,
	DB:       config.REDIS.Db,
}

var db *sql.DB
var User UserModel
var Hub HubModel
var Register RegisterModel
var Message MessageModel
var Invite InviteModel
var redisClient *redis.Client
var ctx = context.Background()

func init() {
	// 初始化logger 紀錄錯誤資訊
	logFile, err := os.OpenFile("./log/model.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	Info = log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stderr, "Warning ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(logFile, "Error ", log.Ldate|log.Ltime|log.Lshortfile)

	// 資料庫連線
	conn := config.DATABASE.Connection
	host := config.DATABASE.Host
	port := config.DATABASE.Port
	user := config.DATABASE.User
	pass := config.DATABASE.Password
	dbname := config.DATABASE.Dbname
	charset := config.DATABASE.Charset
	collation := config.DATABASE.Collation

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&collation=%s&parseTime=true", user, pass, host, port, dbname, charset, collation)
	db, err = sql.Open(conn, dataSource)

	if err != nil {
		panic(err)
	}

	// redis 連線
	redisClient = redis.NewClient(&redisOpt)

	// Initial Model
	User = UserModel{"Users", db}
	Hub = HubModel{"Hubs", db}
	Register = RegisterModel{"Registers", db}
	Message = MessageModel{"Messages", db}
	Invite = InviteModel{"Invites", db}
}

func Destroy() {
	db.Close()
}
