package service

import (
	"context"
	"fmt"
	"unicode/utf8"

	"jus-project/backend/internal/repository"
)

const (
	minQueryLength = 4
	maxResults     = 20
)

// SuggestionService answers prefix searches against the game catalog. It
// has no knowledge of GraphQL or any other transport — it is loaded once
// from a repository.Repository at startup and queried in memory from then
// on.
type SuggestionService struct {
	trie  *trie
	count int
}

func NewSuggestionService(ctx context.Context, repo repository.Repository) (*SuggestionService, error) {
	games, err := repo.ListGames(ctx)
	if err != nil {
		return nil, fmt.Errorf("loading games: %w", err)
	}

	t := newTrie()
	for _, g := range games {
		t.insert(g.Name)
	}

	return &SuggestionService{trie: t, count: len(games)}, nil
}

// Count returns how many games were loaded from the repository at startup.
func (s *SuggestionService) Count() int {
	return s.count
}

// Search returns up to 20 game names starting with term. It returns an
// empty slice for terms shorter than 4 characters or with no matches.
func (s *SuggestionService) Search(term string) []string {
	if utf8.RuneCountInString(term) < minQueryLength {
		return nil
	}
	return s.trie.search(term, maxResults)
}
