package graphql

import (
	"context"
	"slices"
	"testing"

	"jus-project/backend/internal/graphql/generated"
	"jus-project/backend/internal/repository"
	"jus-project/backend/internal/service"
)

type fakeRepository struct {
	games []repository.Game
}

func (f *fakeRepository) ListGames(ctx context.Context) ([]repository.Game, error) {
	return f.games, nil
}

func newTestQueryResolver(t *testing.T, games []repository.Game) generated.QueryResolver {
	t.Helper()

	svc, err := service.NewSuggestionService(context.Background(), &fakeRepository{games: games})
	if err != nil {
		t.Fatalf("NewSuggestionService: %v", err)
	}
	return NewResolver(svc).Query()
}

func TestQueryResolver_Suggestions_ReturnsMatches(t *testing.T) {
	resolver := newTestQueryResolver(t, []repository.Game{
		{ID: 1, Name: "Red Dead Redemption"},
		{ID: 2, Name: "Red Dead Redemption 2"},
		{ID: 3, Name: "The Witcher 3"},
	})

	got, err := resolver.Suggestions(context.Background(), "Red ")
	if err != nil {
		t.Fatalf("Suggestions: %v", err)
	}

	want := []string{"Red Dead Redemption", "Red Dead Redemption 2"}
	if !slices.Equal(got, want) {
		t.Errorf("Suggestions(%q) = %v, want %v", "Red ", got, want)
	}
}

func TestQueryResolver_Suggestions_BelowMinLengthReturnsEmpty(t *testing.T) {
	resolver := newTestQueryResolver(t, []repository.Game{{ID: 1, Name: "Red Dead Redemption"}})

	got, err := resolver.Suggestions(context.Background(), "red")
	if err != nil {
		t.Fatalf("Suggestions: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("Suggestions(%q) = %v, want empty", "red", got)
	}
}
