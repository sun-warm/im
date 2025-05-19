package service

import (
	"context"
	"fmt"
	"message/generated/message"
	"message/utils"
)

func (s *Server) GetMessage(ctx context.Context, req *message.GetMessageRequest) (*message.GetMessageResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request is nil")
	}
	if req.UserName == "" {
		return nil, fmt.Errorf("username is empty")
	}

	conversationID := req.ConversationId
	messagesString, err := utils.GetMessageFromZset(conversationID, req.Cursor, req.Cursor+req.Limit-1)
	if err != nil {
		return nil, err
	}
	fmt.Println("messages:", messagesString)
	messages, err := utils.MarshalMessageStringtoMessage(messagesString)
	if err != nil {
		return nil, err
	}
	return &message.GetMessageResponse{ErrorCode: message.MessageErrorCode_OK, Messages: messages}, nil
}
