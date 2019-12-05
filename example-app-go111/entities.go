package exampleservice

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"google.golang.org/appengine"
)

type Entity struct {
	Key          string `json:"key" datastore:"key"`
	Value        string `json:"value,omitempty" datastore:"value"`
	CreationTime string `json:"creation_time,omitempty" datastore:"creation_time"`
}

func (s *Service) GetEntityHandler(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)
	ctx, _ = appengine.Namespace(ctx, "")

	vars := mux.Vars(r)
	key := vars["key"]

	var entity Entity
	if err := s.db.GetEntity(ctx, key, &entity); err != nil {
		resp := Response{fmt.Sprintf("Error: %s", err.Error())}
		s.writeResponseData(w, http.StatusInternalServerError, &resp)
		return
	}

	s.writeResponseData(w, http.StatusOK, &entity)
	return
}

func (s *Service) PutEntityHandler(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)
	ctx, _ = appengine.Namespace(ctx, "")

	var entity Entity
	if err := s.parseRequest(r, &entity); err != nil {
		resp := Response{fmt.Sprintf("Error: %s", err.Error())}
		s.writeResponseData(w, http.StatusInternalServerError, &resp)
		return
	}

	entity.CreationTime = time.Now().Format(dateLayout)

	if err := s.db.PutEntity(ctx, entity.Key, &entity); err != nil {
		resp := Response{fmt.Sprintf("Error: %s", err.Error())}
		s.writeResponseData(w, http.StatusInternalServerError, &resp)
		return
	}

	return
}

func (s *Service) GetAllEntitiesHandler(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)
	ctx, _ = appengine.Namespace(ctx, "")

	var entities []Entity
	if err := s.db.GetAllEntities(ctx, &entities); err != nil {
		resp := Response{fmt.Sprintf("Error: %s", err.Error())}
		s.writeResponseData(w, http.StatusInternalServerError, &resp)
		return
	}

	s.writeResponseData(w, http.StatusOK, &entities)
	return
}
