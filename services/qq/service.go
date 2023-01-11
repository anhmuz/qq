package qq

import (
	"context"
	"qq/models"
	"qq/pkg/log"
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
}

var _ Service = service{}

func NewService(database qq.Database) (Service, error) {
	return service{
		database: database,
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
	return s.database.Get(key)
}

func (s service) GetAll(ctx context.Context) []models.Entity {
	log.Debug(ctx, "service: get all")
	return s.database.GetAll()
}
