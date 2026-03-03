import { GitCommitHorizontal, GitBranch, GitPullRequest } from 'lucide-react'
import type { AnalysisResponse } from '@/lib/api/types'

interface HistoryItemProps {
  analysis: AnalysisResponse
  onClick: (analysis: AnalysisResponse) => void
}

function typeIcon(type: string) {
  switch (type) {
    case 'single_commit':
      return <GitCommitHorizontal size={16} className="text-blue-500" />
    case 'commit_range':
      return <GitBranch size={16} className="text-purple-500" />
    case 'pull_request':
      return <GitPullRequest size={16} className="text-green-500" />
    default:
      return <GitCommitHorizontal size={16} className="text-gray-400" />
  }
}

function timeAgo(dateStr: string): string {
  const date = new Date(dateStr)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMin = Math.floor(diffMs / 60000)

  if (diffMin < 1) return 'agora'
  if (diffMin < 60) return `há ${diffMin} min`

  const diffHours = Math.floor(diffMin / 60)
  if (diffHours < 24) return `há ${diffHours}h`

  const diffDays = Math.floor(diffHours / 24)
  if (diffDays < 30) return `há ${diffDays}d`

  return date.toLocaleDateString('pt-BR')
}

export function HistoryItem({ analysis, onClick }: HistoryItemProps) {
  const truncated =
    analysis.description.length > 100
      ? analysis.description.slice(0, 100) + '...'
      : analysis.description

  return (
    <button
      onClick={() => onClick(analysis)}
      className="flex w-full items-start gap-3 rounded-md px-3 py-2 text-left transition-colors hover:bg-gray-100 dark:hover:bg-gray-700"
    >
      <div className="mt-0.5 shrink-0">{typeIcon(analysis.type)}</div>
      <div className="min-w-0 flex-1">
        <p className="truncate text-sm text-gray-700 dark:text-gray-300">{truncated}</p>
        <p className="mt-0.5 text-xs text-gray-400 dark:text-gray-500">{timeAgo(analysis.created_at)}</p>
      </div>
    </button>
  )
}
