<template>
  <div class="settings">
    <h1>设置</h1>
    <div class="settings-tabs">
      <button class="tab" :class="{ active: activeTab === 'profile' }" @click="activeTab = 'profile'">个人资料</button>
      <button class="tab" :class="{ active: activeTab === 'security' }" @click="activeTab = 'security'">安全</button>
      <button class="tab" :class="{ active: activeTab === 'notifications' }" @click="activeTab = 'notifications'">通知</button>
      <button class="tab" :class="{ active: activeTab === 'integrations' }" @click="activeTab = 'integrations'">集成</button>
    </div>
    <div class="settings-content">
      <div v-if="activeTab === 'profile'" class="profile-settings">
        <h2>个人资料设置</h2>
        <form class="settings-form">
          <div class="form-group">
            <label for="name">用户名</label>
            <input type="text" id="name" v-model="profile.name" placeholder="请输入用户名">
          </div>
          <div class="form-group">
            <label for="email">邮箱</label>
            <input type="email" id="email" v-model="profile.email" placeholder="请输入邮箱">
          </div>
          <div class="form-group">
            <label for="bio">个人简介</label>
            <textarea id="bio" v-model="profile.bio" placeholder="请输入个人简介" rows="3"></textarea>
          </div>
          <button type="submit" class="btn primary">保存</button>
        </form>
      </div>
      <div v-if="activeTab === 'security'" class="security-settings">
        <h2>安全设置</h2>
        <form class="settings-form">
          <div class="form-group">
            <label for="currentPassword">当前密码</label>
            <input type="password" id="currentPassword" v-model="security.currentPassword" placeholder="请输入当前密码">
          </div>
          <div class="form-group">
            <label for="newPassword">新密码</label>
            <input type="password" id="newPassword" v-model="security.newPassword" placeholder="请输入新密码">
          </div>
          <div class="form-group">
            <label for="confirmPassword">确认密码</label>
            <input type="password" id="confirmPassword" v-model="security.confirmPassword" placeholder="请确认新密码">
          </div>
          <button type="submit" class="btn primary">更改密码</button>
        </form>
      </div>
      <div v-if="activeTab === 'notifications'" class="notifications-settings">
        <h2>通知设置</h2>
        <form class="settings-form">
          <div class="form-group">
            <label>
              <input type="checkbox" v-model="notifications.email">
              邮件通知
            </label>
          </div>
          <div class="form-group">
            <label>
              <input type="checkbox" v-model="notifications.web">
              网页通知
            </label>
          </div>
          <div class="form-group">
            <label>
              <input type="checkbox" v-model="notifications.prs">
              PR 通知
            </label>
          </div>
          <div class="form-group">
            <label>
              <input type="checkbox" v-model="notifications.issues">
              Issue 通知
            </label>
          </div>
          <button type="submit" class="btn primary">保存</button>
        </form>
      </div>
      <div v-if="activeTab === 'integrations'" class="integrations-settings">
        <h2>集成设置</h2>
        <div class="integration-list">
          <div class="integration-item">
            <h3>GitHub</h3>
            <button class="btn">连接</button>
          </div>
          <div class="integration-item">
            <h3>GitLab</h3>
            <button class="btn">连接</button>
          </div>
          <div class="integration-item">
            <h3>Gitea</h3>
            <button class="btn">连接</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const activeTab = ref('profile')

const profile = ref({
  name: 'user1',
  email: 'user1@example.com',
  bio: 'AI 原生代码托管平台开发者'
})

const security = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const notifications = ref({
  email: true,
  web: true,
  prs: true,
  issues: true
})
</script>

<style scoped>
.settings {
  padding: 20px;
}

.settings-tabs {
  display: flex;
  border-bottom: 1px solid #ddd;
  margin: 20px 0;
}

.tab {
  padding: 12px 24px;
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  font-size: 14px;
  font-weight: 600;
  color: #666;
  cursor: pointer;
  transition: all 0.3s;
}

.tab:hover {
  color: #333;
}

.tab.active {
  color: #0366d6;
  border-bottom-color: #0366d6;
}

.settings-content {
  margin-top: 32px;
}

.settings-form {
  max-width: 600px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #333;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  transition: border-color 0.3s;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #0366d6;
  box-shadow: 0 0 0 3px rgba(3, 102, 214, 0.1);
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

.integration-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  max-width: 600px;
}

.integration-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.integration-item h3 {
  margin: 0;
  font-size: 16px;
  color: #333;
}
</style>