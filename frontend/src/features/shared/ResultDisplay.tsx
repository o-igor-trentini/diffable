import { CopyButton } from './CopyButton'
import type { AnalysisResponse } from '@/lib/api/types'

interface ResultDisplayProps {
  result: AnalysisResponse
}

export function ResultDisplay({ result }: ResultDisplayProps) {
  return (
    <div className="mt-6 animate-in fade-in duration-300 rounded-lg border border-gray-200 bg-gray-50 p-5">
      <div className="flex items-start justify-between gap-4">
        <h3 className="text-sm font-semibold text-gray-800">Descrição Gerada</h3>
        <CopyButton text={result.description} />
      </div>

      <div className="mt-3 whitespace-pre-wrap text-sm leading-relaxed text-gray-700">
        {result.description}
      </div>

      <div className="mt-4 flex items-center gap-4 border-t border-gray-200 pt-3 text-xs text-gray-500">
        <span>Modelo: {result.model_used}</span>
        <span>Tokens: {result.tokens_used}</span>
      </div>
    </div>
  )
}
