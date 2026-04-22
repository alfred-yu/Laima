<template>
  <div class="repo-list">
    <h1>仓库列表</h1>
    <div class="repo-actions">
      <button class="create-repo-btn">创建仓库</button>
      <div class="filter-options">
        <select v-model="filter">
          <option value="all">所有仓库</option>
          <option value="owned">我的仓库</option>
          <option value="starred">已星标</option>
        </select>
      </div>
    </div>
    <div class="repo-grid">
      <div v-for="repo in repos" :key="repo.id" class="repo-card">
        <div class="repo-header">
          <h2 class="repo-name">{{ repo.name }}</h2>
          <span class="repo-visibility" :class="repo.visibility">{{ repo.visibility }}</span>
        </div>
        <p class="repo-description">{{ repo.description }}</p>
        <div class="repo-stats">
          <span class="stat">{{ repo.stars }} ⭐</span>
          <span class="stat">{{ repo.forks }} 🍴</span>
          <span class="stat">{{ repo.issues }} 📝</span>
        </div>
        <div class="repo-actions">
          <router-link :to="`/repos/${repo.owner}/${repo.name}/code`" class="btn">查看</router-link>
          <button class="btn secondary">克隆</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const filter = ref('all')
const repos = ref([
  {
    id: 1,
    name: 'laima',
    owner: 'user1',
    description: 'AI 原生代码托管平台',
    visibility: 'public',
    stars: 100,
    forks: 20,
    issues: 5
  },
  {
    id: 2,
    name: 'frontend',
    owner: 'user1',
    description: '前端项目',
    visibility: 'private',
    stars: 50,
    forks: 10,
    issues: 2
  },
  {
    id: 3,
    name: 'backend',
    owner: 'user1',
    description: '后端项目',
    visibility: 'public',
    stars: 80,
    forks: 15,
    issues: 3
  }
])

onMounted(() => {
  // 实际项目中这里会从 API 获取数据
  console.log('RepoList mounted')
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
  background: #0366d6;
  color: #fff;
  border: none;
  border-radius: 4px;
  font-weight: 600;
  cursor: pointer;
}

.filter-options select {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.repo-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.repo-card {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s;
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
  color: #0366d6;
}

.repo-visibility {
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}

.repo-visibility.public {
  background: #d1e7dd;
  color: #0f5132;
}

.repo-visibility.private {
  background: #f8d7da;
  color: #842029;
}

.repo-description {
  margin: 10px 0;
  color: #666;
  font-size: 14px;
}

.repo-stats {
  display: flex;
  gap: 16px;
  margin: 10px 0;
  font-size: 14px;
  color: #666;
}

.repo-actions {
  display: flex;
  gap: 10px;
  margin-top: 15px;
}

.btn {
  padding: 6px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: #fff;
  color: #333;
  text-decoration: none;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.btn:hover {
  background: #f5f5f5;
}

.btn.secondary {
  background: #f5f5f5;
}
</style>