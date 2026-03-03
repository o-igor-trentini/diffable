import { useState, type FormEvent } from 'react'
import { Zap } from 'lucide-react'
import { Button } from '../shared/Button'
import { TextInput } from '../shared/TextInput'
import { TextArea } from '../shared/TextArea'
import { Select } from '../shared/Select'
import { AutocompleteInput } from '../shared/AutocompleteInput'
import { useRepositories } from '@/lib/hooks/useRepositories'
import type { AnalyzePRRequest } from '@/lib/api/types'

const levelOptions = [
  { value: 'functional', label: 'Funcional' },
  { value: 'technical', label: 'Técnico' },
  { value: 'executive', label: 'Executivo' },
]

interface PrFormProps {
  onSubmit: (req: AnalyzePRRequest) => void
  isPending: boolean
}

export function PrForm({ onSubmit, isPending }: PrFormProps) {
  const [workspace, setWorkspace] = useState('')
  const [repoSlug, setRepoSlug] = useState('')
  const [prId, setPrId] = useState('')
  const [rawDiff, setRawDiff] = useState('')
  const [prTitle, setPrTitle] = useState('')
  const [prDescription, setPrDescription] = useState('')
  const [level, setLevel] = useState('functional')
  const [repoQuery, setRepoQuery] = useState('')
  const { data: repos, isLoading: reposLoading } = useRepositories(workspace, repoQuery)

  function handleSubmit(e: FormEvent) {
    e.preventDefault()

    if (rawDiff.trim()) {
      onSubmit({
        raw_diff: rawDiff.trim(),
        pr_title: prTitle.trim(),
        pr_description: prDescription.trim(),
        level,
        ...(prId.trim() && workspace.trim() && repoSlug.trim()
          ? {
              workspace: workspace.trim(),
              repo_slug: repoSlug.trim(),
              pr_id: parseInt(prId, 10),
            }
          : {}),
      })
    } else {
      onSubmit({
        workspace: workspace.trim(),
        repo_slug: repoSlug.trim(),
        pr_id: parseInt(prId, 10),
        level,
      })
    }
  }

  const hasPRID = workspace.trim() && repoSlug.trim() && prId.trim()
  const hasRawDiff = rawDiff.trim() && prTitle.trim()
  const isValid = hasPRID || hasRawDiff

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
        <TextInput
          id="pr-workspace"
          label="Workspace"
          placeholder="meu-workspace"
          value={workspace}
          onChange={(e) => setWorkspace(e.target.value)}
          disabled={isPending}
        />
        <AutocompleteInput
          id="pr-repo"
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
          id="pr-id"
          label="PR ID"
          type="number"
          placeholder="123"
          value={prId}
          onChange={(e) => setPrId(e.target.value)}
          disabled={isPending}
        />
      </div>

      <div className="relative flex items-center py-2">
        <div className="flex-grow border-t border-gray-300 dark:border-gray-600" />
        <span className="mx-4 shrink-0 text-xs text-gray-500 dark:text-gray-400">OU cole o diff manualmente</span>
        <div className="flex-grow border-t border-gray-300 dark:border-gray-600" />
      </div>

      <TextArea
        id="pr-raw-diff"
        label="Diff (raw)"
        placeholder="Cole o diff aqui..."
        rows={6}
        value={rawDiff}
        onChange={(e) => setRawDiff(e.target.value)}
        disabled={isPending}
      />

      <TextInput
        id="pr-title"
        label="Título do PR"
        placeholder="feat: adiciona autenticação OAuth"
        value={prTitle}
        onChange={(e) => setPrTitle(e.target.value)}
        disabled={isPending}
      />

      <TextArea
        id="pr-description"
        label="Descrição do PR (opcional)"
        placeholder="Descrição do pull request..."
        rows={3}
        value={prDescription}
        onChange={(e) => setPrDescription(e.target.value)}
        disabled={isPending}
      />

      <Select
        id="pr-level"
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
