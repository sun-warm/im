package service

import (
	"fmt"
	"net"
	"user/generated/user"

	"google.golang.org/grpc"
)

type server struct {
	user.UnimplementedUserServiceServer
}

func StartUserServer() {
	listen, err := net.Listen("tcp", ":8104")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()
	user.RegisterUserServiceServer(s, &server{})

	defer func() {
		s.Stop()
		listen.Close()
	}()

	fmt.Println("message Serving 8104...")
	err = s.Serve(listen)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
