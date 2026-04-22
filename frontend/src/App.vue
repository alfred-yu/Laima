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
        <div class="user-menu">
          <button class="user-avatar" v-if="isLoggedIn">
            {{ currentUser?.name?.charAt(0) || 'U' }}
          </button>
          <router-link to="/login" v-else class="login-link">登录</router-link>
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
import { computed } from 'vue'
import { useGlobalStore, useUserStore } from './stores'

const globalStore = useGlobalStore()
const userStore = useUserStore()

const isSidebarOpened = computed(() => globalStore.isSidebarOpened)
const isLoggedIn = computed(() => userStore.isLoggedIn)
const currentUser = computed(() => userStore.currentUser)

const toggleSidebar = () => {
  globalStore.toggleSidebar()
}
</script>

<style>
/* 全局样式重置 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  font-size: 14px;
  line-height: 1.5;
  color: #333;
  background-color: #f5f5f5;
}

/* 应用布局 */
.app {
  display: flex;
  min-height: 100vh;
}

/* 侧边栏 */
.sidebar {
  width: 240px;
  background: #2d3748;
  color: #fff;
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
  border-bottom: 1px solid #4a5568;
}

.logo {
  font-size: 20px;
  font-weight: 700;
  color: #fff;
}

.sidebar-nav {
  padding: 20px 0;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 12px 20px;
  color: #e2e8f0;
  text-decoration: none;
  transition: background-color 0.3s;
}

.nav-item:hover {
  background-color: #4a5568;
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
  background: #fff;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

.menu-toggle {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  margin-right: 20px;
  color: #333;
}

.search-box {
  flex: 1;
  max-width: 400px;
}

.search-box input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.user-menu {
  display: flex;
  align-items: center;
}

.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: #0366d6;
  color: #fff;
  border: none;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
}

.login-link {
  color: #0366d6;
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
