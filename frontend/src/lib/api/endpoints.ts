import apiClient from './client'
import type {
  AnalyzeCommitRequest,
  AnalyzeRangeRequest,
  AnalyzePRRequest,
  AnalysisResponse,
} from './types'

export async function analyzeCommit(req: AnalyzeCommitRequest): Promise<AnalysisResponse> {
  const { data } = await apiClient.post<AnalysisResponse>('/analyses/commit', req)
  return data
}

export async function analyzeRange(req: AnalyzeRangeRequest): Promise<AnalysisResponse> {
  const { data } = await apiClient.post<AnalysisResponse>('/analyses/range', req)
  return data
}

export async function analyzePR(req: AnalyzePRRequest): Promise<AnalysisResponse> {
  const { data } = await apiClient.post<AnalysisResponse>('/analyses/pr', req)
  return data
}

export async function getAnalysis(id: string): Promise<AnalysisResponse> {
  const { data } = await apiClient.get<AnalysisResponse>(`/analyses/${id}`)
  return data
}
