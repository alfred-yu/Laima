<template>
  <div class="repo-cicd">
    <h2>CI/CD 流水线</h2>
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
    author: 'user1',
    time: '2小时前',
    branch: 'main',
    status: 'running'
  },
  {
    id: 2,
    author: 'user2',
    time: '1天前',
    branch: 'feature-1',
    status: 'success'
  },
  {
    id: 3,
    author: 'user3',
    time: '3天前',
    branch: 'dev',
    status: 'failed'
  }
])
</script>

<style scoped>
.repo-cicd {
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
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s;
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
  color: #333;
}

.pipeline-meta {
  display: flex;
  gap: 16px;
  font-size: 14px;
  color: #666;
}

.pipeline-status {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}

.pipeline-status.running {
  background: #fff3cd;
  color: #856404;
}

.pipeline-status.success {
  background: #d1e7dd;
  color: #0f5132;
}

.pipeline-status.failed {
  background: #f8d7da;
  color: #842029;
}
</style>