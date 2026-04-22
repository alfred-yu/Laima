<template>
  <div class="register-container">
    <div class="register-form">
      <h1>注册</h1>
      <form @submit.prevent="handleRegister">
        <div class="form-group">
          <label for="username">用户名</label>
          <input 
            type="text" 
            id="username" 
            v-model="form.username" 
            required 
            placeholder="请输入用户名"
          >
        </div>
        <div class="form-group">
          <label for="email">邮箱</label>
          <input 
            type="email" 
            id="email" 
            v-model="form.email" 
            required 
            placeholder="请输入邮箱"
          >
        </div>
        <div class="form-group">
          <label for="password">密码</label>
          <input 
            type="password" 
            id="password" 
            v-model="form.password" 
            required 
            placeholder="请输入密码"
          >
        </div>
        <div class="form-group">
          <label for="confirmPassword">确认密码</label>
          <input 
            type="password" 
            id="confirmPassword" 
            v-model="form.confirmPassword" 
            required 
            placeholder="请确认密码"
          >
        </div>
        <button type="submit" class="register-button" :disabled="isLoading">
          {{ isLoading ? '注册中...' : '注册' }}
        </button>
        <div class="error-message" v-if="error">
          {{ error }}
        </div>
        <div class="login-link">
          已有账号？<router-link to="/login">立即登录</router-link>
        </div>
      </form>
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
  email: '',
  password: '',
  confirmPassword: ''
})

const isLoading = ref(false)
const error = ref('')

const handleRegister = async () => {
  try {
    isLoading.value = true
    error.value = ''
    
    // 验证密码是否一致
    if (form.value.password !== form.value.confirmPassword) {
      error.value = '两次输入的密码不一致'
      return
    }
    
    const response = await authApi.register(form.value.username, form.value.email, form.value.password)
    
    userStore.setToken(response.token)
    userStore.setUser(response.user)
    
    router.push('/')
  } catch (err: any) {
    error.value = err.message || '注册失败，请稍后重试'
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: var(--bg-primary);
  transition: background-color 0.3s;
}

.register-form {
  background: var(--bg-secondary);
  border-radius: 8px;
  padding: 40px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  width: 100%;
  max-width: 400px;
  transition: background-color 0.3s, box-shadow 0.3s;
}

.register-form h1 {
  margin-top: 0;
  margin-bottom: 30px;
  text-align: center;
  color: var(--text-primary);
  transition: color 0.3s;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: var(--text-primary);
  transition: color 0.3s;
}

.form-group input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 14px;
  background: var(--bg-primary);
  color: var(--text-primary);
  transition: border-color 0.3s, background-color 0.3s, color 0.3s;
}

.form-group input:focus {
  outline: none;
  border-color: var(--accent-color);
  box-shadow: 0 0 0 3px rgba(3, 102, 214, 0.1);
}

.register-button {
  width: 100%;
  padding: 12px;
  background: var(--accent-color);
  color: #fff;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.3s;
  margin-top: 10px;
}

.register-button:hover {
  background: var(--accent-color);
  opacity: 0.9;
}

.register-button:disabled {
  background: var(--border-color);
  cursor: not-allowed;
}

.error-message {
  margin-top: 16px;
  padding: 10px;
  background: var(--danger-color);
  color: #fff;
  border-radius: 4px;
  font-size: 14px;
  transition: background-color 0.3s, color 0.3s;
}

.login-link {
  margin-top: 20px;
  text-align: center;
  font-size: 14px;
  color: var(--text-secondary);
  transition: color 0.3s;
}

.login-link a {
  color: var(--accent-color);
  text-decoration: none;
  transition: color 0.3s;
}

.login-link a:hover {
  text-decoration: underline;
}
</style>