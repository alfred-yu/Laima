<template>
  <div class="app-layout">
    <!-- ============ SIDEBAR ============ -->
    <aside class="sidebar" :class="{ 'sidebar-collapsed': !isSidebarOpened }">
      <!-- Logo -->
      <div class="sidebar-header">
        <div class="logo-container">
          <div class="logo-icon">
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 2L4 6V12C4 15.31 6.84 18.17 10.5 18.92V22H13.5V18.92C17.16 18.17 20 15.31 20 12V6L12 2Z" fill="var(--color-primary)" />
              <path d="M12 4.5L18 7.5V12C18 14.8 15.31 17.08 12.5 17.67V19.5H11.5V17.67C8.69 17.08 6 14.8 6 12V7.5L12 4.5Z" fill="white" />
              <circle cx="12" cy="11" r="2.5" fill="var(--color-primary)" />
            </svg>
          </div>
          <span class="logo-text" v-show="isSidebarOpened">Laima</span>
        </div>
      </div>

      <!-- Navigation -->
      <nav class="sidebar-nav">
        <div class="nav-section">
          <div class="nav-label" v-show="isSidebarOpened">
            <small>核心功能</small>
          </div>
          
          <router-link to="/" class="nav-item" active-class="nav-item-active">
            <div class="nav-icon-wrapper">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="3" width="7" height="7" rx="1" />
                <rect x="14" y="3" width="7" height="7" rx="1" />
                <rect x="3" y="14" width="7" height="7" rx="1" />
                <rect x="14" y="14" width="7" height="7" rx="1" />
              </svg>
            </div>
            <span class="nav-text">仪表盘</span>
          </router-link>

          <router-link to="/repos" class="nav-item" active-class="nav-item-active">
            <div class="nav-icon-wrapper">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z" />
                <polyline points="13 2 13 9 20 9" />
              </svg>
            </div>
            <span class="nav-text">仓库</span>
          </router-link>

          <router-link to="/pulls" class="nav-item" active-class="nav-item-active">
            <div class="nav-icon-wrapper">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="18" cy="18" r="3" />
                <circle cx="6" cy="6" r="3" />
                <path d="M13 6h3a2 2 0 0 1 2 2v7" />
                <line x1="6" y1="9" x2="6" y2="21" />
              </svg>
            </div>
            <span class="nav-text">合并请求</span>
          </router-link>

          <router-link to="/issues" class="nav-item" active-class="nav-item-active">
            <div class="nav-icon-wrapper">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10" />
                <line x1="12" y1="8" x2="12" y2="12" />
                <line x1="12" y1="16" x2="12.01" y2="16" />
              </svg>
            </div>
            <span class="nav-text">议题</span>
          </router-link>
        </div>

        <div class="nav-section">
          <div class="nav-label" v-show="isSidebarOpened">
            <small>开发工具</small>
          </div>

          <router-link to="/cicd" class="nav-item" active-class="nav-item-active">
            <div class="nav-icon-wrapper">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10" />
                <polyline points="12 6 12 12 16 14" />
              </svg>
            </div>
            <span class="nav-text">CI/CD</span>
          </router-link>

          <router-link to="/activity" class="nav-item" active-class="nav-item-active">
            <div class="nav-icon-wrapper">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="22 12 18 12 15 21 9 3 6 12 2 12" />
              </svg>
            </div>
            <span class="nav-text">活动记录</span>
          </router-link>

          <router-link to="/settings" class="nav-item" active-class="nav-item-active">
            <div class="nav-icon-wrapper">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="3" />
                <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09a1.65 1.65 0 0 0-1-1.51 1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1-1.51V5a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82.33z" />
              </svg>
            </div>
            <span class="nav-text">设置</span>
          </router-link>
        </div>
      </nav>

      <!-- Sidebar Footer - Collapse Toggle -->
      <div class="sidebar-footer">
        <button class="collapse-btn" @click="toggleSidebar" :title="isSidebarOpened ? '收起侧边栏' : '展开侧边栏'">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="collapse-icon" :style="{ transform: isSidebarOpened ? 'rotate(0deg)' : 'rotate(180deg)' }">
            <path d="M15 18l-6-6 6-6" />
          </svg>
          <span v-show="isSidebarOpened">收起</span>
        </button>
      </div>
    </aside>

    <!-- ============ MAIN CONTENT ============ -->
    <div class="main-container" :class="{ 'sidebar-collapsed': !isSidebarOpened }">
      <!-- Top Header -->
      <header class="top-header">
        <div class="header-left">
          <button class="mobile-menu-btn hide-desktop" @click="toggleSidebar">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="3" y1="12" x2="21" y2="12" />
              <line x1="3" y1="6" x2="21" y2="6" />
              <line x1="3" y1="18" x2="21" y2="18" />
            </svg>
          </button>
          
          <!-- Search -->
          <div class="search-container">
            <svg class="search-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8" />
              <line x1="21" y1="21" x2="16.65" y2="16.65" />
            </svg>
            <input type="text" placeholder="搜索仓库、代码、用户..." class="search-input" />
          </div>
        </div>

        <div class="header-right">
          <!-- Theme Toggle -->
          <div class="header-action">
            <ThemeToggle />
          </div>

          <!-- Notifications -->
          <button class="header-icon-btn" title="通知">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" />
              <path d="M13.73 21a2 2 0 0 1-3.46 0" />
            </svg>
            <span class="notification-badge">3</span>
          </button>

          <!-- User Menu -->
          <div class="user-menu">
            <button class="user-avatar-btn" v-if="isLoggedIn">
              <div class="user-avatar">
                {{ currentUser?.name?.charAt(0) || 'U' }}
              </div>
              <span class="user-name hide-mobile">{{ currentUser?.name || '用户' }}</span>
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="dropdown-arrow">
                <polyline points="6 9 12 15 18 9" />
              </svg>
            </button>
            <router-link to="/login" v-else class="login-btn">
              登录
            </router-link>
          </div>
        </div>
      </header>

      <!-- Page Content -->
      <main class="page-content">
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

