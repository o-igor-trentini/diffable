import { useState, type FormEvent } from 'react'
import { Zap } from 'lucide-react'
import { Button } from '../shared/Button'
import { TextInput } from '../shared/TextInput'
import { TextArea } from '../shared/TextArea'
import { AutocompleteInput } from '../shared/AutocompleteInput'
import { ModeToggle, type SourceMode } from '../shared/ModeToggle'
import { LevelSelector } from '../shared/LevelSelector'
import { useRepositories } from '@/lib/hooks/useRepositories'
import type { AnalyzePRRequest } from '@/lib/api/types'

interface PrFormProps {
  onSubmit: (req: AnalyzePRRequest) => void
  isPending: boolean
}

export function PrForm({ onSubmit, isPending }: PrFormProps) {
  const [mode, setMode] = useState<SourceMode>('bitbucket')
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

    if (mode === 'manual') {
      onSubmit({
        raw_diff: rawDiff.trim(),
        pr_title: prTitle.trim(),
        pr_description: prDescription.trim() || undefined,
        level,
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

  const isValid =
    mode === 'bitbucket'
      ? !!(workspace.trim() && repoSlug.trim() && prId.trim())
      : !!(rawDiff.trim() && prTitle.trim())

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <ModeToggle mode={mode} onModeChange={setMode} disabled={isPending} />

      {mode === 'bitbucket' ? (
        <div className="animate-fade-up space-y-4">
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <TextInput
              id="pr-workspace"
              label="Workspace"
              placeholder="minha-empresa"
              hint="Slug do workspace no Bitbucket"
              value={workspace}
              onChange={(e) => setWorkspace(e.target.value)}
              disabled={isPending}
            />
            <AutocompleteInput
              id="pr-repo"
              label="Repositorio"
              placeholder="Digite para buscar..."
              hint="Digite 2+ caracteres para buscar repositorios"
              value={repoSlug}
              onChange={setRepoSlug}
              onQueryChange={setRepoQuery}
              options={(repos || []).map((r) => ({ value: r.slug, label: r.name }))}
              loading={reposLoading}
              disabled={isPending}
              dependencyMet={workspace.trim().length >= 2}
              dependencyMessage="Preencha o workspace primeiro (min. 2 caracteres)"
            />
          </div>
          <TextInput
            id="pr-id"
            label="Numero do PR"
            type="number"
            placeholder="42"
            hint="Numero identificador do pull request"
            value={prId}
            onChange={(e) => setPrId(e.target.value)}
            disabled={isPending}
          />
        </div>
      ) : (
        <div className="animate-fade-up space-y-4">
          <TextInput
            id="pr-title"
            label="Titulo do PR"
            placeholder="feat: implementa autenticacao OAuth"
            hint="Titulo do pull request (obrigatorio)"
            value={prTitle}
            onChange={(e) => setPrTitle(e.target.value)}
            disabled={isPending}
          />
          <TextArea
            id="pr-description"
            label="Descricao do PR"
            placeholder="Descricao do pull request..."
            hint="Descricao original do PR (opcional, ajuda a gerar resultado melhor)"
            rows={3}
            value={prDescription}
            onChange={(e) => setPrDescription(e.target.value)}
            disabled={isPending}
          />
          <TextArea
            id="pr-raw-diff"
            label="Diff"
            placeholder={'Cole aqui o diff do pull request...\n\ndiff --git a/arquivo.ts b/arquivo.ts\n--- a/arquivo.ts\n+++ b/arquivo.ts'}
            hint="Cole o conteudo do diff (obrigatorio)"
            rows={8}
            value={rawDiff}
            onChange={(e) => setRawDiff(e.target.value)}
            disabled={isPending}
            className="font-mono text-xs"
          />
        </div>
      )}

      <div className="border-t border-stone-100 pt-6 dark:border-white/[0.04]">
        <LevelSelector value={level} onChange={setLevel} disabled={isPending} />
      </div>

      <Button type="submit" disabled={!isValid} loading={isPending}>
        <Zap size={16} />
        Gerar Descricao
      </Button>
    </form>
  )
}
