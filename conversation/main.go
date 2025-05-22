package main

import (
	"conversation/db"
	"conversation/internal/service"
	"conversation/utils/snowflake"
)

func main() {
	db.InitRedis()
	snowflake.NewSnowflake()
	if err := service.StartConversationServer(); err != nil {
		panic(err)
	}
}
