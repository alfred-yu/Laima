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
    
    // 实际项目中这里会调用 API
    // 模拟登录成功
    setTimeout(() => {
      userStore.setToken('mock-token-123')
      userStore.setUser({
        id: 1,
        username: form.value.username,
        email: `${form.value.username}@example.com`,
        name: form.value.username
      })
      router.push('/')
    }, 1000)
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
  background: #f5f5f5;
}

.login-form {
  background: #fff;
  border-radius: 8px;
  padding: 40px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  width: 100%;
  max-width: 400px;
}

.login-form h1 {
  margin-top: 0;
  margin-bottom: 30px;
  text-align: center;
  color: #333;
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

.form-group input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  transition: border-color 0.3s;
}

.form-group input:focus {
  outline: none;
  border-color: #0366d6;
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
  color: #666;
}

.login-button {
  width: 100%;
  padding: 12px;
  background: #0366d6;
  color: #fff;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.3s;
}

.login-button:hover {
  background: #0056b3;
}

.login-button:disabled {
  background: #6c757d;
  cursor: not-allowed;
}

.error-message {
  margin-top: 16px;
  padding: 10px;
  background: #f8d7da;
  color: #842029;
  border-radius: 4px;
  font-size: 14px;
}

.register-link {
  margin-top: 20px;
  text-align: center;
  font-size: 14px;
  color: #666;
}

.register-link a {
  color: #0366d6;
  text-decoration: none;
}

.register-link a:hover {
  text-decoration: underline;
}
</style>