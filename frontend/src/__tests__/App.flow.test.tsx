import { render, screen, fireEvent } from '@testing-library/react'
import { vi, expect, test, beforeEach } from 'vitest'
import App from '../App'
import { fetchSuggestions } from '../api/suggestions'

// Exercises the App -> SuggestionList flow with a mocked fetchSuggestions
// (no network, no debounce-timing assertions — that's covered by manual
// testing, see COMMENTS.md). Real API wiring lives in App.integration.test.tsx.
vi.mock('../api/suggestions', () => ({
  fetchSuggestions: vi.fn(),
}))

const mockedFetchSuggestions = vi.mocked(fetchSuggestions)

beforeEach(() => {
  mockedFetchSuggestions.mockReset()
})

test('typing a term shows the fetched suggestions', async () => {
  mockedFetchSuggestions.mockResolvedValue(['Minecraft', 'Minecraft: Story Mode'])

  render(<App />)
  fireEvent.change(screen.getByRole('textbox', { name: 'Buscar' }), {
    target: { value: 'mine' },
  })

  expect(await screen.findByRole('button', { name: 'Minecraft' })).toBeInTheDocument()
  expect(screen.getByRole('button', { name: 'Minecraft: Story Mode' })).toBeInTheDocument()
})

test('selecting a suggestion fills the field and closes the list', async () => {
  mockedFetchSuggestions.mockResolvedValue(['Minecraft', 'Minecraft: Story Mode'])

  render(<App />)
  const input = screen.getByRole('textbox', { name: 'Buscar' })
  fireEvent.change(input, { target: { value: 'mine' } })

  const option = await screen.findByRole('button', { name: 'Minecraft: Story Mode' })
  fireEvent.click(option)

  expect(input).toHaveValue('Minecraft: Story Mode')
  expect(screen.queryByRole('list')).not.toBeInTheDocument()
})
