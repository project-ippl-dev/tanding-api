package config

import (
	"context"
	"github.com/go-redis/redis/v8"
)

func RedisConnection() (*redis.Client, error) {
	redisConf := Configuration().Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: redisConf.Host + ":" + redisConf.Port,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return rdb, nil
}
