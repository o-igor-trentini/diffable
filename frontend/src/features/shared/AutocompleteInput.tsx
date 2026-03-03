import { useState, useRef, useEffect, type InputHTMLAttributes } from 'react'
import { Loader2, Search, AlertCircle } from 'lucide-react'

interface AutocompleteOption {
  value: string
  label: string
}

interface AutocompleteInputProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'onChange'> {
  label?: string
  hint?: string
  value: string
  onChange: (value: string) => void
  options: AutocompleteOption[]
  loading?: boolean
  onQueryChange: (query: string) => void
  dependencyMet?: boolean
  dependencyMessage?: string
}

export function AutocompleteInput({
  label,
  hint,
  value,
  onChange,
  options,
  loading = false,
  onQueryChange,
  dependencyMet = true,
  dependencyMessage,
  className = '',
  id,
  ...props
}: AutocompleteInputProps) {
  const [isOpen, setIsOpen] = useState(false)
  const [inputValue, setInputValue] = useState(value)
  const [hasSearched, setHasSearched] = useState(false)
  const debounceRef = useRef<ReturnType<typeof setTimeout>>(undefined)
  const containerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    setInputValue(value)
  }, [value])

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
        setIsOpen(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  function handleInputChange(newValue: string) {
    setInputValue(newValue)
    onChange(newValue)

    if (debounceRef.current) {
      clearTimeout(debounceRef.current)
    }

    debounceRef.current = setTimeout(() => {
      onQueryChange(newValue)
      if (newValue.length >= 2 && dependencyMet) {
        setIsOpen(true)
        setHasSearched(true)
      }
    }, 300)
  }

  function handleSelect(option: AutocompleteOption) {
    setInputValue(option.value)
    onChange(option.value)
    setIsOpen(false)
  }

  function handleKeyDown(e: React.KeyboardEvent) {
    if (e.key === 'Escape') {
      setIsOpen(false)
    }
  }

  function handleFocus() {
    if (options.length > 0 && inputValue.length >= 2) {
      setIsOpen(true)
    }
  }

  const showDependencyWarning = !dependencyMet && inputValue.length > 0
  const showNoResults = isOpen && hasSearched && !loading && options.length === 0 && inputValue.length >= 2 && dependencyMet
  const showOptions = isOpen && options.length > 0

  return (
    <div ref={containerRef} className="relative flex flex-col gap-1.5">
      {label && (
        <label htmlFor={id} className="text-sm font-medium text-stone-700 dark:text-stone-300">
          {label}
        </label>
      )}
      <div className="relative">
        <div className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-stone-300 dark:text-stone-600">
          <Search size={14} />
        </div>
        <input
          id={id}
          className={`w-full rounded-lg border bg-white py-2 pl-8 pr-8 text-sm transition-colors placeholder:text-stone-300 focus:border-violet-400 focus:outline-none focus:ring-2 focus:ring-violet-500/20 dark:bg-white/[0.03] dark:text-stone-100 dark:placeholder:text-stone-600 dark:focus:border-violet-500/50 dark:focus:ring-violet-500/10 ${
            showDependencyWarning
              ? 'border-amber-300 dark:border-amber-500/30'
              : 'border-stone-200 dark:border-white/[0.08]'
          } ${className}`}
          value={inputValue}
          onChange={(e) => handleInputChange(e.target.value)}
          onFocus={handleFocus}
          onKeyDown={handleKeyDown}
          autoComplete="off"
          {...props}
        />
        {loading && (
          <div className="absolute right-3 top-1/2 -translate-y-1/2">
            <Loader2 size={14} className="animate-spin text-violet-400" />
          </div>
        )}
      </div>

      {showDependencyWarning && dependencyMessage && (
        <p className="flex items-center gap-1 text-xs text-amber-600 dark:text-amber-400">
          <AlertCircle size={12} />
          {dependencyMessage}
        </p>
      )}

      {hint && !showDependencyWarning && (
        <p className="text-xs text-stone-400 dark:text-stone-500">{hint}</p>
      )}

      {(showOptions || showNoResults) && (
        <ul className="animate-slide-down absolute top-full left-0 right-0 z-20 mt-1 max-h-48 overflow-auto rounded-xl border border-stone-200 bg-white shadow-xl shadow-stone-900/5 dark:border-white/[0.08] dark:bg-[#14141e] dark:shadow-black/30">
          {showNoResults && (
            <li className="px-3 py-3 text-center text-xs text-stone-400 dark:text-stone-500">
              Nenhum repositorio encontrado
            </li>
          )}
          {showOptions &&
            options.map((option) => (
              <li
                key={option.value}
                onClick={() => handleSelect(option)}
                className="cursor-pointer px-3 py-2 text-sm transition-colors hover:bg-violet-50 dark:hover:bg-violet-500/[0.08]"
              >
                <span className="font-mono text-xs font-medium text-stone-700 dark:text-stone-200">
                  {option.value}
                </span>
                {option.label !== option.value && (
                  <span className="ml-2 text-xs text-stone-400 dark:text-stone-500">
                    {option.label}
                  </span>
                )}
              </li>
            ))}
        </ul>
      )}
    </div>
  )
}
