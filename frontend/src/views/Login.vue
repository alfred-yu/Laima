<template>
  <div class="login-page">
    <div class="login-container">
      <!-- Brand Logo -->
      <div class="brand-section">
        <div class="brand-logo">
          <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" class="logo-icon">
            <path d="M12 2L4 6V12C4 15.31 6.84 18.17 10.5 18.92V22H13.5V18.92C17.16 18.17 20 15.31 20 12V6L12 2Z" fill="#fc6d26" />
            <path d="M12 4.5L18 7.5V12C18 14.8 15.31 17.08 12.5 17.67V19.5H11.5V17.67C8.69 17.08 6 14.8 6 12V7.5L12 4.5Z" fill="white" />
            <circle cx="12" cy="11" r="2.5" fill="#fc6d26" />
          </svg>
        </div>
        <h1 class="brand-title">登录到Laima</h1>
      </div>

      <!-- Login Form -->
      <div class="login-form">
        <form @submit.prevent="handleLogin('user')" class="form-content">
          <!-- Username Field -->
          <div class="form-group">
            <label for="username" class="form-label">用户名或主要电子邮件</label>
            <input 
              type="text" 
              id="username" 
              v-model="form.username" 
              required 
              class="form-input"
              :class="{ 'input-error': error && form.username }"
            >
          </div>
          
          <!-- Password Field -->
          <div class="form-group">
            <div class="label-row">
              <label for="password" class="form-label">密码</label>
              <router-link to="/forgot-password" class="forgot-link">忘记密码？</router-link>
            </div>
            <div class="password-input-wrapper">
              <input 
                type="password" 
                id="password" 
                v-model="form.password" 
                required 
                class="form-input"
                :class="{ 'input-error': error && form.password }"
              >
              <button type="button" class="password-toggle">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
                  <circle cx="12" cy="12" r="3" />
                </svg>
              </button>
            </div>
          </div>
          
          <!-- Remember Checkbox -->
          <div class="form-group checkbox-group">
            <input type="checkbox" id="remember" v-model="form.remember" class="form-checkbox">
            <label for="remember" class="checkbox-label">记住账号</label>
          </div>
          
          <!-- Error Message -->
          <div v-if="error" class="error-message">
            <span>{{ error }}</span>
          </div>
          
          <!-- Login Button -->
          <button 
            type="submit" 
            class="login-button"
            :class="{ 'is-loading': isLoading }"
            :disabled="isLoading"
          >
            {{ isLoading ? '登录中...' : '登录' }}
          </button>
          
          <!-- Test Buttons for Role Testing -->
          <div class="test-buttons">
            <button 
              type="button" 
              class="test-button user-test"
              @click="handleLogin('user')"
              :disabled="isLoading"
            >
              测试：普通用户登录
            </button>
            <button 
              type="button" 
              class="test-button admin-test"
              @click="handleLogin('admin')"
              :disabled="isLoading"
            >
              测试：管理员登录
            </button>
          </div>
          
          <!-- Passkey Button -->
          <button type="button" class="passkey-button">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor" class="passkey-icon">
              <path d="M21 10c0-4.4-3.6-8-8-8s-8 3.6-8 8 3.6 8 8 8c1.2 0 2.4-.3 3.5-.8l2.7 2.7c.4.4 1 .4 1.4 0l1.4-1.4c.4-.4.4-1 0-1.4L17 14.5c.5-1.1.8-2.3.8-3.5zm-8 6c-3.3 0-6-2.7-6-6s2.7-6 6-6 6 2.7 6 6-2.7 6-6 6z" />
            </svg>
            通行密钥
          </button>
          
          <!-- Footer Links -->
          <div class="form-footer">
            <p class="terms-text">
              登录即表示您接受使用条款并接受隐私政策和 Cookie 政策。
            </p>
            <p class="register-link">
              还没有账户？
              <router-link to="/register" class="register-button">立即注册</router-link>
            </p>
          </div>
        </form>
        
        <!-- Divider -->
        <div class="divider">
          <span>或使用以下账户登录</span>
        </div>
        
        <!-- Social Login Buttons -->
        <div class="social-login">
          <button type="button" class="social-button">
            <svg width="18" height="18" viewBox="0 0 24 24" class="social-icon google">
              <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
              <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
              <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
              <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
            </svg>
            Google
          </button>
          
          <button type="button" class="social-button">
            <svg width="18" height="18" viewBox="0 0 24 24" class="social-icon github" fill="currentColor">
              <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
            </svg>
            GitHub
          </button>
          
          <button type="button" class="social-button">
            <svg width="18" height="18" viewBox="0 0 24 24" class="social-icon bitbucket" fill="currentColor">
              <path d="M2 0c-.55 0-1 .45-1 1.03l1.71 21.94c.07.55.53.98 1.08.98h15.44c.55 0 1.01-.43 1.08-.98L23 1.03c0-.58-.45-1.03-1-1H2zm12.33 16.01h-3.04l-.91-5.5h4.88l-.93 5.5zm1.71-8.5H7.96l-.37-2h8.82l-.37 2z"/>
            </svg>
            Bitbucket
          </button>
          
          <button type="button" class="social-button">
            <svg width="18" height="18" viewBox="0 0 24 24" class="social-icon salesforce" fill="currentColor">
              <path d="M10.095 18.173a2.068 2.068 0 0 1-2.067 2.068 2.068 2.068 0 0 1-2.067-2.068 2.068 2.068 0 0 1 2.067-2.067 2.068 2.068 0 0 1 2.067 2.067m.689-8.59a2.756 2.756 0 0 1 2.755 2.756 2.756 2.756 0 0 1-2.755 2.755 2.756 2.756 0 0 1-2.756-2.755 2.756 2.756 0 0 1 2.756-2.756m2.631-4.582a3.444 3.444 0 0 1 3.443 3.443 3.444 3.444 0 0 1-3.443 3.443 3.444 3.444 0 0 1-3.443-3.443 3.444 3.444 0 0 1 3.443-3.443m2.522 15.155a2.068 2.068 0 0 1-2.067 2.068 2.068 2.068 0 0 1-2.067-2.068 2.068 2.068 0 0 1 2.067-2.067 2.068 2.068 0 0 1 2.067 2.067m2.589-10.572a2.068 2.068 0 0 1-2.067 2.068 2.068 2.068 0 0 1-2.067-2.068 2.068 2.068 0 0 1 2.067-2.067 2.068 2.068 0 0 1 2.067 2.067M8.395 7.574a2.396 2.396 0 0 1-2.395 2.395 2.396 2.396 0 0 1-2.395-2.395 2.396 2.396 0 0 1 2.395-2.395 2.396 2.396 0 0 1 2.395 2.395M5.764 15.155a2.068 2.068 0 0 1-2.067 2.068A2.068 2.068 0 0 1 1.63 15.155a2.068 2.068 0 0 1 2.067-2.067 2.068 2.068 0 0 1 2.067 2.067"/>
            </svg>
            Salesforce
          </button>
        </div>
        
        <!-- Remember Account -->
        <div class="remember-bottom">
          <input type="checkbox" id="remember-bottom" v-model="form.remember" class="form-checkbox">
          <label for="remember-bottom" class="checkbox-label">记住账号</label>
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

