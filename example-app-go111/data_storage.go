package exampleservice

import (
	"context"

	"google.golang.org/appengine/datastore"
)

type DataStore interface {
	PutEntity(ctx context.Context, key string, entity *Entity) error
	GetEntity(ctx context.Context, key string, entity *Entity) (*Entity, error)
	GetAllEntities(ctx context.Context, entities *[]Entity) (*[]Entity, error)
}

type ds struct {
	kind string
}

func NewDataStoreClient(kind string) DataStore {
	return &ds{
		kind: kind,
	}
}

func (ds *ds) GetEntity(ctx context.Context, key string, entity *Entity) (*Entity, error) {
	k := datastore.NewKey(ctx, ds.kind, key, 0, nil)
	err := datastore.Get(ctx, k, entity)
	return entity, err
}

func (ds *ds) PutEntity(ctx context.Context, key string, entity *Entity) error {
	k := datastore.NewKey(ctx, ds.kind, key, 0, nil)
	_, err := datastore.Put(ctx, k, entity)
	return err
}

func (ds *ds) GetAllEntities(ctx context.Context, entities *[]Entity) (*[]Entity, error) {
	_, err := datastore.NewQuery(ds.kind).GetAll(ctx, entities)
	return entities, err
}
