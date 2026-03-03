import { AlertCircle } from 'lucide-react'

interface ErrorAlertProps {
  message: string
}

export function ErrorAlert({ message }: ErrorAlertProps) {
  return (
    <div
      className="mt-4 flex items-start gap-3 rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm dark:border-red-500/20 dark:bg-red-500/[0.06]"
      role="alert"
    >
      <AlertCircle size={16} className="mt-0.5 shrink-0 text-red-500 dark:text-red-400" />
      <p className="text-red-700 dark:text-red-300">{message}</p>
    </div>
  )
}
