<template>
  <button 
    :class="[
      'btn', 
      `btn-${variant}`, 
      { 'btn-block': block },
      { 'btn-sm': size === 'sm' },
      { 'btn-lg': size === 'lg' }
    ]"
    :disabled="disabled"
    @click="$emit('click')"
  >
    <slot></slot>
  </button>
</template>

<script setup lang="ts">

const props = defineProps({
  variant: {
    type: String,
    default: 'default',
    validator: (value: string) => ['default', 'primary', 'success', 'danger', 'warning', 'info'].includes(value)
  },
  size: {
    type: String,
    default: 'md',
    validator: (value: string) => ['sm', 'md', 'lg'].includes(value)
  },
  block: {
    type: Boolean,
    default: false
  },
  disabled: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['click'])
</script>

<style scoped>
.btn {
  display: inline-block;
  padding: 8px 16px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  text-align: center;
}

.btn:hover {
  opacity: 0.8;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-block {
  display: block;
  width: 100%;
}

.btn-sm {
  padding: 4px 8px;
  font-size: 12px;
}

.btn-lg {
  padding: 12px 24px;
  font-size: 16px;
}

.btn-default {
  background: var(--bg-secondary);
  color: var(--text-primary);
  border-color: var(--border-color);
}

.btn-primary {
  background: var(--accent-color);
  color: #fff;
  border-color: var(--accent-color);
}

.btn-success {
  background: var(--success-color);
  color: #fff;
  border-color: var(--success-color);
}

.btn-danger {
  background: var(--danger-color);
  color: #fff;
  border-color: var(--danger-color);
}

.btn-warning {
  background: var(--warning-color);
  color: var(--text-primary);
  border-color: var(--warning-color);
}

.btn-info {
  background: var(--info-color);
  color: #fff;
  border-color: var(--info-color);
}
</style>