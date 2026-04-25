<template>
  <div :class="['skeleton', `skeleton-${variant}`, { 'skeleton-animate': animate }]"></div>
</template>

<script setup lang="ts">
defineProps({
  variant: {
    type: String,
    default: 'text',
    validator: (value: string) => ['text', 'title', 'avatar', 'button', 'image', 'card'].includes(value)
  },
  animate: {
    type: Boolean,
    default: true
  }
})
</script>

<style scoped>
.skeleton {
  background: var(--bg-tertiary);
  border-radius: 4px;
}

.skeleton-text {
  width: 100%;
  height: 14px;
}

.skeleton-title {
  width: 60%;
  height: 20px;
}

.skeleton-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
}

.skeleton-button {
  width: 80px;
  height: 32px;
  border-radius: 4px;
}

.skeleton-image {
  width: 100%;
  height: 200px;
}

.skeleton-card {
  width: 100%;
  height: 120px;
  border-radius: 8px;
}

.skeleton-animate {
  animation: skeleton-pulse 1.5s ease-in-out infinite;
}

@keyframes skeleton-pulse {
  0% { opacity: 1; }
  50% { opacity: 0.4; }
  100% { opacity: 1; }
}
</style>
