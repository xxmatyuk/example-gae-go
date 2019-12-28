package exampleservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	// DB mocks
	okDBMock = DBMock{
		MockPutEntity: func(key string, entity *Entity) error {
			return nil
		},
		MockGetEntity: func(key string) (*Entity, error) {
			return &Entity{}, nil
		},
		MockGetAllEntities: func() (*[]Entity, error) {
			return &[]Entity{Entity{}}, nil
		},
	}

	// Cache mocks
	okCacheMock = CacheMock{
		MockGetItem: func(key string) (*Entity, error) {
			return &Entity{}, nil
		},
		MockPutItem: func(entity *Entity) error {
			return nil
		},
	}
)

func TestService(t *testing.T) {

	// Define tests
	tests := []struct {
		name             string
		dbMock           DataStoreClient
		cacheMock        CacheClient
		token            string
		withAuth         bool
		expectedCode     int
		expectedResponse Response
	}{
		{
			name:             "Happy-path-no-auth",
			dbMock:           okDBMock,
			cacheMock:        okCacheMock,
			expectedCode:     http.StatusOK,
			expectedResponse: Response{"OK"},
		},
		{
			name:             "Happy-path-with-auth",
			dbMock:           okDBMock,
			cacheMock:        okCacheMock,
			token:            mustEnv("TOKEN"),
			withAuth:         true,
			expectedCode:     http.StatusOK,
			expectedResponse: Response{"OK"},
		},
	}

	// Run positive tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			s := &Service{
				datastoreClient: test.dbMock,
				cacheClient:     test.cacheMock,
				logger:          initLoggerClient("test", "", "", "", ""),
			}

			req, _ := http.NewRequest(http.MethodGet, "/_ah/warmup", &bytes.Buffer{})
			rr := httptest.NewRecorder()

			if test.withAuth {
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", test.token))
			}
			performWarmUpRequest(rr, req, s, test.withAuth)

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

func performWarmUpRequest(rr *httptest.ResponseRecorder, req *http.Request, s *Service, withAuth bool) {

	var handlerFunc http.HandlerFunc
	if withAuth {
		handlerFunc = s.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.WarmupRequestHandler(w, r)
		}))
	} else {
		handlerFunc = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.WarmupRequestHandler(w, r)
		})
	}

	handlerFunc.ServeHTTP(rr, req)
}
