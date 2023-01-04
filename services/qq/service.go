package qq

import (
	"context"
	"log"
	"qq/models"
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
	log.Printf("service: add %+v\n", entity)
	return s.database.Add(entity)
}

func (s service) Remove(ctx context.Context, key string) bool {
	log.Printf("service: remove %s\n", key)
	return s.database.Remove(key)
}

func (s service) Get(ctx context.Context, key string) *models.Entity {
	log.Printf("service: get %s\n", key)
	return s.database.Get(key)
}

func (s service) GetAll(ctx context.Context) []models.Entity {
	log.Println("service: get all")
	return s.database.GetAll()
}
