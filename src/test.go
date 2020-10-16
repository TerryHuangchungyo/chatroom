package main

import (
	"chatroom/config"
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type Person struct {
	name       string
	age        int
	createTime time.Time
}

var ctx = context.Background()

func main() {
	client1 := redis.NewClient(&redis.Options{
		Addr:     config.REDIS.Host + ":" + strconv.FormatInt(int64(config.REDIS.Port), 10),
		Password: config.REDIS.Password,
		DB:       config.REDIS.Db,
	})

	var p = Person{"Terry", 22, time.Now()}

	client1.Publish(ctx, "chan1", p)
	time.Sleep(time.Second * 5)
}
