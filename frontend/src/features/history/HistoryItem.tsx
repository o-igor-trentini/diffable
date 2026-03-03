import { GitCommitHorizontal, GitBranch, GitPullRequest } from 'lucide-react'
import type { AnalysisResponse } from '@/lib/api/types'

interface HistoryItemProps {
  analysis: AnalysisResponse
  onClick: (analysis: AnalysisResponse) => void
}

const typeConfig: Record<string, { icon: React.ReactNode; label: string; color: string }> = {
  single_commit: {
    icon: <GitCommitHorizontal size={14} />,
    label: 'Commit',
    color: 'bg-blue-100 text-blue-600 dark:bg-blue-500/15 dark:text-blue-400',
  },
  commit_range: {
    icon: <GitBranch size={14} />,
    label: 'Range',
    color: 'bg-purple-100 text-purple-600 dark:bg-purple-500/15 dark:text-purple-400',
  },
  pull_request: {
    icon: <GitPullRequest size={14} />,
    label: 'PR',
    color: 'bg-emerald-100 text-emerald-600 dark:bg-emerald-500/15 dark:text-emerald-400',
  },
}

function timeAgo(dateStr: string): string {
  const date = new Date(dateStr)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMin = Math.floor(diffMs / 60000)

  if (diffMin < 1) return 'agora'
  if (diffMin < 60) return `${diffMin}min`

  const diffHours = Math.floor(diffMin / 60)
  if (diffHours < 24) return `${diffHours}h`

  const diffDays = Math.floor(diffHours / 24)
  if (diffDays < 30) return `${diffDays}d`

  return date.toLocaleDateString('pt-BR')
}

export function HistoryItem({ analysis, onClick }: HistoryItemProps) {
  const config = typeConfig[analysis.type] || typeConfig.single_commit
  const truncated =
    analysis.description.length > 120
      ? analysis.description.slice(0, 120) + '...'
      : analysis.description

  return (
    <button
      onClick={() => onClick(analysis)}
      className="flex w-full items-start gap-3 rounded-lg px-3 py-2.5 text-left transition-colors hover:bg-stone-50 dark:hover:bg-white/[0.03]"
    >
      <div className={`mt-0.5 flex shrink-0 items-center gap-1 rounded-md px-1.5 py-0.5 text-[10px] font-semibold ${config.color}`}>
        {config.icon}
        <span>{config.label}</span>
      </div>
      <div className="min-w-0 flex-1">
        <p className="text-sm leading-snug text-stone-600 dark:text-stone-300">{truncated}</p>
        <p className="mt-1 font-mono text-[10px] text-stone-300 dark:text-stone-600">
          {timeAgo(analysis.created_at)}
        </p>
      </div>
    </button>
  )
}
