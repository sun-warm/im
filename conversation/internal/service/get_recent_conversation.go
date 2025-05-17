package service

import (
	"context"
	"conversation/db"
	conversation "conversation/generated/conversation"
	"errors"
)

func (s *Server) GetRecentConversation(ctx context.Context, req *conversation.GetRecentConversationsRequest) (*conversation.GetRecentConversationsResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	if req.UserName == "" {
		return nil, errors.New("username is empty")
	}
	if req.Limit <= 0 {
		return nil, errors.New("limit is invalid")
	}
	conversations, err := db.Rdb.ZRevRange(ctx, "recent_conversation:"+req.UserName, 0, int64(req.Limit)).Result()
	if err != nil {
		return nil, err
	}
	//TODO:这里需要考虑一下，如果没有会话，是否需要返回空的会话列表，同时如何将从zset获取的json string转为结构体返回
	response := conversation.GetRecentConversationsResponse{}
	transferredConversations, err := TransferConversationJsonStringToStruct(conversations)
	if err != nil {
		return nil, err
	}
	response.Conversations = transferredConversations
	response.ErrorCode = conversation.ConversationErrorCode_OK
	return &response, nil
}

// 将对话的JSON字符串转换为结构体
func TransferConversationJsonStringToStruct(conversations []string) ([]*conversation.Conversation, error) {
	// 返回空的结构体和错误
	return nil, nil
}
