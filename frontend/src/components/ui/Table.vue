<template>
  <div class="table-wrapper">
    <table :class="['table', { 'table-striped': striped, 'table-bordered': bordered, 'table-hover': hover }]">
      <thead>
        <tr>
          <th v-for="col in columns" :key="col.key" :style="{ width: col.width }">
            {{ col.label }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="data.length === 0">
          <td :colspan="columns.length" class="table-empty">
            <slot name="empty">No data available</slot>
          </td>
        </tr>
        <tr v-for="(row, index) in data" :key="index" @click="$emit('rowClick', row)">
          <td v-for="col in columns" :key="col.key">
            <slot :name="col.key" :row="row" :index="index">
              {{ row[col.key] }}
            </slot>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
export interface TableColumn {
  key: string
  label: string
  width?: string
}

defineProps<{
  columns: TableColumn[]
  data: Record<string, any>[]
  striped?: boolean
  bordered?: boolean
  hover?: boolean
}>()

defineEmits(['rowClick'])
</script>

<style scoped>
.table-wrapper {
  overflow-x: auto;
  border: 1px solid var(--border-color);
  border-radius: 6px;
}

.table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}

.table th,
.table td {
  padding: 10px 16px;
  text-align: left;
  border-bottom: 1px solid var(--border-color);
}

.table th {
  background: var(--bg-secondary);
  font-weight: 600;
  color: var(--text-secondary);
  white-space: nowrap;
}

.table td {
  color: var(--text-primary);
}

.table-striped tbody tr:nth-child(even) {
  background: var(--bg-secondary);
}

.table-hover tbody tr:hover {
  background: var(--bg-tertiary);
  cursor: pointer;
}

.table-bordered th,
.table-bordered td {
  border: 1px solid var(--border-color);
}

.table-empty {
  text-align: center;
  padding: 40px 16px;
  color: var(--text-secondary);
}
</style>
