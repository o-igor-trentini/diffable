import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import { PrForm } from './PrForm'

describe('PrForm', () => {
  it('renders all fields', () => {
    render(<PrForm onSubmit={vi.fn()} isPending={false} />)

    expect(screen.getByLabelText('Workspace')).toBeInTheDocument()
    expect(screen.getByLabelText('Repositório')).toBeInTheDocument()
    expect(screen.getByLabelText('PR ID')).toBeInTheDocument()
    expect(screen.getByLabelText('Diff (raw)')).toBeInTheDocument()
    expect(screen.getByLabelText('Título do PR')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /gerar descrição/i })).toBeInTheDocument()
  })

  it('disables submit when no fields filled', () => {
    render(<PrForm onSubmit={vi.fn()} isPending={false} />)

    expect(screen.getByRole('button', { name: /gerar descrição/i })).toBeDisabled()
  })

  it('enables submit when workspace/repo/pr_id filled', async () => {
    const user = userEvent.setup()
    render(<PrForm onSubmit={vi.fn()} isPending={false} />)

    await user.type(screen.getByLabelText('Workspace'), 'ws')
    await user.type(screen.getByLabelText('Repositório'), 'repo')
    await user.type(screen.getByLabelText('PR ID'), '42')

    expect(screen.getByRole('button', { name: /gerar descrição/i })).toBeEnabled()
  })

  it('submits with PR ID', async () => {
    const user = userEvent.setup()
    const onSubmit = vi.fn()
    render(<PrForm onSubmit={onSubmit} isPending={false} />)

    await user.type(screen.getByLabelText('Workspace'), 'ws')
    await user.type(screen.getByLabelText('Repositório'), 'repo')
    await user.type(screen.getByLabelText('PR ID'), '42')
    await user.click(screen.getByRole('button', { name: /gerar descrição/i }))

    expect(onSubmit).toHaveBeenCalledWith({
      workspace: 'ws',
      repo_slug: 'repo',
      pr_id: 42,
    })
  })

  it('enables submit when raw diff + PR title filled', async () => {
    const user = userEvent.setup()
    render(<PrForm onSubmit={vi.fn()} isPending={false} />)

    await user.type(screen.getByLabelText('Diff (raw)'), 'diff --git a/main.go')
    await user.type(screen.getByLabelText('Título do PR'), 'My PR')

    expect(screen.getByRole('button', { name: /gerar descrição/i })).toBeEnabled()
  })
})
