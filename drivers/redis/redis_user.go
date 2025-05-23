package redis

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient() *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		Username: "default",
		DB:       0,
	})

	return &RedisClient{
		Client: rdb,
	}
}

func (r *RedisClient) Set(key string, value string, expiration time.Duration) error {
	err := r.Client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return err
}

func (r *RedisClient) Get(key string) (string, error) {
	result, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return result, err
}

func (r *RedisClient) Del(key string) error {
	err := r.Client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return err
}

func (r *RedisClient) GetRecommendationProductsIds(key string) (*[]string, error) {
	ctx := context.Background()

	flag, err := r.Client.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if flag == 0 {
		return nil, redis.Nil
	}

	result, err := r.Client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *RedisClient) SetRecommendationProductsIds(key string, value []string) error {
	ctx := context.Background()

	err := r.Client.RPush(ctx, key, value).Err()
	if err != nil {
		return err
	}

	err = r.Client.Expire(ctx, key, time.Hour * 24).Err()
	if err != nil {
		return err
	}

	return err
}