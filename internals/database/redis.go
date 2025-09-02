package database

import (
	"context"
	"log"

	"github.com/amarjeet-choudhary666/ai_resume_screener/internals/config"
	"github.com/redis/go-redis/v9"
)

var (
	Ctx         = context.Background()
	RedisClient *redis.Client
)

func RedisConnection(cfg config.Config) *redis.Client {
	if cfg.RedisURL == "" {
		log.Fatalf("❌ REDIS_URL not set in environment")
	}

	opt, err := redis.ParseURL(cfg.RedisURL)

	if err != nil {
		log.Fatalf("Could not parse Redis URL: %v", err)
	}
	RedisClient = redis.NewClient(opt)

	if _, err := RedisClient.Ping(Ctx).Result(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	} else {
		log.Println("Redis connection established")
	}

	return RedisClient
}
