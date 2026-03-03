import { useAnalyzeCommit } from '@/lib/hooks/useAnalysis'
import { CommitForm } from './CommitForm'
import { ResultDisplay } from '../shared/ResultDisplay'
import { ErrorAlert } from '../shared/ErrorAlert'
import type { AnalysisResponse } from '@/lib/api/types'

interface CommitAnalysisProps {
  onRefine?: (result: AnalysisResponse) => void
}

export function CommitAnalysis({ onRefine }: CommitAnalysisProps) {
  const { mutate, data, isPending, isError, error } = useAnalyzeCommit()

  return (
    <div>
      <CommitForm onSubmit={mutate} isPending={isPending} />
      {isError && <ErrorAlert message={error.message} />}
      {data && <ResultDisplay result={data} onRefine={onRefine} />}
    </div>
  )
}
