<template>
  <span :class="['badge', `badge-${variant}`, { 'badge-pill': pill }]">
    <slot></slot>
  </span>
</template>

<script setup lang="ts">
defineProps({
  variant: {
    type: String,
    default: 'default',
    validator: (value: string) => ['default', 'primary', 'success', 'danger', 'warning', 'info'].includes(value)
  },
  pill: {
    type: Boolean,
    default: false
  }
})
</script>

<style scoped>
.badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
  line-height: 1.5;
  white-space: nowrap;
}

.badge-pill {
  border-radius: 12px;
}

.badge-default {
  background: var(--bg-tertiary);
  color: var(--text-secondary);
}

.badge-primary {
  background: var(--accent-color);
  color: #fff;
}

.badge-success {
  background: var(--success-color);
  color: #fff;
}

.badge-danger {
  background: var(--danger-color);
  color: #fff;
}

.badge-warning {
  background: var(--warning-color);
  color: var(--text-primary);
}

.badge-info {
  background: var(--info-color);
  color: #fff;
}
</style>
