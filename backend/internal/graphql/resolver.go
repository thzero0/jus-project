package graphql

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"context"

	"jus-project/backend/internal/graphql/generated"
	"jus-project/backend/internal/service"
)

// Resolver holds the dependencies GraphQL field resolvers need. It has one
// field today because the schema has one query; it grows alongside the
// schema, not ahead of it.
type Resolver struct {
	SuggestionService *service.SuggestionService
}

func NewResolver(suggestionService *service.SuggestionService) *Resolver {
	return &Resolver{SuggestionService: suggestionService}
}

// Suggestions is the resolver for the suggestions field.
func (r *queryResolver) Suggestions(ctx context.Context, term string) ([]string, error) {
	names := r.SuggestionService.Search(term)
	return names, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
