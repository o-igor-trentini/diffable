import { Copy, Check } from 'lucide-react'
import { useClipboard } from '@/lib/hooks/useClipboard'

interface CopyButtonProps {
  text: string
}

export function CopyButton({ text }: CopyButtonProps) {
  const { copy, copied } = useClipboard()

  return (
    <button
      onClick={() => copy(text)}
      className="inline-flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm text-gray-600 hover:bg-gray-100 transition-colors"
      title={copied ? 'Copiado!' : 'Copiar'}
    >
      {copied ? (
        <>
          <Check size={16} className="text-green-600" />
          <span className="text-green-600">Copiado!</span>
        </>
      ) : (
        <>
          <Copy size={16} />
          <span>Copiar</span>
        </>
      )}
    </button>
  )
}
