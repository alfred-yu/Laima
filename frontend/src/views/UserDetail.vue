<template>
  <div class="user-detail">
    <div class="user-header">
      <div class="user-avatar">
        {{ user.name.charAt(0) }}
      </div>
      <div class="user-info">
        <h1>{{ user.name }}</h1>
        <p class="user-email">{{ user.email }}</p>
        <p class="user-role">{{ user.role }}</p>
        <div class="user-stats">
          <span class="stat">{{ user.repos }} 仓库</span>
          <span class="stat">{{ user.pulls }} PR</span>
          <span class="stat">{{ user.issues }} Issue</span>
        </div>
      </div>
    </div>
    <div class="user-content">
      <div class="user-repos">
        <h2>用户仓库</h2>
        <div class="repos-list">
          <div v-for="repo in user.reposList" :key="repo.id" class="repo-item">
            <h3 class="repo-name">{{ repo.name }}</h3>
            <p class="repo-description">{{ repo.description }}</p>
            <div class="repo-stats">
              <span class="stat">{{ repo.stars }} ⭐</span>
              <span class="stat">{{ repo.forks }} 🍴</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const user = ref({
  id: 1,
  name: 'user1',
  email: 'user1@example.com',
  role: '管理员',
  repos: 5,
  pulls: 12,
  issues: 8,
  reposList: [
    {
      id: 1,
      name: 'laima',
      description: 'AI 原生代码托管平台',
      stars: 100,
      forks: 20
    },
    {
      id: 2,
      name: 'frontend',
      description: '前端项目',
      stars: 50,
      forks: 10
    },
    {
      id: 3,
      name: 'backend',
      description: '后端项目',
      stars: 80,
      forks: 15
    }
  ]
})
</script>

<style scoped>
.user-detail {
  padding: 20px;
}

.user-header {
  display: flex;
  align-items: flex-start;
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 1px solid #ddd;
}

.user-avatar {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  background: #0366d6;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48px;
  font-weight: 600;
  margin-right: 32px;
}

.user-info h1 {
  margin: 0 0 8px 0;
  font-size: 24px;
  color: #333;
}

.user-email {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: #666;
}

.user-role {
  margin: 0 0 16px 0;
  font-size: 14px;
  color: #999;
}

.user-stats {
  display: flex;
  gap: 24px;
}

.stat {
  font-size: 14px;
  color: #666;
}

.user-content {
  margin-top: 24px;
}

.user-repos h2 {
  margin: 0 0 20px 0;
  font-size: 18px;
  color: #333;
}

.repos-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.repo-item {
  padding: 16px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s;
}

.repo-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.repo-name {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: #0366d6;
}

.repo-description {
  margin: 0 0 12px 0;
  font-size: 14px;
  color: #666;
}

.repo-stats {
  display: flex;
  gap: 16px;
  font-size: 14px;
  color: #666;
}
</style>