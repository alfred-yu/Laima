<template>
  <div class="setting-page">
    <div class="page-header">
      <h1 class="page-title">个人设置</h1>
      <p class="page-subtitle">管理您的个人信息和账户设置</p>
    </div>

    <div class="setting-container">
      <!-- 侧边导航 -->
      <div class="setting-sidebar">
        <nav class="sidebar-nav">
          <a 
            v-for="item in navItems" 
            :key="item.id" 
            :href="`#${item.id}`"
            class="nav-item"
            :class="{ active: activeTab === item.id }"
            @click="activeTab = item.id"
          >
            <span class="nav-icon">{{ item.icon }}</span>
            <span class="nav-text">{{ item.title }}</span>
          </a>
        </nav>
      </div>

      <!-- 主内容区 -->
      <div class="setting-content">
        <!-- 个人资料 -->
        <section id="profile" class="setting-section" v-if="activeTab === 'profile'">
          <div class="section-header">
            <h2 class="section-title">个人资料</h2>
            <p class="section-description">更新您的个人信息和头像</p>
          </div>

          <div class="profile-form">
            <div class="avatar-section">
              <div class="avatar-container">
                <div class="avatar">
                  <img 
                    v-if="user.avatar" 
                    :src="user.avatar" 
                    :alt="user.name"
                    class="avatar-image"
                  >
                  <div v-else class="avatar-placeholder">
                    {{ user.name.charAt(0).toUpperCase() }}
                  </div>
                </div>
                <div class="avatar-actions">
                  <input 
                    type="file" 
                    id="avatar-upload" 
                    accept="image/*" 
                    class="avatar-input"
                    @change="handleAvatarUpload"
                  >
                  <label for="avatar-upload" class="btn btn-ghost btn-sm">
                    更改头像
                  </label>
                  <button 
                    v-if="user.avatar"
                    class="btn btn-ghost btn-sm text-danger"
                    @click="removeAvatar"
                  >
                    移除
                  </button>
                </div>
              </div>
            </div>

            <form @submit.prevent="updateProfile" class="form">
              <div class="form-grid">
                <div class="form-group">
                  <label for="name" class="form-label">用户名</label>
                  <input 
                    type="text" 
                    id="name" 
                    v-model="profileForm.name" 
                    class="form-input"
                    placeholder="输入用户名"
                  >
                </div>

                <div class="form-group">
                  <label for="email" class="form-label">邮箱</label>
                  <input 
                    type="email" 
                    id="email" 
                    v-model="profileForm.email" 
                    class="form-input"
                    placeholder="输入邮箱地址"
                  >
                </div>

                <div class="form-group">
                  <label for="bio" class="form-label">个人简介</label>
                  <textarea 
                    id="bio" 
                    v-model="profileForm.bio" 
                    class="form-input"
                    placeholder="介绍一下自己..."
                    rows="3"
                  ></textarea>
                </div>

                <div class="form-group">
                  <label for="location" class="form-label">位置</label>
                  <input 
                    type="text" 
                    id="location" 
                    v-model="profileForm.location" 
                    class="form-input"
                    placeholder="输入您的位置"
                  >
                </div>

                <div class="form-group">
                  <label for="website" class="form-label">个人网站</label>
                  <input 
                    type="url" 
                    id="website" 
                    v-model="profileForm.website" 
                    class="form-input"
                    placeholder="输入您的网站地址"
                  >
                </div>

                <div class="form-group">
                  <label for="company" class="form-label">公司</label>
                  <input 
                    type="text" 
                    id="company" 
                    v-model="profileForm.company" 
                    class="form-input"
                    placeholder="输入您的公司"
                  >
                </div>
              </div>

              <div class="form-actions">
                <button type="button" class="btn btn-ghost" @click="resetProfileForm">
                  重置
                </button>
                <button type="submit" class="btn btn-primary" :disabled="isUpdating">
                  {{ isUpdating ? '更新中...' : '保存更改' }}
                </button>
              </div>
            </form>
          </div>
        </section>

        <!-- 密码设置 -->
        <section id="password" class="setting-section" v-if="activeTab === 'password'">
          <div class="section-header">
            <h2 class="section-title">密码设置</h2>
            <p class="section-description">更新您的账户密码</p>
          </div>

          <form @submit.prevent="updatePassword" class="form">
            <div class="form-group">
              <label for="current-password" class="form-label">当前密码</label>
              <input 
                type="password" 
                id="current-password" 
                v-model="passwordForm.currentPassword" 
                class="form-input"
                placeholder="输入当前密码"
              >
            </div>

            <div class="form-group">
              <label for="new-password" class="form-label">新密码</label>
              <input 
                type="password" 
                id="new-password" 
                v-model="passwordForm.newPassword" 
                class="form-input"
                placeholder="输入新密码"
              >
              <p class="form-hint">密码长度至少为8个字符，包含字母和数字</p>
            </div>

            <div class="form-group">
              <label for="confirm-password" class="form-label">确认新密码</label>
              <input 
                type="password" 
                id="confirm-password" 
                v-model="passwordForm.confirmPassword" 
                class="form-input"
                placeholder="再次输入新密码"
              >
            </div>

            <div class="form-actions">
              <button type="button" class="btn btn-ghost" @click="resetPasswordForm">
                重置
              </button>
              <button type="submit" class="btn btn-primary" :disabled="isUpdatingPassword">
                {{ isUpdatingPassword ? '更新中...' : '更新密码' }}
              </button>
            </div>
          </form>
        </section>

        <!-- 通知设置 -->
        <section id="notifications" class="setting-section" v-if="activeTab === 'notifications'">
          <div class="section-header">
            <h2 class="section-title">通知设置</h2>
            <p class="section-description">管理您的通知偏好</p>
          </div>

          <div class="notification-settings">
            <div class="setting-item">
              <div class="setting-info">
                <h3 class="setting-title">电子邮件通知</h3>
                <p class="setting-description">接收重要事件的电子邮件通知</p>
              </div>
              <label class="switch">
                <input 
                  type="checkbox" 
                  v-model="notificationForm.email"
                >
                <span class="slider"></span>
              </label>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <h3 class="setting-title">站内通知</h3>
                <p class="setting-description">在网站内接收通知</p>
              </div>
              <label class="switch">
                <input 
                  type="checkbox" 
                  v-model="notificationForm.web"
                >
                <span class="slider"></span>
              </label>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <h3 class="setting-title">PR 通知</h3>
                <p class="setting-description">当有新的 PR 或 PR 更新时通知我</p>
              </div>
              <label class="switch">
                <input 
                  type="checkbox" 
                  v-model="notificationForm.prs"
                >
                <span class="slider"></span>
              </label>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <h3 class="setting-title">Issue 通知</h3>
                <p class="setting-description">当有新的 Issue 或 Issue 更新时通知我</p>
              </div>
              <label class="switch">
                <input 
                  type="checkbox" 
                  v-model="notificationForm.issues"
                >
                <span class="slider"></span>
              </label>
            </div>

            <div class="setting-item">
              <div class="setting-info">
                <h3 class="setting-title">评论通知</h3>
                <p class="setting-description">当有人评论我的代码时通知我</p>
              </div>
              <label class="switch">
                <input 
                  type="checkbox" 
                  v-model="notificationForm.comments"
                >
                <span class="slider"></span>
              </label>
            </div>

            <div class="form-actions">
              <button type="button" class="btn btn-ghost" @click="resetNotificationForm">
                恢复默认
              </button>
              <button type="button" class="btn btn-primary" @click="saveNotificationSettings">
                保存设置
              </button>
            </div>
          </div>
        </section>

        <!-- 账号安全 -->
        <section id="security" class="setting-section" v-if="activeTab === 'security'">
          <div class="section-header">
            <h2 class="section-title">账号安全</h2>
            <p class="section-description">管理您的账号安全设置</p>
          </div>

          <div class="security-settings">
            <div class="security-item">
              <div class="security-info">
                <h3 class="security-title">两步验证</h3>
                <p class="security-description">增加额外的安全层，保护您的账户</p>
              </div>
              <button class="btn btn-ghost btn-sm">
                启用
              </button>
            </div>

            <div class="security-item">
              <div class="security-info">
                <h3 class="security-title">SSH 密钥</h3>
                <p class="security-description">管理您的 SSH 密钥</p>
              </div>
              <button class="btn btn-ghost btn-sm">
                管理
              </button>
            </div>

            <div class="security-item">
              <div class="security-info">
                <h3 class="security-title">应用授权</h3>
                <p class="security-description">管理授权访问您账户的应用</p>
              </div>
              <button class="btn btn-ghost btn-sm">
                管理
              </button>
            </div>

            <div class="security-item">
              <div class="security-info">
                <h3 class="security-title">登录历史</h3>
                <p class="security-description">查看您的最近登录记录</p>
              </div>
              <button class="btn btn-ghost btn-sm">
                查看
              </button>
            </div>
          </div>
        </section>

        <!-- 开发者设置 -->
        <section id="developer" class="setting-section" v-if="activeTab === 'developer'">
          <div class="section-header">
            <h2 class="section-title">开发者设置</h2>
            <p class="section-description">管理您的 API 令牌和开发者设置</p>
          </div>

          <div class="developer-settings">
            <div class="api-tokens">
              <h3 class="token-title">API 令牌</h3>
              <p class="token-description">创建和管理用于 API 访问的令牌</p>
              
              <div class="token-list">
                <div v-if="apiTokens.length === 0" class="empty-state">
                  <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="2" y="2" width="20" height="20" rx="5" ry="5" />
                    <path d="M16 11.37A4 4 0 1 1 12.63 8 4 4 0 0 1 16 11.37z" />
                    <line x1="17.5" y1="6.5" x2="17.51" y2="6.5" />
                  </svg>
                  <p>您还没有创建 API 令牌</p>
                  <button class="btn btn-primary btn-sm">
                    创建令牌
                  </button>
                </div>
                
                <div v-else class="token-item">
                  <div class="token-info">
                    <h4 class="token-name">Personal Access Token</h4>
                    <p class="token-scopes">repo, user</p>
                    <p class="token-created">创建于 2026-04-25</p>
                  </div>
                  <div class="token-actions">
                    <button class="btn btn-ghost btn-sm">
                      编辑
                    </button>
                    <button class="btn btn-ghost btn-sm text-danger">
                      删除
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useUserStore } from '../stores'

