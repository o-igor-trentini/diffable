import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import { AdvancedSettings } from './AdvancedSettings'

// toLocaleString() formats numbers differently by locale (e.g. "2,048" vs "2.048")
// Use formatted value directly to match the component output
const fmt = (n: number) => n.toLocaleString()

describe('AdvancedSettings', () => {
  it('renders toggle button collapsed by default', () => {
    render(<AdvancedSettings value={{}} onChange={vi.fn()} />)

    expect(screen.getByText('Configurações avançadas')).toBeInTheDocument()
    expect(screen.queryByText('Temperature: 0.3')).not.toBeInTheDocument()
  })

  it('expands settings panel on toggle click', async () => {
    const user = userEvent.setup()
    render(<AdvancedSettings value={{}} onChange={vi.fn()} />)

    await user.click(screen.getByText('Configurações avançadas'))

    expect(screen.getByText('Temperature: 0.3')).toBeInTheDocument()
    expect(screen.getByText('Max Tokens')).toBeInTheDocument()
    expect(screen.getByText('Modelo')).toBeInTheDocument()
  })

  it('collapses settings panel on second toggle click', async () => {
    const user = userEvent.setup()
    render(<AdvancedSettings value={{}} onChange={vi.fn()} />)

    await user.click(screen.getByText('Configurações avançadas'))
    expect(screen.getByText('Temperature: 0.3')).toBeInTheDocument()

    await user.click(screen.getByText('Configurações avançadas'))
    expect(screen.queryByText('Temperature: 0.3')).not.toBeInTheDocument()
  })

  it('shows current temperature value from props', async () => {
    const user = userEvent.setup()
    render(<AdvancedSettings value={{ temperature: 0.7 }} onChange={vi.fn()} />)

    await user.click(screen.getByText('Configurações avançadas'))

    expect(screen.getByText('Temperature: 0.7')).toBeInTheDocument()
  })

  it('calls onChange when temperature slider changes', async () => {
    const user = userEvent.setup()
    const onChange = vi.fn()
    render(<AdvancedSettings value={{}} onChange={onChange} />)

    await user.click(screen.getByText('Configurações avançadas'))

    const slider = screen.getByRole('slider')
    Object.getOwnPropertyDescriptor(
      window.HTMLInputElement.prototype,
      'value'
    )!.set!.call(slider, '0.7')
    slider.dispatchEvent(new Event('change', { bubbles: true }))

    expect(onChange).toHaveBeenCalledWith({ temperature: 0.7 })
  })

  it('calls onChange when max_tokens button is clicked', async () => {
    const user = userEvent.setup()
    const onChange = vi.fn()
    render(<AdvancedSettings value={{}} onChange={onChange} />)

    await user.click(screen.getByText('Configurações avançadas'))
    await user.click(screen.getByText(fmt(2048)))

    expect(onChange).toHaveBeenCalledWith({ max_tokens: 2048 })
  })

  it('shows 1024 as default selected max_tokens', async () => {
    const user = userEvent.setup()
    render(<AdvancedSettings value={{}} onChange={vi.fn()} />)

    await user.click(screen.getByText('Configurações avançadas'))

    const btn1024 = screen.getByText(fmt(1024))
    expect(btn1024.className).toContain('violet')
  })

  it('highlights selected max_tokens from props', async () => {
    const user = userEvent.setup()
    render(<AdvancedSettings value={{ max_tokens: 4096 }} onChange={vi.fn()} />)

    await user.click(screen.getByText('Configurações avançadas'))

    const btn4096 = screen.getByText(fmt(4096))
    expect(btn4096.className).toContain('violet')

    const btn1024 = screen.getByText(fmt(1024))
    expect(btn1024.className).not.toContain('violet')
  })

  it('calls onChange when model card is clicked', async () => {
    const user = userEvent.setup()
    const onChange = vi.fn()
    render(<AdvancedSettings value={{}} onChange={onChange} />)

    await user.click(screen.getByText('Configurações avançadas'))
    await user.click(screen.getByText('GPT-4o'))

    expect(onChange).toHaveBeenCalledWith({ model: 'gpt-4o' })
  })

  it('shows Auto as default selected model', async () => {
    const user = userEvent.setup()
    render(<AdvancedSettings value={{}} onChange={vi.fn()} />)

    await user.click(screen.getByText('Configurações avançadas'))

    const autoBtn = screen.getByText('Auto').closest('button')!
    expect(autoBtn.className).toContain('violet')
  })

  it('highlights selected model from props', async () => {
    const user = userEvent.setup()
    render(
      <AdvancedSettings value={{ model: 'gpt-4o-mini' }} onChange={vi.fn()} />
    )

    await user.click(screen.getByText('Configurações avançadas'))

    const miniBtn = screen.getByText('GPT-4o Mini').closest('button')!
    expect(miniBtn.className).toContain('violet')

    const autoBtn = screen.getByText('Auto').closest('button')!
    expect(autoBtn.className).not.toContain('violet')
  })

  it('disables toggle button when disabled prop is true', () => {
    render(<AdvancedSettings value={{}} onChange={vi.fn()} disabled />)

    expect(screen.getByText('Configurações avançadas').closest('button')).toBeDisabled()
  })

  it('disables all inputs when disabled and expanded', async () => {
    const user = userEvent.setup()
    const { rerender } = render(
      <AdvancedSettings value={{}} onChange={vi.fn()} />
    )

    await user.click(screen.getByText('Configurações avançadas'))

    rerender(<AdvancedSettings value={{}} onChange={vi.fn()} disabled />)

    const slider = screen.getByRole('slider')
    expect(slider).toBeDisabled()

    const tokenButtons = screen.getAllByRole('button').filter(
      (btn) => btn.textContent && /^\d/.test(btn.textContent)
    )
    for (const btn of tokenButtons) {
      expect(btn).toBeDisabled()
    }
  })

  it('preserves other overrides when changing one field', async () => {
    const user = userEvent.setup()
    const onChange = vi.fn()
    render(
      <AdvancedSettings
        value={{ temperature: 0.5, max_tokens: 2048 }}
        onChange={onChange}
      />
    )

    await user.click(screen.getByText('Configurações avançadas'))
    await user.click(screen.getByText('GPT-4o Mini'))

    expect(onChange).toHaveBeenCalledWith({
      temperature: 0.5,
      max_tokens: 2048,
      model: 'gpt-4o-mini',
    })
  })
})
