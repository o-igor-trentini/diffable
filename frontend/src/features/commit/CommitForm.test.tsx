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
  it('renders bitbucket mode fields by default', () => {
    renderWithQuery(<CommitForm onSubmit={vi.fn()} isPending={false} />)

    expect(screen.getByLabelText('Workspace')).toBeInTheDocument()
    expect(screen.getByLabelText('Repositorio')).toBeInTheDocument()
    expect(screen.getByLabelText('Hash do Commit')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /gerar descricao/i })).toBeInTheDocument()
  })

  it('disables submit button when no fields filled', () => {
    renderWithQuery(<CommitForm onSubmit={vi.fn()} isPending={false} />)

    expect(screen.getByRole('button', { name: /gerar descricao/i })).toBeDisabled()
  })

  it('shows diff field when switching to manual mode', async () => {
    const user = userEvent.setup()
    renderWithQuery(<CommitForm onSubmit={vi.fn()} isPending={false} />)

    await user.click(screen.getByText('Colar Diff'))

    expect(screen.getByLabelText('Diff')).toBeInTheDocument()
  })

  it('enables submit button when raw diff is filled in manual mode', async () => {
    const user = userEvent.setup()
    renderWithQuery(<CommitForm onSubmit={vi.fn()} isPending={false} />)

    await user.click(screen.getByText('Colar Diff'))
    await user.type(screen.getByLabelText('Diff'), 'diff --git a/main.go')

    expect(screen.getByRole('button', { name: /gerar descricao/i })).toBeEnabled()
  })

  it('submits with raw_diff and level in manual mode', async () => {
    const user = userEvent.setup()
    const onSubmit = vi.fn()
    renderWithQuery(<CommitForm onSubmit={onSubmit} isPending={false} />)

    await user.click(screen.getByText('Colar Diff'))
    await user.type(screen.getByLabelText('Diff'), 'diff --git a/main.go')
    await user.click(screen.getByRole('button', { name: /gerar descricao/i }))

    expect(onSubmit).toHaveBeenCalledWith({ raw_diff: 'diff --git a/main.go', level: 'functional' })
  })

  it('submits with hash and level in bitbucket mode', async () => {
    const user = userEvent.setup()
    const onSubmit = vi.fn()
    renderWithQuery(<CommitForm onSubmit={onSubmit} isPending={false} />)

    await user.type(screen.getByLabelText('Workspace'), 'ws')
    await user.type(screen.getByLabelText('Repositorio'), 'repo')
    await user.type(screen.getByLabelText('Hash do Commit'), 'abc123')
    await user.click(screen.getByRole('button', { name: /gerar descricao/i }))

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
    expect(screen.getByLabelText('Repositorio')).toBeDisabled()
    expect(screen.getByLabelText('Hash do Commit')).toBeDisabled()
  })

  it('does not include overrides when using defaults', async () => {
    const user = userEvent.setup()
    const onSubmit = vi.fn()
    renderWithQuery(<CommitForm onSubmit={onSubmit} isPending={false} />)

    await user.click(screen.getByText('Colar Diff'))
    await user.type(screen.getByLabelText('Diff'), 'diff --git a/main.go')

    // Open advanced settings but don't change anything
    await user.click(screen.getByText('Configurações avançadas'))

    await user.click(screen.getByRole('button', { name: /gerar descricao/i }))

    const payload = onSubmit.mock.calls[0][0]
    expect(payload.overrides).toBeUndefined()
  })

  it('includes overrides when non-default values are set', async () => {
    const user = userEvent.setup()
    const onSubmit = vi.fn()
    renderWithQuery(<CommitForm onSubmit={onSubmit} isPending={false} />)

    await user.click(screen.getByText('Colar Diff'))
    await user.type(screen.getByLabelText('Diff'), 'diff --git a/main.go')

    // Open advanced settings and select GPT-4o model (non-default)
    await user.click(screen.getByText('Configurações avançadas'))
    await user.click(screen.getByText('GPT-4o'))

    await user.click(screen.getByRole('button', { name: /gerar descricao/i }))

    const payload = onSubmit.mock.calls[0][0]
    expect(payload.overrides).toBeDefined()
    expect(payload.overrides.model).toBe('gpt-4o')
  })
})
