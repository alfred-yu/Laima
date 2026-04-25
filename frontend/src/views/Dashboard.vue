<template>
  <div class="dashboard">
    <h1>项目仪表盘</h1>
    <div class="dashboard-cards">
      <div class="card">
        <h2>最近活动</h2>
        <div v-if="activitiesLoading" class="loading-state">
          <Skeleton type="list" :count="4" />
        </div>
        <div v-else-if="activitiesError" class="error-state">
          {{ activitiesError }}
        </div>
        <ul v-else class="activity-list">
          <li v-for="activity in activities" :key="activity.id">
            <span class="activity-actor">{{ activity.actor }}</span>
            <span class="activity-action">{{ activity.action }}</span>
            <span class="activity-target">{{ activity.target }}</span>
            <span class="activity-time">{{ activity.time }}</span>
          </li>
        </ul>
      </div>
      <div class="card">
        <h2>最近 PR</h2>
        <div v-if="prsLoading" class="loading-state">
          <Skeleton type="list" :count="3" />
        </div>
        <div v-else-if="prsError" class="error-state">
          {{ prsError }}
        </div>
        <ul v-else class="pr-list">
          <li v-for="pr in recentPRs" :key="pr.id">
            <span class="pr-title">{{ pr.title }}</span>
            <span class="pr-repo">{{ pr.repo }}</span>
            <span class="pr-status" :class="pr.status">{{ pr.status }}</span>
          </li>
        </ul>
      </div>
      <div class="card">
        <h2>最近 Issue</h2>
        <div v-if="issuesLoading" class="loading-state">
          <Skeleton type="list" :count="3" />
        </div>
        <div v-else-if="issuesError" class="error-state">
          {{ issuesError }}
        </div>
        <ul v-else class="issue-list">
          <li v-for="issue in recentIssues" :key="issue.id">
            <span class="issue-title">{{ issue.title }}</span>
            <span class="issue-repo">{{ issue.repo }}</span>
            <span class="issue-status" :class="issue.status">{{ issue.status }}</span>
          </li>
        </ul>
      </div>
      <div class="card">
        <h2>仓库概览</h2>
        <div v-if="statsLoading" class="loading-state">
          <Skeleton type="text" :count="3" />
        </div>
        <div v-else-if="statsError" class="error-state">
          {{ statsError }}
        </div>
        <div v-else class="repo-stats">
          <div class="stat-item">
            <div class="stat-value">{{ repoCount }}</div>
            <div class="stat-label">总仓库</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">{{ activeRepoCount }}</div>
            <div class="stat-label">活跃仓库</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">{{ privateRepoCount }}</div>
            <div class="stat-label">私有仓库</div>
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

const activities = ref<any[]>([])
const recentPRs = ref<any[]>([])
const recentIssues = ref<any[]>([])
const repoCount = ref(0)
const activeRepoCount = ref(0)
const privateRepoCount = ref(0)

const activitiesLoading = ref(true)
const prsLoading = ref(true)
const issuesLoading = ref(true)
const statsLoading = ref(true)
const activitiesError = ref('')
const prsError = ref('')
const issuesError = ref('')
const statsError = ref('')

const loadActivities = async () => {
  try {
    activitiesLoading.value = true
    activitiesError.value = ''
    // TODO: 当 API 支持活动接口后替换
    // const response = await activityApi.listActivities()
    // activities.value = response.items || []
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
    const response = await prApi.listPRs({ per_page: 3 })
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
    const response = await issueApi.listIssues({ per_page: 3 })
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
    repoCount.value = repos.length
    activeRepoCount.value = repos.length > 0 ? Math.min(3, repos.length) : 0
    privateRepoCount.value = repos.filter((r: any) => r.visibility === 'private').length
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
  padding: 20px;
}

.dashboard-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

.card {
  background: var(--bg-secondary);
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.card h2 {
  margin-top: 0;
  font-size: 18px;
  color: var(--text-primary);
}

.loading-state,
.error-state {
  padding: 20px 0;
  color: var(--text-secondary);
}

.error-state {
  color: var(--danger-color);
}

.activity-list, .pr-list, .issue-list {
  list-style: none;
  padding: 0;
  margin: 10px 0 0 0;
}

.activity-list li, .pr-list li, .issue-list li {
  padding: 10px 0;
  border-bottom: 1px solid var(--border-color);
}

.activity-list li:last-child, .pr-list li:last-child, .issue-list li:last-child {
  border-bottom: none;
}

.activity-actor {
  font-weight: 600;
  margin-right: 8px;
  color: var(--text-primary);
}

.activity-action {
  color: var(--text-secondary);
  margin-right: 8px;
}

.activity-target {
  color: var(--accent-color);
  margin-right: 8px;
}

.activity-time {
  color: var(--text-secondary);
  font-size: 12px;
}

.pr-title, .issue-title {
  display: block;
  font-weight: 600;
  margin-bottom: 4px;
  color: var(--text-primary);
}

.pr-repo, .issue-repo {
  display: block;
  color: var(--text-secondary);
  font-size: 14px;
  margin-bottom: 4px;
}

.pr-status, .issue-status {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}

.pr-status.open, .issue-status.open {
  background: var(--success-color);
  color: #fff;
}

.pr-status.merged {
  background: var(--accent-color);
  color: #fff;
}

.pr-status.closed, .issue-status.closed {
  background: var(--danger-color);
  color: #fff;
}

.repo-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
  margin-top: 10px;
}

.stat-item {
  text-align: center;
  padding: 10px;
  background: var(--bg-primary);
  border-radius: 6px;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
}

.stat-label {
  font-size: 14px;
  color: var(--text-secondary);
  margin-top: 4px;
}
</style>
