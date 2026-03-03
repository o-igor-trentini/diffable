import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { CommitForm } from './CommitForm'

function renderWithQuery(ui: React.ReactElement) {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  })
  return render(
    <QueryClientProvider client={queryClient}>{ui}</QueryClientProvider>
  )
}

describe('CommitForm', () => {
  it('renders all fields', () => {
    renderWithQuery(<CommitForm onSubmit={vi.fn()} isPending={false} />)

    expect(screen.getByLabelText('Workspace')).toBeInTheDocument()
    expect(screen.getByLabelText('Repositório')).toBeInTheDocument()
    expect(screen.getByLabelText('Hash do Commit')).toBeInTheDocument()
    expect(screen.getByLabelText('Diff (raw)')).toBeInTheDocument()
    expect(screen.getByLabelText('Nível da Descrição')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /gerar descrição/i })).toBeInTheDocument()
  })

  it('disables submit button when no fields filled', () => {
    renderWithQuery(<CommitForm onSubmit={vi.fn()} isPending={false} />)

    expect(screen.getByRole('button', { name: /gerar descrição/i })).toBeDisabled()
  })

  it('enables submit button when raw diff is filled', async () => {
    const user = userEvent.setup()
    renderWithQuery(<CommitForm onSubmit={vi.fn()} isPending={false} />)

    await user.type(screen.getByLabelText('Diff (raw)'), 'diff --git a/main.go')

    expect(screen.getByRole('button', { name: /gerar descrição/i })).toBeEnabled()
  })

  it('submits with raw_diff and level when raw diff is filled', async () => {
    const user = userEvent.setup()
    const onSubmit = vi.fn()
    renderWithQuery(<CommitForm onSubmit={onSubmit} isPending={false} />)

    await user.type(screen.getByLabelText('Diff (raw)'), 'diff --git a/main.go')
    await user.click(screen.getByRole('button', { name: /gerar descrição/i }))

    expect(onSubmit).toHaveBeenCalledWith({ raw_diff: 'diff --git a/main.go', level: 'functional' })
  })

  it('submits with hash and level when workspace/repo/hash filled', async () => {
    const user = userEvent.setup()
    const onSubmit = vi.fn()
    renderWithQuery(<CommitForm onSubmit={onSubmit} isPending={false} />)

    await user.type(screen.getByLabelText('Workspace'), 'ws')
    await user.type(screen.getByLabelText('Repositório'), 'repo')
    await user.type(screen.getByLabelText('Hash do Commit'), 'abc123')
    await user.click(screen.getByRole('button', { name: /gerar descrição/i }))

    expect(onSubmit).toHaveBeenCalledWith({
      workspace: 'ws',
      repo_slug: 'repo',
      commit_hash: 'abc123',
      level: 'functional',
    })
  })

  it('disables fields when isPending is true', () => {
    renderWithQuery(<CommitForm onSubmit={vi.fn()} isPending={true} />)

    expect(screen.getByLabelText('Workspace')).toBeDisabled()
    expect(screen.getByLabelText('Repositório')).toBeDisabled()
    expect(screen.getByLabelText('Hash do Commit')).toBeDisabled()
    expect(screen.getByLabelText('Diff (raw)')).toBeDisabled()
  })

  it('level select defaults to functional', () => {
    renderWithQuery(<CommitForm onSubmit={vi.fn()} isPending={false} />)

    const select = screen.getByLabelText('Nível da Descrição') as HTMLSelectElement
    expect(select.value).toBe('functional')
  })
})
