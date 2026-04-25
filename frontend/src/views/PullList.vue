<template>
  <div class="pull-list">
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">合并请求</h1>
        <p class="page-subtitle">管理代码合并请求</p>
      </div>
      <div class="header-right">
        <button class="btn btn-primary">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19" />
            <line x1="5" y1="12" x2="19" y2="12" />
          </svg>
          新建 PR
        </button>
      </div>
    </div>

    <div class="filter-bar">
      <select v-model="filter" @change="loadPulls" class="filter-select">
        <option value="all">所有</option>
        <option value="open">打开</option>
        <option value="merged">已合并</option>
        <option value="closed">已关闭</option>
      </select>
    </div>

    <div class="pulls-container">
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
      <div v-else class="pulls-list">
        <div 
          v-for="(pr, index) in pulls.length ? pulls : mockPulls" 
          :key="pr.id" 
          class="pull-item"
        >
          <div class="pull-icon status-open">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="18" cy="18" r="3" />
              <circle cx="6" cy="6" r="3" />
              <path d="M13 6h3a2 2 0 0 1 2 2v7" />
              <line x1="6" y1="9" x2="6" y2="21" />
            </svg>
          </div>
          <div class="pull-content">
            <h3 class="pull-title">{{ pr.title }}</h3>
            <div class="pull-meta">
              <span class="pull-repo">{{ pr.repo || 'repo-' + index }}</span>
              <span class="pull-author">由 {{ pr.author || 'user-' + index }} 创建</span>
              <span class="pull-time">{{ pr.time || '2天前' }}</span>
            </div>
            <div class="pull-branch">
              <span class="branch-name">{{ pr.source || 'feature-' + index }} </span>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="9 18 15 12 9 6" />
              </svg>
              <span class="branch-name"> {{ pr.target || 'main' }}</span>
            </div>
          </div>
          <span class="status-badge status-open">{{ pr.status || 'OPEN' }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { Skeleton } from '../components';
import { prApi } from '../services/api';

const pulls = ref<any[]>([]);
const filter = ref('all');
const loading = ref(true);
const error = ref('');

const mockPulls = [
  {
    id: 1,
    title: 'feat: 添加用户认证功能',
    repo: 'main-project',
    author: '张三',
    time: '2天前',
    source: 'feature/auth',
    target: 'main',
    status: 'OPEN',
  },
  {
    id: 2,
    title: 'fix: 修复登录页面样式',
    repo: 'frontend-app',
    author: '李四',
    time: '3天前',
    source: 'fix/login-style',
    target: 'main',
    status: 'OPEN',
  },
];

const loadPulls = async () => {
  try {
    loading.value = true;
    error.value = '';
    const response = await prApi.listPRs({ state: filter.value });
    pulls.value = (response as any).items || [];
  } catch (err: any) {
    error.value = err.message || '加载失败';
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  loadPulls();
});
</script>

<style scoped>
.pull-list {
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

.pulls-container {
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

.pulls-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.pull-item {
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

.pull-item:hover {
  border-color: var(--color-primary);
  box-shadow: var(--shadow-md);
}

.pull-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.pull-icon.status-open {
  background: rgba(22, 163, 74, 0.1);
  color: var(--color-success);
}

.pull-content {
  flex: 1;
  min-width: 0;
}

.pull-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 8px 0;
}

.pull-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.pull-branch {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.875rem;
}

.branch-name {
  padding: 2px 8px;
  background: var(--bg-tertiary);
  color: var(--text-secondary);
  border-radius: 4px;
  font-family: var(--font-mono);
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

.status-badge.status-merged {
  background: rgba(79, 70, 229, 0.1);
  color: var(--color-primary);
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
