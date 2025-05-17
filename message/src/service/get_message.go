package service

import (
	"context"
	"fmt"
	"message/db"
	"message/generated/message"

	"github.com/redis/go-redis/v9"
)

func (s *Server) GetMessage(ctx context.Context, req *message.GetMessageRequest) (*message.GetMessageResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request is nil")
	}
	if req.UserName == "" {
		return nil, fmt.Errorf("username is empty")
	}

	conversationID := req.ConversationId
	var messages []string
	messages, err := GetMessageFromZset(conversationID, req.Cursor, req.Limit)
	if err != nil {
		return nil, err
	}
	fmt.Println("messages:", messages)
	return &message.GetMessageResponse{}, nil
}

func GetMessageFromZset(conversationID string, cursor int64, limit int64) ([]string, error) {
	ctx := context.Background()

	// 按数量从ZSet获取消息
	messages, err := db.Rdb.ZRevRange(ctx, "chat:"+conversationID, cursor, limit-1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get recent messages: %v", err)
	}

	return messages, nil
}

func GetMessageFromZsetByTime(conversationID string, startTime float64, endTime float64) ([]string, error) {
	ctx := context.Background()

	// 按数量从ZSet获取消息
	messages, err := db.Rdb.ZRangeByScore(ctx, "chat:"+conversationID, &redis.ZRangeBy{
		Min: fmt.Sprintf("%f", startTime),
		Max: fmt.Sprintf("%f", endTime),
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get recent messages: %v", err)
	}
	return messages, nil
}
