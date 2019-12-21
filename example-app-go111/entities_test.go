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
)

func TestEntities(t *testing.T) {
	testContext, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	testContext, _ = appengine.Namespace(testContext, "")
	defer done()

	testEntity := Entity{
		Key:   "foo",
		Value: "bar",
	}

	okDBMock := DBMock{
		MockPutEntity: func(ctx context.Context, key string, entity *Entity) error {
			return nil
		},
		MockGetEntity: func(ctx context.Context, key string, entity *Entity) (*Entity, error) {
			return &testEntity, nil
		},
		MockGetAllEntities: func(ctx context.Context, entities *[]Entity) (*[]Entity, error) {
			return &[]Entity{testEntity}, nil
		},
	}

	errDBMock := DBMock{
		MockPutEntity: func(ctx context.Context, key string, entity *Entity) error {
			return errors.New("error")
		},
		MockGetEntity: func(ctx context.Context, key string, entity *Entity) (*Entity, error) {
			return nil, errors.New("error")
		},
		MockGetAllEntities: func(ctx context.Context, entities *[]Entity) (*[]Entity, error) {
			return nil, errors.New("error")
		},
	}

	okCacheMock := CacheMock{
		MockGetItem: func(ctx context.Context, key string) (*Entity, error) {
			return &testEntity, nil
		},
		MockPutItem: func(ctx context.Context, entity *Entity) error {
			return nil
		},
	}

	errCacheMock := CacheMock{
		MockGetItem: func(ctx context.Context, key string) (*Entity, error) {
			return nil, errors.New("error")
		},
		MockPutItem: func(ctx context.Context, entity *Entity) error {
			return errors.New("error")
		},
	}

	testsGetEntity := []struct {
		name             string
		dbMock           DataStore
		cacheMock        Cache
		key              string
		expectedCode     int
		expectedResponse Entity
	}{
		{
			name:             "happy-path",
			dbMock:           okDBMock,
			cacheMock:        okCacheMock,
			key:              "foo",
			expectedCode:     http.StatusOK,
			expectedResponse: testEntity,
		},
		{
			name:             "ok-cache-fail-db",
			dbMock:           errDBMock,
			cacheMock:        okCacheMock,
			key:              "foo",
			expectedCode:     http.StatusOK,
			expectedResponse: testEntity,
		},
		{
			name:         "fail-cache-ok-db",
			dbMock:       okDBMock,
			cacheMock:    errCacheMock,
			key:          "foo",
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "fail-cache-fail-db",
			dbMock:       errDBMock,
			cacheMock:    errCacheMock,
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

			performGetEntityRequest(t, testContext, rr, req, s)

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

func performGetEntityRequest(t *testing.T, ctx context.Context, rr *httptest.ResponseRecorder, req *http.Request, s *Service) {

	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.GetEntityHandler(w, r)
	})

	handlerFunc.ServeHTTP(rr, req.WithContext(ctx))
}
