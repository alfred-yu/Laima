<template>
  <div class="repo-detail">
    <div class="repo-header">
      <h1>{{ repoName }}</h1>
      <div class="repo-actions">
        <button class="btn">星标</button>
        <button class="btn">Fork</button>
        <button class="btn primary">克隆</button>
      </div>
    </div>
    <div class="repo-tabs">
      <router-link :to="`/repos/${owner}/${repo}/code`" class="tab">代码</router-link>
      <router-link :to="`/repos/${owner}/${repo}/pulls`" class="tab">PR</router-link>
      <router-link :to="`/repos/${owner}/${repo}/issues`" class="tab">Issue</router-link>
      <router-link :to="`/repos/${owner}/${repo}/cicd`" class="tab">CI/CD</router-link>
    </div>
    <div class="repo-content">
      <router-view />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const owner = computed(() => route.params.owner as string)
const repo = computed(() => route.params.repo as string)
const repoName = computed(() => `${owner.value}/${repo.value}`)
</script>

<style scoped>
.repo-detail {
  padding: 20px;
}

.repo-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.repo-header h1 {
  margin: 0;
  font-size: 24px;
  color: #333;
}

.repo-actions {
  display: flex;
  gap: 10px;
}

.btn {
  padding: 6px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: #fff;
  color: #333;
  cursor: pointer;
  transition: all 0.3s;
}

.btn:hover {
  background: #f5f5f5;
}

.btn.primary {
  background: #0366d6;
  color: #fff;
  border-color: #0366d6;
}

.btn.primary:hover {
  background: #0056b3;
}

.repo-tabs {
  display: flex;
  border-bottom: 1px solid #ddd;
  margin-bottom: 20px;
}

.tab {
  padding: 10px 20px;
  text-decoration: none;
  color: #666;
  border-bottom: 2px solid transparent;
  transition: all 0.3s;
}

.tab:hover {
  color: #333;
  background: #f5f5f5;
}

.tab.router-link-active {
  color: #0366d6;
  border-bottom-color: #0366d6;
  font-weight: 600;
}

.repo-content {
  min-height: 400px;
}
</style>