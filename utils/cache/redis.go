package cache

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/rotisserie/eris"
	"time"
)

type RedisCache struct {
	redisClient *redis.Client
}

func NewRedisCache(address, password string) (rCache *RedisCache, err error) {
	rCache = new(RedisCache)
	rCache.redisClient = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	_, err = rCache.redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, eris.Wrap(err, "Error on ping redis")
	}

	return
}

func (r *RedisCache) Set(ctx context.Context, key string, value []byte, shortTime bool) error {
	var timeTTL time.Duration

	if shortTime {
		timeTTL = time.Hour * 24 * 7
	} else {
		timeTTL = 0
	}
	return eris.Wrap(r.redisClient.Set(ctx, key, value, timeTTL).Err(), "Error on set key")
}

func (r *RedisCache) Get(ctx context.Context, key string) (bytes []byte, err error) {
	if bytes, err = r.redisClient.Get(ctx, key).Bytes(); err != nil {
		if errors.Is(err, redis.Nil) {
			return []byte{}, ErrNoCache
		}
		return []byte{}, eris.Wrap(err, "Erro on get key value")
	}

	return
}

func (r *RedisCache) Down() error {
	return nil
}
