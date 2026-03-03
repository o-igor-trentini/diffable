import { useQuery } from '@tanstack/react-query'
import { listRepositories } from '../api/endpoints'

export function useRepositories(workspace: string, query: string) {
  return useQuery({
    queryKey: ['repositories', workspace, query],
    queryFn: () => listRepositories(workspace, query),
    enabled: workspace.length >= 2 && query.length >= 2,
    staleTime: 5 * 60 * 1000,
  })
}