const userStore = useUserStore()
const activeTab = ref('profile')
const isUpdating = ref(false)
const isUpdatingPassword = ref(false)

// 模拟用户数据
const user = reactive({
  id: 1,
  name: 'user1',
  email: 'user1@example.com',
  avatar: null,
  bio: 'Software developer',
  location: 'Beijing, China',
  website: 'https://example.com',
  company: 'Laima Inc.'
})

// 个人资料表单
const profileForm = reactive({
  name: user.name,
  email: user.email,
  bio: user.bio,
  location: user.location,
  website: user.website,
  company: user.company
})

// 密码表单
const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 通知设置表单
const notificationForm = reactive({
  email: true,
  web: true,
  prs: true,
  issues: true,
  comments: true
})

// API 令牌
const apiTokens = ref([])

// 导航项
const navItems = [
  { id: 'profile', title: '个人资料', icon: '👤' },
  { id: 'password', title: '密码设置', icon: '🔒' },
  { id: 'notifications', title: '通知设置', icon: '🔔' },
  { id: 'security', title: '账号安全', icon: '🛡️' },
  { id: 'developer', title: '开发者设置', icon: '⚙️' }
]

// 处理头像上传
const handleAvatarUpload = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    const file = target.files[0]
    // 这里应该上传文件到服务器，这里只是模拟
    user.avatar = URL.createObjectURL(file)
  }
}

