package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"reddit-clone-backend/internal/auth"
	database "reddit-clone-backend/internal/pkg/db/mysql"

	cors "github.com/rs/cors"

	"github.com/go-chi/chi"
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

	router.Use(auth.Middleware())

	database.InitDB()
	defer database.CloseDB()
	database.Migrate()

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
