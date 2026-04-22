<template>
  <div class="card" :class="{ 'card-shadow': shadow }">
    <div v-if="title" class="card-header">
      <h3 class="card-title">{{ title }}</h3>
    </div>
    <div class="card-body">
      <slot></slot>
    </div>
    <div v-if="$slots.footer" class="card-footer">
      <slot name="footer"></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue'

const props = defineProps({
  title: {
    type: String,
    default: ''
  },
  shadow: {
    type: Boolean,
    default: true
  }
})
</script>

<style scoped>
.card {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #e2e8f0;
}

.card-shadow {
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.card-header {
  padding: 16px 20px;
  border-bottom: 1px solid #e2e8f0;
  background: #f8f9fa;
}

.card-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.card-body {
  padding: 20px;
}

.card-footer {
  padding: 16px 20px;
  border-top: 1px solid #e2e8f0;
  background: #f8f9fa;
}
</style>