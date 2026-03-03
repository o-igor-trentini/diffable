import { useQuery } from '@tanstack/react-query'
import { listAnalyses, getRefinements } from '../api/endpoints'
import type { HistoryFilter } from '../api/types'

export function useAnalysesList(filter?: HistoryFilter) {
  return useQuery({
    queryKey: ['analyses', filter],
    queryFn: () => listAnalyses(filter),
  })
}

export function useRefinements(analysisID: string) {
  return useQuery({
    queryKey: ['refinements', analysisID],
    queryFn: () => getRefinements(analysisID),
    enabled: !!analysisID,
  })
}
