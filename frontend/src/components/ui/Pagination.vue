<template>
  <div class="pagination">
    <button
      class="page-btn"
      :disabled="modelValue <= 1"
      @click="changePage(modelValue - 1)"
    >
      &laquo;
    </button>
    <button
      v-for="page in displayedPages"
      :key="page"
      :class="['page-btn', { 'page-active': page === modelValue, 'page-dots': page === '...' }]"
      :disabled="page === '...'"
      @click="page !== '...' && changePage(page as number)"
    >
      {{ page }}
    </button>
    <button
      class="page-btn"
      :disabled="modelValue >= totalPages"
      @click="changePage(modelValue + 1)"
    >
      &raquo;
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  modelValue: number
  total: number
  pageSize: number
  maxDisplayed?: number
}>()

const emit = defineEmits(['update:modelValue'])

const totalPages = computed(() => Math.ceil(props.total / props.pageSize))

const displayedPages = computed(() => {
  const pages: (number | string)[] = []
  const max = props.maxDisplayed || 7
  const total = totalPages.value
  const current = props.modelValue

  if (total <= max) {
    for (let i = 1; i <= total; i++) pages.push(i)
    return pages
  }

  pages.push(1)

  let start = Math.max(2, current - 1)
  let end = Math.min(total - 1, current + 1)

  if (current <= 3) {
    end = Math.min(max - 1, total - 1)
  }

  if (current >= total - 2) {
    start = Math.max(2, total - max + 2)
  }

  if (start > 2) pages.push('...')

  for (let i = start; i <= end; i++) pages.push(i)

  if (end < total - 1) pages.push('...')

  pages.push(total)

  return pages
})

function changePage(page: number) {
  if (page >= 1 && page <= totalPages.value) {
    emit('update:modelValue', page)
  }
}
</script>

<style scoped>
.pagination {
  display: flex;
  align-items: center;
  gap: 4px;
}

.page-btn {
  min-width: 32px;
  height: 32px;
  padding: 0 8px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.page-btn:hover:not(:disabled):not(.page-active):not(.page-dots) {
  background: var(--bg-secondary);
  border-color: var(--accent-color);
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-active {
  background: var(--accent-color);
  color: #fff;
  border-color: var(--accent-color);
}

.page-dots {
  border: none;
  background: none;
  cursor: default;
}
</style>
