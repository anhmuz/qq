package redisqq

import (
	"context"
	"fmt"
	"qq/models"
	"qq/pkg/rabbitqq"

	"github.com/redis/go-redis/v9"
)

type Database interface {
	Get(ctx context.Context, key string) (*models.Entity, error)
	Set(ctx context.Context, entity *models.Entity) error
}

type database struct {
	client *redis.Client
}

var _ Database = database{}

func NewDatabase() (Database, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     rabbitqq.RedisServerAddr,
		Password: "",
		DB:       0,
	})

	return database{
		client: client,
	}, nil
}

func (d database) Get(ctx context.Context, key string) (*models.Entity, error) {
	value, err := d.client.Get(ctx, key).Result()

	if err == redis.Nil {
		return nil, fmt.Errorf("redis client: key %s does not exist", key)
	}
	if err != nil {
		return nil, fmt.Errorf("redis client: failed to get key %s: %w", key, err)
	}

	return &models.Entity{Key: key, Value: value}, nil
}

func (d database) Set(ctx context.Context, entity *models.Entity) error {
	err := d.client.Set(ctx, entity.Key, entity.Value, 0).Err()
	if err != nil {
		return fmt.Errorf("redis client: failed to set key %s, value %s: %w", entity.Key, entity.Value, err)
	}

	return nil
}
