package main

import (
	"auth/internals/redis"
	"context"
	"os"
)

func main() {
    ctx := context.Background()
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewRedisClient(redisAddr)
	StartWorker(ctx, rdb, "job_queue")
}