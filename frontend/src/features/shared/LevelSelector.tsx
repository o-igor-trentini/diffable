import { ClipboardList, Code2, BarChart3 } from 'lucide-react'

interface LevelOption {
  value: string
  label: string
  description: string
  icon: React.ReactNode
}

const levels: LevelOption[] = [
  {
    value: 'functional',
    label: 'Funcional',
    description: 'Para QA e Product Owners. Foco em mudanças funcionais, sem jargão técnico.',
    icon: <ClipboardList size={18} />,
  },
  {
    value: 'technical',
    label: 'Técnico',
    description: 'Para engenheiros. Linguagem técnica com referências a código e decisões.',
    icon: <Code2 size={18} />,
  },
  {
    value: 'executive',
    label: 'Executivo',
    description: 'Para gestão e stakeholders. Resumo direto em 2-3 frases.',
    icon: <BarChart3 size={18} />,
  },
]

interface LevelSelectorProps {
  value: string
  onChange: (value: string) => void
  disabled?: boolean
}

export function LevelSelector({ value, onChange, disabled }: LevelSelectorProps) {
  return (
    <fieldset className="flex flex-col gap-2" disabled={disabled}>
      <legend className="text-xs font-medium uppercase tracking-wider text-stone-400 dark:text-stone-500">
        Nivel da descricao
      </legend>
      <p className="text-xs text-stone-400 dark:text-stone-500">
        Escolha o publico-alvo da descricao gerada
      </p>
      <div className="mt-1 grid grid-cols-1 gap-2 sm:grid-cols-3">
        {levels.map((level) => {
          const isSelected = value === level.value
          return (
            <button
              key={level.value}
              type="button"
              onClick={() => onChange(level.value)}
              disabled={disabled}
              className={`group relative flex flex-col gap-1.5 rounded-xl border p-3.5 text-left transition-all ${
                isSelected
                  ? 'border-violet-500/50 bg-violet-50 dark:border-violet-400/30 dark:bg-violet-500/[0.07]'
                  : 'border-stone-200 bg-white hover:border-stone-300 dark:border-white/[0.06] dark:bg-white/[0.02] dark:hover:border-white/[0.12]'
              }`}
            >
              <div className="flex items-center gap-2">
                <span
                  className={`${
                    isSelected
                      ? 'text-violet-600 dark:text-violet-400'
                      : 'text-stone-400 dark:text-stone-500'
                  }`}
                >
                  {level.icon}
                </span>
                <span
                  className={`text-sm font-semibold ${
                    isSelected
                      ? 'text-violet-700 dark:text-violet-300'
                      : 'text-stone-700 dark:text-stone-300'
                  }`}
                >
                  {level.label}
                </span>
              </div>
              <p className="text-xs leading-relaxed text-stone-500 dark:text-stone-400">
                {level.description}
              </p>
              {isSelected && (
                <div className="absolute -top-px -right-px -left-px h-0.5 rounded-t-xl bg-gradient-to-r from-violet-600 to-cyan-500" />
              )}
            </button>
          )
        })}
      </div>
    </fieldset>
  )
}
