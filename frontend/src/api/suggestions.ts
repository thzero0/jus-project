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

export async function fetchSuggestions(term: string, signal?: AbortSignal): Promise<string[]> {
  const data = await graphqlClient.request<SuggestionsResponse>({
    document: SUGGESTIONS_QUERY,
    variables: { term },
    signal,
  })
  return data.suggestions
}
