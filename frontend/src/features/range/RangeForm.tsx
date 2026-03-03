import { useState, type FormEvent } from 'react'
import { Zap } from 'lucide-react'
import { Button } from '../shared/Button'
import { TextInput } from '../shared/TextInput'
import { TextArea } from '../shared/TextArea'
import { AutocompleteInput } from '../shared/AutocompleteInput'
import { LevelSelector } from '../shared/LevelSelector'
import { AdvancedSettings } from '../shared/AdvancedSettings'
import { useRepositories } from '@/lib/hooks/useRepositories'
import type { AnalyzeRangeRequest, GenerationOverrides } from '@/lib/api/types'

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
  const [userContext, setUserContext] = useState('')
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

    const trimmedContext = userContext.trim()

    onSubmit({
      workspace: workspace.trim(),
      repo_slug: repoSlug.trim(),
      from_hash: fromHash.trim(),
      to_hash: toHash.trim(),
      level,
      ...(trimmedContext && { user_context: trimmedContext }),
      ...(hasOverrides && { overrides: effectiveOverrides }),
    })
  }

  const isValid =
    workspace.trim() && repoSlug.trim() && fromHash.trim() && toHash.trim()

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="space-y-4">
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <TextInput
            id="range-workspace"
            label="Workspace"
            placeholder="minha-empresa"
            hint="Slug do workspace no Bitbucket"
            value={workspace}
            onChange={(e) => setWorkspace(e.target.value)}
            disabled={isPending}
          />
          <AutocompleteInput
            id="range-repo"
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

        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <TextInput
            id="range-from"
            label="Hash Inicial (from)"
            placeholder="abc1234"
            hint="Commit de inicio do intervalo (exclusivo)"
            value={fromHash}
            onChange={(e) => setFromHash(e.target.value)}
            disabled={isPending}
            className="font-mono"
          />
          <TextInput
            id="range-to"
            label="Hash Final (to)"
            placeholder="def5678"
            hint="Commit final do intervalo (inclusivo)"
            value={toHash}
            onChange={(e) => setToHash(e.target.value)}
            disabled={isPending}
            className="font-mono"
          />
        </div>
      </div>

      <div className="border-t border-stone-100 pt-6 dark:border-white/[0.04]">
        <LevelSelector value={level} onChange={setLevel} disabled={isPending} />
      </div>

      <TextArea
        id="range-user-context"
        label="Contexto Adicional"
        placeholder="Ex: Este PR integra a data de emissão da CNH retornada pela API Nexus. O campo é opcional e pode vir vazio..."
        hint="Opcional. Forneça contexto sobre a tarefa para melhorar a descrição gerada. Especialmente útil no nível QA Detalhado."
        rows={3}
        value={userContext}
        onChange={(e) => setUserContext(e.target.value)}
        disabled={isPending}
      />

      <AdvancedSettings value={overrides} onChange={setOverrides} disabled={isPending} />

      <Button type="submit" disabled={!isValid} loading={isPending}>
        <Zap size={16} />
        Gerar Descricao
      </Button>
    </form>
  )
}
