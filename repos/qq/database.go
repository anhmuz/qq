package qq

import "qq/models"

type Database interface {
	Add(entity models.Entity) bool
	Remove(key string) bool
	Get(key string) *models.Entity
	GetAll() []models.Entity
}

type database struct {
	entities map[string]models.Entity
}

var _ Database = &database{}

func NewDatabase() (Database, error) {
	return &database{
		entities: map[string]models.Entity{},
	}, nil
}

func (d *database) Add(entity models.Entity) bool {
	d.entities[entity.Key] = entity

	return true
}

func (d *database) Remove(key string) bool {
	delete(d.entities, key)

	return true
}

func (d *database) Get(key string) *models.Entity {
	entity, present := d.entities[key]
	if !present {
		return nil
	}
	return &entity
}

func (d *database) GetAll() []models.Entity {
	entities := make([]models.Entity, 0, len(d.entities))
	for _, entity := range d.entities {
		entities = append(entities, entity)
	}

	return entities
}
