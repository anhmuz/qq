package qq

import (
	"context"
	"qq/models"
	"qq/pkg/log"
	"qq/repos/cacheqq"
	"qq/repos/qq"
)

type Service interface {
	Add(ctx context.Context, entity models.Entity) bool
	Remove(ctx context.Context, key string) bool
	Get(ctx context.Context, key string) *models.Entity
	GetAll(ctx context.Context) []models.Entity
}

type service struct {
	database qq.Database
	cache    cacheqq.Cache
}

var _ Service = service{}

func NewService(database qq.Database, cache cacheqq.Cache) (Service, error) {
	return service{
		database: database,
		cache:    cache,
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

	entity, err := s.cache.GetEntity(ctx, key)

	if entity != nil {
		log.Debug(ctx, "service: get from cache", log.Args{"key": key})
		return entity
	}

	if err == nil {
		log.Warning(ctx, "service: key does not exist in cache", log.Args{"key": key})
	} else {
		log.Error(ctx, "service: failed to get from cache", log.Args{"error": err})
	}

	entity = s.database.Get(key)

	if entity == nil {
		return nil
	}

	err = s.cache.SetEntity(ctx, entity)
	if err == nil {
		log.Debug(ctx, "service: set to cache", log.Args{"entity": entity})
	} else {
		log.Error(ctx, "service: failed to set to cache", log.Args{"error": err})
	}

	return entity
}

func (s service) GetAll(ctx context.Context) []models.Entity {
	log.Debug(ctx, "service: get all")
	return s.database.GetAll()
}
