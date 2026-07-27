import { useEffect, useState } from 'react'
import { fetchSuggestions } from './api/suggestions'

const MIN_QUERY_LENGTH = 4
const DEBOUNCE_MS = 250

function isAbortError(error: unknown): boolean {
  return error instanceof DOMException && error.name === 'AbortError'
}

function App() {
  const [term, setTerm] = useState('')
  const [suggestions, setSuggestions] = useState<string[]>([])

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

  const visibleSuggestions = term.length < MIN_QUERY_LENGTH ? [] : suggestions

  return (
    <main>
      <h1>Justarter</h1>
      <input
        type="text"
        value={term}
        onChange={(event) => setTerm(event.target.value)}
        aria-label="Buscar"
      />
      <ul aria-label="Sugestões">
        {visibleSuggestions.map((suggestion) => (
          <li key={suggestion}>{suggestion}</li>
        ))}
      </ul>
    </main>
  )
}

export default App
