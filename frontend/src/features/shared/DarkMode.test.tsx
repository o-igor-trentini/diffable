import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, beforeEach } from 'vitest'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { App } from '@/App'

function renderApp() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  })

  return render(
    <QueryClientProvider client={queryClient}>
      <App />
    </QueryClientProvider>,
  )
}

describe('Dark Mode', () => {
  beforeEach(() => {
    localStorage.clear()
    document.documentElement.classList.remove('dark')
  })

  it('toggles dark class on document element', async () => {
    const user = userEvent.setup()

    renderApp()

    const toggleBtn = screen.getByRole('button', { name: /modo escuro/i })
    expect(document.documentElement.classList.contains('dark')).toBe(false)

    await user.click(toggleBtn)
    expect(document.documentElement.classList.contains('dark')).toBe(true)

    const lightBtn = screen.getByRole('button', { name: /modo claro/i })
    await user.click(lightBtn)
    expect(document.documentElement.classList.contains('dark')).toBe(false)
  })

  it('persists preference in localStorage', async () => {
    const user = userEvent.setup()

    renderApp()

    const toggleBtn = screen.getByRole('button', { name: /modo escuro/i })
    await user.click(toggleBtn)

    expect(localStorage.getItem('diffable-theme')).toBe('dark')
  })

  it('loads persisted dark preference', () => {
    localStorage.setItem('diffable-theme', 'dark')

    renderApp()

    expect(document.documentElement.classList.contains('dark')).toBe(true)
    expect(screen.getByRole('button', { name: /modo claro/i })).toBeInTheDocument()
  })
})
