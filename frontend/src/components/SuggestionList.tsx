import styles from './SuggestionList.module.css'

interface SuggestionListProps {
  suggestions: string[]
  term: string
  onSelect: (suggestion: string) => void
}

// The backend matches by prefix (internal/service/trie.go), so the part of
// each suggestion that matched `term` is always its first term.length
// characters — no substring search needed to know what to bold.
function SuggestionList({ suggestions, term, onSelect }: SuggestionListProps) {
  if (suggestions.length === 0) return null

  return (
    <ul className={styles.list} aria-label="Sugestões">
      {suggestions.map((suggestion) => (
        <li key={suggestion}>
          <button type="button" className={styles.item} onClick={() => onSelect(suggestion)}>
            <strong className={styles.match}>{suggestion.slice(0, term.length)}</strong>
            {suggestion.slice(term.length)}
          </button>
        </li>
      ))}
    </ul>
  )
}

export default SuggestionList
