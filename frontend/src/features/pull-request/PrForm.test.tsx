import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { PrForm } from './PrForm'

function renderWithQuery(ui: React.ReactElement) {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  })
  return render(
    <QueryClientProvider client={queryClient}>{ui}</QueryClientProvider>
  )
}

describe('PrForm', () => {
  it('renders all fields', () => {
    renderWithQuery(<PrForm onSubmit={vi.fn()} isPending={false} />)

    expect(screen.getByLabelText('Workspace')).toBeInTheDocument()
    expect(screen.getByLabelText('Repositório')).toBeInTheDocument()
    expect(screen.getByLabelText('PR ID')).toBeInTheDocument()
    expect(screen.getByLabelText('Diff (raw)')).toBeInTheDocument()
    expect(screen.getByLabelText('Título do PR')).toBeInTheDocument()
    expect(screen.getByLabelText('Nível da Descrição')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /gerar descrição/i })).toBeInTheDocument()
  })

  it('disables submit when no fields filled', () => {
    renderWithQuery(<PrForm onSubmit={vi.fn()} isPending={false} />)

    expect(screen.getByRole('button', { name: /gerar descrição/i })).toBeDisabled()
  })

  it('enables submit when workspace/repo/pr_id filled', async () => {
    const user = userEvent.setup()
    renderWithQuery(<PrForm onSubmit={vi.fn()} isPending={false} />)

    await user.type(screen.getByLabelText('Workspace'), 'ws')
    await user.type(screen.getByLabelText('Repositório'), 'repo')
    await user.type(screen.getByLabelText('PR ID'), '42')

    expect(screen.getByRole('button', { name: /gerar descrição/i })).toBeEnabled()
  })

  it('submits with PR ID and level', async () => {
    const user = userEvent.setup()
    const onSubmit = vi.fn()
    renderWithQuery(<PrForm onSubmit={onSubmit} isPending={false} />)

    await user.type(screen.getByLabelText('Workspace'), 'ws')
    await user.type(screen.getByLabelText('Repositório'), 'repo')
    await user.type(screen.getByLabelText('PR ID'), '42')
    await user.click(screen.getByRole('button', { name: /gerar descrição/i }))

    expect(onSubmit).toHaveBeenCalledWith({
      workspace: 'ws',
      repo_slug: 'repo',
      pr_id: 42,
      level: 'functional',
    })
  })

  it('enables submit when raw diff + PR title filled', async () => {
    const user = userEvent.setup()
    renderWithQuery(<PrForm onSubmit={vi.fn()} isPending={false} />)

    await user.type(screen.getByLabelText('Diff (raw)'), 'diff --git a/main.go')
    await user.type(screen.getByLabelText('Título do PR'), 'My PR')

    expect(screen.getByRole('button', { name: /gerar descrição/i })).toBeEnabled()
  })
})
