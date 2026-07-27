import { render, screen, fireEvent } from '@testing-library/react'
import { vi, expect, test } from 'vitest'
import SuggestionList from '../SuggestionList'

test('renders nothing when there are no suggestions', () => {
  const { container } = render(<SuggestionList suggestions={[]} term="mine" onSelect={() => {}} />)

  expect(container).toBeEmptyDOMElement()
})

test('bolds the part of each suggestion that matched the term', () => {
  render(
    <SuggestionList
      suggestions={['Minecraft', 'Minecraft: Story Mode']}
      term="mine"
      onSelect={() => {}}
    />,
  )

  const [first] = screen.getAllByRole('button')
  expect(first.querySelector('strong')).toHaveTextContent('Mine')
  expect(first).toHaveTextContent('Minecraft')
})

test('calls onSelect with the full suggestion text when clicked', () => {
  const handleSelect = vi.fn()
  render(
    <SuggestionList
      suggestions={['Minecraft', 'Minecraft: Story Mode']}
      term="mine"
      onSelect={handleSelect}
    />,
  )

  fireEvent.click(screen.getByRole('button', { name: 'Minecraft: Story Mode' }))

  expect(handleSelect).toHaveBeenCalledExactlyOnceWith('Minecraft: Story Mode')
})
