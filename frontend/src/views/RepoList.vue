<template>
  <div class="repo-list">
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">仓库</h1>
        <p class="page-subtitle">管理您的项目仓库</p>
      </div>
      <div class="header-right">
        <button class="btn btn-primary" @click="showCreateModal = true">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19" />
            <line x1="5" y1="12" x2="19" y2="12" />
          </svg>
          新建仓库
        </button>
      </div>
    </div>

    <div class="filter-bar">
      <div class="search-box">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="search-icon">
          <circle cx="11" cy="11" r="8" />
          <line x1="21" y1="21" x2="16.65" y2="16.65" />
        </svg>
        <input 
          type="text" 
          v-model="searchQuery" 
          @input="handleSearch" 
          placeholder="搜索仓库..." 
          class="search-input"
        >
        <button v-if="searchQuery" @click="clearSearch" class="clear-search">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18" />
            <line x1="6" y1="6" x2="18" y2="18" />
          </svg>
        </button>
      </div>
      <div class="filter-controls">
        <select v-model="filter" @change="loadRepos" class="filter-select">
          <option value="all">所有仓库</option>
          <option value="my">我的仓库</option>
          <option value="starred">已星标</option>
        </select>
        <select v-model="sortBy" @change="loadRepos" class="sort-select">
          <option value="updated">最近更新</option>
          <option value="created">创建时间</option>
          <option value="stars">星标数</option>
          <option value="forks">Fork 数</option>
        </select>
      </div>
    </div>

    <div class="repos-container">
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
      <div v-else class="repos-grid">
        <div v-for="repo in repos.length ? repos : mockRepos" :key="repo.id" class="repo-card">
          <div class="repo-header">
            <div class="repo-icon">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z" />
                <polyline points="13 2 13 9 20 9" />
              </svg>
            </div>
            <div class="repo-title-section">
              <h3 class="repo-name">{{ repo.name }}</h3>
              <span class="repo-visibility">{{ repo.visibility === 'private' ? '私密' : '公开' }}</span>
            </div>
          </div>
          
          <p class="repo-description">{{ repo.description || '暂无描述' }}</p>
          
          <div class="repo-meta">
            <div class="repo-stat">
              <span class="meta-icon">📄</span>
              <span class="meta-text">{{ repo.language || 'N/A' }}</span>
            </div>
            <div class="repo-stat">
              <span class="meta-icon">⭐</span>
              <span class="meta-text">{{ repo.stars || 0 }}</span>
            </div>
            <div class="repo-stat">
              <span class="meta-icon">🔀</span>
              <span class="meta-text">{{ repo.forks || 0 }}</span>
            </div>
          </div>
          
          <div class="repo-actions">
            <button class="btn btn-ghost btn-sm">查看</button>
            <button class="btn btn-ghost btn-sm">克隆</button>
          </div>
        </div>
      </div>
    </div>

    <Modal v-if="showCreateModal" @close="showCreateModal = false">
      <div class="create-modal">
        <h2>创建新仓库</h2>
        <form @submit.prevent="createRepo">
          <div class="form-group">
            <label for="repo-name">仓库名称</label>
            <input id="repo-name" v-model="newRepo.name" type="text" placeholder="输入仓库名称" />
          </div>
          <div class="form-group">
            <label for="repo-desc">描述</label>
            <textarea id="repo-desc" v-model="newRepo.description" placeholder="仓库描述" />
          </div>
          <div class="form-group">
            <label>可见性</label>
            <div class="radio-group">
              <label>
                <input type="radio" v-model="newRepo.visibility" value="public" />
                公开
              </label>
              <span class="radio-hint">所有人都可以看到这个仓库</span>
            </div>
            <div class="radio-group">
              <label>
                <input type="radio" v-model="newRepo.visibility" value="private" />
                私密
              </label>
              <span class="radio-hint">只有你和你指定的人可以看到</span>
            </div>
          </div>
          <div class="form-group">
            <label>
              <input type="checkbox" v-model="newRepo.initReadme" />
              初始化 README.md 仓库
            </label>
          </div>
          <div class="form-actions">
            <button type="button" class="btn btn-ghost" @click="showCreateModal = false">取消</button>
            <button type="submit" class="btn btn-primary">创建</button>
          </div>
        </form>
      </div>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { repoApi } from '../services/api';
import { Skeleton, Modal } from '../components';

const repos = ref<any[]>([]);
const filter = ref('all');
const sortBy = ref('updated');
const searchQuery = ref('');
const loading = ref(true);
const error = ref('');
const showCreateModal = ref(false);
const newRepo = ref({
  name: '',
  description: '',
  visibility: 'private',
  initReadme: false
});

