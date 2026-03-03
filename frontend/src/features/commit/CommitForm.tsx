import { useState, type FormEvent } from 'react'
import { Zap } from 'lucide-react'
import { Button } from '../shared/Button'
import { TextInput } from '../shared/TextInput'
import { TextArea } from '../shared/TextArea'
import { AutocompleteInput } from '../shared/AutocompleteInput'
import { ModeToggle, type SourceMode } from '../shared/ModeToggle'
import { LevelSelector } from '../shared/LevelSelector'
import { AdvancedSettings } from '../shared/AdvancedSettings'
import { useRepositories } from '@/lib/hooks/useRepositories'
import type { AnalyzeCommitRequest, GenerationOverrides } from '@/lib/api/types'

interface CommitFormProps {
  onSubmit: (req: AnalyzeCommitRequest) => void
  isPending: boolean
}

export function CommitForm({ onSubmit, isPending }: CommitFormProps) {
  const [mode, setMode] = useState<SourceMode>('bitbucket')
  const [workspace, setWorkspace] = useState('')
  const [repoSlug, setRepoSlug] = useState('')
  const [commitHash, setCommitHash] = useState('')
  const [rawDiff, setRawDiff] = useState('')
  const [level, setLevel] = useState('functional')
  const [overrides, setOverrides] = useState<GenerationOverrides>({})
  const [repoQuery, setRepoQuery] = useState('')
  const { data: repos, isLoading: reposLoading } = useRepositories(workspace, repoQuery)

  function handleSubmit(e: FormEvent) {
    e.preventDefault()

    const effectiveOverrides: GenerationOverrides = {}
    if (overrides.temperature !== undefined && overrides.temperature !== 0.3)
      effectiveOverrides.temperature = overrides.temperature
    if (overrides.max_tokens !== undefined && overrides.max_tokens !== 1024)
      effectiveOverrides.max_tokens = overrides.max_tokens
    if (overrides.model !== undefined && overrides.model !== 'auto')
      effectiveOverrides.model = overrides.model

    const hasOverrides = Object.keys(effectiveOverrides).length > 0

    if (mode === 'manual') {
      onSubmit({
        raw_diff: rawDiff.trim(),
        level,
        ...(hasOverrides && { overrides: effectiveOverrides }),
      })
    } else {
      onSubmit({
        workspace: workspace.trim(),
        repo_slug: repoSlug.trim(),
        commit_hash: commitHash.trim(),
        level,
        ...(hasOverrides && { overrides: effectiveOverrides }),
      })
    }
  }

  const isValid =
    mode === 'bitbucket'
      ? !!(workspace.trim() && repoSlug.trim() && commitHash.trim())
      : !!rawDiff.trim()

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <ModeToggle mode={mode} onModeChange={setMode} disabled={isPending} />

      {mode === 'bitbucket' ? (
        <div className="animate-fade-up space-y-4">
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <TextInput
              id="commit-workspace"
              label="Workspace"
              placeholder="minha-empresa"
              hint="Slug do workspace no Bitbucket"
              value={workspace}
              onChange={(e) => setWorkspace(e.target.value)}
              disabled={isPending}
            />
            <AutocompleteInput
              id="commit-repo"
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
            id="commit-hash"
            label="Hash do Commit"
            placeholder="abc1234def5678"
            hint="SHA completo ou abreviado do commit"
            value={commitHash}
            onChange={(e) => setCommitHash(e.target.value)}
            disabled={isPending}
            className="font-mono"
          />
        </div>
      ) : (
        <div className="animate-fade-up">
          <TextArea
            id="commit-raw-diff"
            label="Diff"
            placeholder={'Cole aqui a saida do comando git diff...\n\ndiff --git a/arquivo.ts b/arquivo.ts\n--- a/arquivo.ts\n+++ b/arquivo.ts\n@@ -1,3 +1,4 @@'}
            hint="Cole o conteudo do diff diretamente (saida de git diff, git show, etc.)"
            rows={10}
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

      <AdvancedSettings value={overrides} onChange={setOverrides} disabled={isPending} />

      <Button type="submit" disabled={!isValid} loading={isPending}>
        <Zap size={16} />
        Gerar Descricao
      </Button>
    </form>
  )
}
