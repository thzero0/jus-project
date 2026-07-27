import { useEffect, useState } from 'react'
import { fetchSuggestions } from './api/suggestions'
import SuggestionList from './components/SuggestionList'
import styles from './App.module.css'

const MIN_QUERY_LENGTH = 4
const DEBOUNCE_MS = 250

function isAbortError(error: unknown): boolean {
  return error instanceof DOMException && error.name === 'AbortError'
}

function App() {
  const [term, setTerm] = useState('')
  const [suggestions, setSuggestions] = useState<string[]>([])
  const [isListOpen, setIsListOpen] = useState(false)

  useEffect(() => {
    if (term.length < MIN_QUERY_LENGTH) return

    const controller = new AbortController()

    const timeoutId = setTimeout(() => {
      fetchSuggestions(term, controller.signal)
        .then(setSuggestions)
        .catch((error: unknown) => {
          if (!isAbortError(error)) console.error(error)
        })
    }, DEBOUNCE_MS)

    return () => {
      clearTimeout(timeoutId)
      controller.abort()
    }
  }, [term])

  function handleChange(event: React.ChangeEvent<HTMLInputElement>) {
    setTerm(event.target.value)
    setIsListOpen(true)
  }

  function handleSelect(suggestion: string) {
    setTerm(suggestion)
    setIsListOpen(false)
  }

  const visibleSuggestions = isListOpen && term.length >= MIN_QUERY_LENGTH ? suggestions : []

  return (
    <main className={styles.app}>
      <h1 className={styles.title}>Justarter</h1>
      <div className={styles.searchContainer}>
        <input
          type="text"
          value={term}
          onChange={handleChange}
          aria-label="Buscar"
          placeholder="Buscar jogo..."
          className={styles.input}
        />
        <SuggestionList suggestions={visibleSuggestions} term={term} onSelect={handleSelect} />
      </div>
    </main>
  )
}

export default App
