package qqclient

import (
	"context"
)

type Client interface {
	Add(ctx context.Context, entity Entity) (bool, error)
	Remove(ctx context.Context, key string) (bool, error)
	Get(ctx context.Context, key string) (*Entity, error)
	GetAsync(ctx context.Context, key string) (chan AsyncReply[*Entity], error)
	GetAll(ctx context.Context) ([]Entity, error)
}

type AsyncReply[Result any] struct {
	Result Result
	Err    error
}
