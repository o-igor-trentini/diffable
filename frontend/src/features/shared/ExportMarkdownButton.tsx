import { useState } from 'react'
import { Download, Check } from 'lucide-react'
import type { AnalysisResponse } from '@/lib/api/types'

const levelLabels: Record<string, string> = {
  technical: 'Tecnico',
  functional: 'Funcional',
  executive: 'Executivo',
}

const typeLabels: Record<string, string> = {
  single_commit: 'Commit',
  commit_range: 'Range de Commits',
  pull_request: 'Pull Request',
}

interface ExportMarkdownButtonProps {
  result: AnalysisResponse
}

export function ExportMarkdownButton({ result }: ExportMarkdownButtonProps) {
  const [exported, setExported] = useState(false)

  function handleExport() {
    const date = new Date(result.created_at).toLocaleDateString('pt-BR')
    const level = levelLabels[result.level] || result.level
    const type = typeLabels[result.type] || result.type

    const content = `# Descricao Gerada - Diffable

| Campo | Valor |
|-------|-------|
| Tipo | ${type} |
| Nivel | ${level} |
| Modelo | ${result.model_used} |
| Tokens | ${result.tokens_used} |
| Data | ${date} |

---

${result.description}
`

    const blob = new Blob([content], { type: 'text/markdown' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `diffable-${result.id.slice(0, 8)}.md`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)

    setExported(true)
    setTimeout(() => setExported(false), 2000)
  }

  return (
    <button
      onClick={handleExport}
      className="inline-flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm text-stone-500 transition-colors hover:bg-stone-100 hover:text-stone-700 dark:text-stone-400 dark:hover:bg-white/[0.06] dark:hover:text-stone-200"
      title={exported ? 'Exportado!' : 'Exportar Markdown'}
    >
      {exported ? (
        <>
          <Check size={14} className="text-emerald-500" />
          <span className="text-emerald-600 dark:text-emerald-400">Exportado!</span>
        </>
      ) : (
        <>
          <Download size={14} />
          <span>Exportar</span>
        </>
      )}
    </button>
  )
}
