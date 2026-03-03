import type { ButtonHTMLAttributes } from 'react'
import { LoadingSpinner } from './LoadingSpinner'

type Variant = 'primary' | 'secondary' | 'ghost'

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant
  loading?: boolean
}

const variantStyles: Record<Variant, string> = {
  primary:
    'bg-gradient-to-r from-violet-600 to-violet-500 text-white shadow-md shadow-violet-500/20 hover:from-violet-700 hover:to-violet-600 disabled:from-stone-300 disabled:to-stone-300 disabled:shadow-none dark:from-violet-600 dark:to-violet-500 dark:shadow-violet-500/10 dark:hover:from-violet-500 dark:hover:to-violet-400 dark:disabled:from-stone-700 dark:disabled:to-stone-700 dark:disabled:text-stone-500',
  secondary:
    'bg-stone-100 text-stone-700 hover:bg-stone-200 disabled:bg-stone-50 disabled:text-stone-300 dark:bg-white/[0.06] dark:text-stone-300 dark:hover:bg-white/[0.1] dark:disabled:bg-white/[0.02] dark:disabled:text-stone-600',
  ghost:
    'bg-transparent text-stone-500 hover:bg-stone-100 hover:text-stone-700 disabled:text-stone-300 dark:text-stone-400 dark:hover:bg-white/[0.06] dark:hover:text-stone-200 dark:disabled:text-stone-600',
}

export function Button({
  variant = 'primary',
  loading = false,
  children,
  disabled,
  className = '',
  ...props
}: ButtonProps) {
  return (
    <button
      disabled={disabled || loading}
      className={`inline-flex w-full items-center justify-center gap-2 rounded-xl px-5 py-2.5 text-sm font-semibold transition-all focus:outline-none focus-visible:ring-2 focus-visible:ring-violet-500/50 focus-visible:ring-offset-2 dark:focus-visible:ring-offset-[#08080e] sm:w-auto ${variantStyles[variant]} ${className}`}
      {...props}
    >
      {loading && <LoadingSpinner size="sm" />}
      {children}
    </button>
  )
}
