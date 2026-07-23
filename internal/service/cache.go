package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	TaskCacheTTl = 5 * time.Minute
	UserCacheTTl = 10 * time.Minute
)

type CacheService struct {
	client *redis.Client
}

func NewCacheService(client *redis.Client) *CacheService {
	return &CacheService{client: client}
}

func (cs *CacheService) Get(ctx context.Context, key string) (string, error) {
	if cs.client == nil {
		return "", redis.Nil
	}

	return cs.client.Get(ctx, key).Result()
}

func (cs *CacheService) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	if cs.client == nil {
		return nil
	}

	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return cs.client.Set(ctx, key, jsonValue, ttl).Err()
}

func (cs *CacheService) Delete(ctx context.Context, keys ...string) error {
	if cs.client == nil {
		return nil
	}

	return cs.client.Del(ctx, keys...).Err()
}
