<template>
  <div class="cicd-list">
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">CI/CD</h1>
        <p class="page-subtitle">持续集成和部署流水线</p>
      </div>
      <div class="header-right">
        <button class="btn btn-primary">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 12h-4l-2 2M2 12h4l2-2" />
            <polygon points="8 7 12 3 16 7" />
            <polygon points="16 17 12 21 8 17" />
          </svg>
          运行流水线
        </button>
      </div>
    </div>

    <div class="filter-bar">
      <select v-model="filter" @change="loadPipelines" class="filter-select">
        <option value="all">所有</option>
        <option value="running">运行中</option>
        <option value="success">成功</option>
        <option value="failed">失败</option>
      </select>
    </div>

    <div class="pipelines-container">
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
      </div>
      <div v-else class="pipelines-list">
        <div 
          v-for="(pipeline, index) in pipelines.length ? pipelines : mockPipelines" 
          :key="pipeline.id" 
          class="pipeline-item"
        >
          <div class="pipeline-icon status-success">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10" />
              <polyline points="12 6 12 12 16 14" />
            </svg>
          </div>
          <div class="pipeline-content">
            <h3 class="pipeline-title">流水线 #{{ pipeline.id || index + 1 }}</h3>
            <div class="pipeline-meta">
              <span class="pipeline-repo">{{ pipeline.repo || 'main-project' }}</span>
              <span class="pipeline-author">由 {{ pipeline.author || 'user' }} 触发</span>
              <span class="pipeline-time">{{ pipeline.time || '3小时前' }}</span>
              <span class="pipeline-branch">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="18" r="3" />
                  <circle cx="6" cy="6" r="3" />
                  <path d="M18 9a9 9 0 0 1-9 9" />
                </svg>
                {{ pipeline.branch || 'main' }}
              </span>
            </div>
          </div>
          <span class="status-badge status-success">{{ pipeline.status || 'SUCCESS' }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { Skeleton } from '../components';
import { cicdApi } from '../services/api';

const pipelines = ref<any[]>([]);
const filter = ref('all');
const loading = ref(true);
const error = ref('');

const mockPipelines = [
  {
    id: 1,
    repo: 'main-project',
    author: '张三',
    time: '3小时前',
    branch: 'main',
    status: 'SUCCESS',
  },
  {
    id: 2,
    repo: 'frontend-app',
    author: '李四',
    time: '5小时前',
    branch: 'feature/auth',
    status: 'RUNNING',
  },
];

const loadPipelines = async () => {
  try {
    loading.value = true;
    error.value = '';
    const response = await cicdApi.listPipelines({ status: filter.value });
    pipelines.value = (response as any).items || [];
  } catch (err: any) {
    error.value = err.message || '加载失败';
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  loadPipelines();
});
</script>

<style scoped>
.cicd-list {
  width: 100%;
}

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

.btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 0.9375rem;
  font-weight: 500;
  background: var(--bg-primary);
  color: var(--text-primary);
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.btn-primary {
  background: var(--color-primary);
  border-color: var(--color-primary);
  color: white;
}

.btn-primary:hover {
  background: var(--color-primary-hover);
  border-color: var(--color-primary-hover);
}

.filter-bar {
  display: flex;
  align-items: center;
  margin-bottom: 24px;
}

.filter-select {
  padding: 8px 12px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 0.9375rem;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.pipelines-container {
  min-height: 400px;
}

.loading-state {
  padding: 24px;
}

.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 64px 24px;
  color: var(--text-secondary);
  gap: 16px;
}

.pipelines-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.pipeline-item {
  display: flex;
  gap: 16px;
  padding: 20px;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  align-items: flex-start;
  transition: all 0.15s ease-out;
  cursor: pointer;
}

.pipeline-item:hover {
  border-color: var(--color-primary);
  box-shadow: var(--shadow-md);
}

.pipeline-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.pipeline-icon.status-success {
  background: rgba(22, 163, 74, 0.1);
  color: var(--color-success);
}

.pipeline-icon.status-running {
  background: rgba(217, 119, 6, 0.1);
  color: var(--color-warning);
}

.pipeline-icon.status-failed {
  background: rgba(220, 38, 38, 0.1);
  color: var(--color-danger);
}

.pipeline-content {
  flex: 1;
  min-width: 0;
}

.pipeline-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 8px 0;
}

.pipeline-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.pipeline-branch {
  display: flex;
  align-items: center;
  gap: 4px;
}

.status-badge {
  padding: 6px 12px;
  border-radius: 8px;
  font-size: 0.75rem;
  font-weight: 600;
  flex-shrink: 0;
}

.status-badge.status-success {
  background: rgba(22, 163, 74, 0.1);
  color: var(--color-success);
}

.status-badge.status-running {
  background: rgba(217, 119, 6, 0.1);
  color: var(--color-warning);
}

.status-badge.status-failed {
  background: rgba(220, 38, 38, 0.1);
  color: var(--color-danger);
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
  }
}
</style>
