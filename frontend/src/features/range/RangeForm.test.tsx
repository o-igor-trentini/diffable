import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import { RangeForm } from './RangeForm'

describe('RangeForm', () => {
  it('renders all fields', () => {
    render(<RangeForm onSubmit={vi.fn()} isPending={false} />)

    expect(screen.getByLabelText('Workspace')).toBeInTheDocument()
    expect(screen.getByLabelText('Repositório')).toBeInTheDocument()
    expect(screen.getByLabelText('Hash Inicial (from)')).toBeInTheDocument()
    expect(screen.getByLabelText('Hash Final (to)')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /gerar descrição/i })).toBeInTheDocument()
  })

  it('disables submit when fields are empty', () => {
    render(<RangeForm onSubmit={vi.fn()} isPending={false} />)

    expect(screen.getByRole('button', { name: /gerar descrição/i })).toBeDisabled()
  })

  it('enables submit when all fields are filled', async () => {
    const user = userEvent.setup()
    render(<RangeForm onSubmit={vi.fn()} isPending={false} />)

    await user.type(screen.getByLabelText('Workspace'), 'ws')
    await user.type(screen.getByLabelText('Repositório'), 'repo')
    await user.type(screen.getByLabelText('Hash Inicial (from)'), 'abc')
    await user.type(screen.getByLabelText('Hash Final (to)'), 'def')

    expect(screen.getByRole('button', { name: /gerar descrição/i })).toBeEnabled()
  })

  it('submits with all fields', async () => {
    const user = userEvent.setup()
    const onSubmit = vi.fn()
    render(<RangeForm onSubmit={onSubmit} isPending={false} />)

    await user.type(screen.getByLabelText('Workspace'), 'ws')
    await user.type(screen.getByLabelText('Repositório'), 'repo')
    await user.type(screen.getByLabelText('Hash Inicial (from)'), 'abc123')
    await user.type(screen.getByLabelText('Hash Final (to)'), 'def456')
    await user.click(screen.getByRole('button', { name: /gerar descrição/i }))

    expect(onSubmit).toHaveBeenCalledWith({
      workspace: 'ws',
      repo_slug: 'repo',
      from_hash: 'abc123',
      to_hash: 'def456',
    })
  })
})
