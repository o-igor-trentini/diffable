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
  it('renders bitbucket mode fields by default', () => {
    renderWithQuery(<PrForm onSubmit={vi.fn()} isPending={false} />)

    expect(screen.getByLabelText('Workspace')).toBeInTheDocument()
    expect(screen.getByLabelText('Repositorio')).toBeInTheDocument()
    expect(screen.getByLabelText('Numero do PR')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /gerar descricao/i })).toBeInTheDocument()
  })

  it('disables submit when no fields filled', () => {
    renderWithQuery(<PrForm onSubmit={vi.fn()} isPending={false} />)

    expect(screen.getByRole('button', { name: /gerar descricao/i })).toBeDisabled()
  })

  it('enables submit when workspace/repo/pr_id filled', async () => {
    const user = userEvent.setup()
    renderWithQuery(<PrForm onSubmit={vi.fn()} isPending={false} />)

    await user.type(screen.getByLabelText('Workspace'), 'ws')
    await user.type(screen.getByLabelText('Repositorio'), 'repo')
    await user.type(screen.getByLabelText('Numero do PR'), '42')

    expect(screen.getByRole('button', { name: /gerar descricao/i })).toBeEnabled()
  })

  it('submits with PR ID and level', async () => {
    const user = userEvent.setup()
    const onSubmit = vi.fn()
    renderWithQuery(<PrForm onSubmit={onSubmit} isPending={false} />)

    await user.type(screen.getByLabelText('Workspace'), 'ws')
    await user.type(screen.getByLabelText('Repositorio'), 'repo')
    await user.type(screen.getByLabelText('Numero do PR'), '42')
    await user.click(screen.getByRole('button', { name: /gerar descricao/i }))

    expect(onSubmit).toHaveBeenCalledWith({
      workspace: 'ws',
      repo_slug: 'repo',
      pr_id: 42,
      level: 'functional',
    })
  })

  it('enables submit when raw diff + PR title filled in manual mode', async () => {
    const user = userEvent.setup()
    renderWithQuery(<PrForm onSubmit={vi.fn()} isPending={false} />)

    await user.click(screen.getByText('Colar Diff'))
    await user.type(screen.getByLabelText('Titulo do PR'), 'My PR')
    await user.type(screen.getByLabelText('Diff'), 'diff --git a/main.go')

    expect(screen.getByRole('button', { name: /gerar descricao/i })).toBeEnabled()
  })
})
