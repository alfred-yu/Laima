<template>
  <div class="branch-selector">
    <Dropdown :label="selectedBranch || 'Select branch'" align="left">
      <template #trigger>
        <button class="branch-btn">
          <span class="branch-icon">⌥</span>
          <span class="branch-name">{{ selectedBranch || 'Select branch' }}</span>
          <span class="branch-arrow">▾</span>
        </button>
      </template>
      <div class="branch-dropdown">
        <div class="branch-search">
          <input
            v-model="search"
            class="branch-search-input"
            placeholder="Filter branches..."
          />
        </div>
        <div class="branch-list">
          <div
            v-for="branch in filteredBranches"
            :key="branch"
            :class="['branch-option', { 'branch-active': branch === selectedBranch }]"
            @click="selectBranch(branch)"
          >
            <span class="branch-icon">⌥</span>
            <span>{{ branch }}</span>
            <span v-if="branch === defaultBranch" class="branch-default">default</span>
          </div>
        </div>
      </div>
    </Dropdown>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import Dropdown from '../ui/Dropdown.vue'

const props = defineProps<{
  branches: string[]
  selectedBranch: string
  defaultBranch?: string
}>()

const emit = defineEmits(['update:selectedBranch'])

const search = ref('')

const filteredBranches = computed(() => {
  if (!search.value) return props.branches
  return props.branches.filter(b => b.toLowerCase().includes(search.value.toLowerCase()))
})

function selectBranch(branch: string) {
  emit('update:selectedBranch', branch)
}
</script>

<style scoped>
.branch-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  color: var(--text-primary);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.branch-btn:hover {
  background: var(--bg-tertiary);
}

.branch-icon {
  color: var(--text-secondary);
}

.branch-arrow {
  font-size: 12px;
  color: var(--text-secondary);
}

.branch-dropdown {
  min-width: 240px;
  padding: 8px;
}

.branch-search {
  margin-bottom: 8px;
}

.branch-search-input {
  width: 100%;
  padding: 6px 10px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 13px;
  outline: none;
}

.branch-search-input:focus {
  border-color: var(--accent-color);
}

.branch-list {
  max-height: 240px;
  overflow-y: auto;
}

.branch-option {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
}

.branch-option:hover {
  background: var(--bg-secondary);
}

.branch-active {
  background: var(--accent-color);
  color: #fff;
}

.branch-default {
  font-size: 11px;
  padding: 1px 6px;
  border-radius: 8px;
  background: var(--bg-tertiary);
  color: var(--text-secondary);
  margin-left: auto;
}
</style>
