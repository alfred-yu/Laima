import { useUserStore } from '../stores'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

interface ApiOptions {
  requireAuth?: boolean
}

class ApiClient {
  private baseURL: string
  private token: string | null = null

  constructor(baseURL: string) {
    this.baseURL = baseURL
    const userStore = useUserStore()
    this.token = userStore.token
  }

  private getHeaders(options?: ApiOptions): HeadersInit {
    const headers: HeadersInit = {
      'Content-Type': 'application/json'
    }

    if (options?.requireAuth) {
      const userStore = useUserStore()
      if (userStore.token) {
        headers['Authorization'] = `Bearer ${userStore.token}`
      }
    }

    return headers
  }

  private async handleResponse(response: Response) {
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}))
      throw new Error(errorData.error || `HTTP ${response.status}`)
    }

    return response.json()
  }

  async get<T>(endpoint: string, options?: ApiOptions): Promise<T> {
    const response = await fetch(`${this.baseURL}${endpoint}`, {
      method: 'GET',
      headers: this.getHeaders(options)
    })

    return this.handleResponse(response)
  }

  async post<T>(endpoint: string, data?: any, options?: ApiOptions): Promise<T> {
    const response = await fetch(`${this.baseURL}${endpoint}`, {
      method: 'POST',
      headers: this.getHeaders(options),
      body: data ? JSON.stringify(data) : undefined
    })

    return this.handleResponse(response)
  }

  async put<T>(endpoint: string, data?: any, options?: ApiOptions): Promise<T> {
    const response = await fetch(`${this.baseURL}${endpoint}`, {
      method: 'PUT',
      headers: this.getHeaders(options),
      body: data ? JSON.stringify(data) : undefined
    })

    return this.handleResponse(response)
  }

  async delete<T>(endpoint: string, options?: ApiOptions): Promise<T> {
    const response = await fetch(`${this.baseURL}${endpoint}`, {
      method: 'DELETE',
      headers: this.getHeaders(options)
    })

    return this.handleResponse(response)
  }
}

export const api = new ApiClient(API_BASE_URL)

// 用户认证相关 API
export const authApi = {
  login: async (username: string, password: string) => {
    return api.post('/api/auth/login', { username, password })
  },

  register: async (username: string, email: string, password: string) => {
    return api.post('/api/auth/register', { username, email, password })
  },

  getCurrentUser: async () => {
    return api.get('/api/users/me', { requireAuth: true })
  },

  updateCurrentUser: async (data: any) => {
    return api.put('/api/users/me', data, { requireAuth: true })
  }
}

// 仓库相关 API
export const repoApi = {
  listRepos: async (params?: { q?: string; visibility?: string; page?: number; per_page?: number }) => {
    const searchParams = new URLSearchParams()
    if (params?.q) searchParams.append('q', params.q)
    if (params?.visibility) searchParams.append('visibility', params.visibility)
    if (params?.page) searchParams.append('page', params.page.toString())
    if (params?.per_page) searchParams.append('per_page', params.per_page.toString())

    const queryString = searchParams.toString()
    return api.get(`/api/v1/repos${queryString ? `?${queryString}` : ''}`)
  },

  createRepo: async (data: any) => {
    return api.post('/api/v1/repos', data, { requireAuth: true })
  },

  getRepo: async (owner: string, repo: string) => {
    return api.get(`/api/v1/repos/${owner}/${repo}`)
  },

  updateRepo: async (owner: string, repo: string, data: any) => {
    return api.put(`/api/v1/repos/${owner}/${repo}`, data, { requireAuth: true })
  },

  deleteRepo: async (owner: string, repo: string) => {
    return api.delete(`/api/v1/repos/${owner}/${repo}`, { requireAuth: true })
  },

  listBranches: async (owner: string, repo: string) => {
    return api.get(`/api/v1/repos/${owner}/${repo}/branches`)
  },

  listTags: async (owner: string, repo: string) => {
    return api.get(`/api/v1/repos/${owner}/${repo}/tags`)
  }
}
