package main

import "push/src/service"

func main() {
	if err := service.StartPushServer(); err != nil {
		panic(err)
	}
}
