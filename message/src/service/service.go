package service

import (
	"fmt"
	"message/generated/message"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	message.UnimplementedMessageServiceServer
}

func StartMessageServer() {
	listen, err := net.Listen("tcp", ":8082")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()
	message.RegisterMessageServiceServer(s, &Server{})

	defer func() {
		s.Stop()
		listen.Close()
	}()

	fmt.Println("message Serving 8082...")
	err = s.Serve(listen)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
