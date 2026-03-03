import { useState } from 'react'
import { History, ChevronLeft, ChevronRight, ChevronDown, ChevronUp } from 'lucide-react'
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
  const [open, setOpen] = useState(true)
  const [typeFilter, setTypeFilter] = useState('')
  const [page, setPage] = useState(1)

  const { data, isLoading } = useAnalysesList(
    { type: typeFilter || undefined, page, page_size: PAGE_SIZE },
  )

  const totalPages = data ? Math.ceil(data.total / PAGE_SIZE) : 0

  return (
    <div className="rounded-xl border border-stone-200 bg-white dark:border-white/[0.06] dark:bg-white/[0.02]">
      <button
        onClick={() => setOpen((o) => !o)}
        className="flex w-full items-center justify-between px-5 py-3.5 text-left"
      >
        <div className="flex items-center gap-2.5">
          <History size={16} className="text-stone-400 dark:text-stone-500" />
          <span className="text-sm font-semibold text-stone-700 dark:text-stone-300">
            Historico de Analises
          </span>
          {data && data.total > 0 && (
            <span className="rounded-full bg-stone-100 px-2 py-0.5 text-[10px] font-medium text-stone-500 dark:bg-white/[0.06] dark:text-stone-400">
              {data.total}
            </span>
          )}
        </div>
        {open ? (
          <ChevronUp size={16} className="text-stone-400 dark:text-stone-500" />
        ) : (
          <ChevronDown size={16} className="text-stone-400 dark:text-stone-500" />
        )}
      </button>

      {open && (
        <>
          <div className="border-t border-stone-100 px-5 py-2.5 dark:border-white/[0.04]">
            <div className="flex flex-wrap gap-1">
              {typeOptions.map((opt) => (
                <button
                  key={opt.value}
                  onClick={() => { setTypeFilter(opt.value); setPage(1) }}
                  className={`rounded-lg px-2.5 py-1 text-xs font-medium transition-colors ${
                    typeFilter === opt.value
                      ? 'bg-violet-100 text-violet-700 dark:bg-violet-500/15 dark:text-violet-300'
                      : 'text-stone-400 hover:bg-stone-100 hover:text-stone-600 dark:text-stone-500 dark:hover:bg-white/[0.06] dark:hover:text-stone-300'
                  }`}
                >
                  {opt.label}
                </button>
              ))}
            </div>
          </div>

          <div className="max-h-80 overflow-y-auto border-t border-stone-100 px-2 py-1 dark:border-white/[0.04]">
            {isLoading && (
              <div className="flex justify-center py-10">
                <LoadingSpinner />
              </div>
            )}

            {!isLoading && data && data.data.length === 0 && (
              <p className="py-10 text-center text-sm text-stone-400 dark:text-stone-500">
                Nenhuma analise encontrada
              </p>
            )}

            {!isLoading && data && data.data.map((analysis) => (
              <HistoryItem key={analysis.id} analysis={analysis} onClick={onSelect} />
            ))}
          </div>

          {totalPages > 1 && (
            <div className="flex items-center justify-between border-t border-stone-100 px-5 py-2.5 dark:border-white/[0.04]">
              <button
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                disabled={page <= 1}
                className="inline-flex items-center gap-1 text-xs font-medium text-stone-400 hover:text-stone-600 disabled:text-stone-200 dark:text-stone-500 dark:hover:text-stone-300 dark:disabled:text-stone-700"
              >
                <ChevronLeft size={14} />
                Anterior
              </button>
              <span className="font-mono text-xs text-stone-400 dark:text-stone-500">
                {page} / {totalPages}
              </span>
              <button
                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                disabled={page >= totalPages}
                className="inline-flex items-center gap-1 text-xs font-medium text-stone-400 hover:text-stone-600 disabled:text-stone-200 dark:text-stone-500 dark:hover:text-stone-300 dark:disabled:text-stone-700"
              >
                Proximo
                <ChevronRight size={14} />
              </button>
            </div>
          )}
        </>
      )}
    </div>
  )
}
