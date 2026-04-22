<template>
  <div class="cicd-list">
    <h1>所有 CI/CD 流水线</h1>
    <div class="cicd-actions">
      <button class="btn primary">运行流水线</button>
      <div class="filter-options">
        <select v-model="filter">
          <option value="all">所有流水线</option>
          <option value="running">运行中</option>
          <option value="success">成功</option>
          <option value="failed">失败</option>
        </select>
      </div>
    </div>
    <div class="pipelines-list">
      <div v-for="pipeline in pipelines" :key="pipeline.id" class="pipeline-item">
        <div class="pipeline-info">
          <h3 class="pipeline-name">流水线 #{{ pipeline.id }}</h3>
          <div class="pipeline-meta">
            <span class="pipeline-repo">{{ pipeline.repo }}</span>
            <span class="pipeline-author">由 {{ pipeline.author }} 触发</span>
            <span class="pipeline-time">{{ pipeline.time }}</span>
            <span class="pipeline-branch">{{ pipeline.branch }}</span>
          </div>
        </div>
        <div class="pipeline-status" :class="pipeline.status">
          {{ pipeline.status }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const filter = ref('all')
const pipelines = ref([
  {
    id: 1,
    repo: 'user1/laima',
    author: 'user1',
    time: '2小时前',
    branch: 'main',
    status: 'running'
  },
  {
    id: 2,
    repo: 'user2/frontend',
    author: 'user2',
    time: '1天前',
    branch: 'feature-1',
    status: 'success'
  },
  {
    id: 3,
    repo: 'user3/backend',
    author: 'user3',
    time: '3天前',
    branch: 'dev',
    status: 'failed'
  }
])
</script>

<style scoped>
.cicd-list {
  padding: 20px;
}

.cicd-actions {
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

.pipelines-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.pipeline-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: var(--bg-secondary);
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s, background-color 0.3s, box-shadow 0.3s;
}

.pipeline-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.pipeline-info {
  flex: 1;
}

.pipeline-name {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: var(--text-primary);
  transition: color 0.3s;
}

.pipeline-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  font-size: 14px;
  color: var(--text-secondary);
  transition: color 0.3s;
}

.pipeline-status {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}

.pipeline-status.running {
  background: var(--warning-color);
  color: #fff;
}

.pipeline-status.success {
  background: var(--success-color);
  color: #fff;
}

.pipeline-status.failed {
  background: var(--danger-color);
  color: #fff;
}
</style>