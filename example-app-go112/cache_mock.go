package exampleservice

type CacheMock struct {
	MockGetItem func(key string) (*Entity, error)
	MockPutItem func(entity *Entity) error
}

func (cm CacheMock) GetItem(key string) (*Entity, error) {
	return cm.MockGetItem(key)
}

func (cm CacheMock) PutItem(entity *Entity) error {
	return cm.MockPutItem(entity)
}
