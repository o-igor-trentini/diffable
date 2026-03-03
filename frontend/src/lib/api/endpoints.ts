import apiClient from './client'
import type {
  AnalyzeCommitRequest,
  AnalyzeRangeRequest,
  AnalyzePRRequest,
  AnalysisResponse,
  RefinementResponse,
  PaginatedResponse,
  HistoryFilter,
  RepositoryResponse,
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

export async function refineDescription(id: string, instruction: string): Promise<RefinementResponse> {
  const { data } = await apiClient.post<RefinementResponse>(`/analyses/${id}/refine`, { instruction })
  return data
}

export async function listAnalyses(filter?: HistoryFilter): Promise<PaginatedResponse<AnalysisResponse>> {
  const params = new URLSearchParams()
  if (filter?.type) params.set('type', filter.type)
  if (filter?.page) params.set('page', String(filter.page))
  if (filter?.page_size) params.set('page_size', String(filter.page_size))
  const query = params.toString()
  const { data } = await apiClient.get<PaginatedResponse<AnalysisResponse>>(`/analyses${query ? `?${query}` : ''}`)
  return data
}

export async function getRefinements(id: string): Promise<RefinementResponse[]> {
  const { data } = await apiClient.get<RefinementResponse[]>(`/analyses/${id}/refinements`)
  return data
}

export async function listRepositories(workspace: string, query?: string): Promise<RepositoryResponse[]> {
  const params = new URLSearchParams({ workspace })
  if (query) params.set('q', query)
  const { data } = await apiClient.get<RepositoryResponse[]>(`/bitbucket/repositories?${params.toString()}`)
  return data
}
