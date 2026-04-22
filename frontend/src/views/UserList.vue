<template>
  <div class="user-list">
    <h1>用户列表</h1>
    <div class="users-grid">
      <div v-for="user in users" :key="user.id" class="user-card">
        <div class="user-avatar">
          {{ user.name.charAt(0) }}
        </div>
        <div class="user-info">
          <h3 class="user-name">{{ user.name }}</h3>
          <p class="user-email">{{ user.email }}</p>
          <p class="user-role">{{ user.role }}</p>
        </div>
        <div class="user-actions">
          <router-link :to="`/users/${user.id}`" class="btn">查看</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const users = ref([
  {
    id: 1,
    name: 'user1',
    email: 'user1@example.com',
    role: '管理员'
  },
  {
    id: 2,
    name: 'user2',
    email: 'user2@example.com',
    role: '开发者'
  },
  {
    id: 3,
    name: 'user3',
    email: 'user3@example.com',
    role: '开发者'
  }
])
</script>

<style scoped>
.user-list {
  padding: 20px;
}

.users-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

.user-card {
  display: flex;
  align-items: center;
  padding: 20px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s;
}

.user-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.user-avatar {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: #0366d6;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: 600;
  margin-right: 20px;
}

.user-info {
  flex: 1;
}

.user-name {
  margin: 0 0 4px 0;
  font-size: 16px;
  color: #333;
}

.user-email {
  margin: 0 0 4px 0;
  font-size: 14px;
  color: #666;
}

.user-role {
  margin: 0;
  font-size: 12px;
  color: #999;
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
</style>