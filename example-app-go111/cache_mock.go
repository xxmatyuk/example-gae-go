package exampleservice

import (
	"context"
)

type CacheMock struct {
	MockGetItem func(ctx context.Context, key string) (*Entity, error)
	MockPutItem func(ctx context.Context, entity *Entity) error
}

func (cm CacheMock) GetItem(ctx context.Context, key string) (*Entity, error) {
	return cm.MockGetItem(ctx, key)
}

func (cm CacheMock) PutItem(ctx context.Context, entity *Entity) error {
	return cm.MockPutItem(ctx, entity)
}
