package utils

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"message/db"
	"message/generated/message"
)

func GenerateConversationID(userName, receiver string) string {
	if userName < receiver {
		return userName + receiver
	}
	return receiver + userName
}

func GetMessageFromZset(conversationID string, cursor int64, limit int64) ([]string, error) {
	ctx := context.Background()

	// 按数量从ZSet获取消息
	messages, err := db.Rdb.ZRevRange(ctx, conversationID, cursor, limit-1).Result()
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

func MarshalMessageStringtoMessage(messagesString []string) ([]*message.Message, error) {
	fmt.Println(len(messagesString))
	messages := make([]*message.Message, len(messagesString))
	for i, messageStr := range messagesString {
		message := message.Message{}
		err := sonic.UnmarshalString(messageStr, &message)
		if err != nil {
			fmt.Println("failed to unmarshal message:", err, messageStr)
			return nil, err
		}
		messages[i] = &message
	}
	return messages, nil
}
