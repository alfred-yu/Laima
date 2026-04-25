<template>
  <div class="dashboard">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">仪表盘</h1>
        <p class="page-subtitle">项目概览和最近活动</p>
      </div>
      <div class="header-right">
        <div class="date-range">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
            <line x1="16" y1="2" x2="16" y2="6" />
            <line x1="8" y1="2" x2="8" y2="6" />
            <line x1="3" y1="10" x2="21" y2="10" />
          </svg>
          <span>最近 30 天</span>
        </div>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon-wrapper repo-icon">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z" />
            <polyline points="13 2 13 9 20 9" />
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-value">
            <Skeleton type="text" width="60px" v-if="statsLoading" />
            <span v-else-if="statsError" class="error-text">--</span>
            <span v-else>{{ repoCount }}</span>
          </div>
          <div class="stat-label">总仓库</div>
          <div class="stat-change positive">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="23 6 13.5 15.5 8.5 10.5 1 18" />
              <polyline points="17 6 23 6 23 12" />
            </svg>
            +12%
          </div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon-wrapper pr-icon">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="18" cy="18" r="3" />
            <circle cx="6" cy="6" r="3" />
            <path d="M13 6h3a2 2 0 0 1 2 2v7" />
            <line x1="6" y1="9" x2="6" y2="21" />
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-value">
            <Skeleton type="text" width="60px" v-if="statsLoading" />
            <span v-else-if="statsError" class="error-text">--</span>
            <span v-else>{{ recentPRs.length || 24 }}</span>
          </div>
          <div class="stat-label">合并请求</div>
          <div class="stat-change positive">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="23 6 13.5 15.5 8.5 10.5 1 18" />
              <polyline points="17 6 23 6 23 12" />
            </svg>
            +8%
          </div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon-wrapper issue-icon">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10" />
            <line x1="12" y1="8" x2="12" y2="12" />
            <line x1="12" y1="16" x2="12.01" y2="16" />
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-value">
            <Skeleton type="text" width="60px" v-if="statsLoading" />
            <span v-else-if="statsError" class="error-text">--</span>
            <span v-else>{{ recentIssues.length || 56 }}</span>
          </div>
          <div class="stat-label">议题</div>
          <div class="stat-change negative">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="23 18 13.5 8.5 8.5 13.5 1 6" />
              <polyline points="17 18 23 18 23 12" />
            </svg>
            -3%
          </div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon-wrapper build-icon">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10" />
            <polyline points="12 6 12 12 16 14" />
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-value">
            <Skeleton type="text" width="60px" v-if="statsLoading" />
            <span v-else-if="statsError" class="error-text">--</span>
            <span v-else>128</span>
          </div>
          <div class="stat-label">构建数</div>
          <div class="stat-change positive">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="23 6 13.5 15.5 8.5 10.5 1 18" />
              <polyline points="17 6 23 6 23 12" />
            </svg>
            +15%
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content Grid -->
    <div class="content-grid">
      <!-- Recent Activity -->
      <div class="card activity-card">
        <div class="card-header">
          <h2 class="card-title">最近活动</h2>
          <button class="btn btn-link">查看全部</button>
        </div>
        <div class="card-body">
          <div v-if="activitiesLoading" class="loading-container">
            <Skeleton type="list" :count="5" />
          </div>
          <div v-else-if="activitiesError" class="error-container">
            <div class="error-icon">
              <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10" />
                <line x1="12" y1="8" x2="12" y2="12" />
                <line x1="12" y1="16" x2="12.01" y2="16" />
              </svg>
            </div>
            <p>{{ activitiesError }}</p>
          </div>
          <div v-else class="activity-list">
            <div 
              v-for="(activity, index) in (activities.length ? activities : mockActivities)" 
              :key="activity.id || index" 
              class="activity-item"
            >
              <div class="activity-avatar">
                {{ (activity.actor || 'U').charAt(0).toUpperCase() }}
              </div>
              <div class="activity-details">
                <div class="activity-title">
                  <span class="actor-name">{{ activity.actor || '用户' }}</span>
                  <span class="action-text">{{ activity.action || '更新了' }}</span>
                  <span class="target-link">{{ activity.target || '项目文件' }}</span>
                </div>
                <div class="activity-time">{{ activity.time || '刚刚' }}</div>
              </div>
              <div class="activity-badge">
                <span class="badge-type">{{ getActivityType(activity.action) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Column -->
      <div class="right-column">
        <!-- Recent PRs -->
        <div class="card">
          <div class="card-header">
            <h2 class="card-title">合并请求</h2>
            <button class="btn btn-link">全部</button>
          </div>
          <div class="card-body">
            <div v-if="prsLoading" class="loading-container">
              <Skeleton type="list" :count="3" />
            </div>
            <div v-else-if="prsError" class="empty-state">
              <p>{{ prsError }}</p>
            </div>
            <div v-else class="pr-list">
              <div 
                v-for="(pr, index) in (recentPRs.length ? recentPRs : mockPRs)" 
                :key="pr.id || index" 
                class="pr-item"
              >
                <div class="pr-icon status-open">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10" />
                  </svg>
                </div>
                <div class="pr-content">
                  <div class="pr-title">{{ pr.title || '新功能开发' }}</div>
                  <div class="pr-meta">
                    <span class="pr-repo">{{ pr.repo || 'main-project' }}</span>
                    <span class="pr-branch">#{{ pr.id || index + 1 }}</span>
                  </div>
                </div>
                <span class="status-badge status-open">{{ pr.status || 'OPEN' }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Recent Issues -->
        <div class="card">
          <div class="card-header">
            <h2 class="card-title">议题</h2>
            <button class="btn btn-link">全部</button>
          </div>
          <div class="card-body">
            <div v-if="issuesLoading" class="loading-container">
              <Skeleton type="list" :count="3" />
            </div>
            <div v-else-if="issuesError" class="empty-state">
              <p>{{ issuesError }}</p>
            </div>
            <div v-else class="issue-list">
              <div 
                v-for="(issue, index) in (recentIssues.length ? recentIssues : mockIssues)" 
                :key="issue.id || index" 
                class="issue-item"
              >
                <div class="issue-icon status-open">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10" />
                    <line x1="12" y1="8" x2="12" y2="12" />
                    <line x1="12" y1="16" x2="12.01" y2="16" />
                  </svg>
                </div>
                <div class="issue-content">
                  <div class="issue-title">{{ issue.title || 'Bug 修复' }}</div>
                  <div class="issue-meta">
                    <span class="issue-repo">{{ issue.repo || 'frontend-app' }}</span>
                    <span v-if="issue.labels" class="issue-labels">
                      <span 
                        v-for="(label, labelIndex) in issue.labels" 
                        :key="labelIndex" 
                        class="label-pill"
                      >
                        {{ label }}
                      </span>
                    </span>
                  </div>
                </div>
                <span class="status-badge status-open">{{ issue.status || 'OPEN' }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Skeleton } from '../components'
import { repoApi, prApi, issueApi } from '../services/api'

// Data
const activities = ref<any[]>([])
const recentPRs = ref<any[]>([])
const recentIssues = ref<any[]>([])
const repoCount = ref(0)
const activeRepoCount = ref(0)
const privateRepoCount = ref(0)

// Loading & Error States
const activitiesLoading = ref(true)
const prsLoading = ref(true)
const issuesLoading = ref(true)
const statsLoading = ref(true)
const activitiesError = ref('')
const prsError = ref('')
const issuesError = ref('')
const statsError = ref('')

// Mock Data for Display
const mockActivities = [
  { id: 1, actor: '张三', action: '提交代码到', target: 'repo-main', time: '2分钟前' },
  { id: 2, actor: '李四', action: '创建了', target: 'PR #45', time: '15分钟前' },
  { id: 3, actor: '王五', action: '评论了', target: 'Issue #12', time: '1小时前' },
  { id: 4, actor: '赵六', action: '合并了', target: 'PR #42', time: '2小时前' },
  { id: 5, actor: '钱七', action: '创建了', target: 'repo-new', time: '3小时前' },
]

const mockPRs = [
  { id: 1, title: 'Feat: 添加用户认证功能', repo: 'main-project', status: 'OPEN' },
  { id: 2, title: 'Fix: 修复登录页面样式问题', repo: 'frontend-app', status: 'OPEN' },
  { id: 3, title: 'Refactor: 优化组件性能', repo: 'ui-lib', status: 'MERGED' },
]

const mockIssues = [
  { id: 1, title: '登录超时问题', repo: 'main-project', labels: ['bug', 'high'], status: 'OPEN' },
  { id: 2, title: '添加暗黑模式支持', repo: 'frontend-app', labels: ['feature'], status: 'OPEN' },
  { id: 3, title: '文档错误', repo: 'docs', labels: ['doc'], status: 'CLOSED' },
]

// Helper Functions
const getActivityType = (action?: string) => {
  if (!action) return 'update'
  if (action.includes('提交')) return 'commit'
  if (action.includes('创建') || action.includes('新建')) return 'create'
  if (action.includes('评论')) return 'comment'
  if (action.includes('合并')) return 'merge'
  return 'update'
}

// API Calls
const loadActivities = async () => {
  try {
    activitiesLoading.value = true
    activitiesError.value = ''
    // TODO: Activity API
  } catch (err: any) {
    activitiesError.value = err.message || '加载活动失败'
  } finally {
    activitiesLoading.value = false
  }
}

const loadRecentPRs = async () => {
  try {
    prsLoading.value = true
    prsError.value = ''
    const response = await prApi.listPRs({ per_page: 5 })
    recentPRs.value = (response as any).items || []
  } catch (err: any) {
    prsError.value = err.message || '加载 PR 失败'
  } finally {
    prsLoading.value = false
  }
}

const loadRecentIssues = async () => {
  try {
    issuesLoading.value = true
    issuesError.value = ''
    const response = await issueApi.listIssues({ per_page: 5 })
    recentIssues.value = (response as any).items || []
  } catch (err: any) {
    issuesError.value = err.message || '加载 Issue 失败'
  } finally {
    issuesLoading.value = false
  }
}

const loadStats = async () => {
  try {
    statsLoading.value = true
    statsError.value = ''
    const response = await repoApi.listRepos({ per_page: 100 })
    const repos = (response as any).items || []
    repoCount.value = repos.length || 8
    activeRepoCount.value = repos.length > 0 ? Math.min(3, repos.length) : 3
    privateRepoCount.value = repos.filter((r: any) => r.visibility === 'private').length || 2
  } catch (err: any) {
    statsError.value = err.message || '加载统计失败'
  } finally {
    statsLoading.value = false
  }
}

onMounted(() => {
  loadActivities()
  loadRecentPRs()
  loadRecentIssues()
  loadStats()
})
</script>

<style scoped>
.dashboard {
  width: 100%;
}

/* Page Header */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
}

.header-left {
  flex: 1;
}

.page-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 4px 0;
  letter-spacing: -0.025em;
}

.page-subtitle {
  font-size: 0.9375rem;
  color: var(--text-secondary);
  margin: 0;
}

.date-range {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  color: var(--text-secondary);
  font-size: 0.875rem;
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  display: flex;
  gap: 16px;
  padding: 20px;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  transition: all 0.2s ease-out;
}

.stat-card:hover {
  border-color: var(--color-primary);
  box-shadow: var(--shadow-md);
}

.stat-icon-wrapper {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.repo-icon {
  background: rgba(79, 70, 229, 0.1);
  color: var(--color-primary);
}

.pr-icon {
  background: rgba(22, 163, 74, 0.1);
  color: var(--color-success);
}

.issue-icon {
  background: rgba(217, 119, 6, 0.1);
  color: var(--color-warning);
}

.build-icon {
  background: rgba(8, 145, 178, 0.1);
  color: var(--color-info);
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1.2;
}

.error-text {
  color: var(--text-muted);
}

.stat-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin-top: 4px;
}

.stat-change {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 0.8125rem;
  font-weight: 500;
  margin-top: 8px;
}

.stat-change.positive {
  color: var(--color-success);
}

.stat-change.negative {
  color: var(--color-danger);
}

/* Content Grid */
.content-grid {
  display: grid;
  grid-template-columns: 1fr 400px;
  gap: 24px;
}

.right-column {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* Card */
.card {
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  overflow: hidden;
}

.activity-card {
  min-height: 400px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color);
}

.card-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.card-body {
  padding: 0;
}

/* Buttons */
.btn {
  padding: 6px 12px;
  border: none;
  border-radius: 6px;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.btn-link {
  background: transparent;
  color: var(--color-primary);
  padding: 0;
  font-weight: 500;
}

.btn-link:hover {
  text-decoration: underline;
}

/* Loading & Error States */
.loading-container {
  padding: 20px;
}

.error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  color: var(--text-secondary);
}

.error-icon {
  color: var(--color-danger);
  margin-bottom: 12px;
  opacity: 0.5;
}

.empty-state {
  padding: 20px;
  color: var(--text-secondary);
}

/* Activity List */
.activity-list {
  display: flex;
  flex-direction: column;
}

.activity-item {
  display: flex;
  gap: 12px;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color-light);
  transition: background-color 0.15s ease-out;
}

