package init

import (
	"context"
	"fmt"
	"log"
	"time"

	"blog/internal/global"

	"github.com/go-redis/redis/v8"
)

// InitRedis 初始化Redis连接
func InitRedis() error {
	cfg := global.Config
	if cfg == nil {
		return fmt.Errorf("config not initialized")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	global.Redis = rdb
	log.Println("Redis initialized successfully")
	return nil
}
