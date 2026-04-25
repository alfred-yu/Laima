<template>
  <div class="activity-page">
    <div class="page-header">
      <h1 class="page-title">活动记录</h1>
      <p class="page-subtitle">查看您的最近操作和贡献</p>
    </div>

    <div class="activity-container">
      <!-- 筛选器 -->
      <div class="activity-filter">
        <div class="filter-tabs">
          <button 
            v-for="tab in filterTabs" 
            :key="tab.value"
            :class="['filter-tab', { active: activeFilter === tab.value }]"
            @click="activeFilter = tab.value"
          >
            <span class="tab-icon">{{ tab.icon }}</span>
            <span class="tab-text">{{ tab.label }}</span>
          </button>
        </div>
        <div class="date-range">
          <select v-model="dateRange" class="date-select">
            <option value="7">最近 7 天</option>
            <option value="30">最近 30 天</option>
            <option value="90">最近 90 天</option>
            <option value="all">全部</option>
          </select>
        </div>
      </div>

      <!-- 活动统计 -->
      <div class="activity-stats">
        <div class="stat-card">
          <div class="stat-icon contribute">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
              <polyline points="14 2 14 8 20 8" />
              <line x1="16" y1="13" x2="8" y2="13" />
              <line x1="16" y1="17" x2="8" y2="17" />
              <polyline points="10 9 9 9 8 9" />
            </svg>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.commits }}</div>
            <div class="stat-label">代码提交</div>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon pr">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="18" cy="18" r="3" />
              <circle cx="6" cy="6" r="3" />
              <path d="M13 6h3a2 2 0 0 1 2 2v7" />
              <line x1="6" y1="9" x2="6" y2="21" />
            </svg>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.prs }}</div>
            <div class="stat-label">Pull Request</div>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon issue">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10" />
              <line x1="12" y1="8" x2="12" y2="12" />
              <line x1="12" y1="16" x2="12.01" y2="16" />
            </svg>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.issues }}</div>
            <div class="stat-label">Issue</div>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon star">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" />
            </svg>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.stars }}</div>
            <div class="stat-label">星标</div>
          </div>
        </div>
      </div>

      <!-- 活动时间线 -->
      <div class="activity-timeline">
        <div v-if="loading" class="loading-state">
          <Skeleton type="list" :count="5" />
        </div>
        <div v-else-if="error" class="error-state">
          <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10" />
            <line x1="12" y1="8" x2="12" y2="12" />
            <line x1="12" y1="16" x2="12.01" y2="16" />
          </svg>
          <p>{{ error }}</p>
          <button class="btn btn-primary" @click="loadActivities">重试</button>
        </div>
        <div v-else-if="activities.length === 0" class="empty-state">
          <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <circle cx="12" cy="12" r="10" />
            <polyline points="12 6 12 12 16 14" />
          </svg>
          <h3>暂无活动记录</h3>
          <p>开始使用 Laima 进行代码开发，您的活动将显示在这里</p>
        </div>
        <div v-else class="timeline">
          <div v-for="(activity, index) in activities" :key="activity.id" class="timeline-item">
            <div class="timeline-marker" :class="activity.type"></div>
            <div class="timeline-content">
              <div class="activity-header">
                <div class="activity-type" :class="activity.type">
                  <span class="type-icon">{{ activity.icon }}</span>
                  <span class="type-text">{{ activity.action }}</span>
                </div>
                <div class="activity-time">{{ activity.time }}</div>
              </div>
              <div class="activity-details">
                <p class="activity-description">{{ activity.description }}</p>
                <div v-if="activity.repo" class="activity-repo">
                  <span class="repo-icon">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z" />
                      <polyline points="13 2 13 9 20 9" />
                    </svg>
                  </span>
                  <span class="repo-name">{{ activity.repo }}</span>
                </div>
                <div v-if="activity.commit" class="activity-commit">
                  <span class="commit-hash">{{ activity.commit.hash }}</span>
                  <span class="commit-message">{{ activity.commit.message }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div v-if="!loading && activities.length > 0" class="pagination">
        <button class="btn btn-ghost btn-sm" :disabled="page === 1" @click="page--; loadActivities()">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="15 18 9 12 15 6" />
          </svg>
          上一页
        </button>
        <span class="page-info">第 {{ page }} 页，共 {{ totalPages }} 页</span>
        <button class="btn btn-ghost btn-sm" :disabled="page === totalPages" @click="page++; loadActivities()">
          下一页
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="9 18 15 12 9 6" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Skeleton } from '../components'

