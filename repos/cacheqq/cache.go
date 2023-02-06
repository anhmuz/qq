package cacheqq

import (
	"context"
	"fmt"
	"qq/models"
	"qq/pkg/rabbitqq"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	GetEntity(ctx context.Context, key string) (*models.Entity, error)
	SetEntity(ctx context.Context, entity *models.Entity) error
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
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get key %s: %w", key, err)
	}

	return &models.Entity{Key: key, Value: value}, nil
}

func (c cache) SetEntity(ctx context.Context, entity *models.Entity) error {
	err := c.redisClient.Set(ctx, entity.Key, entity.Value, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s, value %s: %w", entity.Key, entity.Value, err)
	}

	return nil
}
