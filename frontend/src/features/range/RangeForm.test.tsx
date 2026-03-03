import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { RangeForm } from './RangeForm'

function renderWithQuery(ui: React.ReactElement) {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  })
  return render(
    <QueryClientProvider client={queryClient}>{ui}</QueryClientProvider>
  )
}

describe('RangeForm', () => {
  it('renders all fields', () => {
    renderWithQuery(<RangeForm onSubmit={vi.fn()} isPending={false} />)

    expect(screen.getByLabelText('Workspace')).toBeInTheDocument()
    expect(screen.getByLabelText('Repositorio')).toBeInTheDocument()
    expect(screen.getByLabelText('Hash Inicial (from)')).toBeInTheDocument()
    expect(screen.getByLabelText('Hash Final (to)')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /gerar descricao/i })).toBeInTheDocument()
  })

  it('disables submit when fields are empty', () => {
    renderWithQuery(<RangeForm onSubmit={vi.fn()} isPending={false} />)

    expect(screen.getByRole('button', { name: /gerar descricao/i })).toBeDisabled()
  })

  it('enables submit when all fields are filled', async () => {
    const user = userEvent.setup()
    renderWithQuery(<RangeForm onSubmit={vi.fn()} isPending={false} />)

    await user.type(screen.getByLabelText('Workspace'), 'ws')
    await user.type(screen.getByLabelText('Repositorio'), 'repo')
    await user.type(screen.getByLabelText('Hash Inicial (from)'), 'abc')
    await user.type(screen.getByLabelText('Hash Final (to)'), 'def')

    expect(screen.getByRole('button', { name: /gerar descricao/i })).toBeEnabled()
  })

  it('submits with all fields including level', async () => {
    const user = userEvent.setup()
    const onSubmit = vi.fn()
    renderWithQuery(<RangeForm onSubmit={onSubmit} isPending={false} />)

    await user.type(screen.getByLabelText('Workspace'), 'ws')
    await user.type(screen.getByLabelText('Repositorio'), 'repo')
    await user.type(screen.getByLabelText('Hash Inicial (from)'), 'abc123')
    await user.type(screen.getByLabelText('Hash Final (to)'), 'def456')
    await user.click(screen.getByRole('button', { name: /gerar descricao/i }))

    expect(onSubmit).toHaveBeenCalledWith({
      workspace: 'ws',
      repo_slug: 'repo',
      from_hash: 'abc123',
      to_hash: 'def456',
      level: 'functional',
    })
  })
})
