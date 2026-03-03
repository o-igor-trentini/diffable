import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { ResultDisplay } from './ResultDisplay'

describe('ResultDisplay', () => {
  const mockResult = {
    id: 'test-id',
    type: 'single_commit',
    level: 'functional',
    description: 'Generated test description',
    model_used: 'gpt-4o-mini',
    tokens_used: 150,
    created_at: '2025-01-01T00:00:00Z',
  }

  it('renders the description', () => {
    render(<ResultDisplay result={mockResult} />)

    expect(screen.getByText('Generated test description')).toBeInTheDocument()
  })

  it('shows level, model and token info', () => {
    render(<ResultDisplay result={mockResult} />)

    expect(screen.getByText('Nível: Funcional')).toBeInTheDocument()
    expect(screen.getByText('Modelo: gpt-4o-mini')).toBeInTheDocument()
    expect(screen.getByText('Tokens: 150')).toBeInTheDocument()
  })

  it('renders copy and export buttons', () => {
    render(<ResultDisplay result={mockResult} />)

    expect(screen.getByText('Copiar')).toBeInTheDocument()
    expect(screen.getByText('Exportar')).toBeInTheDocument()
  })
})
