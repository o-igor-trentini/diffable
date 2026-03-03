import { RefreshCw, Sparkles } from 'lucide-react'
import { CopyButton } from './CopyButton'
import { ExportMarkdownButton } from './ExportMarkdownButton'
import type { AnalysisResponse } from '@/lib/api/types'

const levelLabels: Record<string, string> = {
  technical: 'Tecnico',
  functional: 'Funcional',
  executive: 'Executivo',
}

interface ResultDisplayProps {
  result: AnalysisResponse
  onRefine?: (result: AnalysisResponse) => void
}

export function ResultDisplay({ result, onRefine }: ResultDisplayProps) {
  return (
    <div className="animate-fade-up mt-6 rounded-xl border border-stone-200 bg-white p-5 dark:border-white/[0.06] dark:bg-white/[0.02]">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-2">
          <div className="flex h-6 w-6 items-center justify-center rounded-md bg-gradient-to-br from-violet-600 to-cyan-500">
            <Sparkles size={12} className="text-white" />
          </div>
          <h3 className="text-sm font-semibold text-stone-800 dark:text-stone-200">
            Descricao Gerada
          </h3>
        </div>
        <div className="flex items-center gap-1">
          {onRefine && (
            <button
              onClick={() => onRefine(result)}
              className="inline-flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm text-stone-500 transition-colors hover:bg-stone-100 hover:text-stone-700 dark:text-stone-400 dark:hover:bg-white/[0.06] dark:hover:text-stone-200"
              title="Refinar descricao"
            >
              <RefreshCw size={14} />
              <span>Refinar</span>
            </button>
          )}
          <CopyButton text={result.description} />
          <ExportMarkdownButton result={result} />
        </div>
      </div>

      <div className="mt-4 whitespace-pre-wrap rounded-lg bg-stone-50 p-4 font-mono text-[13px] leading-relaxed text-stone-700 dark:bg-white/[0.03] dark:text-stone-300">
        {result.description}
      </div>

      <div className="mt-4 flex flex-wrap items-center gap-3 border-t border-stone-100 pt-3 text-xs text-stone-400 dark:border-white/[0.04] dark:text-stone-500">
        <span className="rounded-md bg-stone-100 px-2 py-0.5 dark:bg-white/[0.06]">
          {levelLabels[result.level] || result.level}
        </span>
        <span className="font-mono">{result.model_used}</span>
        <span>{result.tokens_used} tokens</span>
      </div>
    </div>
  )
}
