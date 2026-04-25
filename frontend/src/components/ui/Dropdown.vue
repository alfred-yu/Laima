<template>
  <div :class="['dropdown', { 'dropdown-open': isOpen }]" ref="dropdownRef">
    <div class="dropdown-trigger" @click="toggle">
      <slot name="trigger">
        <button class="dropdown-button">{{ label }}</button>
      </slot>
    </div>
    <Transition name="dropdown">
      <div v-if="isOpen" :class="['dropdown-menu', `dropdown-menu-${align}`]">
        <slot></slot>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  label: {
    type: String,
    default: 'Options'
  },
  align: {
    type: String,
    default: 'left',
    validator: (value: string) => ['left', 'right'].includes(value)
  }
})

const isOpen = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)

function toggle() {
  isOpen.value = !isOpen.value
}

function close(event: MouseEvent) {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target as Node)) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', close)
})

onUnmounted(() => {
  document.removeEventListener('click', close)
})
</script>

<style scoped>
.dropdown {
  position: relative;
  display: inline-block;
}

.dropdown-button {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  color: var(--text-primary);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.dropdown-button:hover {
  background: var(--bg-tertiary);
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  margin-top: 4px;
  min-width: 160px;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  z-index: 100;
  padding: 4px 0;
}

.dropdown-menu-left { left: 0; }
.dropdown-menu-right { right: 0; }

.dropdown-enter-active,
.dropdown-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