// 移除头像
const removeAvatar = () => {
  user.avatar = null
}

// 更新个人资料
const updateProfile = async () => {
  try {
    isUpdating.value = true
    // 模拟 API 调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 更新用户数据
    Object.assign(user, profileForm)
    
    // 显示成功消息
    alert('个人资料更新成功！')
  } catch (error) {
    alert('更新失败，请重试')
  } finally {
    isUpdating.value = false
  }
}

// 重置个人资料表单
const resetProfileForm = () => {
  Object.assign(profileForm, {
    name: user.name,
    email: user.email,
    bio: user.bio,
    location: user.location,
    website: user.website,
    company: user.company
  })
}

// 更新密码
const updatePassword = async () => {
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    alert('两次输入的密码不一致')
    return
  }

  if (passwordForm.newPassword.length < 8) {
    alert('密码长度至少为8个字符')
    return
  }

  try {
    isUpdatingPassword.value = true
    // 模拟 API 调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 显示成功消息
    alert('密码更新成功！')
    resetPasswordForm()
  } catch (error) {
    alert('更新失败，请重试')
  } finally {
    isUpdatingPassword.value = false
  }
}

// 重置密码表单
const resetPasswordForm = () => {
  Object.assign(passwordForm, {
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
  })
}

// 保存通知设置
const saveNotificationSettings = () => {
  // 模拟 API 调用
  setTimeout(() => {
    alert('通知设置保存成功！')
  }, 500)
}