const mockRepos = [
  {
    id: 1,
    name: 'frontend-app',
    description: '前端应用项目，使用 Vue 3 开发',
    visibility: 'public',
    language: 'Vue',
    stars: 128,
    forks: 34,
    updated_at: '2026-04-25T10:00:00Z',
    created_at: '2026-04-01T00:00:00Z'
  },
  {
    id: 2,
    name: 'backend-api',
    description: '后端 API 服务，使用 Go 开发',
    visibility: 'private',
    language: 'Go',
    stars: 87,
    forks: 15,
    updated_at: '2026-04-24T15:30:00Z',
    created_at: '2026-03-15T00:00:00Z'
  },
  {
    id: 3,
    name: 'ui-components',
    description: '通用 UI 组件库',
    visibility: 'public',
    language: 'TypeScript',
    stars: 256,
    forks: 67,
    updated_at: '2026-04-23T09:15:00Z',
    created_at: '2026-02-10T00:00:00Z'
  },
  {
    id: 4,
    name: 'mobile-app',
    description: '移动应用项目，使用 React Native 开发',
    visibility: 'public',
    language: 'JavaScript',
    stars: 45,
    forks: 8,
    updated_at: '2026-04-22T14:45:00Z',
    created_at: '2026-01-05T00:00:00Z'
  },
  {
    id: 5,
    name: 'api-gateway',
    description: 'API 网关服务',
    visibility: 'private',
    language: 'Node.js',
    stars: 23,
    forks: 3,
    updated_at: '2026-04-21T11:20:00Z',
    created_at: '2025-12-20T00:00:00Z'
  }
];

// 处理搜索
const handleSearch = () => {
  // 防抖处理
  clearTimeout(window.searchTimeout);
  window.searchTimeout = setTimeout(() => {
    loadRepos();
  }, 300);
};

// 清除搜索
const clearSearch = () => {
  searchQuery.value = '';
  loadRepos();
};

// 加载仓库
const loadRepos = async () => {
  try {
    loading.value = true;
    error.value = '';
    
    // 模拟 API 调用
    // const response = await repoApi.listRepos();
    // repos.value = (response as any).items || [];
    
    // 模拟搜索和排序
    let filteredRepos = [...mockRepos];
    
    // 搜索过滤
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase();
      filteredRepos = filteredRepos.filter(repo => 
        repo.name.toLowerCase().includes(query) || 
        repo.description.toLowerCase().includes(query)
      );
    }
    
    // 排序
    filteredRepos.sort((a, b) => {
      switch (sortBy.value) {
        case 'updated':
          return new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime();
        case 'created':
          return new Date(b.created_at).getTime() - new Date(a.created_at).getTime();
        case 'stars':
          return b.stars - a.stars;
        case 'forks':
          return b.forks - a.forks;
        default:
          return 0;
      }
    });
    
    repos.value = filteredRepos;
  } catch (err: any) {
    error.value = err.message || '加载失败';
  } finally {
    loading.value = false;
  }
};

const createRepo = async () => {
  showCreateModal.value = false;
  newRepo.value = { name: '', description: '', visibility: 'private', initReadme: false };
  loadRepos();
};

onMounted(() => {
  loadRepos();
});
</script>

<style scoped>
.repo-list {
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

.btn-ghost {
  background: transparent;
  border-color: var(--border-color);
}

.btn-ghost:hover {
  background: var(--bg-tertiary);
}

.btn-sm {
  padding: 6px 12px;
  font-size: 0.875rem;
}

.filter-bar {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.search-box {
  flex: 1;
  min-width: 300px;
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 12px;
  color: var(--text-muted);
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: 10px 12px 10px 40px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 0.9375rem;
  transition: all 0.15s ease-out;
}

.search-input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px var(--color-primary-bg);
}

.clear-search {
  position: absolute;
  right: 12px;
  background: none;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  transition: all 0.15s ease-out;
}

.clear-search:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.filter-controls {
  display: flex;
  gap: 12px;
  align-items: center;
}

.filter-select,
.sort-select {
  padding: 8px 12px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 0.9375rem;
  cursor: pointer;
  transition: all 0.15s ease-out;
  min-width: 120px;
}

.filter-select:hover,
.sort-select:hover {
  border-color: var(--color-primary);
}

.filter-select:focus,
.sort-select:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px var(--color-primary-bg);
}

.repos-container {
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

.repos-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 16px;
}

.repo-card {
  display: flex;
  flex-direction: column;
  padding: 20px;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  transition: all 0.15s ease-out;
}

.repo-card:hover {
  border-color: var(--color-primary);
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

.repo-header {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 12px;
}

.repo-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  background: var(--color-primary-bg);
  border-radius: 8px;
  color: var(--color-primary);
  flex-shrink: 0;
}

.repo-title-section {
  flex: 1;
  min-width: 0;
}

.repo-name {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.repo-visibility {
  display: inline-block;
  padding: 2px 8px;
  margin-top: 4px;
  font-size: 0.75rem;
  color: var(--text-muted);
  background: var(--bg-tertiary);
  border-radius: 4px;
}

.repo-description {
  font-size: 0.9375rem;
  color: var(--text-secondary);
  line-height: 1.5;
  margin-bottom: 16px;
}

.repo-meta {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
}

.repo-stat {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.repo-actions {
  display: flex;
  gap: 8px;
  margin-top: auto;
  padding-top: 16px;
  border-top: 1px solid var(--border-color-light);
}

.create-modal {
  width: 100%;
  max-width: 480px;
}

.create-modal h2 {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 24px 0;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-size: 0.9375rem;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 0.9375rem;
  transition: all 0.15s ease-out;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px var(--color-primary-bg);
}

.form-group textarea {
  resize: vertical;
  min-height: 80px;
}

.radio-group {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 12px;
  background: var(--bg-tertiary);
  border-radius: 8px;
  margin-bottom: 8px;
}

.radio-group input {
  margin-top: 3px;
}

.radio-group label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.radio-hint {
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 8px;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
  }

  .repos-grid {
    grid-template-columns: 1fr;
  }
}
</style>
