import { useAnalyzeCommit } from '@/lib/hooks/useAnalysis'
import { CommitForm } from './CommitForm'
import { ResultDisplay } from '../shared/ResultDisplay'
import { ErrorAlert } from '../shared/ErrorAlert'

export function CommitAnalysis() {
  const { mutate, data, isPending, isError, error } = useAnalyzeCommit()

  return (
    <div>
      <CommitForm onSubmit={mutate} isPending={isPending} />
      {isError && <ErrorAlert message={error.message} />}
      {data && <ResultDisplay result={data} />}
    </div>
  )
}
