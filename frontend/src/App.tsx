import { useState } from 'react'
import { fetchSuggestions } from './api/suggestions'

const MIN_QUERY_LENGTH = 4

function App() {
  const [term, setTerm] = useState('')
  const [suggestions, setSuggestions] = useState<string[]>([])

  async function handleChange(event: React.ChangeEvent<HTMLInputElement>) {
    const value = event.target.value
    setTerm(value)

    if (value.length < MIN_QUERY_LENGTH) {
      setSuggestions([])
      return
    }

    const results = await fetchSuggestions(value)
    setSuggestions(results)
  }

  return (
    <main>
      <h1>Justarter</h1>
      <input type="text" value={term} onChange={handleChange} aria-label="Buscar" />
      <ul aria-label="Sugestões">
        {suggestions.map((suggestion) => (
          <li key={suggestion}>{suggestion}</li>
        ))}
      </ul>
    </main>
  )
}

export default App
