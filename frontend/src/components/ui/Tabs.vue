<template>
  <div class="tabs">
    <div class="tabs-nav">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        :class="['tab-item', { 'tab-active': modelValue === tab.key }]"
        @click="selectTab(tab.key)"
      >
        {{ tab.label }}
        <span v-if="tab.count !== undefined" class="tab-count">{{ tab.count }}</span>
      </button>
    </div>
    <div class="tabs-content">
      <slot></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
export interface Tab {
  key: string
  label: string
  count?: number
}

const props = defineProps<{
  tabs: Tab[]
  modelValue: string
}>()

const emit = defineEmits(['update:modelValue'])

function selectTab(key: string) {
  emit('update:modelValue', key)
}
</script>

<style scoped>
.tabs {
  width: 100%;
}

.tabs-nav {
  display: flex;
  border-bottom: 2px solid var(--border-color);
  gap: 0;
}

.tab-item {
  padding: 10px 20px;
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  margin-bottom: -2px;
  color: var(--text-secondary);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 6px;
}

.tab-item:hover {
  color: var(--text-primary);
}

.tab-active {
  color: var(--accent-color);
  border-bottom-color: var(--accent-color);
}

.tab-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 20px;
  height: 20px;
  padding: 0 6px;
  border-radius: 10px;
  background: var(--bg-tertiary);
  font-size: 12px;
  color: var(--text-secondary);
}

.tab-active .tab-count {
  background: var(--accent-color);
  color: #fff;
}

.tabs-content {
  padding: 16px 0;
}
</style>
