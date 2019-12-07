package exampleservice

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
)

type Entity struct {
	Key   string `json:"key" datastore:"key"`
	Value string `json:"value,omitempty" datastore:"value"`
}

type DataStoreClient interface {
	PutEntity(key string, entity *Entity) error
	GetEntity(key string, entity *Entity) error
	GetAllEntities() ([]Entity, error)
}

type ds struct {
	client *datastore.Client
	kind   string
}

func NewDataStoreClient(projectID, kind string) DataStoreClient {

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)

	if err != nil {
		panic(fmt.Sprintf("Error while datastore client init: %s", err.Error()))
	}

	return &ds{
		kind:   kind,
		client: client,
	}
}

func (ds *ds) GetEntity(key string, entity *Entity) error {
	ctx := context.Background()
	k := datastore.NameKey(ds.kind, key, nil)
	return ds.client.Get(ctx, k, &entity)
}

func (ds *ds) PutEntity(key string, entity *Entity) error {
	ctx := context.Background()
	k := datastore.NameKey(ds.kind, key, nil)
	_, err := ds.client.Put(ctx, k, entity)
	return err
}

func (ds *ds) GetAllEntities() ([]Entity, error) {
	ctx := context.Background()
	var entities []Entity
	_, err := ds.client.GetAll(ctx, datastore.NewQuery(ds.kind), &entities)
	return entities, err
}
