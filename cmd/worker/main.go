package main

import (
	"auth/internals/db"
	"auth/internals/redis"
	"auth/internals/store"
	"context"
	"fmt"
	"os"
)

func main() {
    ctx := context.Background()
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb := redis.NewRedisClient(redisAddr)
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}
	db, err := db.New(dbURL, 10, 5, "5m")
	if err != nil {
		fmt.Println("Error connecting to database: ", err)
		return
	}
	store := store.NewStore(db)
	
	StartWorker(ctx, rdb, "job_queue",store)
}