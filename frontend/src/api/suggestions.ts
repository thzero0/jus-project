import { gql } from 'graphql-request'
import { graphqlClient } from '../lib/graphqlClient'

const SUGGESTIONS_QUERY = gql`
  query Suggestions($term: String!) {
    suggestions(term: $term)
  }
`

interface SuggestionsResponse {
  suggestions: string[]
}

export async function fetchSuggestions(term: string): Promise<string[]> {
  const data = await graphqlClient.request<SuggestionsResponse>(SUGGESTIONS_QUERY, { term })
  return data.suggestions
}
