<template>
  <div :class="['merge-status', `merge-${status}`]">
    <span class="merge-icon">{{ iconMap[status] }}</span>
    <span class="merge-text">{{ labelMap[status] }}</span>
    <slot></slot>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  status: 'mergeable' | 'conflict' | 'checking' | 'blocked'
}>()

const iconMap: Record<string, string> = {
  mergeable: '✓',
  conflict: '✕',
  checking: '⟳',
  blocked: '⊘'
}

const labelMap: Record<string, string> = {
  mergeable: 'No conflicts, can be merged',
  conflict: 'Merge conflicts detected',
  checking: 'Checking merge status...',
  blocked: 'Merge is blocked'
}
</script>

<style scoped>
.merge-status {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
}

.merge-icon {
  font-size: 16px;
}

.merge-mergeable {
  background: rgba(46, 160, 67, 0.1);
  color: var(--success-color);
  border: 1px solid var(--success-color);
}

.merge-conflict {
  background: rgba(248, 81, 73, 0.1);
  color: var(--danger-color);
  border: 1px solid var(--danger-color);
}

.merge-checking {
  background: rgba(56, 132, 255, 0.1);
  color: var(--info-color);
  border: 1px solid var(--info-color);
}

.merge-checking .merge-icon {
  animation: spin 1s linear infinite;
}

.merge-blocked {
  background: rgba(158, 142, 106, 0.1);
  color: var(--warning-color);
  border: 1px solid var(--warning-color);
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
