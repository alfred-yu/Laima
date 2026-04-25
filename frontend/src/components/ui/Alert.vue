<template>
  <Transition name="alert">
    <div v-if="visible" :class="['alert', `alert-${variant}`]">
      <span class="alert-icon">{{ iconMap[variant] }}</span>
      <div class="alert-content">
        <div v-if="title" class="alert-title">{{ title }}</div>
        <div class="alert-message"><slot></slot></div>
      </div>
      <button v-if="closable" class="alert-close" @click="close">&times;</button>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps({
  variant: {
    type: String,
    default: 'info',
    validator: (value: string) => ['info', 'success', 'warning', 'danger'].includes(value)
  },
  title: {
    type: String,
    default: ''
  },
  closable: {
    type: Boolean,
    default: false
  }
})

const visible = ref(true)

const iconMap: Record<string, string> = {
  info: 'ℹ',
  success: '✓',
  warning: '⚠',
  danger: '✕'
}

const emit = defineEmits(['close'])

function close() {
  visible.value = false
  emit('close')
}
</script>

<style scoped>
.alert {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 12px 16px;
  border-radius: 6px;
  border: 1px solid;
  font-size: 14px;
}

.alert-icon {
  font-size: 16px;
  line-height: 1.5;
  flex-shrink: 0;
}

.alert-content {
  flex: 1;
}

.alert-title {
  font-weight: 600;
  margin-bottom: 4px;
}

.alert-message {
  line-height: 1.5;
}

.alert-close {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  color: inherit;
  opacity: 0.6;
  padding: 0;
  line-height: 1;
}

.alert-close:hover {
  opacity: 1;
}

.alert-info {
  background: var(--info-color);
  border-color: var(--info-color);
  color: #fff;
  opacity: 0.9;
}

.alert-success {
  background: var(--success-color);
  border-color: var(--success-color);
  color: #fff;
  opacity: 0.9;
}

.alert-warning {
  background: var(--warning-color);
  border-color: var(--warning-color);
  color: var(--text-primary);
  opacity: 0.9;
}

.alert-danger {
  background: var(--danger-color);
  border-color: var(--danger-color);
  color: #fff;
  opacity: 0.9;
}

.alert-enter-active,
.alert-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.alert-enter-from,
.alert-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
