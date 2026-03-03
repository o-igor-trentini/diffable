interface PlaceholderProps {
  title: string
  subtitle: string
}

export function Placeholder({ title, subtitle }: PlaceholderProps) {
  return (
    <div className="text-center text-gray-500 py-12">
      <p className="text-lg font-medium">{title}</p>
      <p className="text-sm mt-1">{subtitle}</p>
    </div>
  )
}
