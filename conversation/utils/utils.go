package utils

import (
	"context"
	"conversation/db"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func CheckUsersValid(users []string) bool {
	if len(users) == 0 {
		return false
	}
	//TODO:像用户这种数据应该放在mysql 还是 redis？
	for _, user := range users {
		if user == "" {
			return false
		}
	}
	return true
}

func GenerateConversationID(u1 string, u2 string) string {
	if u1 < u2 {
		return u1 + "_" + u2
	}
	return u2 + "_" + u1
}

func CreateEmptyZset(conversationID string) error {
	ctx := context.Background()
	key := "chat:" + conversationID

	// 检查键是否存在，如果不存在则创建空ZSet
	//TODO：其实这里应该直接用exist创建字段，之后修改
	exists, err := db.Rdb.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to check key existence: %v", err)
	}

	if exists == 0 {
		// 添加一个虚拟成员然后立即删除，确保创建了空ZSet
		err := db.Rdb.ZAdd(ctx, key, redis.Z{Score: 0, Member: "__temp__"}).Err()
		if err != nil {
			return fmt.Errorf("failed to initialize ZSet: %v", err)
		}
		_, err = db.Rdb.ZRem(ctx, key, "__temp__").Result()
		if err != nil {
			return fmt.Errorf("failed to remove temp member: %v", err)
		}
	}

	return nil
}

func GenerateGroupConversationID() (string, error) {
	ctx := context.Background()
	// 使用 Redis 的自增命令生成唯一 ID
	id, err := db.Rdb.Incr(ctx, "global:conversation_id").Result()
	if err != nil {
		return "", fmt.Errorf("failed to generate conversation ID: %v", err)
	}
	return fmt.Sprintf("%d", id), nil
}

func UpdateUserRecentConversations(userID []string, conversationID string) error {

	ctx := context.Background()
	timestamp := float64(time.Now().Unix()) // 使用当前时间戳作为 Score

	// 使用tx开始pipeline事务
	tx := db.Rdb.TxPipeline()
	// 添加消息到 ZSet recent_conversation list
	for _, uID := range userID {
		tx.ZAdd(ctx, "rc:"+uID, redis.Z{
			Score:  timestamp,
			Member: conversationID,
		})
	}
	//TODO: 应该加个补偿？？？？如果更新失败是重试还是丢入某个消费队列上？
	_, err := tx.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update recentConversationList to ZSet: %v", err)
	}

	return nil

}
