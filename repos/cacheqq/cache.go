package cacheqq

import (
	"context"
	"fmt"
	"qq/models"
	"qq/pkg/rabbitqq"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	GetEntity(ctx context.Context, key string) (*models.Entity, error)
	SetEntity(ctx context.Context, key string, entity *models.Entity) error
	DeleteEntity(ctx context.Context, key string) error
}

type cache struct {
	redisClient *redis.Client
}

var _ Cache = cache{}

func NewRedisCache() Cache {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     rabbitqq.RedisServerAddr,
		Password: "",
		DB:       0,
	})

	return cache{
		redisClient: redisClient,
	}
}

func (c cache) GetEntity(ctx context.Context, key string) (*models.Entity, error) {
	value, err := c.redisClient.Get(ctx, key).Result()

	if err == redis.Nil {
		return nil, fmt.Errorf("key %s does not exist", key)

	}
	if err != nil {
		return nil, fmt.Errorf("failed to get key %s: %w", key, err)
	}

	if value == "" {
		return nil, nil
	}

	return &models.Entity{Key: key, Value: value}, nil
}

func (c cache) SetEntity(ctx context.Context, key string, entity *models.Entity) error {
	var value string

	if entity != nil {
		value = entity.Value
	}

	err := c.redisClient.Set(ctx, key, value, 10*time.Minute).Err()

	if err != nil {
		return fmt.Errorf("failed to set key %s, value %s", key, value)
	}

	return nil
}

func (c cache) DeleteEntity(ctx context.Context, key string) error {
	err := c.redisClient.Del(ctx, key).Err()

	if err != nil {
		return fmt.Errorf("failed to delete key %s", key)
	}

	return nil
}
