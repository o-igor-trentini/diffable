import { useState } from 'react'
import { Layout } from './features/shared/Layout'
import { TabNavigation, tabs, type TabId } from './features/shared/TabNavigation'
import { Placeholder } from './features/shared/Placeholder'

export function App() {
  const [activeTab, setActiveTab] = useState<TabId>('commit')

  const currentTab = tabs.find((t) => t.id === activeTab)!

  return (
    <Layout>
      <div className="rounded-lg bg-white shadow">
        <TabNavigation activeTab={activeTab} onTabChange={setActiveTab} />

        <div className="p-6">
          <Placeholder
            title={currentTab.placeholder.title}
            subtitle={currentTab.placeholder.subtitle}
          />
        </div>
      </div>
    </Layout>
  )
}
