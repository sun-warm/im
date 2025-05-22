package main

import (
	"user/client"
	"user/dao"
	"user/src/service"
)

func main() {
	if err := dao.InitDB(); err != nil {
		panic(err)
	}
	_, err := client.InitPushClient()
	if err != nil {
		panic(err)
	}
	service.StartUserServer()

}
