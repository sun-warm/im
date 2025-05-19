package main

import (
	"context"
	"fmt"
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
	fmt.Println("Redis connection result:", res)
	service.StartMessageServer()
}
