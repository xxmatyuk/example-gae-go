package exampleservice

import (
	"net/http"
	"os"
)

const (
	appLogName = "app-log"
)

type Service struct {
	logger *loggingClient
}

func NewService() *Service {

	var (
		l   *loggingClient
		err error
	)

	if l, err = initLoggerClient(
		appLogName,
		os.Getenv("GOOGLE_CLOUD_PROJECT"),
		os.Getenv("GAE_VERSION"),
		os.Getenv("GAE_INSTANCE"),
		os.Getenv("GAE_SERVICE"),
	); err != nil {
		panic("No logger!")
	}

	return &Service{
		logger: l,
	}
}

func (s *Service) WarmupRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}
