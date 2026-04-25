<template>
  <div class="review-panel">
    <div class="review-header">
      <h3 class="review-title">Reviews</h3>
      <slot name="actions"></slot>
    </div>
    <div class="review-list">
      <div v-for="review in reviews" :key="review.id" :class="['review-item', `review-${review.state}`]">
        <div class="review-author">
          <div class="review-avatar">{{ review.reviewer.charAt(0).toUpperCase() }}</div>
          <span class="review-name">{{ review.reviewer }}</span>
        </div>
        <span :class="['review-state', `state-${review.state}`]">
          {{ stateLabel(review.state) }}
        </span>
        <div v-if="review.body" class="review-body">{{ review.body }}</div>
        <div class="review-time">{{ review.time }}</div>
      </div>
    </div>
    <div v-if="reviews.length === 0" class="review-empty">
      No reviews yet
    </div>
  </div>
</template>

<script setup lang="ts">
export interface Review {
  id: number
  reviewer: string
  state: 'approved' | 'changes_requested' | 'commented' | 'pending'
  body: string
  time: string
}

defineProps<{
  reviews: Review[]
}>()

function stateLabel(state: string): string {
  const map: Record<string, string> = {
    approved: 'Approved',
    changes_requested: 'Changes Requested',
    commented: 'Commented',
    pending: 'Pending'
  }
  return map[state] || state
}
</script>

<style scoped>
.review-panel {
  border: 1px solid var(--border-color);
  border-radius: 6px;
  overflow: hidden;
}

.review-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
}

.review-title {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
}

.review-item {
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
}

.review-item:last-child { border-bottom: none; }

.review-author {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.review-avatar {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: var(--accent-color);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
}

.review-name {
  font-size: 14px;
  font-weight: 500;
}

.review-state {
  display: inline-block;
  font-size: 12px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 10px;
  margin-top: 4px;
}

.state-approved { background: var(--success-color); color: #fff; }
.state-changes_requested { background: var(--danger-color); color: #fff; }
.state-commented { background: var(--info-color); color: #fff; }
.state-pending { background: var(--bg-tertiary); color: var(--text-secondary); }

.review-body {
  font-size: 14px;
  color: var(--text-primary);
  margin-top: 8px;
  padding: 8px 12px;
  background: var(--bg-secondary);
  border-radius: 4px;
}

.review-time {
  font-size: 12px;
  color: var(--text-secondary);
  margin-top: 4px;
}

.review-empty {
  padding: 24px;
  text-align: center;
  color: var(--text-secondary);
  font-size: 14px;
}
</style>
