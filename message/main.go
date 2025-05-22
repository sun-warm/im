package main

import (
	"context"
	"fmt"
	"message/client"
	"message/db"
	"message/src/service"
)

func main() {
	if err := db.InitRedis(); err != nil {
		fmt.Println("failed to connect to redis:", err)
		return
	}
	ctx := context.Background()
	res, err := db.Rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println(err)
	}
	_, err = client.InitPushClient()
	if err != nil {
		fmt.Println("failed to connect to push server:", err)
		return
	}

	fmt.Println("Redis connection result:", res)
	if err := service.StartMessageServer(); err != nil {
		panic(err)
	}
}
