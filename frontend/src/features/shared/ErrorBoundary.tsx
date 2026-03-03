import { Component, type ErrorInfo, type ReactNode } from 'react'
import { AlertCircle, RefreshCw } from 'lucide-react'

interface ErrorBoundaryProps {
  children: ReactNode
}

interface ErrorBoundaryState {
  hasError: boolean
  error: Error | null
}

export class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: ErrorBoundaryProps) {
    super(props)
    this.state = { hasError: false, error: null }
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('ErrorBoundary caught an error:', error, errorInfo)
  }

  handleRetry = () => {
    this.setState({ hasError: false, error: null })
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4 dark:bg-gray-900">
          <div className="w-full max-w-md rounded-lg bg-white p-8 text-center shadow-lg dark:bg-gray-800">
            <AlertCircle size={48} className="mx-auto text-red-500" />
            <h2 className="mt-4 text-lg font-semibold text-gray-900 dark:text-gray-100">
              Algo deu errado
            </h2>
            <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">
              Ocorreu um erro inesperado. Tente novamente.
            </p>

            {import.meta.env.DEV && this.state.error && (
              <pre className="mt-4 max-h-32 overflow-auto rounded bg-gray-100 p-3 text-left text-xs text-red-700 dark:bg-gray-700 dark:text-red-400">
                {this.state.error.message}
              </pre>
            )}

            <button
              onClick={this.handleRetry}
              className="mt-6 inline-flex items-center gap-2 rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 transition-colors"
            >
              <RefreshCw size={16} />
              Tentar novamente
            </button>
          </div>
        </div>
      )
    }

    return this.props.children
  }
}
