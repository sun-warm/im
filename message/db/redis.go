package db

import (
	"context"
	"fmt"

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
	//FIXME:这里面如果defer close会导致当InitRedis后直接关闭连接，导致后续无法请求。
	//defer Rdb.Close()
	ctx := context.Background()

	// 测试连接
	res, err := Rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}
	fmt.Println("Redis connection result:", res)
	return nil
}
