package qq

import (
	"log"
	"qq/models"
	"qq/repos/qq"
)

type Service interface {
	Add(entity models.Entity) bool
	Remove(key string) bool
	Get(key string) *models.Entity
	GetAll() []models.Entity
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

func (s service) Add(entity models.Entity) bool {
	log.Printf("service: add %+v\n", entity)
	return s.database.Add(entity)
}

func (s service) Remove(key string) bool {
	log.Printf("service: remove %s\n", key)
	return s.database.Remove(key)
}

func (s service) Get(key string) *models.Entity {
	log.Printf("service: get %s\n", key)
	return s.database.Get(key)
}

func (s service) GetAll() []models.Entity {
	log.Println("service: get all")
	return s.database.GetAll()
}
