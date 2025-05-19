package service

import (
	"conversation/generated/conversation"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	conversation.UnimplementedConversationServiceServer
}

func StartConversationServer() {
	listen, err := net.Listen("tcp", ":8103")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()
	conversation.RegisterConversationServiceServer(s, &Server{})

	defer func() {
		s.Stop()
		listen.Close()
	}()

	fmt.Println("message Serving 8103...")
	err = s.Serve(listen)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
