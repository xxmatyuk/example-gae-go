package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	service "github.com/xxmatyuk/example-gae-go/example_app_go112"
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	http.ListenAndServe(":"+port, r)
}
