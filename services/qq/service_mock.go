package qq

import (
	"context"
	"qq/models"
)

type serviceMock struct {
	AddMock    func(ctx context.Context, entity models.Entity) bool
	RemoveMock func(ctx context.Context, key string) bool
	GetMock    func(ctx context.Context, key string) *models.Entity
	GetAllMock func(ctx context.Context) []models.Entity
}

var _ Service = serviceMock{}

func NewServiceMock() serviceMock {
	return serviceMock{}
}

func (s serviceMock) Add(ctx context.Context, entity models.Entity) bool {
	return s.AddMock(ctx, entity)
}

func (s serviceMock) Remove(ctx context.Context, key string) bool {
	return s.RemoveMock(ctx, key)
}

func (s serviceMock) Get(ctx context.Context, key string) *models.Entity {
	return s.GetMock(ctx, key)
}

func (s serviceMock) GetAll(ctx context.Context) []models.Entity {
	return s.GetAllMock(ctx)
}
