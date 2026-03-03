import { useState, type FormEvent } from 'react'
import { RefreshCw } from 'lucide-react'
import { Button } from '../shared/Button'
import { TextArea } from '../shared/TextArea'

const suggestions = [
  'Simplifique a linguagem',
  'Mais tecnico e detalhado',
  'Resuma em 2-3 frases',
  'Foque nos impactos para QA',
  'Foque no impacto de negocio',
]

interface RefineFormProps {
  initialDescription?: string
  onSubmit: (instruction: string) => void
  isPending: boolean
}

export function RefineForm({ initialDescription, onSubmit, isPending }: RefineFormProps) {
  const [instruction, setInstruction] = useState('')

  function handleSubmit(e: FormEvent) {
    e.preventDefault()
    if (instruction.trim()) {
      onSubmit(instruction.trim())
    }
  }

  function handleSuggestion(text: string) {
    setInstruction(text)
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-5">
      <div className="flex flex-col gap-1.5">
        <label className="text-sm font-medium text-stone-700 dark:text-stone-300">
          Descricao Original
        </label>
        <div className="max-h-48 overflow-y-auto rounded-lg border border-stone-200 bg-stone-50 p-4 font-mono text-xs leading-relaxed whitespace-pre-wrap text-stone-600 dark:border-white/[0.06] dark:bg-white/[0.02] dark:text-stone-400">
          {initialDescription || 'Nenhuma descricao carregada'}
        </div>
      </div>

      <div className="flex flex-col gap-2">
        <label className="text-sm font-medium text-stone-700 dark:text-stone-300">
          Sugestoes rapidas
        </label>
        <div className="flex flex-wrap gap-1.5">
          {suggestions.map((s) => (
            <button
              key={s}
              type="button"
              onClick={() => handleSuggestion(s)}
              disabled={isPending}
              className={`rounded-lg border px-3 py-1.5 text-xs font-medium transition-all ${
                instruction === s
                  ? 'border-violet-300 bg-violet-50 text-violet-700 dark:border-violet-500/30 dark:bg-violet-500/10 dark:text-violet-300'
                  : 'border-stone-200 text-stone-500 hover:border-stone-300 hover:text-stone-700 dark:border-white/[0.08] dark:text-stone-400 dark:hover:border-white/[0.15] dark:hover:text-stone-200'
              }`}
            >
              {s}
            </button>
          ))}
        </div>
      </div>

      <TextArea
        id="refine-instruction"
        label="Instrucao personalizada"
        placeholder="Descreva como voce quer ajustar a descricao..."
        hint="Escreva instrucoes livres ou selecione uma sugestao acima"
        rows={3}
        value={instruction}
        onChange={(e) => setInstruction(e.target.value)}
        disabled={isPending}
      />

      <Button type="submit" disabled={!instruction.trim() || !initialDescription} loading={isPending}>
        <RefreshCw size={16} />
        Refinar Descricao
      </Button>
    </form>
  )
}
