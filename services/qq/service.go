package qq

import (
	"context"
	"qq/models"
	"qq/pkg/log"
	"qq/repos/qq"
	"qq/repos/redisqq"
)

type Service interface {
	Add(ctx context.Context, entity models.Entity) bool
	Remove(ctx context.Context, key string) bool
	Get(ctx context.Context, key string) *models.Entity
	GetAll(ctx context.Context) []models.Entity
}

type service struct {
	database      qq.Database
	redisDatabase redisqq.Database
}

var _ Service = service{}

func NewService(database qq.Database, redisDatabase redisqq.Database) (Service, error) {
	return service{
		database:      database,
		redisDatabase: redisDatabase,
	}, nil
}

func (s service) Add(ctx context.Context, entity models.Entity) bool {
	log.Debug(ctx, "service: add", log.Args{"entity": entity})
	return s.database.Add(entity)
}

func (s service) Remove(ctx context.Context, key string) bool {
	log.Debug(ctx, "service: remove", log.Args{"key": key})
	return s.database.Remove(key)
}

func (s service) Get(ctx context.Context, key string) *models.Entity {
	log.Debug(ctx, "service: get", log.Args{"key": key})

	entity, err := s.redisDatabase.Get(ctx, key)
	if err == nil {
		log.Debug(ctx, "service: get from redis cache", log.Args{"key": key})
		return entity
	}
	log.Debug(ctx, "service: failed to get from redis cache", log.Args{"error": err})

	entity = s.database.Get(key)

	if entity == nil {
		return nil
	}

	err = s.redisDatabase.Set(ctx, entity)
	if err == nil {
		log.Debug(ctx, "service: set to redis cache", log.Args{"entity": entity})
	} else {
		log.Debug(ctx, "service: failed to set to redis cache", log.Args{"error": err})
	}

	return entity
}

func (s service) GetAll(ctx context.Context) []models.Entity {
	log.Debug(ctx, "service: get all")
	return s.database.GetAll()
}
