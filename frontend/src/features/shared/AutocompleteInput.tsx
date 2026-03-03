import { useState, useRef, useEffect, type InputHTMLAttributes } from 'react'
import { Loader2 } from 'lucide-react'

interface AutocompleteOption {
  value: string
  label: string
}

interface AutocompleteInputProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'onChange'> {
  label?: string
  value: string
  onChange: (value: string) => void
  options: AutocompleteOption[]
  loading?: boolean
  onQueryChange: (query: string) => void
}

export function AutocompleteInput({
  label,
  value,
  onChange,
  options,
  loading = false,
  onQueryChange,
  className = '',
  id,
  ...props
}: AutocompleteInputProps) {
  const [isOpen, setIsOpen] = useState(false)
  const [inputValue, setInputValue] = useState(value)
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
      if (newValue.length >= 2) {
        setIsOpen(true)
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

  return (
    <div ref={containerRef} className="relative flex flex-col gap-1">
      {label && (
        <label htmlFor={id} className="text-sm font-medium text-gray-700 dark:text-gray-300">
          {label}
        </label>
      )}
      <div className="relative">
        <input
          id={id}
          className={`w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100 dark:placeholder-gray-400 ${className}`}
          value={inputValue}
          onChange={(e) => handleInputChange(e.target.value)}
          onFocus={() => options.length > 0 && setIsOpen(true)}
          onKeyDown={handleKeyDown}
          autoComplete="off"
          {...props}
        />
        {loading && (
          <div className="absolute right-2 top-1/2 -translate-y-1/2">
            <Loader2 size={16} className="animate-spin text-gray-400" />
          </div>
        )}
      </div>
      {isOpen && options.length > 0 && (
        <ul className="absolute top-full left-0 right-0 z-10 mt-1 max-h-48 overflow-auto rounded-md border border-gray-200 bg-white shadow-lg dark:border-gray-600 dark:bg-gray-700">
          {options.map((option) => (
            <li
              key={option.value}
              onClick={() => handleSelect(option)}
              className="cursor-pointer px-3 py-2 text-sm text-gray-700 hover:bg-blue-50 dark:text-gray-200 dark:hover:bg-gray-600"
            >
              <span className="font-medium">{option.value}</span>
              {option.label !== option.value && (
                <span className="ml-2 text-gray-400 dark:text-gray-500">{option.label}</span>
              )}
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
