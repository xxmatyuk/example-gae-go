package exmapleservice

import (
	"net/http"
)

type service struct{}

func NewService() *service {
	return &service{}
}

func (s *service) WarmupRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}
