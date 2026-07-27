package repository

import (
	"context"
	"os"
	"testing"
)

// TestPostgresRepository_ListGames exercises the real Postgres wiring used
// by cmd/main.go (NewPostgresRepository + ListGames) against a database
// seeded from db/games.csv via db/seed.sh. It requires DATABASE_URL to
// point at a running, seeded database:
//
//	docker compose up -d db seed
//	export DATABASE_URL="postgres://admin:123@localhost:5432/games-db?sslmode=disable"
//	go test ./internal/repository/... -v
//
// It is skipped (not failed) when DATABASE_URL is unset, so `go test ./...`
// stays green in CI, which doesn't provision a database for this job.
func TestPostgresRepository_ListGames(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set; see docs/testing-backend.md to run this test")
	}

	repo, err := NewPostgresRepository(databaseURL)
	if err != nil {
		t.Fatalf("NewPostgresRepository: %v", err)
	}
	defer func() { _ = repo.Close() }()

	games, err := repo.ListGames(context.Background())
	if err != nil {
		t.Fatalf("ListGames: %v", err)
	}

	// db/games.csv after dedup by id in db/clean_data.py. Update this if the
	// versioned dataset changes deliberately.
	const wantCount = 8978
	if len(games) != wantCount {
		t.Errorf("ListGames() returned %d games, want %d", len(games), wantCount)
	}

	var found bool
	for _, g := range games {
		if g.Name == "The Elder Scrolls VI" {
			found = true
			break
		}
	}
	if !found {
		t.Error(`ListGames() missing expected game "The Elder Scrolls VI"`)
	}
}
