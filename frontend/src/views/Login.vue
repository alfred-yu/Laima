<template>
  <div class="login-page">
    <div class="login-container">
      <!-- Brand Logo -->
      <div class="login-brand">
        <div class="brand-logo">
          <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" class="logo-icon">
            <path d="M12 2L4 6V12C4 15.31 6.84 18.17 10.5 18.92V22H13.5V18.92C17.16 18.17 20 15.31 20 12V6L12 2Z" fill="var(--color-primary)" />
            <path d="M12 4.5L18 7.5V12C18 14.8 15.31 17.08 12.5 17.67V19.5H11.5V17.67C8.69 17.08 6 14.8 6 12V7.5L12 4.5Z" fill="white" />
            <circle cx="12" cy="11" r="2.5" fill="var(--color-primary)" />
          </svg>
        </div>
        <h1 class="brand-title">Laima</h1>
        <p class="brand-subtitle">专业的代码托管平台</p>
      </div>

      <!-- Login Form Card -->
      <div class="login-form-card">
        <div class="form-header">
          <h2 class="form-title">登录</h2>
          <p class="form-subtitle">欢迎回来，继续您的开发工作</p>
        </div>
        
        <form @submit.prevent="handleLogin" class="form-content">
          <div class="form-group">
            <label for="username" class="form-label">用户名或邮箱</label>
            <div class="input-wrapper">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="input-icon">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
                <circle cx="12" cy="7" r="4" />
              </svg>
              <input 
                type="text" 
                id="username" 
                v-model="form.username" 
                required 
                placeholder="请输入用户名或邮箱"
                class="form-input"
                :class="{ 'input-error': error && form.username }"
              >
            </div>
          </div>
          
          <div class="form-group">
            <div class="label-row">
              <label for="password" class="form-label">密码</label>
              <a href="#" class="forgot-link">忘记密码？</a>
            </div>
            <div class="input-wrapper">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="input-icon">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
                <path d="M7 11V7a5 5 0 0 1 10 0v4" />
              </svg>
              <input 
                type="password" 
                id="password" 
                v-model="form.password" 
                required 
                placeholder="请输入密码"
                class="form-input"
                :class="{ 'input-error': error && form.password }"
              >
            </div>
          </div>
          
          <div class="form-group checkbox-group">
            <input type="checkbox" id="remember" v-model="form.remember" class="form-checkbox">
            <label for="remember" class="checkbox-label">记住我</label>
          </div>
          
          <div v-if="error" class="error-message">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10" />
              <line x1="12" y1="8" x2="12" y2="12" />
              <line x1="12" y1="16" x2="12.01" y2="16" />
            </svg>
            <span>{{ error }}</span>
          </div>
          
          <button 
            type="submit" 
            class="login-button"
            :class="{ 'is-loading': isLoading }"
            :disabled="isLoading"
          >
            <svg v-if="isLoading" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="loading-spinner">
              <circle cx="12" cy="12" r="10" stroke-dasharray="32" stroke-dashoffset="32" />
            </svg>
            {{ isLoading ? '登录中...' : '登录' }}
          </button>
          
          <div class="divider">
            <span>或</span>
          </div>
          
          <div class="social-login">
            <button type="button" class="social-button github">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
              </svg>
              <span>使用 GitHub 登录</span>
            </button>
          </div>
        </form>
        
        <div class="form-footer">
          <p class="register-link">
            还没有账号？
            <router-link to="/register" class="register-button">立即注册</router-link>
          </p>
        </div>
      </div>

      <!-- Footer Links -->
      <div class="login-footer">
        <div class="footer-links">
          <a href="#" class="footer-link">关于我们</a>
          <a href="#" class="footer-link">隐私政策</a>
          <a href="#" class="footer-link">服务条款</a>
          <a href="#" class="footer-link">帮助中心</a>
        </div>
        <div class="footer-copyright">
          © 2026 Laima. 保留所有权利
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores'
import { authApi } from '../services/api'

const router = useRouter()
const userStore = useUserStore()

const form = ref({
  username: '',
  password: '',
  remember: false
})

const isLoading = ref(false)
const error = ref('')

const handleLogin = async () => {
  try {
    isLoading.value = true
    error.value = ''
    
    const response = await authApi.login(form.value.username, form.value.password) as { token: string; user: any }
    
    userStore.setToken(response.token)
    userStore.setUser(response.user)
    
    // 登录成功后跳转到用户页面或仪表盘
    router.push('/dashboard')
  } catch (err: any) {
    error.value = err.message || '登录失败，请检查用户名和密码'
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: var(--bg-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}

.login-container {
  width: 100%;
  max-width: 480px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 32px;
}

/* Brand Logo */
.login-brand {
  text-align: center;
  margin-bottom: 16px;
}

.brand-logo {
  margin-bottom: 16px;
}

.logo-icon {
  width: 48px;
  height: 48px;
  margin: 0 auto;
}

.brand-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 8px 0;
  letter-spacing: -0.025em;
}

.brand-subtitle {
  font-size: 1rem;
  color: var(--text-secondary);
  margin: 0;
  line-height: 1.5;
}

/* Login Form Card */
.login-form-card {
  width: 100%;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 32px;
  box-shadow: var(--shadow-md);
}

.form-header {
  margin-bottom: 24px;
}

.form-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 8px 0;
}

