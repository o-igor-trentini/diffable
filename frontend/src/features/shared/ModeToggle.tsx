import { Cloud, FileCode2 } from 'lucide-react'

export type SourceMode = 'bitbucket' | 'manual'

interface ModeToggleProps {
  mode: SourceMode
  onModeChange: (mode: SourceMode) => void
  disabled?: boolean
}

export function ModeToggle({ mode, onModeChange, disabled }: ModeToggleProps) {
  return (
    <div className="flex flex-col gap-2">
      <span className="text-xs font-medium uppercase tracking-wider text-stone-400 dark:text-stone-500">
        Fonte do diff
      </span>
      <div className="inline-flex rounded-xl bg-stone-100 p-1 dark:bg-white/[0.04]">
        <button
          type="button"
          disabled={disabled}
          onClick={() => onModeChange('bitbucket')}
          className={`flex items-center gap-2 rounded-lg px-4 py-2 text-sm font-medium transition-all ${
            mode === 'bitbucket'
              ? 'bg-white text-violet-700 shadow-sm dark:bg-white/[0.1] dark:text-violet-300'
              : 'text-stone-500 hover:text-stone-700 dark:text-stone-400 dark:hover:text-stone-200'
          }`}
        >
          <Cloud size={15} />
          Buscar do Bitbucket
        </button>
        <button
          type="button"
          disabled={disabled}
          onClick={() => onModeChange('manual')}
          className={`flex items-center gap-2 rounded-lg px-4 py-2 text-sm font-medium transition-all ${
            mode === 'manual'
              ? 'bg-white text-violet-700 shadow-sm dark:bg-white/[0.1] dark:text-violet-300'
              : 'text-stone-500 hover:text-stone-700 dark:text-stone-400 dark:hover:text-stone-200'
          }`}
        >
          <FileCode2 size={15} />
          Colar Diff
        </button>
      </div>
    </div>
  )
}
