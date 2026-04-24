<template>
  <div class="dashboard">
    <h1>项目仪表盘</h1>
    <div class="dashboard-cards">
      <div class="card">
        <h2>最近活动</h2>
        <ul class="activity-list">
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
        <ul class="pr-list">
          <li v-for="pr in recentPRs" :key="pr.id">
            <span class="pr-title">{{ pr.title }}</span>
            <span class="pr-repo">{{ pr.repo }}</span>
            <span class="pr-status" :class="pr.status">{{ pr.status }}</span>
          </li>
        </ul>
      </div>
      <div class="card">
        <h2>最近 Issue</h2>
        <ul class="issue-list">
          <li v-for="issue in recentIssues" :key="issue.id">
            <span class="issue-title">{{ issue.title }}</span>
            <span class="issue-repo">{{ issue.repo }}</span>
            <span class="issue-status" :class="issue.status">{{ issue.status }}</span>
          </li>
        </ul>
      </div>
      <div class="card">
        <h2>仓库概览</h2>
        <div class="repo-stats">
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

// 模拟数据
const activities = ref([
  { id: 1, actor: 'user1', action: 'push', target: 'repo1', time: '2分钟前' },
  { id: 2, actor: 'user2', action: 'create', target: 'PR #123', time: '5分钟前' },
  { id: 3, actor: 'user3', action: 'comment', target: 'Issue #456', time: '10分钟前' },
  { id: 4, actor: 'user4', action: 'merge', target: 'PR #122', time: '15分钟前' }
])

const recentPRs = ref([
  { id: 123, title: 'Add new feature', repo: 'repo1', status: 'open' },
  { id: 122, title: 'Fix bug', repo: 'repo2', status: 'merged' },
  { id: 121, title: 'Update documentation', repo: 'repo3', status: 'closed' }
])

const recentIssues = ref([
  { id: 456, title: 'Bug in login', repo: 'repo1', status: 'open' },
  { id: 455, title: 'Feature request', repo: 'repo2', status: 'open' },
  { id: 454, title: 'Performance issue', repo: 'repo3', status: 'closed' }
])

const repoCount = ref(10)
const activeRepoCount = ref(5)
const privateRepoCount = ref(3)

onMounted(() => {
  // 实际项目中这里会从 API 获取数据
  console.log('Dashboard mounted')
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
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.card h2 {
  margin-top: 0;
  font-size: 18px;
  color: #333;
}

.activity-list, .pr-list, .issue-list {
  list-style: none;
  padding: 0;
  margin: 10px 0 0 0;
}

.activity-list li, .pr-list li, .issue-list li {
  padding: 10px 0;
  border-bottom: 1px solid #eee;
}

.activity-list li:last-child, .pr-list li:last-child, .issue-list li:last-child {
  border-bottom: none;
}

.activity-actor {
  font-weight: 600;
  margin-right: 8px;
}

.activity-action {
  color: #666;
  margin-right: 8px;
}

.activity-target {
  color: #0366d6;
  margin-right: 8px;
}

.activity-time {
  color: #999;
  font-size: 12px;
}

.pr-title, .issue-title {
  display: block;
  font-weight: 600;
  margin-bottom: 4px;
}

.pr-repo, .issue-repo {
  display: block;
  color: #666;
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
  background: #d1e7dd;
  color: #0f5132;
}

.pr-status.merged {
  background: #cfe2ff;
  color: #084298;
}

.pr-status.closed, .issue-status.closed {
  background: #f8d7da;
  color: #842029;
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
  background: #f8f9fa;
  border-radius: 6px;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #333;
}

.stat-label {
  font-size: 14px;
  color: #666;
  margin-top: 4px;
}
</style>