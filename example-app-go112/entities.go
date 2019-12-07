package exampleservice

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Service) GetEntityHandler(w http.ResponseWriter, r *http.Request) {

	var (
		entity *Entity
		err    error
	)

	vars := mux.Vars(r)
	key := vars["key"]

	if entity, err = s.cacheClient.GetItem(key); err != nil {
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

	if err := s.datastoreClient.PutEntity(entity.Key, &entity); err != nil {
		resp := Response{fmt.Sprintf("Error: %s", err.Error())}
		s.writeResponseData(w, http.StatusInternalServerError, &resp)
		return
	}

	_ = s.cacheClient.PutItem(&entity)

	return
}

func (s *Service) GetAllEntitiesHandler(w http.ResponseWriter, r *http.Request) {

	var (
		entities []Entity
		err      error
	)

	if entities, err = s.datastoreClient.GetAllEntities(); err != nil {
		resp := Response{fmt.Sprintf("Error: %s", err.Error())}
		s.writeResponseData(w, http.StatusInternalServerError, &resp)
		return
	}

	s.writeResponseData(w, http.StatusOK, &entities)
	return
}