const activeFilter = ref('all')
const dateRange = ref('30')
const activities = ref<any[]>([])
const loading = ref(true)
const error = ref('')
const page = ref(1)
const totalPages = ref(5)

// 活动统计
const stats = ref({
  commits: 42,
  prs: 18,
  issues: 12,
  stars: 25
})

// 筛选标签
const filterTabs = [
  { value: 'all', label: '所有活动', icon: '📋' },
  { value: 'commit', label: '代码提交', icon: '📝' },
  { value: 'pr', label: 'Pull Request', icon: '🔄' },
  { value: 'issue', label: 'Issue', icon: '🐛' },
  { value: 'star', label: '星标', icon: '⭐' }
]

// 模拟活动数据
const mockActivities = [
  {
    id: 1,
    type: 'commit',
    icon: '📝',
    action: '提交代码',
    description: '在 frontend-app 仓库中提交了新的功能代码',
    repo: 'frontend-app',
    commit: {
      hash: 'a1b2c3d',
      message: 'Add new dashboard component'
    },
    time: '2分钟前'
  },
  {
    id: 2,
    type: 'pr',
    icon: '🔄',
    action: '创建 PR',
    description: '在 backend-api 仓库中创建了新的 Pull Request',
    repo: 'backend-api',
    time: '15分钟前'
  },
  {
    id: 3,
    type: 'issue',
    icon: '🐛',
    action: '创建 Issue',
    description: '在 ui-components 仓库中创建了新的 Issue',
    repo: 'ui-components',
    time: '30分钟前'
  },
  {
    id: 4,
    type: 'star',
    icon: '⭐',
    action: '星标仓库',
    description: '星标了 awesome-vue 仓库',
    repo: 'awesome-vue',
    time: '1小时前'
  },
  {
    id: 5,
    type: 'commit',
    icon: '📝',
    action: '提交代码',
    description: '在 mobile-app 仓库中修复了 bug',
    repo: 'mobile-app',
    commit: {
      hash: 'e4f5g6h',
      message: 'Fix login bug'
    },
    time: '2小时前'
  },
  {
    id: 6,
    type: 'pr',
    icon: '🔄',
    action: '合并 PR',
    description: '合并了 frontend-app 仓库中的 PR #123',
    repo: 'frontend-app',
    time: '3小时前'
  },
  {
    id: 7,
    type: 'issue',
    icon: '🐛',
    action: '关闭 Issue',
    description: '关闭了 backend-api 仓库中的 Issue #456',
    repo: 'backend-api',
    time: '4小时前'
  },
  {
    id: 8,
    type: 'commit',
    icon: '📝',
    action: '提交代码',
    description: '在 api-gateway 仓库中更新了依赖',
    repo: 'api-gateway',
    commit: {
      hash: 'i7j8k9l',
      message: 'Update dependencies'
    },
    time: '5小时前'
  }
]

// 加载活动数据
const loadActivities = async () => {
  try {
    loading.value = true
    error.value = ''
    
    // 模拟 API 调用
    await new Promise(resolve => setTimeout(resolve, 500))
    
    // 模拟筛选
    let filteredActivities = [...mockActivities]
    
    if (activeFilter.value !== 'all') {
      filteredActivities = filteredActivities.filter(activity => activity.type === activeFilter.value)
    }
    
    activities.value = filteredActivities
  } catch (err: any) {
    error.value = err.message || '加载失败'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadActivities()
})
</script>

<style scoped>
.activity-page {
  padding: 24px;
  min-height: 100vh;
}

.page-header {
  margin-bottom: 32px;
}

.page-title {
  font-size: 2rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 8px 0;
  letter-spacing: -0.025em;
}

.page-subtitle {
  font-size: 1rem;
  color: var(--text-secondary);
  margin: 0;
}

.activity-container {
  max-width: 1200px;
  margin: 0 auto;
}

