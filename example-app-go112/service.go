package exampleservice

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-redis/redis"
)

const (
	contentTypeHeader   = "Content-type"
	applicationJSONType = "application/json"
	appLogName          = "app_log"
)

type Response struct {
	Message string `json:"message"`
}

type Service struct {
	datastoreClient DataStoreClient
	cacheClient     CacheClient
	redisClient     *redis.Client
	logger          *loggingClient
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic(k + " environment variable is not set")
	}
	return v
}

func (s *Service) parseRequest(r *http.Request, marshalledRequest interface{}) error {
	var (
		err         error
		requestBody []byte
	)

	requestBody, err = ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}

	err = json.Unmarshal(requestBody, &marshalledRequest)
	if err != nil {
		return err
	}

	return err
}

func (s *Service) writeResponseData(w http.ResponseWriter, code int, data interface{}) error {

	jsonResponse, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.WriteHeader(code)
	w.Header().Set(contentTypeHeader, applicationJSONType)
	w.Write(jsonResponse)

	return nil
}

func NewService() *Service {

	logger := initLoggerClient(
		appLogName,
		mustEnv("GOOGLE_CLOUD_PROJECT"),
		os.Getenv("GAE_VERSION"),
		os.Getenv("GAE_INSTANCE"),
		os.Getenv("GAE_SERVICE"),
	)

	datastoreClient := NewDataStoreClient(
		mustEnv("GOOGLE_CLOUD_PROJECT"),
		mustEnv("DS_KIND"),
	)

	cacheClient := NewCacheDBClient(
		datastoreClient,
		mustEnv("REDISHOST"),
		mustEnv("REDISPORT"),
		mustEnv("REDISDB"),
	)

	return &Service{
		datastoreClient: datastoreClient,
		logger:          logger,
		cacheClient:     cacheClient,
	}
}

func (s *Service) WarmupRequestHandler(w http.ResponseWriter, r *http.Request) {
	s.writeResponseData(w, http.StatusOK, &Response{"OK"})
	return
}
