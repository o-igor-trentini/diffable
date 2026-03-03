import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import { LevelSelector } from './LevelSelector'

describe('LevelSelector', () => {
  it('renders 4 level options including QA Detalhado', () => {
    render(<LevelSelector value="functional" onChange={vi.fn()} />)

    expect(screen.getByText('Funcional')).toBeInTheDocument()
    expect(screen.getByText('QA Detalhado')).toBeInTheDocument()
    expect(screen.getByText('Técnico')).toBeInTheDocument()
    expect(screen.getByText('Executivo')).toBeInTheDocument()

    const buttons = screen.getAllByRole('button')
    expect(buttons).toHaveLength(4)
  })

  it('calls onChange with qa_detailed when QA Detalhado is clicked', async () => {
    const user = userEvent.setup()
    const onChange = vi.fn()
    render(<LevelSelector value="functional" onChange={onChange} />)

    await user.click(screen.getByText('QA Detalhado'))

    expect(onChange).toHaveBeenCalledWith('qa_detailed')
  })
})
