<template>
  <div class="forgot-password-page">
    <div class="forgot-password-container">
      <!-- Brand Logo -->
      <div class="forgot-password-brand">
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

      <!-- Forgot Password Form Card -->
      <div class="forgot-password-card">
        <div class="form-header">
          <h2 class="form-title">忘记密码</h2>
          <p class="form-subtitle">输入您的邮箱地址，我们将发送重置密码的链接</p>
        </div>
        
        <form @submit.prevent="handleForgotPassword" class="form-content">
          <div class="form-group">
            <label for="email" class="form-label">邮箱地址</label>
            <div class="input-wrapper">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="input-icon">
                <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z" />
                <polyline points="22,6 12,13 2,6" />
              </svg>
              <input 
                type="email" 
                id="email" 
                v-model="form.email" 
                required 
                placeholder="请输入您的邮箱地址"
                class="form-input"
                :class="{ 'input-error': error && form.email }"
              >
            </div>
          </div>
          
          <div v-if="error" class="error-message">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10" />
              <line x1="12" y1="8" x2="12" y2="12" />
              <line x1="12" y1="16" x2="12.01" y2="16" />
            </svg>
            <span>{{ error }}</span>
          </div>
          
          <div v-if="success" class="success-message">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" />
              <polyline points="22,4 12,14.01 9,11.01" />
            </svg>
            <span>{{ success }}</span>
          </div>
          
          <button 
            type="submit" 
            class="forgot-password-button"
            :class="{ 'is-loading': isLoading }"
            :disabled="isLoading"
          >
            <svg v-if="isLoading" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="loading-spinner">
              <circle cx="12" cy="12" r="10" stroke-dasharray="32" stroke-dashoffset="32" />
            </svg>
            {{ isLoading ? '发送中...' : '发送重置链接' }}
          </button>
        </form>
        
        <div class="form-footer">
          <router-link to="/login" class="back-to-login">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="15,18 9,12 15,6" />
            </svg>
            返回登录
          </router-link>
        </div>
      </div>

      <!-- Footer Links -->
      <div class="forgot-password-footer">
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
import { authApi } from '../services/api'

const router = useRouter()

const form = ref({
  email: ''
})

const isLoading = ref(false)
const error = ref('')
const success = ref('')

const handleForgotPassword = async () => {
  try {
    isLoading.value = true
    error.value = ''
    success.value = ''
    
    // 模拟 API 调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 模拟成功响应
    success.value = '重置密码链接已发送到您的邮箱，请查收'
    
    // 3秒后跳转到登录页面
    setTimeout(() => {
      router.push('/login')
    }, 3000)
  } catch (err: any) {
    error.value = err.message || '发送失败，请稍后重试'
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.forgot-password-page {
  min-height: 100vh;
  background: var(--bg-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}

.forgot-password-container {
  width: 100%;
  max-width: 480px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 32px;
}

/* Brand Logo */
.forgot-password-brand {
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

/* Forgot Password Form Card */
.forgot-password-card {
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

.form-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
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

.success-message {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px;
  background: var(--color-success-bg);
  color: var(--color-success);
  border-radius: 8px;
  font-size: 0.875rem;
  border: 1px solid rgba(16, 185, 129, 0.2);
}

.forgot-password-button {
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

.forgot-password-button:hover:not(:disabled) {
  background: var(--color-primary-hover);
  border-color: var(--color-primary-hover);
  transform: translateY(-1px);
}

.forgot-password-button:disabled {
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

.form-footer {
  margin-top: 24px;
  text-align: center;
}

.back-to-login {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.875rem;
  color: var(--color-primary);
  text-decoration: none;
  font-weight: 500;
  transition: all 0.15s ease-out;
}

.back-to-login:hover {
  color: var(--color-primary-hover);
  text-decoration: underline;
}

/* Footer Links */
.forgot-password-footer {
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
  .forgot-password-page {
    padding: 20px;
  }
  
  .forgot-password-container {
    gap: 24px;
  }
  
  .forgot-password-card {
    padding: 24px;
  }
  
  .brand-title {
    font-size: 1.5rem;
  }
  
  .form-title {
    font-size: 1.125rem;
  }
  
  .footer-links {
    gap: 16px;
  }
}

@media (max-width: 480px) {
  .forgot-password-page {
    padding: 16px;
  }
  
  .forgot-password-card {
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
