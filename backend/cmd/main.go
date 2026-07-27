package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"jus-project/backend/internal/repository"
	"jus-project/backend/internal/service"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	repo, err := repository.NewPostgresRepository(databaseURL)
	if err != nil {
		log.Fatalf("connecting to database: %v", err)
	}
	defer func() { _ = repo.Close() }()

	svc, err := service.NewSuggestionService(context.Background(), repo)
	if err != nil {
		log.Fatalf("loading suggestion service: %v", err)
	}
	log.Printf("suggestion service ready: %d games loaded", svc.Count())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "ok")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
