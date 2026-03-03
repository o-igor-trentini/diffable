import { RefreshCw } from 'lucide-react'
import { CopyButton } from './CopyButton'
import type { AnalysisResponse } from '@/lib/api/types'

interface ResultDisplayProps {
  result: AnalysisResponse
  onRefine?: (result: AnalysisResponse) => void
}

export function ResultDisplay({ result, onRefine }: ResultDisplayProps) {
  return (
    <div className="mt-6 animate-in fade-in duration-300 rounded-lg border border-gray-200 bg-gray-50 p-5 dark:border-gray-700 dark:bg-gray-700/50">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between sm:gap-4">
        <h3 className="text-sm font-semibold text-gray-800 dark:text-gray-200">Descrição Gerada</h3>
        <div className="flex items-center gap-2">
          {onRefine && (
            <button
              onClick={() => onRefine(result)}
              className="inline-flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm text-gray-600 hover:bg-gray-100 transition-colors dark:text-gray-300 dark:hover:bg-gray-600"
              title="Refinar descrição"
            >
              <RefreshCw size={16} />
              <span>Refinar</span>
            </button>
          )}
          <CopyButton text={result.description} />
        </div>
      </div>

      <div className="mt-3 whitespace-pre-wrap text-sm leading-relaxed text-gray-700 dark:text-gray-300">
        {result.description}
      </div>

      <div className="mt-4 flex items-center gap-4 border-t border-gray-200 pt-3 text-xs text-gray-500 dark:border-gray-600 dark:text-gray-400">
        <span>Modelo: {result.model_used}</span>
        <span>Tokens: {result.tokens_used}</span>
      </div>
    </div>
  )
}
