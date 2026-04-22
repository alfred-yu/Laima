<template>
  <div class="app">
    <div class="sidebar" :class="{ 'sidebar-collapsed': !isSidebarOpened }">
      <div class="sidebar-header">
        <h1 class="logo">Laima</h1>
      </div>
      <nav class="sidebar-nav">
        <router-link to="/" class="nav-item">
          <span class="nav-icon">📊</span>
          <span class="nav-text">仪表盘</span>
        </router-link>
        <router-link to="/repos" class="nav-item">
          <span class="nav-icon">📁</span>
          <span class="nav-text">仓库</span>
        </router-link>
        <router-link to="/pulls" class="nav-item">
          <span class="nav-icon">🔄</span>
          <span class="nav-text">PR</span>
        </router-link>
        <router-link to="/issues" class="nav-item">
          <span class="nav-icon">📝</span>
          <span class="nav-text">Issue</span>
        </router-link>
        <router-link to="/cicd" class="nav-item">
          <span class="nav-icon">⚙️</span>
          <span class="nav-text">CI/CD</span>
        </router-link>
        <router-link to="/users" class="nav-item">
          <span class="nav-icon">👥</span>
          <span class="nav-text">用户</span>
        </router-link>
        <router-link to="/settings" class="nav-item">
          <span class="nav-icon">⚙️</span>
          <span class="nav-text">设置</span>
        </router-link>
      </nav>
    </div>
    <div class="main-content">
      <header class="top-nav">
        <button class="menu-toggle" @click="toggleSidebar">
          ☰
        </button>
        <div class="search-box">
          <input type="text" placeholder="搜索..." />
        </div>
        <div class="top-actions">
          <ThemeToggle />
          <div class="user-menu">
            <button class="user-avatar" v-if="isLoggedIn">
              {{ currentUser?.name?.charAt(0) || 'U' }}
            </button>
            <router-link to="/login" v-else class="login-link">登录</router-link>
          </div>
        </div>
      </header>
      <main class="content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useGlobalStore, useUserStore } from './stores'
import { ThemeToggle } from './components'

const globalStore = useGlobalStore()
const userStore = useUserStore()

const isSidebarOpened = computed(() => globalStore.isSidebarOpened)
const isLoggedIn = computed(() => userStore.isLoggedIn)
const currentUser = computed(() => userStore.currentUser)

const toggleSidebar = () => {
  globalStore.toggleSidebar()
}

// 初始化主题
onMounted(() => {
  const currentTheme = globalStore.currentTheme
  document.documentElement.setAttribute('data-theme', currentTheme)
})
</script>

<style>
/* 全局样式重置 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

:root {
  --bg-primary: #f5f5f5;
  --bg-secondary: #fff;
  --text-primary: #333;
  --text-secondary: #666;
  --border-color: #ddd;
  --sidebar-bg: #2d3748;
  --sidebar-text: #e2e8f0;
  --sidebar-hover: #4a5568;
  --accent-color: #0366d6;
  --success-color: #28a745;
  --danger-color: #dc3545;
  --warning-color: #ffc107;
  --info-color: #17a2b8;
}

/* 暗色主题变量 */
[data-theme="dark"] {
  --bg-primary: #1a1a1a;
  --bg-secondary: #2d2d2d;
  --text-primary: #e2e8f0;
  --text-secondary: #a0aec0;
  --border-color: #4a5568;
  --sidebar-bg: #1a202c;
  --sidebar-text: #e2e8f0;
  --sidebar-hover: #2d3748;
  --accent-color: #4299e1;
  --success-color: #48bb78;
  --danger-color: #f56565;
  --warning-color: #ed8936;
  --info-color: #38b2ac;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  font-size: 14px;
  line-height: 1.5;
  color: var(--text-primary);
  background-color: var(--bg-primary);
  transition: background-color 0.3s, color 0.3s;
}

/* 应用布局 */
.app {
  display: flex;
  min-height: 100vh;
}

/* 侧边栏 */
.sidebar {
  width: 240px;
  background: var(--sidebar-bg);
  color: var(--sidebar-text);
  transition: width 0.3s ease;
  height: 100vh;
  position: fixed;
  overflow-y: auto;
}

.sidebar.sidebar-collapsed {
  width: 64px;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid var(--border-color);
}

.logo {
  font-size: 20px;
  font-weight: 700;
  color: var(--sidebar-text);
}

.sidebar-nav {
  padding: 20px 0;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 12px 20px;
  color: var(--sidebar-text);
  text-decoration: none;
  transition: background-color 0.3s;
}

.nav-item:hover {
  background-color: var(--sidebar-hover);
}

.nav-icon {
  margin-right: 12px;
  font-size: 18px;
}

.sidebar.sidebar-collapsed .nav-text {
  display: none;
}

.sidebar.sidebar-collapsed .nav-icon {
  margin-right: 0;
}

/* 主内容区域 */
.main-content {
  flex: 1;
  margin-left: 240px;
  transition: margin-left 0.3s ease;
  min-height: 100vh;
}

.sidebar.sidebar-collapsed + .main-content {
  margin-left: 64px;
}

/* 顶部导航栏 */
.top-nav {
  display: flex;
  align-items: center;
  padding: 0 20px;
  height: 64px;
  background: var(--bg-secondary);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
  transition: background-color 0.3s;
}

.menu-toggle {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  margin-right: 20px;
  color: var(--text-primary);
}

.search-box {
  flex: 1;
  max-width: 400px;
}

.search-box input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 14px;
  background: var(--bg-primary);
  color: var(--text-primary);
  transition: all 0.3s;
}

.top-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-menu {
  display: flex;
  align-items: center;
}

.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--accent-color);
  color: #fff;
  border: none;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
}

.login-link {
  color: var(--accent-color);
  text-decoration: none;
  font-weight: 600;
}

/* 内容区域 */
.content {
  padding: 20px;
  min-height: calc(100vh - 64px);
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    transform: translateX(-100%);
  }
  
  .sidebar.sidebar-opened {
    transform: translateX(0);
  }
  
  .main-content {
    margin-left: 0;
  }
}
</style>