// 重置通知设置
const resetNotificationForm = () => {
  Object.assign(notificationForm, {
    email: true,
    web: true,
    prs: true,
    issues: true,
    comments: true
  })
}

onMounted(() => {
  // 加载用户数据
  // 这里应该从 API 获取用户数据
})
</script>

<style scoped>
.setting-page {
  padding: 24px;
  min-height: 100vh;
}

.page-header {
  margin-bottom: 32px;
}

.page-title {
  font-size: 2rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 8px 0;
  letter-spacing: -0.025em;
}

.page-subtitle {
  font-size: 1rem;
  color: var(--text-secondary);
  margin: 0;
}

.setting-container {
  display: flex;
  gap: 24px;
  max-width: 1200px;
  margin: 0 auto;
}

/* 侧边导航 */
.setting-sidebar {
  flex: 0 0 240px;
}

.sidebar-nav {
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 8px;
  box-shadow: var(--shadow-sm);
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 8px;
  color: var(--text-secondary);
  text-decoration: none;
  transition: all 0.15s ease-out;
  font-size: 0.9375rem;
  font-weight: 500;
}

.nav-item:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.nav-item.active {
  background: var(--color-primary-bg);
  color: var(--color-primary);
  font-weight: 600;
}

.nav-icon {
  font-size: 1.125rem;
  width: 20px;
  text-align: center;
}

/* 主内容区 */
.setting-content {
  flex: 1;
  min-width: 0;
}

.setting-section {
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 32px;
  box-shadow: var(--shadow-sm);
  margin-bottom: 24px;
  animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.section-header {
  margin-bottom: 24px;
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 8px 0;
}

.section-description {
  font-size: 0.9375rem;
  color: var(--text-secondary);
  margin: 0;
}

/* 个人资料表单 */
.profile-form {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.avatar-section {
  display: flex;
  justify-content: flex-start;
  align-items: center;
}

.avatar-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  margin-right: 40px;
}

.avatar {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  overflow: hidden;
  border: 3px solid var(--border-color);
  box-shadow: var(--shadow-md);
  transition: transform 0.2s ease-out;
}

.avatar:hover {
  transform: scale(1.05);
}

.avatar-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  width: 100%;
  height: 100%;
  background: var(--color-primary);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48px;
  font-weight: 600;
}

.avatar-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: center;
}

.avatar-input {
  display: none;
}

/* 表单样式 */
.form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
}

