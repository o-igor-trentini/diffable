import { Moon, Sun } from 'lucide-react'

interface LayoutProps {
  children: React.ReactNode
  dark: boolean
  onToggleDark: () => void
}

export function Layout({ children, dark, onToggleDark }: LayoutProps) {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors">
      <header className="bg-white shadow-sm dark:bg-gray-800">
        <div className="mx-auto flex max-w-4xl items-center justify-between px-4 py-4">
          <div>
            <h1 className="text-xl font-bold text-gray-900 dark:text-gray-100">Diffable</h1>
            <p className="text-sm text-gray-500 dark:text-gray-400">
              Gere descrições automáticas para seus commits e PRs
            </p>
          </div>
          <button
            onClick={onToggleDark}
            className="rounded-md p-2 text-gray-500 hover:bg-gray-100 transition-colors dark:text-gray-400 dark:hover:bg-gray-700"
            title={dark ? 'Modo claro' : 'Modo escuro'}
            aria-label={dark ? 'Ativar modo claro' : 'Ativar modo escuro'}
          >
            {dark ? <Sun size={20} /> : <Moon size={20} />}
          </button>
        </div>
      </header>

      <main className="mx-auto max-w-4xl px-4 py-6">
        {children}
      </main>
    </div>
  )
}
