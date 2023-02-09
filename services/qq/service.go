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

	s.database.Add(entity)

	cachedEntity, err := s.cache.GetEntity(ctx, entity.Key)
	_ = cachedEntity

	if err != nil {
		return true
	}

	err = s.cache.SetEntity(ctx, entity.Key, &entity)
	if err == nil {
		log.Debug(ctx, "set to cache", log.Args{"key": entity.Key, "entity": entity})
	} else {
		log.Warning(ctx, "failed to set to cache", log.Args{"error": err})
	}

	return true
}

func (s service) Remove(ctx context.Context, key string) bool {
	log.Debug(ctx, "service: remove", log.Args{"key": key})

	err := s.cache.DeleteEntity(ctx, key)
	if err != nil {
		log.Warning(ctx, "failed to delete from cache", log.Args{"error": err})
	}

	return s.database.Remove(key)
}

func (s service) Get(ctx context.Context, key string) *models.Entity {
	log.Debug(ctx, "service: get", log.Args{"key": key})

	entity, err := s.cache.GetEntity(ctx, key)

	log.Debug(ctx, "get from cache", log.Args{"key": key, "entity": entity, "error": err})

	if err == nil {
		return entity
	}

	entity = s.database.Get(key)

	err = s.cache.SetEntity(ctx, key, entity)
	if err == nil {
		log.Debug(ctx, "set to cache", log.Args{"key": key, "entity": entity})
	} else {
		log.Warning(ctx, "failed to set to cache", log.Args{"error": err})
	}

	return entity
}

func (s service) GetAll(ctx context.Context) []models.Entity {
	log.Debug(ctx, "service: get all")
	return s.database.GetAll()
}
