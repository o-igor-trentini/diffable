import { useAnalyzeRange } from '@/lib/hooks/useAnalysis'
import { RangeForm } from './RangeForm'
import { ResultDisplay } from '../shared/ResultDisplay'
import { ErrorAlert } from '../shared/ErrorAlert'

export function RangeAnalysis() {
  const { mutate, data, isPending, isError, error } = useAnalyzeRange()

  return (
    <div>
      <RangeForm onSubmit={mutate} isPending={isPending} />
      {isError && <ErrorAlert message={error.message} />}
      {data && <ResultDisplay result={data} />}
    </div>
  )
}
