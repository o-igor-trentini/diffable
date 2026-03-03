import { RefreshCw } from 'lucide-react'
import { RefineForm } from './RefineForm'
import { CopyButton } from '../shared/CopyButton'
import { ErrorAlert } from '../shared/ErrorAlert'
import { useRefineDescription } from '@/lib/hooks/useAnalysis'
import type { AnalysisResponse } from '@/lib/api/types'

interface RefineDescriptionProps {
  analysis?: AnalysisResponse | null
}

export function RefineDescription({ analysis }: RefineDescriptionProps) {
  const { mutate, data, isPending, isError, error } = useRefineDescription()

  function handleRefine(instruction: string) {
    if (!analysis) return
    mutate({ id: analysis.id, instruction })
  }

  if (!analysis) {
    return (
      <div className="flex flex-col items-center justify-center py-16 text-center">
        <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-stone-100 dark:bg-white/[0.04]">
          <RefreshCw size={20} className="text-stone-300 dark:text-stone-600" />
        </div>
        <p className="mt-4 text-sm font-medium text-stone-500 dark:text-stone-400">
          Nenhuma descricao para refinar
        </p>
        <p className="mt-1 max-w-xs text-xs leading-relaxed text-stone-400 dark:text-stone-500">
          Gere uma descricao em qualquer aba e clique em &quot;Refinar&quot; para ajusta-la aqui.
        </p>
      </div>
    )
  }

  return (
    <div>
      <RefineForm
        initialDescription={analysis.description}
        onSubmit={handleRefine}
        isPending={isPending}
      />

      {isError && <ErrorAlert message={error.message} />}

      {data && (
        <div className="animate-fade-up mt-6 rounded-xl border border-stone-200 bg-white p-5 dark:border-white/[0.06] dark:bg-white/[0.02]">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <h3 className="text-sm font-semibold text-stone-800 dark:text-stone-200">
              Descricao Refinada
            </h3>
            <CopyButton text={data.refined_description} />
          </div>

          <div className="mt-4 whitespace-pre-wrap rounded-lg bg-stone-50 p-4 font-mono text-[13px] leading-relaxed text-stone-700 dark:bg-white/[0.03] dark:text-stone-300">
            {data.refined_description}
          </div>

          <div className="mt-4 flex items-center gap-3 border-t border-stone-100 pt-3 text-xs text-stone-400 dark:border-white/[0.04] dark:text-stone-500">
            <span className="font-mono">{data.model_used}</span>
            <span>{data.tokens_used} tokens</span>
          </div>
        </div>
      )}
    </div>
  )
}
