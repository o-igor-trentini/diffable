export interface GenerationOverrides {
  max_tokens?: number
  temperature?: number
  model?: string
}

export interface AnalyzeCommitRequest {
  workspace?: string
  repo_slug?: string
  commit_hash?: string
  raw_diff?: string
  level?: string
  user_context?: string
  overrides?: GenerationOverrides
}

export interface AnalyzeRangeRequest {
  workspace: string
  repo_slug: string
  from_hash: string
  to_hash: string
  level?: string
  user_context?: string
  overrides?: GenerationOverrides
}

export interface AnalyzePRRequest {
  workspace?: string
  repo_slug?: string
  pr_id?: number
  raw_diff?: string
  pr_title?: string
  pr_description?: string
  level?: string
  user_context?: string
  overrides?: GenerationOverrides
}

export interface AnalysisResponse {
  id: string
  type: string
  level: string
  description: string
  model_used: string
  tokens_used: number
  created_at: string
}

export interface ErrorResponse {
  error: string
  message: string
  details?: string
}

export interface RefineRequest {
  instruction: string
}

export interface RefinementResponse {
  id: string
  analysis_id: string
  instruction: string
  refined_description: string
  model_used: string
  tokens_used: number
  created_at: string
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  page_size: number
}

export interface RepositoryResponse {
  slug: string
  name: string
  full_name: string
}

export interface HistoryFilter {
  type?: string
  page?: number
  page_size?: number
}
