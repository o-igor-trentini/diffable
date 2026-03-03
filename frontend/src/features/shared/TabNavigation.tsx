import { GitCommitHorizontal, GitBranch, GitPullRequest, RefreshCw } from 'lucide-react'
import { TabButton } from './TabButton'

export type TabId = 'commit' | 'range' | 'pr' | 'refine'

export interface TabConfig {
  id: TabId
  label: string
  icon: React.ReactNode
  placeholder: { title: string; subtitle: string }
}

export const tabs: TabConfig[] = [
  { id: 'commit', label: 'Commit', icon: <GitCommitHorizontal size={18} />, placeholder: { title: 'Análise de Commit', subtitle: 'Em breve — Fase 4' } },
  { id: 'range', label: 'Range', icon: <GitBranch size={18} />, placeholder: { title: 'Análise de Range', subtitle: 'Em breve — Fase 4' } },
  { id: 'pr', label: 'PR', icon: <GitPullRequest size={18} />, placeholder: { title: 'Análise de PR', subtitle: 'Em breve — Fase 4' } },
  { id: 'refine', label: 'Refinar', icon: <RefreshCw size={18} />, placeholder: { title: 'Refinar Descrição', subtitle: 'Em breve — Fase 5' } },
]

interface TabNavigationProps {
  activeTab: TabId
  onTabChange: (tab: TabId) => void
}

export function TabNavigation({ activeTab, onTabChange }: TabNavigationProps) {
  return (
    <nav className="flex overflow-x-auto border-b border-gray-200 dark:border-gray-700" role="tablist">
      {tabs.map((tab) => (
        <TabButton
          key={tab.id}
          active={activeTab === tab.id}
          icon={tab.icon}
          label={tab.label}
          onClick={() => onTabChange(tab.id)}
        />
      ))}
    </nav>
  )
}
