import { useMutation } from '@tanstack/react-query'
import { analyzeCommit, analyzeRange, analyzePR } from '../api/endpoints'
import type {
  AnalyzeCommitRequest,
  AnalyzeRangeRequest,
  AnalyzePRRequest,
} from '../api/types'

export function useAnalyzeCommit() {
  return useMutation({
    mutationFn: (req: AnalyzeCommitRequest) => analyzeCommit(req),
  })
}

export function useAnalyzeRange() {
  return useMutation({
    mutationFn: (req: AnalyzeRangeRequest) => analyzeRange(req),
  })
}

export function useAnalyzePR() {
  return useMutation({
    mutationFn: (req: AnalyzePRRequest) => analyzePR(req),
  })
}
