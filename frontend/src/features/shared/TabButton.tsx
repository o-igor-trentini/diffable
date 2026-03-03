interface TabButtonProps {
  active: boolean
  icon: React.ReactNode
  label: string
  onClick: () => void
}

export function TabButton({ active, icon, label, onClick }: TabButtonProps) {
  return (
    <button
      role="tab"
      aria-selected={active}
      onClick={onClick}
      className={`relative flex shrink-0 items-center gap-2 px-4 py-3 text-sm font-medium transition-colors ${
        active
          ? 'text-violet-700 dark:text-violet-300'
          : 'text-stone-400 hover:text-stone-600 dark:text-stone-500 dark:hover:text-stone-300'
      }`}
    >
      <span className={active ? 'text-violet-600 dark:text-violet-400' : ''}>{icon}</span>
      {label}
      {active && (
        <span className="absolute bottom-0 left-2 right-2 h-0.5 rounded-full bg-gradient-to-r from-violet-600 to-cyan-500" />
      )}
    </button>
  )
}
