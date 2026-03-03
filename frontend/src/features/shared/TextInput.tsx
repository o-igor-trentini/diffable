import type { InputHTMLAttributes } from 'react'

interface TextInputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string
  hint?: string
  error?: string
}

export function TextInput({ label, hint, error, className = '', id, ...props }: TextInputProps) {
  return (
    <div className="flex flex-col gap-1.5">
      {label && (
        <label htmlFor={id} className="text-sm font-medium text-stone-700 dark:text-stone-300">
          {label}
        </label>
      )}
      <input
        id={id}
        className={`w-full rounded-lg border bg-white px-3 py-2 text-sm transition-colors placeholder:text-stone-300 focus:border-violet-400 focus:outline-none focus:ring-2 focus:ring-violet-500/20 dark:bg-white/[0.03] dark:text-stone-100 dark:placeholder:text-stone-600 dark:focus:border-violet-500/50 dark:focus:ring-violet-500/10 ${
          error
            ? 'border-red-400 dark:border-red-500/50'
            : 'border-stone-200 dark:border-white/[0.08]'
        } ${className}`}
        {...props}
      />
      {error && <p className="text-xs text-red-500 dark:text-red-400">{error}</p>}
      {hint && !error && (
        <p className="text-xs text-stone-400 dark:text-stone-500">{hint}</p>
      )}
    </div>
  )
}
