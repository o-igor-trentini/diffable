import { useState, type FormEvent } from 'react'
import { Zap } from 'lucide-react'
import { Button } from '../shared/Button'
import { TextInput } from '../shared/TextInput'
import { TextArea } from '../shared/TextArea'
import { Select } from '../shared/Select'
import { AutocompleteInput } from '../shared/AutocompleteInput'
import { useRepositories } from '@/lib/hooks/useRepositories'
import type { AnalyzeCommitRequest } from '@/lib/api/types'

const levelOptions = [
  { value: 'functional', label: 'Funcional' },
  { value: 'technical', label: 'Técnico' },
  { value: 'executive', label: 'Executivo' },
]

interface CommitFormProps {
  onSubmit: (req: AnalyzeCommitRequest) => void
  isPending: boolean
}

export function CommitForm({ onSubmit, isPending }: CommitFormProps) {
  const [workspace, setWorkspace] = useState('')
  const [repoSlug, setRepoSlug] = useState('')
  const [commitHash, setCommitHash] = useState('')
  const [rawDiff, setRawDiff] = useState('')
  const [level, setLevel] = useState('functional')
  const [repoQuery, setRepoQuery] = useState('')
  const { data: repos, isLoading: reposLoading } = useRepositories(workspace, repoQuery)

  function handleSubmit(e: FormEvent) {
    e.preventDefault()

    if (rawDiff.trim()) {
      onSubmit({ raw_diff: rawDiff.trim(), level })
    } else {
      onSubmit({
        workspace: workspace.trim(),
        repo_slug: repoSlug.trim(),
        commit_hash: commitHash.trim(),
        level,
      })
    }
  }

  const hasHash = workspace.trim() && repoSlug.trim() && commitHash.trim()
  const hasRawDiff = rawDiff.trim()
  const isValid = hasHash || hasRawDiff

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
        <TextInput
          id="commit-workspace"
          label="Workspace"
          placeholder="meu-workspace"
          value={workspace}
          onChange={(e) => setWorkspace(e.target.value)}
          disabled={isPending}
        />
        <AutocompleteInput
          id="commit-repo"
          label="Repositório"
          placeholder="meu-repo"
          value={repoSlug}
          onChange={setRepoSlug}
          onQueryChange={setRepoQuery}
          options={(repos || []).map((r) => ({ value: r.slug, label: r.name }))}
          loading={reposLoading}
          disabled={isPending}
        />
        <TextInput
          id="commit-hash"
          label="Hash do Commit"
          placeholder="abc1234"
          value={commitHash}
          onChange={(e) => setCommitHash(e.target.value)}
          disabled={isPending}
        />
      </div>

      <div className="relative flex items-center py-2">
        <div className="flex-grow border-t border-gray-300 dark:border-gray-600" />
        <span className="mx-4 shrink-0 text-xs text-gray-500 dark:text-gray-400">OU cole o diff manualmente</span>
        <div className="flex-grow border-t border-gray-300 dark:border-gray-600" />
      </div>

      <TextArea
        id="commit-raw-diff"
        label="Diff (raw)"
        placeholder="Cole o diff aqui..."
        rows={6}
        value={rawDiff}
        onChange={(e) => setRawDiff(e.target.value)}
        disabled={isPending}
      />

      <Select
        id="commit-level"
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