.activity-item:last-child {
  border-bottom: none;
}

.activity-item:hover {
  background: var(--bg-secondary);
}

.activity-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--color-primary), #7c3aed);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 0.875rem;
  flex-shrink: 0;
}

.activity-details {
  flex: 1;
  min-width: 0;
}

.activity-title {
  font-size: 0.9375rem;
  color: var(--text-primary);
  line-height: 1.4;
}

.actor-name {
  font-weight: 600;
}

.action-text {
  color: var(--text-secondary);
}

.target-link {
  color: var(--color-primary);
  font-weight: 500;
}

.activity-time {
  font-size: 0.8125rem;
  color: var(--text-muted);
  margin-top: 4px;
}

.activity-badge {
  flex-shrink: 0;
}

.badge-type {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 500;
  background: var(--bg-tertiary);
  color: var(--text-secondary);
}

/* PR & Issue Lists */
.pr-list,
.issue-list {
  display: flex;
  flex-direction: column;
}

.pr-item,
.issue-item {
  display: flex;
  gap: 12px;
  padding: 14px 20px;
  align-items: center;
  border-bottom: 1px solid var(--border-color-light);
  transition: background-color 0.15s ease-out;
  cursor: pointer;
}

.pr-item:last-child,
.issue-item:last-child {
  border-bottom: none;
}

