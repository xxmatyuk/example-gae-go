package exampleservice

type DBMock struct {
	MockPutEntity      func(key string, entity *Entity) error
	MockGetEntity      func(key string) (*Entity, error)
	MockGetAllEntities func() (*[]Entity, error)
}

func (dm DBMock) PutEntity(key string, entity *Entity) error {
	return dm.MockPutEntity(key, entity)
}

func (dm DBMock) GetEntity(key string) (*Entity, error) {
	return dm.MockGetEntity(key)
}

func (dm DBMock) GetAllEntities() (*[]Entity, error) {
	return dm.MockGetAllEntities()
}
