<template>
  <div class="login-container">
    <div class="login-form">
      <h1>登录</h1>
      <form @submit.prevent="handleLogin">
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
          <label for="password">密码</label>
          <input 
            type="password" 
            id="password" 
            v-model="form.password" 
            required 
            placeholder="请输入密码"
          >
        </div>
        <div class="form-group remember">
          <input type="checkbox" id="remember" v-model="form.remember">
          <label for="remember">记住我</label>
        </div>
        <button type="submit" class="login-button" :disabled="isLoading">
          {{ isLoading ? '登录中...' : '登录' }}
        </button>
        <div class="error-message" v-if="error">
          {{ error }}
        </div>
        <div class="register-link">
          还没有账号？<router-link to="/register">立即注册</router-link>
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
    
    router.push('/')
  } catch (err: any) {
    error.value = err.message || '登录失败，请检查用户名和密码'
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: var(--bg-primary);
  transition: background-color 0.3s;
}

.login-form {
  background: var(--bg-secondary);
  border-radius: 8px;
  padding: 40px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  width: 100%;
  max-width: 400px;
  transition: background-color 0.3s, box-shadow 0.3s;
}

.login-form h1 {
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

.form-group.remember {
  display: flex;
  align-items: center;
  margin-bottom: 24px;
}

.form-group.remember input {
  width: auto;
  margin-right: 8px;
}

.form-group.remember label {
  margin-bottom: 0;
  font-weight: normal;
  color: var(--text-secondary);
  transition: color 0.3s;
}

.login-button {
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
}

.login-button:hover {
  background: var(--accent-color);
  opacity: 0.9;
}

.login-button:disabled {
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

.register-link {
  margin-top: 20px;
  text-align: center;
  font-size: 14px;
  color: var(--text-secondary);
  transition: color 0.3s;
}

.register-link a {
  color: var(--accent-color);
  text-decoration: none;
  transition: color 0.3s;
}

.register-link a:hover {
  text-decoration: underline;
}
</style>