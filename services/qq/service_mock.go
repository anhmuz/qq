package qq

import (
	"context"
	"qq/models"
)

type ServiceMock struct {
	AddMock    func(ctx context.Context, entity models.Entity) bool
	RemoveMock func(ctx context.Context, key string) bool
	GetMock    func(ctx context.Context, key string) *models.Entity
	GetAllMock func(ctx context.Context) []models.Entity
}

var _ Service = ServiceMock{}

func (s ServiceMock) Add(ctx context.Context, entity models.Entity) bool {
	return s.AddMock(ctx, entity)
}

func (s ServiceMock) Remove(ctx context.Context, key string) bool {
	return s.RemoveMock(ctx, key)
}

func (s ServiceMock) Get(ctx context.Context, key string) *models.Entity {
	return s.GetMock(ctx, key)
}

func (s ServiceMock) GetAll(ctx context.Context) []models.Entity {
	return s.GetAllMock(ctx)
}
