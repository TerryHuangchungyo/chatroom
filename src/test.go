package main

import (
	"chatroom/config"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	client1 := redis.NewClient(&redis.Options{
		Addr:     config.REDIS.Host + ":" + strconv.FormatInt(int64(config.REDIS.Port), 10),
		Password: config.REDIS.Password,
		DB:       config.REDIS.Db,
	})

	pong, err := client1.Ping(ctx).Result()
	if err != nil {
		fmt.Println(pong)
		fmt.Println(err.Error())
	}

	client2 := redis.NewClient(&redis.Options{
		Addr:     config.REDIS.Host + ":" + strconv.FormatInt(int64(config.REDIS.Port), 10),
		Password: config.REDIS.Password,
		DB:       config.REDIS.Db,
	})

	quit := make(chan int)

	sub := client1.PSubscribe(ctx)
	go func() {
		for {
			select {
			case msg := <-sub.Channel():
				fmt.Printf("Message %s from channel %s\n", msg.Payload, msg.Channel)
			case <-quit:
				fmt.Println("Exit")
				break
			}
		}
	}()

	sub.Subscribe(ctx, "chan1")

	client2.Publish(ctx, "chan1", "Hello! I'm client2")

	sub.Subscribe(ctx, "chan2")
	client2.Publish(ctx, "chan2", "@")
	client2.Publish(ctx, "chan1", "@")
	quit <- 1
	time.Sleep(time.Second * 10)
}
