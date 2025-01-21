package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(addr string) (*RedisRepository, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return &RedisRepository{client: rdb}, nil
}

func (r *RedisRepository) Set(ctx context.Context, key, val string) error {
	return r.client.Set(ctx, key, val, 0).Err()
}

func (r *RedisRepository) Get(ctx context.Context, key string) (string, error) {
	res, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return res, nil
}
