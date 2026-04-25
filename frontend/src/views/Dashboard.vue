<template>
  <div class="dashboard-page">
    <div class="page-header">
      <h1 class="page-title">
        {{ isAdmin ? '管理仪表盘' : '仪表盘' }}
      </h1>
      <p class="page-subtitle">
        {{ isAdmin ? '查看系统概览和管理功能' : '查看您的项目概览和活动统计' }}
      </p>
    </div>

    <div class="dashboard-container">
      <!-- 管理员专属内容 -->
      <div v-if="isAdmin" class="admin-section">
        <!-- 系统概览 -->
        <div class="admin-overview">
          <div class="section-header">
            <h3>系统概览</h3>
          </div>
          <div class="admin-stats">
            <div class="admin-stat-card">
              <div class="stat-icon users">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
                  <circle cx="9" cy="7" r="4" />
                  <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
                  <path d="M16 3.13a4 4 0 0 1 0 7.75" />
                </svg>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ adminStats.totalUsers }}</div>
                <div class="stat-label">总用户数</div>
              </div>
            </div>
            <div class="admin-stat-card">
              <div class="stat-icon repos">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z" />
                  <polyline points="13 2 13 9 20 9" />
                </svg>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ adminStats.totalRepos }}</div>
                <div class="stat-label">总仓库数</div>
              </div>
            </div>
            <div class="admin-stat-card">
              <div class="stat-icon commits">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
                  <polyline points="14 2 14 8 20 8" />
                  <line x1="16" y1="13" x2="8" y2="13" />
                  <line x1="16" y1="17" x2="8" y2="17" />
                  <polyline points="10 9 9 9 8 9" />
                </svg>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ adminStats.totalCommits }}</div>
                <div class="stat-label">总提交数</div>
              </div>
            </div>
            <div class="admin-stat-card">
              <div class="stat-icon health">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
                </svg>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ adminStats.systemHealth }}%</div>
                <div class="stat-label">系统健康</div>
              </div>
            </div>
          </div>
        </div>

        <!-- 管理操作 -->
        <div class="admin-actions">
          <div class="section-header">
            <h3>管理操作</h3>
          </div>
          <div class="action-grid">
            <button class="admin-action-card">
              <div class="action-icon">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
                  <circle cx="9" cy="7" r="4" />
                  <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
                  <path d="M16 3.13a4 4 0 0 1 0 7.75" />
                </svg>
              </div>
              <div class="action-content">
                <h3>用户管理</h3>
                <p>管理系统用户和权限</p>
              </div>
            </button>
            <button class="admin-action-card">
              <div class="action-icon">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z" />
                  <polyline points="13 2 13 9 20 9" />
                </svg>
              </div>
              <div class="action-content">
                <h3>仓库管理</h3>
                <p>管理系统仓库和项目</p>
              </div>
            </button>
            <button class="admin-action-card">
              <div class="action-icon">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="3" />
                  <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z" />
                </svg>
              </div>
              <div class="action-content">
                <h3>系统设置</h3>
                <p>配置系统参数和功能</p>
              </div>
            </button>
            <button class="admin-action-card">
              <div class="action-icon">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M14 9a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h6a2 2 0 0 1 2 2v3z" />
                  <path d="M20 17a2 2 0 0 1-2 2h-6a2 2 0 0 1-2-2v-3a2 2 0 0 1 2-2h6a2 2 0 0 1 2 2v3z" />
                </svg>
              </div>
              <div class="action-content">
                <h3>日志管理</h3>
                <p>查看系统日志和审计</p>
              </div>
            </button>
          </div>
        </div>

        <!-- 系统健康 -->
        <div class="system-health">
          <div class="section-header">
            <h3>系统健康状态</h3>
          </div>
          <div class="health-grid">
            <div v-for="(service, index) in systemHealthData" :key="index" class="health-card">
              <div class="health-header">
                <span class="service-name">{{ service.name }}</span>
                <span class="service-status" :class="service.status">
                  {{ service.status === 'online' ? '在线' : '离线' }}
                </span>
              </div>
              <div class="health-uptime">
                <span class="uptime-label">运行时间:</span>
                <span class="uptime-value">{{ service.uptime }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 通用内容（所有用户可见） -->
      <div class="common-section">
        <!-- 快速操作 -->
        <div class="quick-actions">
          <button class="action-card create-repo">
            <div class="action-icon">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 2L4 6V12C4 15.31 6.84 18.17 10.5 18.92V22H13.5V18.92C17.16 18.17 20 15.31 20 12V6L12 2Z" />
                <path d="M12 4.5L18 7.5V12C18 14.8 15.31 17.08 12.5 17.67V19.5H11.5V17.67C8.69 17.08 6 14.8 6 12V7.5L12 4.5Z" />
                <circle cx="12" cy="11" r="2.5" fill="currentColor" />
              </svg>
            </div>
            <div class="action-content">
              <h3>创建仓库</h3>
              <p>开始一个新的代码项目</p>
            </div>
          </button>
          <button class="action-card import-repo">
            <div class="action-icon">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                <polyline points="7 10 12 15 17 10" />
                <line x1="12" y1="15" x2="12" y2="3" />
              </svg>
            </div>
            <div class="action-content">
              <h3>导入仓库</h3>
              <p>从其他平台导入代码</p>
            </div>
          </button>
          <button class="action-card new-pr">
            <div class="action-icon">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="18" cy="18" r="3" />
                <circle cx="6" cy="6" r="3" />
                <path d="M13 6h3a2 2 0 0 1 2 2v7" />
                <line x1="6" y1="9" x2="6" y2="21" />
              </svg>
            </div>
            <div class="action-content">
              <h3>新建 PR</h3>
              <p>提交代码合并请求</p>
            </div>
          </button>
          <button class="action-card new-issue">
            <div class="action-icon">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10" />
                <line x1="12" y1="8" x2="12" y2="12" />
                <line x1="12" y1="16" x2="12.01" y2="16" />
              </svg>
            </div>
            <div class="action-content">
              <h3>新建 Issue</h3>
              <p>报告问题或提出建议</p>
            </div>
          </button>
        </div>

        <!-- 数据概览 -->
        <div class="stats-overview">
          <div class="stat-card">
            <div class="stat-header">
              <h3>总仓库数</h3>
              <span class="stat-icon repo">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z" />
                  <polyline points="13 2 13 9 20 9" />
                </svg>
              </span>
            </div>
            <div class="stat-value">{{ stats.repos }}</div>
            <div class="stat-change positive">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="22 12 18 12 15 21 9 3 6 12 2 12" />
              </svg>
              +{{ stats.reposChange }}% 较上月
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-header">
              <h3>代码提交</h3>
              <span class="stat-icon commit">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
                  <polyline points="14 2 14 8 20 8" />
                  <line x1="16" y1="13" x2="8" y2="13" />
                  <line x1="16" y1="17" x2="8" y2="17" />
                  <polyline points="10 9 9 9 8 9" />
                </svg>
              </span>
            </div>
            <div class="stat-value">{{ stats.commits }}</div>
            <div class="stat-change positive">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="22 12 18 12 15 21 9 3 6 12 2 12" />
              </svg>
              +{{ stats.commitsChange }}% 较上月
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-header">
              <h3>合并请求</h3>
              <span class="stat-icon pr">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="18" cy="18" r="3" />
                  <circle cx="6" cy="6" r="3" />
                  <path d="M13 6h3a2 2 0 0 1 2 2v7" />
                  <line x1="6" y1="9" x2="6" y2="21" />
                </svg>
              </span>
            </div>
            <div class="stat-value">{{ stats.prs }}</div>
            <div class="stat-change negative">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="2 12 6 12 9 3 15 21 18 12 22 12" />
              </svg>
              -{{ stats.prsChange }}% 较上月
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-header">
              <h3>议题</h3>
              <span class="stat-icon issue">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10" />
                  <line x1="12" y1="8" x2="12" y2="12" />
                  <line x1="12" y1="16" x2="12.01" y2="16" />
                </svg>
              </span>
            </div>
            <div class="stat-value">{{ stats.issues }}</div>
            <div class="stat-change positive">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="22 12 18 12 15 21 9 3 6 12 2 12" />
              </svg>
              +{{ stats.issuesChange }}% 较上月
            </div>
          </div>
        </div>

        <!-- 图表区域 -->
        <div class="charts-section">
          <!-- 代码提交趋势 -->
          <div class="chart-card">
            <div class="chart-header">
              <h3>代码提交趋势</h3>
              <div class="chart-controls">
                <select v-model="commitChartRange" class="chart-select">
                  <option value="7">最近 7 天</option>
                  <option value="30">最近 30 天</option>
                  <option value="90">最近 90 天</option>
                  <option value="365">最近一年</option>
                </select>
              </div>
            </div>
            <div class="chart-content">
              <div v-if="loading" class="chart-loading">
                <Skeleton type="bar" :count="7" />
              </div>
              <div v-else class="commit-chart">
                <div class="chart-bars">
                  <div v-for="(day, index) in commitData" :key="index" class="chart-bar">
                    <div class="bar-container">
                      <div class="bar" :style="{ height: (day.commits / maxCommits) * 100 + '%' }"></div>
                    </div>
                    <div class="bar-label">{{ day.date }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 项目活跃度 -->
          <div class="chart-card">
            <div class="chart-header">
              <h3>项目活跃度</h3>
              <div class="chart-controls">
                <select v-model="activityChartRange" class="chart-select">
                  <option value="7">最近 7 天</option>
                  <option value="30">最近 30 天</option>
                  <option value="90">最近 90 天</option>
                </select>
              </div>
            </div>
            <div class="chart-content">
              <div v-if="loading" class="chart-loading">
                <Skeleton type="bar" :count="5" />
              </div>
              <div v-else class="activity-chart">
                <div class="activity-list">
                  <div v-for="(repo, index) in activityData" :key="index" class="activity-item">
                    <div class="activity-info">
                      <div class="repo-name">{{ repo.name }}</div>
                      <div class="repo-stats">
                        <span class="stat-item">{{ repo.commits }} 次提交</span>
                        <span class="stat-item">{{ repo.prs }} 个 PR</span>
                        <span class="stat-item">{{ repo.issues }} 个 Issue</span>
                      </div>
                    </div>
                    <div class="activity-progress">
                      <div class="progress-bar">
                        <div class="progress-fill" :style="{ width: repo.activity + '%' }"></div>
                      </div>
                      <div class="progress-value">{{ repo.activity }}%</div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 最近活动 -->
        <div class="recent-activity">
          <div class="section-header">
            <h3>最近活动</h3>
            <router-link to="/activity" class="view-all">查看全部</router-link>
          </div>
          <div class="activity-list">
            <div v-if="loading" class="activity-loading">
              <Skeleton type="list" :count="5" />
            </div>
            <div v-else v-for="(activity, index) in recentActivities" :key="index" class="activity-item">
              <div class="activity-avatar" :class="activity.type">
                <span class="avatar-icon">{{ activity.icon }}</span>
              </div>
              <div class="activity-content">
                <div class="activity-header">
                  <span class="activity-type">{{ activity.typeText }}</span>
                  <span class="activity-time">{{ activity.time }}</span>
                </div>
                <div class="activity-description">{{ activity.description }}</div>
                <div v-if="activity.repo" class="activity-repo">
                  <span class="repo-icon">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z" />
                      <polyline points="13 2 13 9 20 9" />
                    </svg>
                  </span>
                  <span class="repo-name">{{ activity.repo }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 热门仓库 -->
        <div class="popular-repos">
          <div class="section-header">
            <h3>热门仓库</h3>
            <router-link to="/repos" class="view-all">查看全部</router-link>
          </div>
          <div class="repo-grid">
            <div v-if="loading" class="repo-loading">
              <Skeleton type="card" :count="4" />
            </div>
            <div v-else v-for="(repo, index) in popularRepos" :key="index" class="repo-card">
              <div class="repo-header">
                <div class="repo-visibility" :class="repo.visibility">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
                  </svg>
                </div>
                <div class="repo-name">{{ repo.name }}</div>
              </div>
              <div class="repo-description">{{ repo.description }}</div>
              <div class="repo-stats">
                <div class="stat-item">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="22 12 18 12 15 21 9 3 6 12 2 12" />
                  </svg>
                  {{ repo.stars }}
                </div>
                <div class="stat-item">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M5 5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2v16l-7-3.5L5 21V5z" />
                  </svg>
                  {{ repo.forks }}
                </div>
                <div class="stat-item">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10" />
                    <line x1="12" y1="8" x2="12" y2="12" />
                    <line x1="12" y1="16" x2="12.01" y2="16" />
                  </svg>
                  {{ repo.issues }}
                </div>
              </div>
              <div class="repo-language">
                <div class="language-color" :style="{ backgroundColor: repo.languageColor }"></div>
                <span class="language-name">{{ repo.language }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Skeleton } from '../components'
import { useUserStore } from '../stores'

// 用户状态
const userStore = useUserStore()
const isAdmin = computed(() => userStore.isAdmin)
const isUser = computed(() => userStore.isUser)

// 加载状态
const loading = ref(true)

// 统计数据
const stats = ref({
  repos: 12,
  reposChange: 15,
  commits: 142,
  commitsChange: 23,
  prs: 28,
  prsChange: 8,
  issues: 16,
  issuesChange: 12
})

// 管理员统计数据
const adminStats = ref({
  totalUsers: 128,
  totalRepos: 456,
  totalCommits: 3240,
  totalPrs: 528,
  activeUsers: 86,
  systemHealth: 98
})

// 图表范围
const commitChartRange = ref('30')
const activityChartRange = ref('30')

// 代码提交数据
const commitData = ref([
  { date: '1月', commits: 12 },
  { date: '2月', commits: 18 },
  { date: '3月', commits: 25 },
  { date: '4月', commits: 32 },
  { date: '5月', commits: 28 },
  { date: '6月', commits: 35 },
  { date: '7月', commits: 42 }
])

// 最大提交数（用于计算图表高度）
const maxCommits = computed(() => {
  return Math.max(...commitData.value.map(item => item.commits))
})

// 项目活跃度数据
const activityData = ref([
  { name: 'frontend-app', commits: 42, prs: 8, issues: 3, activity: 85 },
  { name: 'backend-api', commits: 38, prs: 12, issues: 5, activity: 78 },
  { name: 'mobile-app', commits: 25, prs: 5, issues: 2, activity: 65 },
  { name: 'ui-components', commits: 18, prs: 3, issues: 1, activity: 45 },
  { name: 'api-gateway', commits: 12, prs: 2, issues: 0, activity: 32 }
])

// 最近活动数据
const recentActivities = ref([
  {
    type: 'commit',
    typeText: '代码提交',
    icon: '📝',
    description: '在 frontend-app 仓库中提交了新的功能代码',
    repo: 'frontend-app',
    time: '2分钟前'
  },
  {
    type: 'pr',
    typeText: '创建 PR',
    icon: '🔄',
    description: '在 backend-api 仓库中创建了新的 Pull Request',
    repo: 'backend-api',
    time: '15分钟前'
  },
  {
    type: 'issue',
    typeText: '创建 Issue',
    icon: '🐛',
    description: '在 ui-components 仓库中创建了新的 Issue',
    repo: 'ui-components',
    time: '30分钟前'
  },
  {
    type: 'star',
    typeText: '星标仓库',
    icon: '⭐',
    description: '星标了 awesome-vue 仓库',
    repo: 'awesome-vue',
    time: '1小时前'
  },
  {
    type: 'commit',
    typeText: '代码提交',
    icon: '📝',
    description: '在 mobile-app 仓库中修复了 bug',
    repo: 'mobile-app',
    time: '2小时前'
  }
])

// 热门仓库数据
const popularRepos = ref([
  {
    name: 'frontend-app',
    description: '前端应用主仓库，包含用户界面和交互逻辑',
    stars: 45,
    forks: 12,
    issues: 3,
    language: 'Vue',
    languageColor: '#4fc08d',
    visibility: 'public'
  },
  {
    name: 'backend-api',
    description: '后端 API 服务，提供数据接口和业务逻辑',
    stars: 38,
    forks: 8,
    issues: 5,
    language: 'Node.js',
    languageColor: '#339933',
    visibility: 'public'
  },
  {
    name: 'mobile-app',
    description: '移动应用客户端，支持 iOS 和 Android',
    stars: 25,
    forks: 5,
    issues: 2,
    language: 'React Native',
    languageColor: '#007acc',
    visibility: 'private'
  },
  {
    name: 'ui-components',
    description: '可复用的 UI 组件库',
    stars: 18,
    forks: 3,
    issues: 1,
    language: 'TypeScript',
    languageColor: '#3178c6',
    visibility: 'public'
  }
])

// 系统健康数据
const systemHealthData = ref([
  { name: 'API 服务', status: 'online', uptime: '99.9%' },
  { name: '数据库', status: 'online', uptime: '99.8%' },
  { name: '存储服务', status: 'online', uptime: '100%' },
  { name: '缓存服务', status: 'online', uptime: '99.9%' }
])

// 模拟加载数据
const loadData = async () => {
  try {
    loading.value = true
    // 模拟 API 调用
    await new Promise(resolve => setTimeout(resolve, 800))
  } catch (error) {
    console.error('加载数据失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.dashboard-page {
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

.dashboard-container {
  max-width: 1400px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 32px;
}

/* 快速操作 */
.quick-actions {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 16px;
}

.action-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  box-shadow: var(--shadow-sm);
  cursor: pointer;
  transition: all 0.15s ease-out;
  text-align: left;
}

.action-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
  border-color: var(--color-primary);
}

.action-card.create-repo {
  border-left: 4px solid var(--color-primary);
}

.action-card.import-repo {
  border-left: 4px solid var(--color-success);
}

.action-card.new-pr {
  border-left: 4px solid var(--color-warning);
}

.action-card.new-issue {
  border-left: 4px solid var(--color-danger);
}

.action-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.create-repo .action-icon {
  background: var(--color-primary-bg);
  color: var(--color-primary);
}

.import-repo .action-icon {
  background: var(--color-success-bg);
  color: var(--color-success);
}

.new-pr .action-icon {
  background: var(--color-warning-bg);
  color: var(--color-warning);
}

.new-issue .action-icon {
  background: var(--color-danger-bg);
  color: var(--color-danger);
}

.action-content h3 {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 4px 0;
}

.action-content p {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin: 0;
}

/* 数据概览 */
.stats-overview {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 16px;
}

.stat-card {
  padding: 24px;
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

.stat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.stat-header h3 {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-secondary);
  margin: 0;
}

.stat-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-icon.repo {
  background: var(--color-primary-bg);
  color: var(--color-primary);
}

.stat-icon.commit {
  background: var(--color-success-bg);
  color: var(--color-success);
}

.stat-icon.pr {
  background: var(--color-warning-bg);
  color: var(--color-warning);
}

.stat-icon.issue {
  background: var(--color-danger-bg);
  color: var(--color-danger);
}

.stat-value {
  font-size: 2rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 8px 0;
  letter-spacing: -0.025em;
}

.stat-change {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 0.8125rem;
  font-weight: 500;
}

.stat-change.positive {
  color: var(--color-success);
}

.stat-change.negative {
  color: var(--color-danger);
}

/* 图表区域 */
.charts-section {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
}

.chart-card {
  padding: 24px;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  box-shadow: var(--shadow-sm);
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.chart-header h3 {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.chart-select {
  padding: 6px 12px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  background: var(--bg-secondary);
  color: var(--text-primary);
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.chart-select:hover {
  border-color: var(--color-primary);
}

.chart-content {
  min-height: 200px;
}

.chart-loading {
  padding: 24px 0;
}

/* 代码提交图表 */
.commit-chart {
  height: 200px;
}

.chart-bars {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  height: 100%;
  gap: 8px;
}

.chart-bar {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.bar-container {
  flex: 1;
  width: 100%;
  min-height: 20px;
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.bar {
  width: 100%;
  max-width: 32px;
  background: linear-gradient(180deg, var(--color-primary), var(--color-primary-dark));
  border-radius: 4px 4px 0 0;
  transition: height 0.3s ease-out;
}

.bar-label {
  font-size: 0.75rem;
  color: var(--text-muted);
  text-align: center;
  white-space: nowrap;
}

/* 项目活跃度图表 */
.activity-chart {
  height: 200px;
}

.activity-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.activity-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.activity-info {
  flex: 1;
  min-width: 0;
}

.repo-name {
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.repo-stats {
  display: flex;
  gap: 16px;
  font-size: 0.75rem;
  color: var(--text-muted);
}

.activity-progress {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 120px;
}

.progress-bar {
  flex: 1;
  height: 6px;
  background: var(--bg-secondary);
  border-radius: 3px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--color-primary), var(--color-primary-light));
  border-radius: 3px;
  transition: width 0.3s ease-out;
}

.progress-value {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--text-secondary);
  min-width: 36px;
  text-align: right;
}

/* 最近活动 */
.recent-activity {
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  box-shadow: var(--shadow-sm);
  padding: 24px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.section-header h3 {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.view-all {
  font-size: 0.875rem;
  color: var(--color-primary);
  text-decoration: none;
  font-weight: 500;
  transition: all 0.15s ease-out;
}

.view-all:hover {
  text-decoration: underline;
  color: var(--color-primary-hover);
}

.activity-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.activity-loading {
  padding: 16px 0;
}

.activity-item {
  display: flex;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid var(--border-color-light);
}

.activity-item:last-child {
  border-bottom: none;
}

.activity-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  color: white;
  font-size: 14px;
}

.activity-avatar.commit {
  background: var(--color-primary);
}

.activity-avatar.pr {
  background: var(--color-success);
}

.activity-avatar.issue {
  background: var(--color-warning);
}

.activity-avatar.star {
  background: var(--color-danger);
}

.activity-content {
  flex: 1;
  min-width: 0;
}

.activity-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.activity-type {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
}

.activity-time {
  font-size: 0.75rem;
  color: var(--text-muted);
  white-space: nowrap;
}

.activity-description {
  font-size: 0.9375rem;
  color: var(--text-secondary);
  margin: 0 0 8px 0;
  line-height: 1.4;
}

.activity-repo {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.8125rem;
  color: var(--color-primary);
}

.repo-icon {
  flex-shrink: 0;
}

/* 热门仓库 */
.popular-repos {
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  box-shadow: var(--shadow-sm);
  padding: 24px;
}

.repo-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.repo-loading {
  grid-column: 1 / -1;
}

.repo-card {
  padding: 20px;
  border: 1px solid var(--border-color-light);
  border-radius: 8px;
  background: var(--bg-secondary);
  transition: all 0.15s ease-out;
}

.repo-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-sm);
  border-color: var(--color-primary);
}

.repo-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.repo-visibility {
  width: 20px;
  height: 20px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
}

.repo-visibility.public {
  background: var(--color-success-bg);
  color: var(--color-success);
}

.repo-visibility.private {
  background: var(--color-warning-bg);
  color: var(--color-warning);
}

.repo-name {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 1rem;
}

.repo-description {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin: 0 0 16px 0;
  line-height: 1.4;
  min-height: 40px;
}

.repo-stats {
  display: flex;
  gap: 16px;
  margin-bottom: 12px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 0.75rem;
  color: var(--text-muted);
}

.repo-language {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.75rem;
  color: var(--text-muted);
}

.language-color {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  flex-shrink: 0;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .charts-section {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .dashboard-page {
    padding: 16px;
  }
  
  .quick-actions {
    grid-template-columns: 1fr;
  }
  
  .stats-overview {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .repo-grid {
    grid-template-columns: 1fr;
  }
  
  .action-card {
    padding: 20px;
  }
  
  .stat-card {
    padding: 20px;
  }
  
  .chart-card {
    padding: 20px;
  }
  
  .recent-activity,
  .popular-repos {
    padding: 20px;
  }
}

/* 管理员相关样式 */
.admin-section {
  display: flex;
  flex-direction: column;
  gap: 32px;
  margin-bottom: 32px;
}

.admin-overview {
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  box-shadow: var(--shadow-sm);
  padding: 24px;
}

.admin-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-top: 16px;
}

.admin-stat-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 20px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color-light);
  border-radius: 8px;
  transition: all 0.15s ease-out;
}

.admin-stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-sm);
  border-color: var(--color-primary);
}

.admin-stat-card .stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.admin-stat-card .stat-icon.users {
  background: var(--color-primary-bg);
  color: var(--color-primary);
}

.admin-stat-card .stat-icon.repos {
  background: var(--color-success-bg);
  color: var(--color-success);
}

.admin-stat-card .stat-icon.commits {
  background: var(--color-warning-bg);
  color: var(--color-warning);
}

.admin-stat-card .stat-icon.health {
  background: var(--color-danger-bg);
  color: var(--color-danger);
}

.admin-stat-card .stat-info {
  flex: 1;
}

.admin-stat-card .stat-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.admin-stat-card .stat-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.admin-actions {
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  box-shadow: var(--shadow-sm);
  padding: 24px;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 16px;
  margin-top: 16px;
}

.admin-action-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color-light);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.15s ease-out;
  text-align: left;
  border-left: 4px solid var(--color-primary);
}

.admin-action-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
  border-color: var(--color-primary);
}

