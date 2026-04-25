<template>
  <div class="commit-list">
    <div v-for="commit in commits" :key="commit.sha" class="commit-item" @click="$emit('select', commit)">
      <div class="commit-avatar">
        <img v-if="commit.authorAvatar" :src="commit.authorAvatar" :alt="commit.author" class="avatar-img" />
        <div v-else class="avatar-placeholder">{{ commit.author?.charAt(0)?.toUpperCase() }}</div>
      </div>
      <div class="commit-info">
        <div class="commit-message">{{ commit.message }}</div>
        <div class="commit-meta">
          <span class="commit-author">{{ commit.author }}</span>
          <span class="commit-time">{{ commit.time }}</span>
        </div>
      </div>
      <div class="commit-sha">
        <code>{{ commit.sha.slice(0, 7) }}</code>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
export interface Commit {
  sha: string
  message: string
  author: string
  authorAvatar?: string
  time: string
}

defineProps<{
  commits: Commit[]
}>()

defineEmits(['select'])
</script>

<style scoped>
.commit-list {
  border: 1px solid var(--border-color);
  border-radius: 6px;
  overflow: hidden;
}

.commit-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  cursor: pointer;
  transition: background 0.15s;
}

.commit-item + .commit-item {
  border-top: 1px solid var(--border-color);
}

.commit-item:hover {
  background: var(--bg-secondary);
}

.commit-avatar {
  flex-shrink: 0;
}

.avatar-img {
  width: 32px;
  height: 32px;
  border-radius: 50%;
}

.avatar-placeholder {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--accent-color);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
}

.commit-info {
  flex: 1;
  min-width: 0;
}

.commit-message {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.commit-meta {
  display: flex;
  gap: 8px;
  font-size: 12px;
  color: var(--text-secondary);
  margin-top: 2px;
}

.commit-sha code {
  font-size: 12px;
  padding: 2px 6px;
  background: var(--bg-secondary);
  border-radius: 4px;
  color: var(--accent-color);
}
</style>
