import { useState } from 'react'
import { History, ChevronLeft, ChevronRight } from 'lucide-react'
import { HistoryItem } from './HistoryItem'
import { LoadingSpinner } from '../shared/LoadingSpinner'
import { useAnalysesList } from '@/lib/hooks/useHistory'
import type { AnalysisResponse } from '@/lib/api/types'

interface HistoryPanelProps {
  onSelect: (analysis: AnalysisResponse) => void
}

const typeOptions = [
  { label: 'Todos', value: '' },
  { label: 'Commit', value: 'single_commit' },
  { label: 'Range', value: 'commit_range' },
  { label: 'PR', value: 'pull_request' },
]

const PAGE_SIZE = 10

export function HistoryPanel({ onSelect }: HistoryPanelProps) {
  const [open, setOpen] = useState(false)
  const [typeFilter, setTypeFilter] = useState('')
  const [page, setPage] = useState(1)

  const { data, isLoading } = useAnalysesList(
    open ? { type: typeFilter || undefined, page, page_size: PAGE_SIZE } : undefined,
  )

  const totalPages = data ? Math.ceil(data.total / PAGE_SIZE) : 0

  if (!open) {
    return (
      <button
        onClick={() => setOpen(true)}
        className="inline-flex items-center gap-2 rounded-md px-3 py-2 text-sm text-gray-600 hover:bg-gray-100 transition-colors"
      >
        <History size={16} />
        Histórico
      </button>
    )
  }

  return (
    <div className="rounded-lg border border-gray-200 bg-white">
      <div className="flex items-center justify-between border-b border-gray-200 px-4 py-3">
        <div className="flex items-center gap-2">
          <History size={16} className="text-gray-500" />
          <span className="text-sm font-medium text-gray-700">Histórico</span>
        </div>
        <button
          onClick={() => setOpen(false)}
          className="text-xs text-gray-400 hover:text-gray-600"
        >
          Fechar
        </button>
      </div>

      <div className="border-b border-gray-200 px-4 py-2">
        <div className="flex gap-1">
          {typeOptions.map((opt) => (
            <button
              key={opt.value}
              onClick={() => { setTypeFilter(opt.value); setPage(1) }}
              className={`rounded-md px-2 py-1 text-xs transition-colors ${
                typeFilter === opt.value
                  ? 'bg-blue-100 text-blue-700'
                  : 'text-gray-500 hover:bg-gray-100'
              }`}
            >
              {opt.label}
            </button>
          ))}
        </div>
      </div>

      <div className="max-h-80 overflow-y-auto p-2">
        {isLoading && (
          <div className="flex justify-center py-8">
            <LoadingSpinner />
          </div>
        )}

        {!isLoading && data && data.data.length === 0 && (
          <p className="py-8 text-center text-sm text-gray-400">
            Nenhuma análise encontrada
          </p>
        )}

        {!isLoading && data && data.data.map((analysis) => (
          <HistoryItem key={analysis.id} analysis={analysis} onClick={onSelect} />
        ))}
      </div>

      {totalPages > 1 && (
        <div className="flex items-center justify-between border-t border-gray-200 px-4 py-2">
          <button
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            disabled={page <= 1}
            className="inline-flex items-center gap-1 text-xs text-gray-500 hover:text-gray-700 disabled:text-gray-300"
          >
            <ChevronLeft size={14} />
            Anterior
          </button>
          <span className="text-xs text-gray-400">
            {page} / {totalPages}
          </span>
          <button
            onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
            disabled={page >= totalPages}
            className="inline-flex items-center gap-1 text-xs text-gray-500 hover:text-gray-700 disabled:text-gray-300"
          >
            Próximo
            <ChevronRight size={14} />
          </button>
        </div>
      )}
    </div>
  )
}