.pr-item:hover,
.issue-item:hover {
  background: var(--bg-secondary);
}

.pr-icon,
.issue-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pr-icon.status-open,
.issue-icon.status-open {
  color: var(--color-success);
}

.pr-content,
.issue-content {
  flex: 1;
  min-width: 0;
}

.pr-title,
.issue-title {
  font-size: 0.9375rem;
  color: var(--text-primary);
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.pr-meta,
.issue-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 4px;
}

.pr-repo,
.issue-repo {
  font-size: 0.8125rem;
  color: var(--text-muted);
}

.pr-branch {
  font-size: 0.8125rem;
  color: var(--text-muted);
}

.issue-labels {
  display: flex;
  gap: 6px;
}

.label-pill {
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 0.75rem;
  background: var(--bg-tertiary);
  color: var(--text-secondary);
}

.status-badge {
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 600;
  flex-shrink: 0;
}

.status-badge.status-open {
  background: rgba(22, 163, 74, 0.1);
  color: var(--color-success);
}

.status-badge.status-merged {
  background: rgba(79, 70, 229, 0.1);
  color: var(--color-primary);
}

.status-badge.status-closed {
  background: rgba(220, 38, 38, 0.1);
  color: var(--color-danger);
}

/* Responsive Design */
@media (max-width: 1024px) {
  .content-grid {
    grid-template-columns: 1fr;
  }
  
  .right-column {
    flex-direction: row;
    overflow-x: auto;
  }
  
  .right-column > .card {
    min-width: 320px;
  }
}

@media (max-width: 640px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .right-column {
    flex-direction: column;
  }
  
  .right-column > .card {
    min-width: 0;
  }
}
</style>
