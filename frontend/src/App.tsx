import { useState, useEffect, useCallback } from 'react'
import { Layout } from './features/shared/Layout'
import { TabNavigation, type TabId } from './features/shared/TabNavigation'
import { CommitAnalysis } from './features/commit/CommitAnalysis'
import { RangeAnalysis } from './features/range/RangeAnalysis'
import { PrAnalysis } from './features/pull-request/PrAnalysis'
import { RefineDescription } from './features/refine/RefineDescription'
import { HistoryPanel } from './features/history/HistoryPanel'
import type { AnalysisResponse } from './lib/api/types'

function useDarkMode() {
  const [dark, setDark] = useState(() => {
    const saved = localStorage.getItem('diffable-theme')
    if (saved) return saved === 'dark'
    return window.matchMedia('(prefers-color-scheme: dark)').matches
  })

  useEffect(() => {
    const root = document.documentElement
    if (dark) {
      root.classList.add('dark')
    } else {
      root.classList.remove('dark')
    }
    localStorage.setItem('diffable-theme', dark ? 'dark' : 'light')
  }, [dark])

  return [dark, setDark] as const
}

export function App() {
  const [activeTab, setActiveTab] = useState<TabId>('commit')
  const [currentAnalysis, setCurrentAnalysis] = useState<AnalysisResponse | null>(null)
  const [dark, setDark] = useDarkMode()

  function handleRefine(result: AnalysisResponse) {
    setCurrentAnalysis(result)
    setActiveTab('refine')
  }

  function handleHistorySelect(analysis: AnalysisResponse) {
    setCurrentAnalysis(analysis)
    setActiveTab('refine')
  }

  const handleKeyDown = useCallback((e: KeyboardEvent) => {
    // Ctrl+Enter: submit active form
    if (e.ctrlKey && e.key === 'Enter') {
      const form = document.querySelector('form')
      if (form) {
        const submitBtn = form.querySelector('button[type="submit"]') as HTMLButtonElement | null
        if (submitBtn && !submitBtn.disabled) {
          form.requestSubmit(submitBtn)
        }
      }
    }

    // Ctrl+Shift+C: copy result
    if (e.ctrlKey && e.shiftKey && e.key === 'C') {
      const copyBtn = document.querySelector('[data-copy-button]') as HTMLButtonElement | null
      if (copyBtn) {
        copyBtn.click()
      }
    }
  }, [])

  useEffect(() => {
    document.addEventListener('keydown', handleKeyDown)
    return () => document.removeEventListener('keydown', handleKeyDown)
  }, [handleKeyDown])

  function renderTab() {
    switch (activeTab) {
      case 'commit':
        return <CommitAnalysis onRefine={handleRefine} />
      case 'range':
        return <RangeAnalysis onRefine={handleRefine} />
      case 'pr':
        return <PrAnalysis onRefine={handleRefine} />
      case 'refine':
        return <RefineDescription analysis={currentAnalysis} />
    }
  }

  return (
    <Layout dark={dark} onToggleDark={() => setDark((d) => !d)}>
      <div className="rounded-lg bg-white shadow dark:bg-gray-800">
        <TabNavigation activeTab={activeTab} onTabChange={setActiveTab} />

        <div className="p-4 sm:p-6">
          {renderTab()}
        </div>
      </div>

      <div className="mt-4">
        <HistoryPanel onSelect={handleHistorySelect} />
      </div>
    </Layout>
  )
}
