package config

import (
	"context"
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(redisConf RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisConf.Host + ":" + redisConf.Port,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return rdb, nil
}
