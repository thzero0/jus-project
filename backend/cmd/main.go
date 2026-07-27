package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"

	graphqlapi "jus-project/backend/internal/graphql"
	"jus-project/backend/internal/graphql/generated"
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

	resolver := graphqlapi.NewResolver(svc)
	schema := generated.NewExecutableSchema(generated.Config{Resolvers: resolver})

	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		corsOrigin = "http://localhost:3000"
	}
	graphqlHandler := cors.New(cors.Options{
		AllowedOrigins: []string{corsOrigin},
		AllowedMethods: []string{http.MethodPost},
		AllowedHeaders: []string{"Content-Type"},
	}).Handler(handler.NewDefaultServer(schema))

	http.Handle("/graphql", graphqlHandler)
	http.Handle("/", playground.Handler("GraphQL Playground", "/graphql"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
