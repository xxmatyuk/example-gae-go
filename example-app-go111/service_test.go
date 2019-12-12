package exampleservice

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
)

func TestService(t *testing.T) {
	testContext, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	testContext, _ = appengine.Namespace(testContext, "")
	defer done()

	// DB mocks
	okDBMock := DBMock{
		MockPutEntity: func(ctx context.Context, key string, entity *Entity) error {
			return nil
		},
		MockGetEntity: func(ctx context.Context, key string, entity *Entity) error {
			return nil
		},
		MockGetAllEntities: func(ctx context.Context, entities *[]Entity) error {
			return nil
		},
	}

	// Cache mocks
	okCacheMock := CacheMock{
		MockGetItem: func(ctx context.Context, key string) (*Entity, error) {
			return &Entity{}, nil
		},
		MockPutItem: func(ctx context.Context, entity *Entity) error {
			return nil
		},
	}

	// Define tests
	tests := []struct {
		name             string
		dbMock           DataStore
		cacheMock        Cache
		expectedCode     int
		expectedResponse Response
	}{
		{
			name:             "Happy-path",
			dbMock:           okDBMock,
			cacheMock:        okCacheMock,
			expectedCode:     http.StatusOK,
			expectedResponse: Response{"OK"},
		},
	}

	// Run positive tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			s := &Service{
				db:    test.dbMock,
				cache: test.cacheMock,
			}

			req, _ := http.NewRequest(http.MethodGet, "/_ah/warmup", &bytes.Buffer{})
			rr := httptest.NewRecorder()

			performWarmUpRequest(t, testContext, rr, req, s)

			// Status code check
			if statusCode := rr.Code; statusCode != test.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", statusCode, test.expectedCode)
			}

			// Body check
			var actualResponse Response
			json.NewDecoder(rr.Body).Decode(&actualResponse)
			if actualResponse != test.expectedResponse {
				t.Errorf("handler returned wrong response code: got %v want %v", actualResponse, test.expectedResponse)
			}

		})
	}

}

func performWarmUpRequest(t *testing.T, ctx context.Context, rr *httptest.ResponseRecorder, req *http.Request, s *Service) {

	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.WarmupRequestHandler(w, r.WithContext(ctx))
	})

	handlerFunc.ServeHTTP(rr, req)
}
