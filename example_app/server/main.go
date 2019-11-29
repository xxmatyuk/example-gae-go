package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"google.golang.org/appengine"

	service "github.com/xxmatyuk/example-gae-go111/example_app"
)

func main() {

	// Router
	r := mux.NewRouter()

	// Service
	s := service.NewService()

	// Handlers
	r.HandleFunc("/_ah/warmup", s.WarmupRequestHandler).Methods("GET")
	r.HandleFunc("/test-admin", s.WarmupRequestHandler).Methods("GET")
	r.HandleFunc("/test-no-admin", s.WarmupRequestHandler).Methods("GET")
	r.HandleFunc("/test-custom-auth", s.Auth(s.WarmupRequestHandler)).Methods("GET")

	http.Handle("/", r)

	appengine.Main()

}
