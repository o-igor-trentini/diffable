import axios from 'axios'

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor: add X-Request-ID
apiClient.interceptors.request.use((config) => {
  config.headers['X-Request-ID'] = crypto.randomUUID()
  return config
})

// Response interceptor: retry once on 5xx, extract error messages
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const config = error.config
    if (
      error.response &&
      error.response.status >= 500 &&
      !config._retried
    ) {
      config._retried = true
      return apiClient(config)
    }

    const message =
      error.response?.data?.message ||
      error.response?.data?.error ||
      error.message ||
      'An unexpected error occurred'
    return Promise.reject(new Error(message))
  },
)

export default apiClient
