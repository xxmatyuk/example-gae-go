package exampleservice

import (
	"context"
)

type DBMock struct {
	MockPutEntity      func(ctx context.Context, key string, entity *Entity) error
	MockGetEntity      func(ctx context.Context, key string, entity *Entity) error
	MockGetAllEntities func(ctx context.Context, entities *[]Entity) error
}

func (dm DBMock) PutEntity(ctx context.Context, key string, entity *Entity) error {
	return dm.MockPutEntity(ctx, key, entity)
}

func (dm DBMock) GetEntity(ctx context.Context, key string, entity *Entity) error {
	return dm.MockGetEntity(ctx, key, entity)
}

func (dm DBMock) GetAllEntities(ctx context.Context, entities *[]Entity) error {
	return dm.MockGetAllEntities(ctx, entities)
}