.form-input {
  padding: 12px 16px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 0.9375rem;
  background: var(--bg-primary);
  color: var(--text-primary);
  transition: all 0.15s ease-out;
}

.form-input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px var(--color-primary-bg);
}

.form-hint {
  font-size: 0.8125rem;
  color: var(--text-muted);
  margin: 4px 0 0 0;
}

.form-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 8px;
}

/* 按钮样式 */
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

.btn-primary:hover:not(:disabled) {
  background: var(--color-primary-hover);
  border-color: var(--color-primary-hover);
  transform: translateY(-1px);
}

.btn-primary:disabled {
  background: var(--bg-tertiary);
  border-color: var(--border-color);
  color: var(--text-muted);
  cursor: not-allowed;
  transform: none;
}

.btn-ghost {
  background: transparent;
  border-color: var(--border-color);
}

.btn-ghost:hover {
  background: var(--bg-secondary);
  border-color: var(--color-gray-300);
}

.btn-sm {
  padding: 6px 12px;
  font-size: 0.875rem;
}

.text-danger {
  color: var(--color-danger);
}

.text-danger:hover {
  background: var(--color-danger-bg);
  border-color: var(--color-danger);
}

/* 通知设置 */
.notification-settings {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.setting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: var(--bg-secondary);
  border-radius: 8px;
  transition: all 0.15s ease-out;
}

.setting-item:hover {
  background: var(--bg-tertiary);
}

.setting-info {
  flex: 1;
  min-width: 0;
}

.setting-title {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 4px 0;
}

.setting-description {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin: 0;
}

/* 开关样式 */
.switch {
  position: relative;
  display: inline-block;
  width: 50px;
  height: 24px;
}

.switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: var(--bg-tertiary);
  transition: 0.4s;
  border-radius: 24px;
}

.slider:before {
  position: absolute;
  content: "";
  height: 18px;
  width: 18px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  transition: 0.4s;
  border-radius: 50%;
}

input:checked + .slider {
  background-color: var(--color-primary);
}

input:checked + .slider:before {
  transform: translateX(26px);
}

/* 安全设置 */
.security-settings {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.security-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: var(--bg-secondary);
  border-radius: 8px;
  transition: all 0.15s ease-out;
}

.security-item:hover {
  background: var(--bg-tertiary);
}

.security-info {
  flex: 1;
  min-width: 0;
}

.security-title {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 4px 0;
}

.security-description {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin: 0;
}

/* 开发者设置 */
.developer-settings {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.api-tokens {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.token-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.token-description {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin: 0;
}

.token-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 24px;
  background: var(--bg-secondary);
  border-radius: 8px;
  text-align: center;
  color: var(--text-secondary);
  gap: 16px;
}

.token-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: var(--bg-secondary);
  border-radius: 8px;
  transition: all 0.15s ease-out;
}

.token-item:hover {
  background: var(--bg-tertiary);
}

.token-info {
  flex: 1;
  min-width: 0;
}

.token-name {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 4px 0;
}

.token-scopes {
  font-size: 0.8125rem;
  color: var(--text-secondary);
  margin: 0 0 4px 0;
}

.token-created {
  font-size: 0.75rem;
  color: var(--text-muted);
  margin: 0;
}

.token-actions {
  display: flex;
  gap: 8px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .setting-container {
    flex-direction: column;
  }
  
  .setting-sidebar {
    flex: none;
  }
  
  .sidebar-nav {
    display: flex;
    overflow-x: auto;
    padding: 8px 12px;
  }
  
  .nav-item {
    white-space: nowrap;
  }
  
  .form-grid {
    grid-template-columns: 1fr;
  }
  
  .avatar-section {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .avatar-container {
    margin-right: 0;
    margin-bottom: 24px;
  }
}

@media (max-width: 480px) {
  .setting-page {
    padding: 16px;
  }
  
  .setting-section {
    padding: 24px;
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  .btn {
    width: 100%;
    justify-content: center;
  }
}
</style>
