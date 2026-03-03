import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import { TabNavigation } from './TabNavigation'

describe('TabNavigation', () => {
  it('renders all 4 tabs', () => {
    render(<TabNavigation activeTab="commit" onTabChange={() => {}} />)

    expect(screen.getByText('Commit')).toBeInTheDocument()
    expect(screen.getByText('Range')).toBeInTheDocument()
    expect(screen.getByText('Pull Request')).toBeInTheDocument()
    expect(screen.getByText('Refinar')).toBeInTheDocument()
  })

  it('marks the active tab with aria-selected', () => {
    render(<TabNavigation activeTab="pr" onTabChange={() => {}} />)

    const prTab = screen.getByText('Pull Request').closest('button')
    const commitTab = screen.getByText('Commit').closest('button')

    expect(prTab).toHaveAttribute('aria-selected', 'true')
    expect(commitTab).toHaveAttribute('aria-selected', 'false')
  })

  it('calls onTabChange when clicking a tab', async () => {
    const user = userEvent.setup()
    const onChange = vi.fn()
    render(<TabNavigation activeTab="commit" onTabChange={onChange} />)

    await user.click(screen.getByText('Range'))

    expect(onChange).toHaveBeenCalledWith('range')
  })
})
