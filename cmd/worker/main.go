package main

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

func main() {
    ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	StartWorker(ctx, rdb, "job_queue")
}