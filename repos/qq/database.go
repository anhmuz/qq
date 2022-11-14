package qq

import "qq/models"

type Database interface {
	Add(entity models.Entity) bool
	Remove(key string) bool
	Get(key string) *models.Entity
	GetAll() []models.Entity
}

type database struct {
	entities []models.Entity
}

var _ Database = database{}

func NewDatabase() (Database, error) {
	return database{
		entities: []models.Entity{},
	}, nil
}

func (d database) Add(entity models.Entity) bool {
	d.entities = append(d.entities, entity)

	return true
}

func (d database) Remove(key string) bool {
	for i, entity := range d.entities {
		if entity.Key == key {
			d.entities = append(d.entities[:i], d.entities[i+1:]...)
			return true
		}
	}

	return false
}

func (d database) Get(key string) *models.Entity {
	for _, entity := range d.entities {
		if entity.Key == key {
			return &entity
		}
	}
	return nil
}

func (d database) GetAll() []models.Entity {
	return d.entities
}
