package repository

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// PostgresRepository is the production Repository implementation. The
// connection string is passed in by the caller (sourced from the
// DATABASE_URL environment variable), keeping this package free of any
// assumption about how configuration is loaded.
type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(databaseURL string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("opening database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) ListGames(ctx context.Context) ([]Game, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM games ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("querying games: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var games []Game
	for rows.Next() {
		var g Game
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, fmt.Errorf("scanning game row: %w", err)
		}
		games = append(games, g)
	}

	return games, rows.Err()
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}
