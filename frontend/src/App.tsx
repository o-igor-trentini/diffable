import { useState } from 'react'
import { Layout } from './features/shared/Layout'
import { TabNavigation, type TabId } from './features/shared/TabNavigation'
import { CommitAnalysis } from './features/commit/CommitAnalysis'
import { RangeAnalysis } from './features/range/RangeAnalysis'
import { PrAnalysis } from './features/pull-request/PrAnalysis'
import { RefineDescription } from './features/refine/RefineDescription'
import { HistoryPanel } from './features/history/HistoryPanel'
import type { AnalysisResponse } from './lib/api/types'

export function App() {
  const [activeTab, setActiveTab] = useState<TabId>('commit')
  const [currentAnalysis, setCurrentAnalysis] = useState<AnalysisResponse | null>(null)

  function handleRefine(result: AnalysisResponse) {
    setCurrentAnalysis(result)
    setActiveTab('refine')
  }

  function handleHistorySelect(analysis: AnalysisResponse) {
    setCurrentAnalysis(analysis)
    setActiveTab('refine')
  }

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
    <Layout>
      <div className="rounded-lg bg-white shadow">
        <TabNavigation activeTab={activeTab} onTabChange={setActiveTab} />

        <div className="p-6">
          {renderTab()}
        </div>
      </div>

      <div className="mt-4">
        <HistoryPanel onSelect={handleHistorySelect} />
      </div>
    </Layout>
  )
}
