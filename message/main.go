package main

import (
	"fmt"
	"message/db"
	"message/src/service"
)

func main() {
	if err := db.InitRedis(); err != nil {
		fmt.Println("failed to connect to redis:", err)
		return
	}
	service.StartMessageServer()
}
