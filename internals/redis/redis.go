package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func NewRedisClient(address string) *Redis {
    return &Redis{
		Client: redis.NewClient(&redis.Options{
			Addr: address,
		}),
	}
}

func (rdb *Redis) EnqueueJob(ctx context.Context, jobName string, jobData []byte) error {
	err := rdb.Client.LPush(ctx, jobName, jobData).Err()
	if err != nil {
		return err	
	}
	return nil
}

func(rdb *Redis) DequeueJob(ctx context.Context, jobName string) ([]byte, error) {
	job, err := rdb.Client.RPop(ctx, jobName).Result()
	if err != nil {
		return nil, err
	}
	return []byte(job), nil
}