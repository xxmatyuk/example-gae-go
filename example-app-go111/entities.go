package exampleservice

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Entity struct {
	Key   string `json:"key" datastore:"key"`
	Value string `json:"value,omitempty" datastore:"value"`
}

func (s *Service) GetEntityHandler(w http.ResponseWriter, r *http.Request) {

	var (
		entity *Entity
		err    error
	)

	vars := mux.Vars(r)
	key := vars["key"]

	if entity, err = s.cache.GetItem(r.Context(), key); err != nil {
		resp := Response{fmt.Sprintf("Error: %s", err.Error())}
		s.writeResponseData(w, http.StatusInternalServerError, &resp)
		return
	}

	s.writeResponseData(w, http.StatusOK, &entity)

	return
}

func (s *Service) PutEntityHandler(w http.ResponseWriter, r *http.Request) {

	var entity Entity
	if err := s.parseRequest(r, &entity); err != nil {
		resp := Response{fmt.Sprintf("Error: %s", err.Error())}
		s.writeResponseData(w, http.StatusInternalServerError, &resp)
		return
	}

	if err := s.db.PutEntity(r.Context(), entity.Key, &entity); err != nil {
		resp := Response{fmt.Sprintf("Error: %s", err.Error())}
		s.writeResponseData(w, http.StatusInternalServerError, &resp)
		return
	}

	_ = s.cache.PutItem(r.Context(), &entity)

	return
}

func (s *Service) GetAllEntitiesHandler(w http.ResponseWriter, r *http.Request) {

	var entities []Entity
	if _, err := s.db.GetAllEntities(r.Context(), &entities); err != nil {
		resp := Response{fmt.Sprintf("Error: %s", err.Error())}
		s.writeResponseData(w, http.StatusInternalServerError, &resp)
		return
	}

	s.writeResponseData(w, http.StatusOK, &entities)
	return
}
