package exampleservice

import (
	"context"
	"time"

	"google.golang.org/appengine/memcache"
)

type Cache interface {
	GetItem(ctx context.Context, key string) (*Entity, error)
	PutItem(ctx context.Context, entity *Entity) error
}

type cachedDB struct {
	DB DataStore
}

func NewCacheDB(db DataStore) Cache {
	return &cachedDB{
		DB: db,
	}
}

func (c *cachedDB) GetItem(ctx context.Context, key string) (*Entity, error) {

	item, err := memcache.Get(ctx, key)
	if err == nil {
		return &Entity{
			Key:   key,
			Value: string(item.Value),
		}, nil
	}

	var e Entity
	if _, err = c.DB.GetEntity(ctx, key, &e); err != nil {
		return nil, err
	}

	return &e, c.PutItem(ctx, &e)
}

func (c *cachedDB) PutItem(ctx context.Context, entity *Entity) error {
	if err := memcache.Set(ctx, &memcache.Item{
		Key:        entity.Key,
		Value:      []byte(entity.Value),
		Expiration: 15 * time.Minute,
	}); err != nil {
		return err
	}
	return nil
}