/* 筛选器 */
.activity-filter {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.filter-tabs {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.filter-tab {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border: 1px solid var(--border-color);
  border-radius: 20px;
  background: var(--bg-primary);
  color: var(--text-secondary);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.filter-tab:hover {
  border-color: var(--color-primary);
  color: var(--color-primary);
}

.filter-tab.active {
  background: var(--color-primary);
  border-color: var(--color-primary);
  color: white;
}

.tab-icon {
  font-size: 1rem;
}

.date-range {
  display: flex;
  align-items: center;
}

.date-select {
  padding: 8px 12px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.date-select:hover {
  border-color: var(--color-primary);
}

.date-select:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px var(--color-primary-bg);
}

/* 活动统计 */
.activity-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 32px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  box-shadow: var(--shadow-sm);
  transition: all 0.15s ease-out;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
  border-color: var(--color-primary);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.stat-icon.contribute {
  background: linear-gradient(135deg, #4f46e5, #8b5cf6);
}

.stat-icon.pr {
  background: linear-gradient(135deg, #10b981, #34d399);
}

.stat-icon.issue {
  background: linear-gradient(135deg, #f59e0b, #fbbf24);
}

.stat-icon.star {
  background: linear-gradient(135deg, #ef4444, #f87171);
}

.stat-content {
  flex: 1;
  min-width: 0;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 4px 0;
}

.stat-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin: 0;
}

/* 活动时间线 */
.activity-timeline {
  margin-bottom: 32px;
}

.loading-state {
  padding: 24px;
}

.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 64px 24px;
  background: var(--bg-secondary);
  border-radius: 12px;
  text-align: center;
  color: var(--text-secondary);
  gap: 16px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 24px;
  background: var(--bg-secondary);
  border-radius: 12px;
  text-align: center;
  color: var(--text-secondary);
  gap: 16px;
}

.empty-state h3 {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.empty-state p {
  max-width: 400px;
  margin: 0;
  line-height: 1.5;
}

.timeline {
  position: relative;
  padding-left: 32px;
}

.timeline::before {
  content: '';
  position: absolute;
  left: 11px;
  top: 0;
  bottom: 0;
  width: 2px;
  background: var(--border-color);
}

.timeline-item {
  position: relative;
  margin-bottom: 24px;
  animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.timeline-marker {
  position: absolute;
  left: -32px;
  top: 4px;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 12px;
  z-index: 1;
  border: 3px solid var(--bg-primary);
}

.timeline-marker.commit {
  background: #4f46e5;
}

.timeline-marker.pr {
  background: #10b981;
}

.timeline-marker.issue {
  background: #f59e0b;
}

.timeline-marker.star {
  background: #ef4444;
}

.timeline-content {
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 20px;
  box-shadow: var(--shadow-sm);
  transition: all 0.15s ease-out;
}

.timeline-content:hover {
  transform: translateX(4px);
  box-shadow: var(--shadow-md);
  border-color: var(--color-primary);
}

.activity-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.activity-type {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  font-size: 0.9375rem;
}

.activity-type.commit {
  color: #4f46e5;
}

.activity-type.pr {
  color: #10b981;
}

.activity-type.issue {
  color: #f59e0b;
}

.activity-type.star {
  color: #ef4444;
}

.type-icon {
  font-size: 1rem;
}

.activity-time {
  font-size: 0.8125rem;
  color: var(--text-muted);
  white-space: nowrap;
}

.activity-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.activity-description {
  margin: 0;
  font-size: 0.9375rem;
  color: var(--text-primary);
  line-height: 1.5;
}

.activity-repo {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.875rem;
  color: var(--text-secondary);
  padding: 8px 12px;
  background: var(--bg-secondary);
  border-radius: 8px;
}

.repo-icon {
  color: var(--color-primary);
}

.activity-commit {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 8px 12px;
  background: var(--bg-secondary);
  border-radius: 8px;
  font-size: 0.875rem;
}

.commit-hash {
  color: var(--color-primary);
  font-family: monospace;
  font-size: 0.8125rem;
}

.commit-message {
  color: var(--text-secondary);
}

/* 分页 */
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 24px 0;
}

.page-info {
  font-size: 0.875rem;
  color: var(--text-secondary);
  min-width: 120px;
  text-align: center;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .activity-filter {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .filter-tabs {
    width: 100%;
  }
  
  .filter-tab {
    flex: 1;
    justify-content: center;
  }
  
  .activity-stats {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .timeline {
    padding-left: 24px;
  }
  
  .timeline-marker {
    left: -24px;
    width: 20px;
    height: 20px;
    font-size: 10px;
  }
  
  .timeline::before {
    left: 9px;
  }
}

@media (max-width: 480px) {
  .activity-page {
    padding: 16px;
  }
  
  .activity-stats {
    grid-template-columns: 1fr;
  }
  
  .filter-tab {
    padding: 6px 12px;
    font-size: 0.8125rem;
  }
  
  .timeline-content {
    padding: 16px;
  }
  
  .activity-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
}
</style>
