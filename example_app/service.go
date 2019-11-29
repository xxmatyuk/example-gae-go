package exampleservice

import (
	"net/http"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) WarmupRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}