.admin-action-card .action-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: var(--color-primary-bg);
  color: var(--color-primary);
}

.admin-action-card .action-content h3 {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 4px 0;
}

.admin-action-card .action-content p {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin: 0;
}

.system-health {
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  box-shadow: var(--shadow-sm);
  padding: 24px;
}

.health-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-top: 16px;
}

.health-card {
  padding: 20px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color-light);
  border-radius: 8px;
  transition: all 0.15s ease-out;
}

.health-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-sm);
  border-color: var(--color-primary);
}

.health-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.service-name {
  font-weight: 500;
  color: var(--text-primary);
}

.service-status {
  font-size: 0.75rem;
  font-weight: 500;
  padding: 4px 8px;
  border-radius: 12px;
}

.service-status.online {
  background: var(--color-success-bg);
  color: var(--color-success);
}

.service-status.offline {
  background: var(--color-danger-bg);
  color: var(--color-danger);
}

.health-uptime {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.875rem;
}

.uptime-label {
  color: var(--text-secondary);
}

.uptime-value {
  font-weight: 500;
  color: var(--text-primary);
}

.common-section {
  display: flex;
  flex-direction: column;
  gap: 32px;
}

@media (max-width: 1200px) {
  .charts-section {
    grid-template-columns: 1fr;
  }
  
  .admin-stats {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .health-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .dashboard-page {
    padding: 16px;
  }
  
  .quick-actions {
    grid-template-columns: 1fr;
  }
  
  .stats-overview {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .repo-grid {
    grid-template-columns: 1fr;
  }
  
  .action-card {
    padding: 20px;
  }
  
  .stat-card {
    padding: 20px;
  }
  
  .chart-card {
    padding: 20px;
  }
  
  .recent-activity,
  .popular-repos,
  .admin-overview,
  .admin-actions,
  .system-health {
    padding: 20px;
  }
  
  .admin-stats {
    grid-template-columns: 1fr;
  }
  
  .health-grid {
    grid-template-columns: 1fr;
  }
  
  .action-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 480px) {
  .stats-overview {
    grid-template-columns: 1fr;
  }
  
  .activity-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .activity-progress {
    width: 100%;
    min-width: unset;
  }
  
  .action-card {
    flex-direction: column;
    text-align: center;
  }
  
  .action-icon {
    align-self: center;
  }
  
  .admin-stat-card {
    flex-direction: column;
    text-align: center;
  }
  
  .admin-stat-card .stat-icon {
    align-self: center;
  }
  
  .health-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .health-uptime {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
}
</style>
