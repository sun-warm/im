package service

import (
	"context"
	"fmt"
	"push/proto"

	"github.com/bytedance/sonic"

	"github.com/gorilla/websocket"
)

func (s *server) Push(ctx context.Context, req proto.PushMessageRequest) (*proto.PushMessageResponse, error) {
	conn, ok := s.connections.Load(req.Receiver)
	if !ok {
		return &proto.PushMessageResponse{
			ErrorCode:    proto.PushErrorCode_USER_OFFLINE,
			ErrorMessage: "User is offline. Message stored.",
		}, nil
	}
	//fixme：如何处理离线的各种消息？
	// 发送消息
	wsConn := conn.(*websocket.Conn)
	messageString, err := sonic.MarshalString(req.PushMessage)
	if err != nil {
		//Fixme：这里不需要做补偿，直接丢回去，告诉他传的有问题。
		return nil, err
	}
	err = wsConn.WriteMessage(websocket.TextMessage, []byte(messageString))
	if err != nil {
		fmt.Printf("failed to send message to user %s: %v\n", req.Receiver, err)
		return &proto.PushMessageResponse{
			ErrorCode:    proto.PushErrorCode_SOCKET_ERROR,
			ErrorMessage: "Failed to send message",
		}, nil
	}
	//fixme:如果这次丢失败了，有什么办法补偿嘛？
	fmt.Println("Message sent successfully to user:", req.Receiver)
	return &proto.PushMessageResponse{
		ErrorCode: proto.PushErrorCode_OK,
		Content:   messageString,
	}, nil
}
