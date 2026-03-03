import { renderHook, act } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useClipboard } from './useClipboard'

describe('useClipboard', () => {
  beforeEach(() => {
    Object.assign(navigator, {
      clipboard: {
        writeText: vi.fn().mockResolvedValue(undefined),
      },
    })
  })

  it('starts with copied as false', () => {
    const { result } = renderHook(() => useClipboard())

    expect(result.current.copied).toBe(false)
  })

  it('copies text and sets copied to true', async () => {
    const { result } = renderHook(() => useClipboard())

    await act(async () => {
      result.current.copy('test text')
    })

    expect(navigator.clipboard.writeText).toHaveBeenCalledWith('test text')
    expect(result.current.copied).toBe(true)
  })

  it('resets copied to false after 2 seconds', async () => {
    vi.useFakeTimers()
    const { result } = renderHook(() => useClipboard())

    await act(async () => {
      result.current.copy('test text')
    })

    expect(result.current.copied).toBe(true)

    act(() => {
      vi.advanceTimersByTime(2000)
    })

    expect(result.current.copied).toBe(false)
    vi.useRealTimers()
  })
})
