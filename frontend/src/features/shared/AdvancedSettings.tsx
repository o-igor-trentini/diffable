import { useState } from 'react'
import { Settings2, ChevronDown } from 'lucide-react'
import type { GenerationOverrides } from '../../lib/api/types'

const TOKEN_OPTIONS = [256, 512, 1024, 2048, 4096] as const

interface ModelOption {
  value: string
  label: string
  description: string
}

const modelOptions: ModelOption[] = [
  {
    value: 'auto',
    label: 'Auto',
    description: 'Seleção automática baseada no tamanho e tipo',
  },
  {
    value: 'gpt-4o-mini',
    label: 'GPT-4o Mini',
    description: 'Rápido e eficiente para diffs simples',
  },
  {
    value: 'gpt-4o',
    label: 'GPT-4o',
    description: 'Mais capaz, ideal para diffs complexos',
  },
]

interface AdvancedSettingsProps {
  value: GenerationOverrides
  onChange: (value: GenerationOverrides) => void
  disabled?: boolean
}

export function AdvancedSettings({ value, onChange, disabled }: AdvancedSettingsProps) {
  const [open, setOpen] = useState(false)

  const temperature = value.temperature ?? 0.3
  const maxTokens = value.max_tokens ?? 1024
  const model = value.model ?? 'auto'

  return (
    <div className="flex flex-col gap-2">
      <button
        type="button"
        onClick={() => setOpen((prev) => !prev)}
        disabled={disabled}
        className="flex items-center gap-2 text-xs font-medium uppercase tracking-wider text-stone-400 transition-colors hover:text-stone-600 dark:text-stone-500 dark:hover:text-stone-300"
      >
        <Settings2 size={14} />
        <span>Configurações avançadas</span>
        <ChevronDown
          size={14}
          className={`transition-transform ${open ? 'rotate-180' : ''}`}
        />
      </button>

      {open && (
        <div className="mt-1 flex flex-col gap-6 rounded-xl border border-stone-200 bg-white p-4 dark:border-white/[0.06] dark:bg-white/[0.02]">
          {/* Temperature */}
          <fieldset className="flex flex-col gap-2" disabled={disabled}>
            <legend className="text-xs font-medium uppercase tracking-wider text-stone-400 dark:text-stone-500">
              Temperature: {temperature.toFixed(1)}
            </legend>
            <p className="text-xs text-stone-400 dark:text-stone-500">
              Controla a criatividade. Valores baixos = mais preciso, altos = mais criativo.
            </p>
            <div className="flex flex-col gap-1">
              <input
                type="range"
                min={0}
                max={1}
                step={0.1}
                value={temperature}
                onChange={(e) =>
                  onChange({ ...value, temperature: parseFloat(e.target.value) })
                }
                disabled={disabled}
                className="h-2 w-full cursor-pointer appearance-none rounded-lg bg-stone-200 accent-violet-500 dark:bg-white/[0.06]"
              />
              <div className="flex justify-between text-[10px] text-stone-400 dark:text-stone-500">
                <span>Preciso (0.0)</span>
                <span>Criativo (1.0)</span>
              </div>
            </div>
          </fieldset>

          {/* Max Tokens */}
          <fieldset className="flex flex-col gap-2" disabled={disabled}>
            <legend className="text-xs font-medium uppercase tracking-wider text-stone-400 dark:text-stone-500">
              Max Tokens
            </legend>
            <p className="text-xs text-stone-400 dark:text-stone-500">
              Limite de tokens na resposta gerada. Mais tokens = descrição mais longa.
            </p>
            <div className="mt-1 flex flex-wrap gap-2">
              {TOKEN_OPTIONS.map((opt) => {
                const isSelected = maxTokens === opt
                return (
                  <button
                    key={opt}
                    type="button"
                    onClick={() => onChange({ ...value, max_tokens: opt })}
                    disabled={disabled}
                    className={`rounded-xl border px-3.5 py-2 text-sm font-medium transition-all ${
                      isSelected
                        ? 'border-violet-500/50 bg-violet-50 text-violet-700 dark:border-violet-400/30 dark:bg-violet-500/[0.07] dark:text-violet-300'
                        : 'border-stone-200 bg-white text-stone-600 hover:border-stone-300 dark:border-white/[0.06] dark:bg-white/[0.02] dark:text-stone-400 dark:hover:border-white/[0.12]'
                    }`}
                  >
                    {opt.toLocaleString()}
                  </button>
                )
              })}
            </div>
          </fieldset>

          {/* Modelo */}
          <fieldset className="flex flex-col gap-2" disabled={disabled}>
            <legend className="text-xs font-medium uppercase tracking-wider text-stone-400 dark:text-stone-500">
              Modelo
            </legend>
            <p className="text-xs text-stone-400 dark:text-stone-500">
              Escolha o modelo de IA. &apos;Auto&apos; seleciona automaticamente com base no tipo
              de análise.
            </p>
            <div className="mt-1 grid grid-cols-1 gap-2 sm:grid-cols-3">
              {modelOptions.map((opt) => {
                const isSelected = model === opt.value
                return (
                  <button
                    key={opt.value}
                    type="button"
                    onClick={() => onChange({ ...value, model: opt.value })}
                    disabled={disabled}
                    className={`group relative flex flex-col gap-1.5 rounded-xl border p-3.5 text-left transition-all ${
                      isSelected
                        ? 'border-violet-500/50 bg-violet-50 dark:border-violet-400/30 dark:bg-violet-500/[0.07]'
                        : 'border-stone-200 bg-white hover:border-stone-300 dark:border-white/[0.06] dark:bg-white/[0.02] dark:hover:border-white/[0.12]'
                    }`}
                  >
                    <span
                      className={`text-sm font-semibold ${
                        isSelected
                          ? 'text-violet-700 dark:text-violet-300'
                          : 'text-stone-700 dark:text-stone-300'
                      }`}
                    >
                      {opt.label}
                    </span>
                    <p className="text-xs leading-relaxed text-stone-500 dark:text-stone-400">
                      {opt.description}
                    </p>
                    {isSelected && (
                      <div className="absolute -top-px -right-px -left-px h-0.5 rounded-t-xl bg-gradient-to-r from-violet-600 to-cyan-500" />
                    )}
                  </button>
                )
              })}
            </div>
          </fieldset>
        </div>
      )}
    </div>
  )
}
