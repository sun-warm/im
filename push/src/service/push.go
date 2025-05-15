package service

import (
	"context"
	"fmt"
	"push/proto"

	"github.com/gorilla/websocket"
)

func (s *server) Push(ctx context.Context, req proto.PushMessageRequest) (*proto.PushMessageResponse, error) {
	conn, ok := s.connections.Load(req.Receiver)
	if !ok {
		return &proto.PushMessageResponse{
			ErrorCode:    proto.PushErrorCode_SOCKET_ERROR,
			ErrorMessage: "User not connected",
		}, nil
	}

	// 发送消息
	wsConn := conn.(*websocket.Conn)
	err := wsConn.WriteMessage(websocket.TextMessage, []byte(req.Content))
	if err != nil {
		fmt.Printf("failed to send message to user %s: %v\n", req.Receiver, err)
		return &proto.PushMessageResponse{
			ErrorCode:    proto.PushErrorCode_SOCKET_ERROR,
			ErrorMessage: "Failed to send message",
		}, nil
	}
	fmt.Println("Message sent successfully to user:", req.Receiver)
	return &proto.PushMessageResponse{
		ErrorCode: proto.PushErrorCode_OK,
		Content:   req.Content,
	}, nil
}
