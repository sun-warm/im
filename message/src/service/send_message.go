package service

import (
	"context"
	"fmt"
	"message/client"
	"message/db"
	"message/generated/message"
	"message/generated/push_service"
	"message/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
1、写入消息入conversation
2、推送消息给接收者
3、更新发送者和接收者的最近会话列表
*/

// FIXME:如果是单聊的话，其实可以当A向B发消息后，再建立一个Conversation
func (s *Server) SendMessage(ctx context.Context, req *message.SendMessageRequest) (*message.SendMessageResponse, error) {
	fmt.Println(1111)

	if req == nil || req.Message == nil || req.Receiver == "" {
		return nil, fmt.Errorf("request is wrong")
	}

	conversationID := utils.GenerateConversationID(req.UserName, req.Receiver)
	//1、写入消息入conversation
	if err := AddMessageToZSet(conversationID, req.Content); err != nil {
		return nil, err
	}

	// 2、推送消息给接收者
	_, err := client.PushServiceClient.Client.PushMessage(ctx, &push_service.PushMessageRequest{UserName: req.Receiver, Content: req.Content})
	if err != nil {
		//TODO:记录错误， 具体处理是客户端延迟拉取消息 or 如何处理？
		fmt.Println(err.Error())
		return nil, err
	}
	fmt.Println(1111)
	// 3、更新发送者和接收者的最近会话列表
	if err := UpdateRecentConversation(req.UserName, req.Receiver, conversationID); err != nil {
		return nil, err
	}
	return &message.SendMessageResponse{Content: "Hello "}, nil
}

func AddMessageToZSet(conversationID, message string) error {
	ctx := context.Background()
	timestamp := float64(time.Now().Unix()) // 使用当前时间戳作为 Score

	// 添加消息到 ZSet
	err := db.Rdb.ZAdd(ctx, "chat:"+conversationID, redis.Z{
		Score:  timestamp,
		Member: message,
	}).Err()
	if err != nil {
		return fmt.Errorf("failed to add message to ZSet: %v", err)
	}
	return nil
}

// 考虑加pipeline，还是加事务pipeline
func UpdateRecentConversation(sender, receiver, conversationID string) error {
	ctx := context.Background()
	timestamp := float64(time.Now().Unix()) // 使用当前时间戳作为 Score

	// 使用tx开始pipeline事务
	tx := db.Rdb.TxPipeline()
	// 添加消息到 ZSet recent_conversation list
	tx.ZAdd(ctx, "rc:"+sender, redis.Z{
		Score:  timestamp,
		Member: conversationID,
	})
	tx.ZAdd(ctx, "rc:"+receiver, redis.Z{
		Score:  timestamp,
		Member: conversationID,
	})
	//TODO: 应该加个补偿？？？？如果更新失败是重试还是丢入某个消费队列上？
	_, err := tx.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update recentConversationList to ZSet: %v", err)
	}

	return nil
}
