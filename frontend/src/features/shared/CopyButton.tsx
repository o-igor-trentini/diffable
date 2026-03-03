import { Copy, Check } from 'lucide-react'
import { useClipboard } from '@/lib/hooks/useClipboard'

interface CopyButtonProps {
  text: string
}

export function CopyButton({ text }: CopyButtonProps) {
  const { copy, copied } = useClipboard()

  return (
    <button
      data-copy-button
      onClick={() => copy(text)}
      className="inline-flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm text-stone-500 transition-colors hover:bg-stone-100 hover:text-stone-700 dark:text-stone-400 dark:hover:bg-white/[0.06] dark:hover:text-stone-200"
      title={copied ? 'Copiado!' : 'Copiar (Ctrl+Shift+C)'}
    >
      {copied ? (
        <>
          <Check size={14} className="text-emerald-500" />
          <span className="text-emerald-600 dark:text-emerald-400">Copiado!</span>
        </>
      ) : (
        <>
          <Copy size={14} />
          <span>Copiar</span>
        </>
      )}
    </button>
  )
}
