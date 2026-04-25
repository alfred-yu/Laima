import { defineStore } from 'pinia'

// 用户状态管理
export const useUserStore = defineStore('user', {
  state: () => ({
    user: null as any,
    token: localStorage.getItem('token') || '',
    role: localStorage.getItem('role') || 'user', // 默认普通用户
    isLoading: false,
    error: null as string | null
  }),
  getters: {
    isLoggedIn: (state) => !!state.token,
    currentUser: (state) => state.user,
    isAdmin: (state) => state.role === 'admin',
    isUser: (state) => state.role === 'user'
  },
  actions: {
    setToken(token: string) {
      this.token = token
      localStorage.setItem('token', token)
    },
    setUser(user: any) {
      this.user = user
    },
    setRole(role: string) {
      this.role = role
      localStorage.setItem('role', role)
    },
    logout() {
      this.token = ''
      this.user = null
      this.role = 'user'
      localStorage.removeItem('token')
      localStorage.removeItem('role')
    },
    setLoading(loading: boolean) {
      this.isLoading = loading
    },
    setError(error: string | null) {
      this.error = error
    }
  }
})

// 仓库状态管理
export const useRepoStore = defineStore('repo', {
  state: () => ({
    repos: [] as any[],
    currentRepo: null as any,
    isLoading: false,
    error: null as string | null
  }),
  getters: {
    repoList: (state) => state.repos,
    currentRepository: (state) => state.currentRepo
  },
  actions: {
    setRepos(repos: any[]) {
      this.repos = repos
    },
    setCurrentRepo(repo: any) {
      this.currentRepo = repo
    },
    setLoading(loading: boolean) {
      this.isLoading = loading
    },
    setError(error: string | null) {
      this.error = error
    }
  }
})

// PR 状态管理
export const usePRStore = defineStore('pr', {
  state: () => ({
    prs: [] as any[],
    currentPR: null as any,
    isLoading: false,
    error: null as string | null
  }),
  getters: {
    prList: (state) => state.prs,
    currentPullRequest: (state) => state.currentPR
  },
  actions: {
    setPRs(prs: any[]) {
      this.prs = prs
    },
    setCurrentPR(pr: any) {
      this.currentPR = pr
    },
    setLoading(loading: boolean) {
      this.isLoading = loading
    },
    setError(error: string | null) {
      this.error = error
    }
  }
})

// Issue 状态管理
export const useIssueStore = defineStore('issue', {
  state: () => ({
    issues: [] as any[],
    currentIssue: null as any,
    isLoading: false,
    error: null as string | null
  }),
  getters: {
    issueList: (state) => state.issues,
    currentIssueDetail: (state) => state.currentIssue
  },
  actions: {
    setIssues(issues: any[]) {
      this.issues = issues
    },
    setCurrentIssue(issue: any) {
      this.currentIssue = issue
    },
    setLoading(loading: boolean) {
      this.isLoading = loading
    },
    setError(error: string | null) {
      this.error = error
    }
  }
})

// CI/CD 状态管理
export const useCICDStore = defineStore('cicd', {
  state: () => ({
    pipelines: [] as any[],
    currentPipeline: null as any,
    isLoading: false,
    error: null as string | null
  }),
  getters: {
    pipelineList: (state) => state.pipelines,
    currentPipelineDetail: (state) => state.currentPipeline
  },
  actions: {
    setPipelines(pipelines: any[]) {
      this.pipelines = pipelines
    },
    setCurrentPipeline(pipeline: any) {
      this.currentPipeline = pipeline
    },
    setLoading(loading: boolean) {
      this.isLoading = loading
    },
    setError(error: string | null) {
      this.error = error
    }
  }
})

// 全局状态管理
export const useGlobalStore = defineStore('global', {
  state: () => ({
    sidebar: {
      opened: true
    },
    theme: localStorage.getItem('theme') || 'light',
    notifications: [] as any[]
  }),
  getters: {
    isSidebarOpened: (state) => state.sidebar.opened,
    currentTheme: (state) => state.theme
  },
  actions: {
    toggleSidebar() {
      this.sidebar.opened = !this.sidebar.opened
    },
    setTheme(theme: string) {
      this.theme = theme
      localStorage.setItem('theme', theme)
    },
    addNotification(notification: any) {
      this.notifications.push(notification)
      // 5秒后自动移除通知
      setTimeout(() => {
        this.removeNotification(notification.id)
      }, 5000)
    },
    removeNotification(id: string) {
      this.notifications = this.notifications.filter(n => n.id !== id)
    }
  }
})