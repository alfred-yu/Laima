<template>
  <div class="repo-list">
    <h1>仓库列表</h1>
    <div class="repo-actions">
      <button class="create-repo-btn" @click="showCreateModal = true">创建仓库</button>
      <div class="filter-options">
        <select v-model="filter">
          <option value="all">所有仓库</option>
          <option value="owned">我的仓库</option>
          <option value="starred">已星标</option>
        </select>
      </div>
    </div>

    <div v-if="isLoading" class="loading">加载中...</div>
    <div v-else-if="error" class="error-message">{{ error }}</div>
    <div v-else class="repo-grid">
      <div v-for="repo in repos" :key="repo.id" class="repo-card">
        <div class="repo-header">
          <h2 class="repo-name">{{ repo.name }}</h2>
          <span class="repo-visibility" :class="repo.visibility">{{ repo.visibility }}</span>
        </div>
        <p class="repo-description">{{ repo.description || '暂无描述' }}</p>
        <div class="repo-stats">
          <span class="stat">📁 {{ repo.default_branch || 'main' }}</span>
        </div>
        <div class="repo-actions">
          <router-link :to="`/repos/${repo.full_path?.split('/')[0] || ''}/${repo.name}/code`" class="btn">查看</router-link>
          <button class="btn secondary">克隆</button>
        </div>
      </div>
    </div>

    <!-- 创建仓库模态框 -->
    <div v-if="showCreateModal" class="modal-overlay" @click="showCreateModal = false">
      <div class="modal" @click.stop>
        <h2>创建新仓库</h2>
        <form @submit.prevent="handleCreateRepo">
          <div class="form-group">
            <label for="repo-name">仓库名称 *</label>
            <input 
              type="text" 
              id="repo-name" 
              v-model="createForm.name" 
              required 
              placeholder="输入仓库名称"
            >
          </div>
          <div class="form-group">
            <label for="repo-description">描述</label>
            <textarea 
              id="repo-description" 
              v-model="createForm.description" 
              placeholder="输入仓库描述"
              rows="3"
            ></textarea>
          </div>
          <div class="form-group">
            <label>可见性</label>
            <div class="radio-group">
              <label>
                <input type="radio" v-model="createForm.visibility" value="public">
                公开 (Public)
              </label>
              <label>
                <input type="radio" v-model="createForm.visibility" value="private">
                私有 (Private)
              </label>
            </div>
          </div>
          <div class="form-group checkbox-group">
            <label>
              <input type="checkbox" v-model="createForm.auto_init">
              使用 README 初始化仓库
            </label>
          </div>
          <div class="modal-actions">
            <button type="button" class="btn secondary" @click="showCreateModal = false">取消</button>
            <button type="submit" class="btn primary" :disabled="createLoading">
              {{ createLoading ? '创建中...' : '创建仓库' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRepoStore } from '../stores'
import { repoApi } from '../services/api'

const repoStore = useRepoStore()

const filter = ref('all')
const repos = ref<any[]>([])
const isLoading = ref(false)
const error = ref('')
const showCreateModal = ref(false)
const createLoading = ref(false)

const createForm = ref({
  name: '',
  description: '',
  visibility: 'private',
  auto_init: false
})

const loadRepos = async () => {
  try {
    isLoading.value = true
    error.value = ''
    
    const response = await repoApi.listRepos() as { items: any[] }
    repos.value = response.items || []
    repoStore.setRepos(repos.value)
  } catch (err: any) {
    error.value = err.message || '加载仓库失败'
  } finally {
    isLoading.value = false
  }
}

const handleCreateRepo = async () => {
  try {
    createLoading.value = true
    
    await repoApi.createRepo({
      name: createForm.value.name,
      description: createForm.value.description,
      visibility: createForm.value.visibility,
      auto_init: createForm.value.auto_init
    })
    
    showCreateModal.value = false
    createForm.value = {
      name: '',
      description: '',
      visibility: 'private',
      auto_init: false
    }
    
    // 重新加载仓库列表
    await loadRepos()
  } catch (err: any) {
    alert(err.message || '创建仓库失败')
  } finally {
    createLoading.value = false
  }
}

onMounted(() => {
  loadRepos()
})
</script>

<style scoped>
.repo-list {
  padding: 20px;
}

.repo-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.create-repo-btn {
  padding: 8px 16px;
  background: var(--accent-color);
  color: #fff;
  border: none;
  border-radius: 4px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.3s;
}

.create-repo-btn:hover {
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

.loading {
  text-align: center;
  padding: 40px;
  color: var(--text-secondary);
}

.repo-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.repo-card {
  background: var(--bg-secondary);
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s, background-color 0.3s, box-shadow 0.3s;
}

.repo-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.repo-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 10px;
}

.repo-name {
  margin: 0;
  font-size: 18px;
  color: var(--accent-color);
  transition: color 0.3s;
}

.repo-visibility {
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}

.repo-visibility.public {
  background: var(--success-color);
  color: #fff;
}

.repo-visibility.private {
  background: var(--danger-color);
  color: #fff;
}

.repo-description {
  margin: 10px 0;
  color: var(--text-secondary);
  font-size: 14px;
  transition: color 0.3s;
}

.repo-stats {
  display: flex;
  gap: 16px;
  margin: 10px 0;
  font-size: 14px;
  color: var(--text-secondary);
  transition: color 0.3s;
}

.repo-actions {
  display: flex;
  gap: 10px;
  margin-top: 15px;
}

.btn {
  padding: 6px 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--bg-primary);
  color: var(--text-primary);
  text-decoration: none;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.btn:hover {
  background: var(--bg-secondary);
}

.btn.secondary {
  background: var(--bg-secondary);
}

.btn.primary {
  background: var(--accent-color);
  color: #fff;
  border-color: var(--accent-color);
}

/* 模态框样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal {
  background: var(--bg-secondary);
  border-radius: 8px;
  padding: 30px;
  width: 100%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal h2 {
  margin-top: 0;
  margin-bottom: 20px;
  color: var(--text-primary);
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: var(--text-primary);
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 14px;
  background: var(--bg-primary);
  color: var(--text-primary);
  box-sizing: border-box;
}

.form-group textarea {
  resize: vertical;
}

.radio-group {
  display: flex;
  gap: 20px;
}

.radio-group label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: normal;
}

.checkbox-group label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: normal;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
}

.error-message {
  margin: 20px 0;
  padding: 12px;
  background: var(--danger-color);
  color: #fff;
  border-radius: 4px;
  text-align: center;
}
</style>