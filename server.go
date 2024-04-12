package main

import (
	"log"
	"net/http"
	"os"

	cors "github.com/rs/cors"

	chi "github.com/go-chi/chi/v5"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		Debug:            true,
	}).Handler)
	log.Printf("Vote chain is listening on http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
