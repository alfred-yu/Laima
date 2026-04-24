<template>
  <div class="pull-list">
    <h1>所有 PR</h1>
    <div class="pulls-actions">
      <button class="btn primary">新建 PR</button>
      <div class="filter-options">
        <select v-model="filter">
          <option value="all">所有 PR</option>
          <option value="open">打开</option>
          <option value="merged">已合并</option>
          <option value="closed">已关闭</option>
        </select>
      </div>
    </div>
    <div class="pulls-list">
      <div v-for="pr in pulls" :key="pr.id" class="pull-item">
        <div class="pull-info">
          <h3 class="pull-title">{{ pr.title }}</h3>
          <div class="pull-meta">
            <span class="pull-repo">{{ pr.repo }}</span>
            <span class="pull-author">由 {{ pr.author }} 创建</span>
            <span class="pull-time">{{ pr.time }}</span>
            <span class="pull-branch">{{ pr.source }} → {{ pr.target }}</span>
          </div>
        </div>
        <div class="pull-status" :class="pr.status">
          {{ pr.status }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const filter = ref('all')
const pulls = ref([
  {
    id: 1,
    title: 'Add new feature',
    repo: 'user1/laima',
    author: 'user1',
    time: '2小时前',
    source: 'feature-1',
    target: 'main',
    status: 'open'
  },
  {
    id: 2,
    title: 'Fix bug',
    repo: 'user2/frontend',
    author: 'user2',
    time: '1天前',
    source: 'bugfix-1',
    target: 'main',
    status: 'merged'
  },
  {
    id: 3,
    title: 'Update documentation',
    repo: 'user3/backend',
    author: 'user3',
    time: '3天前',
    source: 'docs-update',
    target: 'main',
    status: 'closed'
  }
])
</script>

<style scoped>
.pull-list {
  padding: 20px;
}

.pulls-actions {
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

.pulls-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.pull-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: var(--bg-secondary);
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s, background-color 0.3s, box-shadow 0.3s;
}

.pull-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.pull-info {
  flex: 1;
}

.pull-title {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: var(--text-primary);
  transition: color 0.3s;
}

.pull-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  font-size: 14px;
  color: var(--text-secondary);
  transition: color 0.3s;
}

.pull-status {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}

.pull-status.open {
  background: var(--success-color);
  color: #fff;
}

.pull-status.merged {
  background: var(--accent-color);
  color: #fff;
}

.pull-status.closed {
  background: var(--danger-color);
  color: #fff;
}
</style>