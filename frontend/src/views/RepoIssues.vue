<template>
  <div class="repo-issues">
    <h2>Issue 列表</h2>
    <div class="issues-actions">
      <button class="btn primary">新建 Issue</button>
      <div class="filter-options">
        <select v-model="filter">
          <option value="all">所有 Issue</option>
          <option value="open">打开</option>
          <option value="closed">已关闭</option>
        </select>
      </div>
    </div>
    <div class="issues-list">
      <div v-for="issue in issues" :key="issue.id" class="issue-item">
        <div class="issue-info">
          <h3 class="issue-title">{{ issue.title }}</h3>
          <div class="issue-meta">
            <span class="issue-author">由 {{ issue.author }} 创建</span>
            <span class="issue-time">{{ issue.time }}</span>
            <span v-for="label in issue.labels" :key="label" class="issue-label">{{ label }}</span>
          </div>
        </div>
        <div class="issue-status" :class="issue.status">
          {{ issue.status }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const filter = ref('all')
const issues = ref([
  {
    id: 1,
    title: 'Bug in login',
    author: 'user1',
    time: '2小时前',
    labels: ['bug', 'high'],
    status: 'open'
  },
  {
    id: 2,
    title: 'Feature request',
    author: 'user2',
    time: '1天前',
    labels: ['feature', 'medium'],
    status: 'open'
  },
  {
    id: 3,
    title: 'Performance issue',
    author: 'user3',
    time: '3天前',
    labels: ['performance', 'low'],
    status: 'closed'
  }
])
</script>

<style scoped>
.repo-issues {
  padding: 20px;
}

.issues-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.btn {
  padding: 8px 16px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--bg-primary);
  color: var(--text-primary);
  cursor: pointer;
  transition: all 0.3s;
}

.btn.primary {
  background: var(--accent-color);
  color: #fff;
  border-color: var(--accent-color);
}

.btn.primary:hover {
  opacity: 0.9;
}

.filter-options select {
  padding: 8px 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 14px;
  background: var(--bg-primary);
  color: var(--text-primary);
  transition: border-color 0.3s, background-color 0.3s, color 0.3s;
}

.issues-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.issue-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: var(--bg-secondary);
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s, background-color 0.3s, box-shadow 0.3s;
}

.issue-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.issue-info {
  flex: 1;
}

.issue-title {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: var(--text-primary);
  transition: color 0.3s;
}

.issue-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  font-size: 14px;
  color: var(--text-secondary);
  transition: color 0.3s;
}

.issue-label {
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  background: var(--border-color);
  color: var(--text-secondary);
  transition: background-color 0.3s, color 0.3s;
}

.issue-status {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}

.issue-status.open {
  background: var(--success-color);
  color: #fff;
}

.issue-status.closed {
  background: var(--danger-color);
  color: #fff;
}
</style>