package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

// RedisClient 全局 Redis 客户端实例
var RedisClient *redis.Client

// ConnectRedis 连接 Redis 并进行初始化检查。
// 参数：
//
//	host     Redis 主机地址，例如 "127.0.0.1"
//	port     Redis 服务端口，例如 6379
//	password Redis 密码，如果没有则传空字符串
//	db       要使用的 Redis 数据库序号，默认 0
func ConnectRedis(host string, port int, password string, db int) {
	// 1. 构造地址
	addr := fmt.Sprintf("%s:%d", host, port)

	// 2. 创建客户端
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// 3. 带超时的 Context，用于 Ping 健康检查
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 4. 尝试 PING，确保 Redis 可用
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("尝试连接 Redis 失败 (%s): %v", addr, err)
	}

	// 5. 全局保存客户端实例
	RedisClient = client

	fmt.Println("成功连接至 Redis:", addr)
}
