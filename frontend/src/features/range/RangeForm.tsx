import { useState, type FormEvent } from 'react'
import { Zap } from 'lucide-react'
import { Button } from '../shared/Button'
import { TextInput } from '../shared/TextInput'
import { Select } from '../shared/Select'
import { AutocompleteInput } from '../shared/AutocompleteInput'
import { useRepositories } from '@/lib/hooks/useRepositories'
import type { AnalyzeRangeRequest } from '@/lib/api/types'

const levelOptions = [
  { value: 'functional', label: 'Funcional' },
  { value: 'technical', label: 'Técnico' },
  { value: 'executive', label: 'Executivo' },
]

interface RangeFormProps {
  onSubmit: (req: AnalyzeRangeRequest) => void
  isPending: boolean
}

export function RangeForm({ onSubmit, isPending }: RangeFormProps) {
  const [workspace, setWorkspace] = useState('')
  const [repoSlug, setRepoSlug] = useState('')
  const [fromHash, setFromHash] = useState('')
  const [toHash, setToHash] = useState('')
  const [level, setLevel] = useState('functional')
  const [repoQuery, setRepoQuery] = useState('')
  const { data: repos, isLoading: reposLoading } = useRepositories(workspace, repoQuery)

  function handleSubmit(e: FormEvent) {
    e.preventDefault()
    onSubmit({
      workspace: workspace.trim(),
      repo_slug: repoSlug.trim(),
      from_hash: fromHash.trim(),
      to_hash: toHash.trim(),
      level,
    })
  }

  const isValid =
    workspace.trim() && repoSlug.trim() && fromHash.trim() && toHash.trim()

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <TextInput
          id="range-workspace"
          label="Workspace"
          placeholder="meu-workspace"
          value={workspace}
          onChange={(e) => setWorkspace(e.target.value)}
          disabled={isPending}
        />
        <AutocompleteInput
          id="range-repo"
          label="Repositório"
          placeholder="meu-repo"
          value={repoSlug}
          onChange={setRepoSlug}
          onQueryChange={setRepoQuery}
          options={(repos || []).map((r) => ({ value: r.slug, label: r.name }))}
          loading={reposLoading}
          disabled={isPending}
        />
      </div>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <TextInput
          id="range-from"
          label="Hash Inicial (from)"
          placeholder="abc1234"
          value={fromHash}
          onChange={(e) => setFromHash(e.target.value)}
          disabled={isPending}
        />
        <TextInput
          id="range-to"
          label="Hash Final (to)"
          placeholder="def5678"
          value={toHash}
          onChange={(e) => setToHash(e.target.value)}
          disabled={isPending}
        />
      </div>

      <Select
        id="range-level"
        label="Nível da Descrição"
        options={levelOptions}
        value={level}
        onChange={(e) => setLevel(e.target.value)}
        disabled={isPending}
      />

      <Button type="submit" disabled={!isValid} loading={isPending}>
        <Zap size={16} />
        Gerar Descrição
      </Button>
    </form>
  )
}