.form-subtitle {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin: 0;
}

.form-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.label-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.form-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
}

.forgot-link {
  font-size: 0.875rem;
  color: var(--color-primary);
  text-decoration: none;
  font-weight: 500;
  transition: color 0.15s ease-out;
}

.forgot-link:hover {
  color: var(--color-primary-hover);
  text-decoration: underline;
}

.input-wrapper {
  position: relative;
}

.input-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-muted);
  pointer-events: none;
}

.form-input {
  width: 100%;
  padding: 12px 12px 12px 40px;
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

.input-error {
  border-color: var(--color-danger);
}

.checkbox-group {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-direction: row;
}

.form-checkbox {
  width: 16px;
  height: 16px;
  accent-color: var(--color-primary);
}

.checkbox-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
  cursor: pointer;
}

.error-message {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px;
  background: var(--color-danger-bg);
  color: var(--color-danger);
  border-radius: 8px;
  font-size: 0.875rem;
  border: 1px solid rgba(220, 38, 38, 0.2);
}

.login-button {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px;
  background: var(--color-primary);
  color: white;
  border: 1px solid var(--color-primary);
  border-radius: 8px;
  font-size: 0.9375rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.login-button:hover:not(:disabled) {
  background: var(--color-primary-hover);
  border-color: var(--color-primary-hover);
  transform: translateY(-1px);
}

.login-button:disabled {
  background: var(--bg-tertiary);
  border-color: var(--border-color);
  color: var(--text-muted);
  cursor: not-allowed;
  transform: none;
}

.is-loading {
  opacity: 0.8;
}

.loading-spinner {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loading-spinner circle {
  animation: dash 1.5s ease-in-out infinite;
}

@keyframes dash {
  0% { stroke-dashoffset: 32; }
  50% { stroke-dashoffset: 8; }
  100% { stroke-dashoffset: 32; }
}

.divider {
  position: relative;
  display: flex;
  align-items: center;
  margin: 20px 0;
}

.divider::before {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--border-color);
}

.divider span {
  padding: 0 16px;
  font-size: 0.8125rem;
  color: var(--text-muted);
  background: var(--bg-primary);
}

.divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--border-color);
}

.social-login {
  display: flex;
  gap: 12px;
}

.social-button {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.social-button:hover {
  background: var(--bg-secondary);
  border-color: var(--color-gray-300);
  transform: translateY(-1px);
}

.social-button.github {
  border-color: var(--color-gray-700);
  color: var(--color-gray-700);
  background: var(--bg-primary);
}

.social-button.github:hover {
  background: var(--color-gray-50);
  border-color: var(--color-gray-600);
}

.form-footer {
  margin-top: 24px;
  text-align: center;
}

.register-link {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin: 0;
}

.register-button {
  color: var(--color-primary);
  font-weight: 600;
  text-decoration: none;
  transition: color 0.15s ease-out;
}

.register-button:hover {
  color: var(--color-primary-hover);
  text-decoration: underline;
}

/* Footer Links */
.login-footer {
  width: 100%;
  text-align: center;
  margin-top: 16px;
}

.footer-links {
  display: flex;
  justify-content: center;
  gap: 24px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.footer-link {
  font-size: 0.75rem;
  color: var(--text-secondary);
  text-decoration: none;
  transition: color 0.15s ease-out;
}

.footer-link:hover {
  color: var(--color-primary);
  text-decoration: underline;
}

.footer-copyright {
  font-size: 0.75rem;
  color: var(--text-muted);
}

/* Responsive Design */
@media (max-width: 768px) {
  .login-page {
    padding: 20px;
  }
  
  .login-container {
    gap: 24px;
  }
  
  .login-form-card {
    padding: 24px;
  }
  
  .brand-title {
    font-size: 1.5rem;
  }
  
  .form-title {
    font-size: 1.125rem;
  }
  
  .social-login {
    flex-direction: column;
  }
  
  .footer-links {
    gap: 16px;
  }
}

@media (max-width: 480px) {
  .login-page {
    padding: 16px;
  }
  
  .login-form-card {
    padding: 20px;
  }
  
  .brand-title {
    font-size: 1.25rem;
  }
  
  .form-content {
    gap: 12px;
  }
  
  .footer-links {
    gap: 12px;
  }
}
</style>
