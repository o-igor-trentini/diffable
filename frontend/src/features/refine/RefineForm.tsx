import { useState, type FormEvent } from 'react'
import { RefreshCw } from 'lucide-react'
import { Button } from '../shared/Button'
import { TextArea } from '../shared/TextArea'

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

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <TextArea
        id="refine-original"
        label="Descrição Original"
        rows={6}
        value={initialDescription ?? ''}
        disabled
      />

      <TextArea
        id="refine-instruction"
        label="Instrução de Refinamento"
        placeholder="ex: simplifique, mais técnico, mais resumido, foque nos impactos para o QA..."
        rows={3}
        value={instruction}
        onChange={(e) => setInstruction(e.target.value)}
        disabled={isPending}
      />

      <Button type="submit" disabled={!instruction.trim() || !initialDescription} loading={isPending}>
        <RefreshCw size={16} />
        Refinar
      </Button>
    </form>
  )
}
