package exampleservice

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	contentTypeHeader   = "Content-type"
	applicationJSONType = "application/json"
	dateLayout          = "2006-01-02T15:04:05"
)

type Response struct {
	Message string `json:"message"`
}

type Service struct {
	db DataStore
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
	return &Service{
		db: NewDataStoreClient(mustEnv("DS_KIND")),
	}
}

func (s *Service) WarmupRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}
