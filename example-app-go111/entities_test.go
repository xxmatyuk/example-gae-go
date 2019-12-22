package exampleservice

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
)

var (
	testEntity = Entity{
		Key:   "foo",
		Value: "bar",
	}

	// DB Mocks
	okDBPut = func(ctx context.Context, key string, entity *Entity) error {
		return nil
	}
	errDBPut = func(ctx context.Context, key string, entity *Entity) error {
		return errors.New("fail")
	}
	okDBGet = func(ctx context.Context, key string, entity *Entity) (*Entity, error) {
		return &testEntity, nil
	}
	errDDGet = func(ctx context.Context, key string, entity *Entity) (*Entity, error) {
		return nil, errors.New("fail")
	}
	missDBGet = func(ctx context.Context, key string, entity *Entity) (*Entity, error) {
		return nil, datastore.ErrNoSuchEntity
	}
	okDBGetAll = func(ctx context.Context, entities *[]Entity) (*[]Entity, error) {
		return &[]Entity{testEntity}, nil
	}
	errDBGetAll = func(ctx context.Context, entities *[]Entity) (*[]Entity, error) {
		return nil, errors.New("fail")
	}

	// Cache mocks
	okCacheGet = func(ctx context.Context, key string) (*Entity, error) {
		return &testEntity, nil
	}
	errCacheGet = func(ctx context.Context, key string) (*Entity, error) {
		return nil, errors.New("fail")
	}
	missCacheGet = func(ctx context.Context, key string) (*Entity, error) {
		return nil, datastore.ErrNoSuchEntity
	}
	okCachePut = func(ctx context.Context, entity *Entity) error {
		return nil
	}
)

func TestEntities(t *testing.T) {
	testContext, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	testContext, _ = appengine.Namespace(testContext, "")
	defer done()

	testsGetEntity := []struct {
		name             string
		dbMock           DataStore
		cacheMock        Cache
		key              string
		expectedCode     int
		expectedResponse Entity
	}{
		{
			name: "happy-path",
			dbMock: DBMock{
				MockPutEntity:      okDBPut,
				MockGetEntity:      okDBGet,
				MockGetAllEntities: okDBGetAll,
			},
			cacheMock: CacheMock{
				MockGetItem: okCacheGet,
				MockPutItem: okCachePut,
			},
			key:              "foo",
			expectedCode:     http.StatusOK,
			expectedResponse: testEntity,
		},
		{
			name: "ok-cache-fail-db",
			dbMock: DBMock{
				MockPutEntity:      okDBPut,
				MockGetEntity:      errDDGet,
				MockGetAllEntities: okDBGetAll,
			},
			cacheMock: CacheMock{
				MockGetItem: okCacheGet,
				MockPutItem: okCachePut,
			},
			key:              "foo",
			expectedCode:     http.StatusOK,
			expectedResponse: testEntity,
		},
		{
			name: "cache-miss-ok-miss",
			dbMock: DBMock{
				MockPutEntity:      okDBPut,
				MockGetEntity:      okDBGet,
				MockGetAllEntities: okDBGetAll,
			},
			cacheMock: CacheMock{
				MockGetItem: missCacheGet,
				MockPutItem: okCachePut,
			},
			key:          "foo",
			expectedCode: http.StatusNotFound,
		},
		{
			name:   "fail-cache-ok-db",
			dbMock: okDBMock,
			cacheMock: CacheMock{
				MockGetItem: errCacheGet,
				MockPutItem: okCachePut,
			},
			key:          "foo",
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "fail-cache-fail-db",
			dbMock: DBMock{
				MockPutEntity:      okDBPut,
				MockGetEntity:      errDDGet,
				MockGetAllEntities: okDBGetAll,
			},
			cacheMock: CacheMock{
				MockGetItem: errCacheGet,
				MockPutItem: okCachePut,
			},
			key:          "foo",
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, test := range testsGetEntity {
		t.Run(test.name, func(t *testing.T) {

			s := &Service{
				db:    test.dbMock,
				cache: test.cacheMock,
			}

			url := fmt.Sprintf("/get-entity/%s", test.key)
			req, _ := http.NewRequest(http.MethodGet, url, &bytes.Buffer{})
			rr := httptest.NewRecorder()

			performGetEntityRequest(testContext, rr, req, s)

			// Status code check
			if statusCode := rr.Code; statusCode != test.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", statusCode, test.expectedCode)
			}

			// Body check
			var actualResponse Entity
			json.NewDecoder(rr.Body).Decode(&actualResponse)
			if actualResponse != test.expectedResponse {
				t.Errorf("handler returned wrong response code: got %v want %v", actualResponse, test.expectedResponse)
			}

		})
	}
}

func performGetEntityRequest(ctx context.Context, rr *httptest.ResponseRecorder, req *http.Request, s *Service) {

	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.GetEntityHandler(w, r)
	})

	handlerFunc.ServeHTTP(rr, req.WithContext(ctx))
}
