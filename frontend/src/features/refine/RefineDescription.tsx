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
      <div className="py-12 text-center text-gray-500">
        <p className="text-lg font-medium">Refinar Descrição</p>
        <p className="mt-1 text-sm">
          Gere uma descrição primeiro e clique em &quot;Refinar&quot; para ajustá-la.
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
        <div className="mt-6 animate-in fade-in duration-300 rounded-lg border border-gray-200 bg-gray-50 p-5">
          <div className="flex items-start justify-between gap-4">
            <h3 className="text-sm font-semibold text-gray-800">Descrição Refinada</h3>
            <CopyButton text={data.refined_description} />
          </div>

          <div className="mt-3 whitespace-pre-wrap text-sm leading-relaxed text-gray-700">
            {data.refined_description}
          </div>

          <div className="mt-4 flex items-center gap-4 border-t border-gray-200 pt-3 text-xs text-gray-500">
            <span>Modelo: {data.model_used}</span>
            <span>Tokens: {data.tokens_used}</span>
          </div>
        </div>
      )}
    </div>
  )
}
