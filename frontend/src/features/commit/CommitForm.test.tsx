import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import { CommitForm } from './CommitForm'

describe('CommitForm', () => {
  it('renders all fields', () => {
    render(<CommitForm onSubmit={vi.fn()} isPending={false} />)

    expect(screen.getByLabelText('Workspace')).toBeInTheDocument()
    expect(screen.getByLabelText('Repositório')).toBeInTheDocument()
    expect(screen.getByLabelText('Hash do Commit')).toBeInTheDocument()
    expect(screen.getByLabelText('Diff (raw)')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /gerar descrição/i })).toBeInTheDocument()
  })

  it('disables submit button when no fields filled', () => {
    render(<CommitForm onSubmit={vi.fn()} isPending={false} />)

    expect(screen.getByRole('button', { name: /gerar descrição/i })).toBeDisabled()
  })

  it('enables submit button when raw diff is filled', async () => {
    const user = userEvent.setup()
    render(<CommitForm onSubmit={vi.fn()} isPending={false} />)

    await user.type(screen.getByLabelText('Diff (raw)'), 'diff --git a/main.go')

    expect(screen.getByRole('button', { name: /gerar descrição/i })).toBeEnabled()
  })

  it('submits with raw_diff when raw diff is filled', async () => {
    const user = userEvent.setup()
    const onSubmit = vi.fn()
    render(<CommitForm onSubmit={onSubmit} isPending={false} />)

    await user.type(screen.getByLabelText('Diff (raw)'), 'diff --git a/main.go')
    await user.click(screen.getByRole('button', { name: /gerar descrição/i }))

    expect(onSubmit).toHaveBeenCalledWith({ raw_diff: 'diff --git a/main.go' })
  })

  it('submits with hash when workspace/repo/hash filled', async () => {
    const user = userEvent.setup()
    const onSubmit = vi.fn()
    render(<CommitForm onSubmit={onSubmit} isPending={false} />)

    await user.type(screen.getByLabelText('Workspace'), 'ws')
    await user.type(screen.getByLabelText('Repositório'), 'repo')
    await user.type(screen.getByLabelText('Hash do Commit'), 'abc123')
    await user.click(screen.getByRole('button', { name: /gerar descrição/i }))

    expect(onSubmit).toHaveBeenCalledWith({
      workspace: 'ws',
      repo_slug: 'repo',
      commit_hash: 'abc123',
    })
  })

  it('disables fields when isPending is true', () => {
    render(<CommitForm onSubmit={vi.fn()} isPending={true} />)

    expect(screen.getByLabelText('Workspace')).toBeDisabled()
    expect(screen.getByLabelText('Repositório')).toBeDisabled()
    expect(screen.getByLabelText('Hash do Commit')).toBeDisabled()
    expect(screen.getByLabelText('Diff (raw)')).toBeDisabled()
  })
})
