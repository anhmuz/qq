package qqclient

import (
	"context"
	"qq/pkg/protocol"
)

type Client interface {
	Add(ctx context.Context, entity protocol.Entity) (bool, error)
	Remove(ctx context.Context, key string) (bool, error)
	Get(ctx context.Context, key string) (*protocol.Entity, error)
	GetAsync(ctx context.Context, key string) (chan AsyncReply[*protocol.Entity], error)
	GetAll(ctx context.Context) ([]protocol.Entity, error)
}

type AsyncReply[Result any] struct {
	Result Result
	Err    error
}
