import { useState } from 'react'
import { Download, Check } from 'lucide-react'
import type { AnalysisResponse } from '@/lib/api/types'

const levelLabels: Record<string, string> = {
  technical: 'Técnico',
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

    const content = `# Descrição Gerada - Diffable

| Campo | Valor |
|-------|-------|
| Tipo | ${type} |
| Nível | ${level} |
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
      className="inline-flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm text-gray-600 hover:bg-gray-100 transition-colors dark:text-gray-300 dark:hover:bg-gray-600"
      title={exported ? 'Exportado!' : 'Exportar Markdown'}
    >
      {exported ? (
        <>
          <Check size={16} className="text-green-600 dark:text-green-400" />
          <span className="text-green-600 dark:text-green-400">Exportado!</span>
        </>
      ) : (
        <>
          <Download size={16} />
          <span>Exportar</span>
        </>
      )}
    </button>
  )
}