// Initialize theme
onMounted(() => {
  const currentTheme = globalStore.currentTheme
  document.documentElement.setAttribute('data-theme', currentTheme)
})
</script>

<style scoped>
/* ============================================
   APP LAYOUT
   ============================================ */
.app-layout {
  display: flex;
  min-height: 100vh;
  background: var(--bg-secondary);
}

/* ============================================
   SIDEBAR
   ============================================ */
.sidebar {
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  width: 260px;
  background: var(--bg-primary);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  z-index: 400;
  transition: width 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.sidebar-collapsed {
  width: 72px;
}

/* Sidebar Header */
.sidebar-header {
  padding: 20px 16px;
  border-bottom: 1px solid var(--border-color-light);
}

.logo-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-icon {
  width: 36px;
  height: 36px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-text {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: -0.025em;
  white-space: nowrap;
}

/* Sidebar Navigation */
.sidebar-nav {
  flex: 1;
  padding: 16px 8px;
  overflow-y: auto;
}

.nav-section {
  margin-bottom: 24px;
}

.nav-label {
  padding: 8px 12px 4px;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--text-muted);
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  margin: 2px 0;
  border-radius: 8px;
  color: var(--text-secondary);
  text-decoration: none;
  transition: all 0.15s ease-out;
}

.nav-item:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
  text-decoration: none;
}

.nav-item-active {
  background: var(--color-primary-bg);
  color: var(--color-primary);
  font-weight: 500;
}

.nav-item-active:hover {
  background: rgba(79, 70, 229, 0.12);
}

.nav-icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  flex-shrink: 0;
}

.nav-text {
  white-space: nowrap;
  font-size: 0.9375rem;
}

.sidebar-collapsed .nav-text,
.sidebar-collapsed .nav-label {
  display: none;
}

/* Sidebar Footer */
.sidebar-footer {
  padding: 12px;
  border-top: 1px solid var(--border-color-light);
}

.collapse-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 10px 12px;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.9375rem;
  transition: all 0.15s ease-out;
}

.collapse-btn:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.collapse-icon {
  flex-shrink: 0;
  transition: transform 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

/* ============================================
   MAIN CONTAINER
   ============================================ */
.main-container {
  flex: 1;
  margin-left: 260px;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  transition: margin-left 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.main-container.sidebar-collapsed {
  margin-left: 72px;
}

/* ============================================
   TOP HEADER
   ============================================ */
.top-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
  padding: 0 24px;
  background: var(--bg-primary);
  border-bottom: 1px solid var(--border-color);
  position: sticky;
  top: 0;
  z-index: 300;
}

.header-left,
.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

/* Mobile Menu */
.mobile-menu-btn {
  display: none;
  background: transparent;
  border: none;
  padding: 8px;
  border-radius: 6px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.mobile-menu-btn:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

/* Search */
.search-container {
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
  width: 320px;
  padding: 10px 12px 10px 40px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 0.9375rem;
  background: var(--bg-secondary);
  color: var(--text-primary);
  transition: all 0.15s ease-out;
}

.search-input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px var(--color-primary-bg);
  background: var(--bg-primary);
}

.search-input::placeholder {
  color: var(--text-muted);
}

/* Header Actions */
.header-icon-btn {
  position: relative;
  width: 38px;
  height: 38px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.header-icon-btn:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.notification-badge {
  position: absolute;
  top: 6px;
  right: 6px;
  width: 16px;
  height: 16px;
  background: var(--color-danger);
  color: white;
  font-size: 0.6875rem;
  font-weight: 600;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 0 0 2px var(--bg-primary);
}

/* User Menu */
.user-avatar-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 12px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--bg-primary);
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.user-avatar-btn:hover {
  background: var(--bg-secondary);
  border-color: var(--color-gray-300);
}

.user-avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--color-primary), #7c3aed);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 0.875rem;
}

.user-name {
  font-size: 0.9375rem;
  color: var(--text-primary);
  font-weight: 500;
}

.dropdown-arrow {
  color: var(--text-muted);
}

.login-btn {
  padding: 8px 16px;
  background: var(--color-primary);
  color: white;
  text-decoration: none;
  border-radius: 6px;
  font-weight: 500;
  transition: all 0.15s ease-out;
}

.login-btn:hover {
  background: var(--color-primary-hover);
  text-decoration: none;
}

/* ============================================
   PAGE CONTENT
   ============================================ */
.page-content {
  flex: 1;
  padding: 24px;
  min-height: calc(100vh - 64px);
}

/* ============================================
   TRANSITIONS
   ============================================ */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease-out;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* ============================================
   RESPONSIVE DESIGN
   ============================================ */
@media (max-width: 1024px) {
  .search-input {
    width: 240px;
  }
}

@media (max-width: 768px) {
  .mobile-menu-btn {
    display: flex;
  }
  
  .sidebar {
    transform: translateX(-100%);
    box-shadow: var(--shadow-lg);
  }
  
  .sidebar:not(.sidebar-collapsed) {
    transform: translateX(0);
  }
  
  .sidebar-collapsed {
    width: 260px;
  }
  
  .main-container {
    margin-left: 0;
  }
  
  .main-container.sidebar-collapsed {
    margin-left: 0;
  }
  
  .search-input {
    width: 200px;
  }
  
  .page-content {
    padding: 16px;
  }
}

@media (max-width: 640px) {
  .search-container {
    display: none;
  }
  
  .top-header {
    padding: 0 16px;
  }
}
</style>
