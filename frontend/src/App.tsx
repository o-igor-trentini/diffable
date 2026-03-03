import { useState } from 'react'
import { Layout } from './features/shared/Layout'
import { TabNavigation, type TabId } from './features/shared/TabNavigation'
import { Placeholder } from './features/shared/Placeholder'
import { CommitAnalysis } from './features/commit/CommitAnalysis'
import { RangeAnalysis } from './features/range/RangeAnalysis'
import { PrAnalysis } from './features/pull-request/PrAnalysis'

export function App() {
  const [activeTab, setActiveTab] = useState<TabId>('commit')

  function renderTab() {
    switch (activeTab) {
      case 'commit':
        return <CommitAnalysis />
      case 'range':
        return <RangeAnalysis />
      case 'pr':
        return <PrAnalysis />
      case 'refine':
        return (
          <Placeholder
            title="Refinar Descrição"
            subtitle="Em breve — Fase 5"
          />
        )
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
    </Layout>
  )
}
