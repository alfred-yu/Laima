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
  color: var(--text-primary);
  transition: color 0.3s;
}

.repo-actions {
  display: flex;
  gap: 10px;
}

.btn {
  padding: 6px 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--bg-primary);
  color: var(--text-primary);
  cursor: pointer;
  transition: all 0.3s;
}

.btn:hover {
  background: var(--bg-secondary);
}

.btn.primary {
  background: var(--accent-color);
  color: #fff;
  border-color: var(--accent-color);
}

.btn.primary:hover {
  opacity: 0.9;
}

.repo-tabs {
  display: flex;
  border-bottom: 1px solid var(--border-color);
  margin-bottom: 20px;
  transition: border-color 0.3s;
}

.tab {
  padding: 10px 20px;
  text-decoration: none;
  color: var(--text-secondary);
  border-bottom: 2px solid transparent;
  transition: all 0.3s;
}

.tab:hover {
  color: var(--text-primary);
  background: var(--bg-primary);
}

.tab.router-link-active {
  color: var(--accent-color);
  border-bottom-color: var(--accent-color);
  font-weight: 600;
}

.repo-content {
  min-height: 400px;
}
</style>