package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"jus-project/backend/internal/repository"
)

type fakeRepository struct {
	games []repository.Game
	err   error
}

func (f *fakeRepository) ListGames(ctx context.Context) ([]repository.Game, error) {
	return f.games, f.err
}

func newTestService(t *testing.T, games []repository.Game) *SuggestionService {
	t.Helper()

	svc, err := NewSuggestionService(context.Background(), &fakeRepository{games: games})
	if err != nil {
		t.Fatalf("NewSuggestionService: %v", err)
	}
	return svc
}

func TestNewSuggestionService_RepositoryError(t *testing.T) {
	_, err := NewSuggestionService(context.Background(), &fakeRepository{err: errors.New("boom")})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestSearch_BelowMinLengthReturnsEmpty(t *testing.T) {
	svc := newTestService(t, []repository.Game{{ID: 1, Name: "Red Dead Redemption"}})

	for _, term := range []string{"", "r", "re", "red"} {
		if got := svc.Search(term); len(got) != 0 {
			t.Errorf("Search(%q) = %v, want empty", term, got)
		}
	}
}

func TestSearch_CapsAt20Results(t *testing.T) {
	games := make([]repository.Game, 0, 25)
	for i := 0; i < 25; i++ {
		games = append(games, repository.Game{ID: i, Name: fmt.Sprintf("Test Game %02d", i)})
	}
	svc := newTestService(t, games)

	got := svc.Search("Test")
	if len(got) != maxResults {
		t.Errorf("Search(%q) returned %d results, want %d", "Test", len(got), maxResults)
	}
}
