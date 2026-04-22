<template>
  <div class="issue-list">
    <h1>所有 Issue</h1>
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
            <span class="issue-repo">{{ issue.repo }}</span>
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
    repo: 'user1/laima',
    author: 'user1',
    time: '2小时前',
    labels: ['bug', 'high'],
    status: 'open'
  },
  {
    id: 2,
    title: 'Feature request',
    repo: 'user2/frontend',
    author: 'user2',
    time: '1天前',
    labels: ['feature', 'medium'],
    status: 'open'
  },
  {
    id: 3,
    title: 'Performance issue',
    repo: 'user3/backend',
    author: 'user3',
    time: '3天前',
    labels: ['performance', 'low'],
    status: 'closed'
  }
])
</script>

<style scoped>
.issue-list {
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
  border: 1px solid #ddd;
  border-radius: 4px;
  background: #fff;
  color: #333;
  cursor: pointer;
  transition: all 0.3s;
}

.btn.primary {
  background: #0366d6;
  color: #fff;
  border-color: #0366d6;
}

.btn.primary:hover {
  background: #0056b3;
}

.filter-options select {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
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
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s;
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
  color: #333;
}

.issue-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  font-size: 14px;
  color: #666;
}

.issue-label {
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  background: #e9ecef;
  color: #495057;
}

.issue-status {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}

.issue-status.open {
  background: #d1e7dd;
  color: #0f5132;
}

.issue-status.closed {
  background: #f8d7da;
  color: #842029;
}
</style>