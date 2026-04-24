<template>
  <button class="theme-toggle" @click="toggleTheme" title="切换主题">
    <span v-if="currentTheme === 'light'">🌙</span>
    <span v-else>☀️</span>
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useGlobalStore } from '../stores'

const globalStore = useGlobalStore()
const currentTheme = computed(() => globalStore.currentTheme)

const toggleTheme = () => {
  const newTheme = currentTheme.value === 'light' ? 'dark' : 'light'
  globalStore.setTheme(newTheme)
  document.documentElement.setAttribute('data-theme', newTheme)
}
</script>

<style scoped>
.theme-toggle {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  padding: 8px;
  border-radius: 50%;
  transition: background-color 0.3s;
}

.theme-toggle:hover {
  background-color: rgba(0, 0, 0, 0.1);
}
</style>