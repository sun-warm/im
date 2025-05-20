package main

import (
	"user/dao"
	"user/src/service"
)

func main() {
	if err := dao.InitDB(); err != nil {
		panic(err)
	}

	service.StartUserServer()

}