const handleLogin = async (role = 'user') => {
  try {
    isLoading.value = true
    error.value = ''
    
    // 模拟登录请求
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 模拟登录成功
    const response = {
      token: 'mock-token-123',
      user: {
        id: 1,
        username: form.value.username || 'testuser',
        email: `${form.value.username || 'testuser'}@example.com`,
        name: role === 'admin' ? '管理员用户' : '测试用户'
      }
    }
    
    userStore.setToken(response.token)
    userStore.setUser(response.user)
    userStore.setRole(role)
    
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
  display: flex;
  align-items: center;
  justify-content: center;
  background: #ffffff;
  padding: 48px 16px;
}

.login-container {
  width: 100%;
  max-width: 400px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

/* Brand Section */
.brand-section {
  text-align: center;
  margin-bottom: 32px;
}

.brand-logo {
  margin-bottom: 16px;
}

.logo-icon {
  width: 64px;
  height: 64px;
}

.brand-title {
  font-size: 1.25rem;
  font-weight: 700;
  color: #333;
  margin: 0;
}

/* Login Form */
.login-form {
  width: 100%;
}

.form-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-label {
  font-size: 0.875rem;
  font-weight: 600;
  color: #333;
}

.label-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.forgot-link {
  font-size: 0.75rem;
  color: #1f78d1;
  text-decoration: none;
}

.forgot-link:hover {
  text-decoration: underline;
}

.form-input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #bfbfbf;
  border-radius: 6px;
  font-size: 0.9375rem;
  background: #fff;
  color: #333;
  transition: border-color 0.1s ease;
}

.form-input:focus {
  outline: none;
  border-color: #1f78d1;
  box-shadow: 0 0 0 3px rgba(31, 120, 209, 0.15);
}

.input-error {
  border-color: #d93025;
}

.password-input-wrapper {
  position: relative;
}

.password-toggle {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  padding: 4px;
  cursor: pointer;
  color: #666;
  display: flex;
  align-items: center;
  justify-content: center;
}

.checkbox-group {
  flex-direction: row;
  align-items: center;
  gap: 8px;
  margin-top: 4px;
}

.form-checkbox {
  width: 16px;
  height: 16px;
  margin: 0;
  accent-color: #1f78d1;
}

.checkbox-label {
  font-size: 0.875rem;
  color: #666;
  cursor: pointer;
  margin: 0;
}

.error-message {
  padding: 10px 12px;
  background: #fef0f0;
  color: #d93025;
  border: 1px solid #fde2e2;
  border-radius: 6px;
  font-size: 0.875rem;
  text-align: center;
}

.login-button {
  width: 100%;
  padding: 12px 16px;
  background: #1f78d1;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 0.9375rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.1s ease;
  margin-top: 4px;
}

.login-button:hover:not(:disabled) {
  background: #1a5fa5;
}

.login-button:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.is-loading {
  opacity: 0.8;
}

.passkey-button {
  width: 100%;
  padding: 10px 16px;
  background: white;
  color: #333;
  border: 1px solid #bfbfbf;
  border-radius: 6px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.1s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.passkey-button:hover {
  background: #f5f5f5;
  border-color: #999;
}

.passkey-icon {
  color: #333;
}

/* Test Buttons */
.test-buttons {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 12px;
}

.test-button {
  width: 100%;
  padding: 10px 16px;
  border: 1px solid #e5e5e5;
  border-radius: 6px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.1s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.test-button.user-test {
  background: #f0f8ff;
  color: #1f78d1;
  border-color: #d0e3ff;
}

.test-button.user-test:hover:not(:disabled) {
  background: #e0f0ff;
  border-color: #b0d0ff;
}

.test-button.admin-test {
  background: #f8f0ff;
  color: #6f42c1;
  border-color: #e6d0ff;
}

.test-button.admin-test:hover:not(:disabled) {
  background: #f0e0ff;
  border-color: #d6b0ff;
}

.test-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.form-footer {
  text-align: center;
  margin-top: 8px;
}

.terms-text {
  font-size: 0.75rem;
  color: #666;
  margin: 0 0 12px 0;
  line-height: 1.4;
}

.register-link {
  font-size: 0.875rem;
  color: #666;
  margin: 0;
}

.register-button {
  color: #1f78d1;
  font-weight: 500;
  text-decoration: none;
}

.register-button:hover {
  text-decoration: underline;
}

.divider {
  position: relative;
  display: flex;
  align-items: center;
  margin: 24px 0;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: #e5e5e5;
}

.divider span {
  padding: 0 16px;
  font-size: 0.75rem;
  color: #666;
  background: white;
}

.social-login {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.social-button {
  width: 100%;
  padding: 10px 16px;
  background: white;
  color: #333;
  border: 1px solid #bfbfbf;
  border-radius: 6px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.1s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
}

.social-button:hover {
  background: #f5f5f5;
  border-color: #999;
}

.social-icon {
  flex-shrink: 0;
}

.remember-bottom {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 16px;
  justify-content: center;
}

/* Responsive Design */
@media (max-width: 480px) {
  .login-page {
    padding: 32px 16px;
  }
  
  .login-container {
    max-width: 100%;
  }
}
</style>
