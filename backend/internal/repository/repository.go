package repository

import "context"

// Game is the data-layer representation of a single game record.
type Game struct {
	ID   int
	Name string
}

// Repository loads the full game catalog. Implementations back it with
// whatever storage is configured (Postgres in production, an in-memory
// fake in tests).
type Repository interface {
	ListGames(ctx context.Context) ([]Game, error)
}
