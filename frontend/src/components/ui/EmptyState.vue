<template>
  <div class="empty-state">
    <div class="empty-icon">
      <slot name="icon">
        <span class="empty-icon-default">📭</span>
      </slot>
    </div>
    <h3 class="empty-title">{{ title }}</h3>
    <p v-if="description" class="empty-description">{{ description }}</p>
    <div v-if="$slots.action" class="empty-action">
      <slot name="action"></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  description?: string
}>()
</script>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.empty-icon-default {
  font-size: 48px;
}

.empty-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 8px;
}

.empty-description {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0;
  max-width: 400px;
}

.empty-action {
  margin-top: 20px;
}
</style>
