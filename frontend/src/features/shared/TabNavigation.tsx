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
  { id: 'commit', label: 'Commit', icon: <GitCommitHorizontal size={16} />, placeholder: { title: 'Analise de Commit', subtitle: 'Commit individual' } },
  { id: 'range', label: 'Range', icon: <GitBranch size={16} />, placeholder: { title: 'Analise de Range', subtitle: 'Intervalo de commits' } },
  { id: 'pr', label: 'Pull Request', icon: <GitPullRequest size={16} />, placeholder: { title: 'Analise de PR', subtitle: 'Pull request' } },
  { id: 'refine', label: 'Refinar', icon: <RefreshCw size={16} />, placeholder: { title: 'Refinar Descricao', subtitle: 'Ajustar resultado' } },
]

interface TabNavigationProps {
  activeTab: TabId
  onTabChange: (tab: TabId) => void
}

export function TabNavigation({ activeTab, onTabChange }: TabNavigationProps) {
  return (
    <nav className="flex gap-0.5 overflow-x-auto border-b border-stone-200 px-2 dark:border-white/[0.06]" role="tablist">
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
