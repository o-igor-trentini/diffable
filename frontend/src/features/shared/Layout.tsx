import { Moon, Sun, Sparkles } from 'lucide-react'

interface LayoutProps {
  children: React.ReactNode
  dark: boolean
  onToggleDark: () => void
}

export function Layout({ children, dark, onToggleDark }: LayoutProps) {
  return (
    <div className="min-h-screen bg-stone-50 transition-colors dark:bg-[#08080e]">
      <header className="sticky top-0 z-30 border-b border-stone-200/80 bg-white/80 backdrop-blur-xl dark:border-white/[0.06] dark:bg-[#08080e]/80">
        <div className="mx-auto flex max-w-5xl items-center justify-between px-6 py-3.5">
          <div className="flex items-center gap-3">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-violet-600 to-cyan-500 shadow-lg shadow-violet-500/20">
              <Sparkles size={16} className="text-white" />
            </div>
            <div>
              <h1 className="text-gradient text-base font-bold tracking-tight">
                Diffable
              </h1>
              <p className="text-[11px] text-stone-400 dark:text-stone-500">
                Descricoes inteligentes para seus diffs
              </p>
            </div>
          </div>
          <button
            onClick={onToggleDark}
            className="rounded-lg p-2 text-stone-400 transition-colors hover:bg-stone-100 hover:text-stone-600 dark:text-stone-500 dark:hover:bg-white/[0.06] dark:hover:text-stone-300"
            title={dark ? 'Modo claro' : 'Modo escuro'}
            aria-label={dark ? 'Ativar modo claro' : 'Ativar modo escuro'}
          >
            {dark ? <Sun size={18} /> : <Moon size={18} />}
          </button>
        </div>
      </header>

      <main className="mx-auto max-w-5xl px-6 py-8">
        {children}
      </main>
    </div>
  )
}
