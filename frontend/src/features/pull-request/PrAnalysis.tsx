import { useAnalyzePR } from '@/lib/hooks/useAnalysis'
import { PrForm } from './PrForm'
import { ResultDisplay } from '../shared/ResultDisplay'
import { ErrorAlert } from '../shared/ErrorAlert'
import type { AnalysisResponse } from '@/lib/api/types'

interface PrAnalysisProps {
  onRefine?: (result: AnalysisResponse) => void
}

export function PrAnalysis({ onRefine }: PrAnalysisProps) {
  const { mutate, data, isPending, isError, error } = useAnalyzePR()

  return (
    <div>
      <PrForm onSubmit={mutate} isPending={isPending} />
      {isError && <ErrorAlert message={error.message} />}
      {data && <ResultDisplay result={data} onRefine={onRefine} />}
    </div>
  )
}
