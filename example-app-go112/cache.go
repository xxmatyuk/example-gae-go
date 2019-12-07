package exampleservice

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type CacheClient interface {
	GetItem(key string) (*Entity, error)
	PutItem(entity *Entity) error
}

type cachedDB struct {
	redisClient *redis.Client
	dbClient    DataStoreClient
}

func NewCacheDBClient(dbClient DataStoreClient, redisHost string, redisPort string, db string) CacheClient {

	idb, err := strconv.Atoi(db)
	if err != nil {
		panic(fmt.Sprintf("Error while Redis client init: %s", err.Error()))
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password:   "",
		DB:         idb,
		MaxRetries: 2,
	})

	return &cachedDB{
		redisClient: redisClient,
		dbClient:    dbClient,
	}
}

func (c *cachedDB) GetItem(key string) (*Entity, error) {

	status := c.redisClient.Get(key)
	_, err := status.Result()
	if err == nil && status.Val() != "" {
		return &Entity{
			Key:   key,
			Value: string(status.Val()),
		}, nil
	}

	var e Entity
	if err = c.dbClient.GetEntity(key, &e); err != nil {
		return nil, err
	}

	return &e, c.PutItem(&e)
}

func (c *cachedDB) PutItem(entity *Entity) error {

	status := c.redisClient.Set(entity.Key, entity.Value, 15*time.Minute)
	_, err := status.Result()

	return err
}
