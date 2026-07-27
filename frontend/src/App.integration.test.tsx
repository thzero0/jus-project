import { describe, test, expect } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import App from './App'

// Exercises the real GraphQL client against a running `docker compose up`
// stack — see docs/testing-backend.md for how the backend gets seeded.
// Skipped by default so `npm test` stays fast and dependency-free in CI;
// run with `RUN_INTEGRATION_TESTS=1 npm test` locally.
describe.skipIf(process.env.RUN_INTEGRATION_TESTS !== '1')(
  'App (integração com a API real)',
  () => {
    test('busca sugestões reais via GraphQL', async () => {
      render(<App />)

      const input = screen.getByRole('textbox', { name: 'Buscar' })
      fireEvent.change(input, { target: { value: 'minecraf' } })

      await waitFor(() => {
        expect(screen.getByText('Minecraft')).toBeInTheDocument()
      })
    })
  },
)
