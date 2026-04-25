<template>
  <div class="issue-list">
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">议题</h1>
        <p class="page-subtitle">管理项目问题和任务</p>
      </div>
      <div class="header-right">
        <button class="btn btn-primary">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19" />
            <line x1="5" y1="12" x2="19" y2="12" />
          </svg>
          新建 Issue
        </button>
      </div>
    </div>

    <div class="filter-bar">
      <select v-model="filter" @change="loadIssues" class="filter-select">
        <option value="all">所有</option>
        <option value="open">打开</option>
        <option value="closed">已关闭</option>
      </select>
    </div>

    <div class="issues-container">
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
      <div v-else class="issues-list">
        <div 
          v-for="(issue, index) in issues.length ? issues : mockIssues" 
          :key="issue.id" 
          class="issue-item"
        >
          <div class="issue-icon status-open">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10" />
              <line x1="12" y1="8" x2="12" y2="12" />
              <line x1="12" y1="16" x2="12.01" y2="16" />
            </svg>
          </div>
          <div class="issue-content">
            <h3 class="issue-title">{{ issue.title }}</h3>
            <div class="issue-meta">
              <span class="issue-repo">{{ issue.repo || 'repo-' + index }}</span>
              <span class="issue-author">由 {{ issue.author || 'user-' + index }} 创建</span>
              <span class="issue-time">{{ issue.time || '1天前' }}</span>
            </div>
            <div class="issue-labels">
              <span 
                v-for="(label, labelIndex) in (issue.labels || ['bug', 'priority']) " 
                :key="labelIndex" 
                class="label-pill"
              >
                {{ label }}
              </span>
            </div>
          </div>
          <span class="status-badge status-open">{{ issue.status || 'OPEN' }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { Skeleton } from '../components';
import { issueApi } from '../services/api';

const issues = ref<any[]>([]);
const filter = ref('all');
const loading = ref(true);
const error = ref('');

const mockIssues = [
  {
    id: 1,
    title: '登录超时问题修复',
    repo: 'main-project',
    author: '张三',
    time: '1天前',
    labels: ['bug', 'high priority'],
    status: 'OPEN',
  },
  {
    id: 2,
    title: '添加暗黑模式支持',
    repo: 'frontend-app',
    author: '李四',
    time: '3天前',
    labels: ['feature', 'enhancement'],
    status: 'OPEN',
  },
];

const loadIssues = async () => {
  try {
    loading.value = true;
    error.value = '';
    const response = await issueApi.listIssues({ state: filter.value });
    issues.value = (response as any).items || [];
  } catch (err: any) {
    error.value = err.message || '加载失败';
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  loadIssues();
});
</script>

<style scoped>
.issue-list {
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

.issues-container {
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

.issues-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.issue-item {
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

.issue-item:hover {
  border-color: var(--color-primary);
  box-shadow: var(--shadow-md);
}

.issue-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.issue-icon.status-open {
  background: rgba(22, 163, 74, 0.1);
  color: var(--color-success);
}

.issue-content {
  flex: 1;
  min-width: 0;
}

.issue-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 8px 0;
}

.issue-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin-bottom: 12px;
}

.issue-labels {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.label-pill {
  padding: 4px 10px;
  background: var(--bg-tertiary);
  color: var(--text-secondary);
  border-radius: 999px;
  font-size: 0.8125rem;
}

.status-badge {
  padding: 6px 12px;
  border-radius: 8px;
  font-size: 0.75rem;
  font-weight: 600;
  flex-shrink: 0;
}

.status-badge.status-open {
  background: rgba(22, 163, 74, 0.1);
  color: var(--color-success);
}

.status-badge.status-closed {
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
