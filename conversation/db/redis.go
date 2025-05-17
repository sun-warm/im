package db

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func InitRedis() error {
	// 创建 Redis 客户端
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 地址
		Password: "",               // 密码
		DB:       0,                // 使用默认 DB
	})
	defer Rdb.Close()
	ctx := context.Background()

	// 测试连接
	_, err := Rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}
