<template>
  <Teleport to="body">
    <div class="toast-container">
      <TransitionGroup name="toast">
        <div
          v-for="toast in toasts"
          :key="toast.id"
          :class="['toast', `toast-${toast.type}`]"
        >
          <span class="toast-icon">{{ iconMap[toast.type] }}</span>
          <span class="toast-message">{{ toast.message }}</span>
          <button class="toast-close" @click="removeToast(toast.id)">&times;</button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface Toast {
  id: number
  type: 'success' | 'danger' | 'warning' | 'info'
  message: string
}

const toasts = ref<Toast[]>([])
let nextId = 0

const iconMap: Record<string, string> = {
  success: '✓',
  danger: '✕',
  warning: '⚠',
  info: 'ℹ'
}

function addToast(type: Toast['type'], message: string, duration = 3000) {
  const id = nextId++
  toasts.value.push({ id, type, message })
  if (duration > 0) {
    setTimeout(() => removeToast(id), duration)
  }
}

function removeToast(id: number) {
  toasts.value = toasts.value.filter(t => t.id !== id)
}

function success(message: string, duration?: number) { addToast('success', message, duration) }
function danger(message: string, duration?: number) { addToast('danger', message, duration) }
function warning(message: string, duration?: number) { addToast('warning', message, duration) }
function info(message: string, duration?: number) { addToast('info', message, duration) }

defineExpose({ success, danger, warning, info, addToast, removeToast })
</script>

<style scoped>
.toast-container {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 2000;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.toast {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  border-radius: 6px;
  min-width: 280px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  font-size: 14px;
}

.toast-icon { font-size: 14px; flex-shrink: 0; }
.toast-message { flex: 1; }

.toast-close {
  background: none;
  border: none;
  font-size: 16px;
  cursor: pointer;
  color: inherit;
  opacity: 0.6;
  padding: 0;
}

.toast-close:hover { opacity: 1; }

.toast-success { background: var(--success-color); color: #fff; }
.toast-danger { background: var(--danger-color); color: #fff; }
.toast-warning { background: var(--warning-color); color: var(--text-primary); }
.toast-info { background: var(--info-color); color: #fff; }

.toast-enter-active { transition: all 0.3s ease; }
.toast-leave-active { transition: all 0.2s ease; }
.toast-enter-from { opacity: 0; transform: translateX(40px); }
.toast-leave-to { opacity: 0; transform: translateX(40px); }
</style>
