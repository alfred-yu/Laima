<template>
  <div class="issue-detail">
    <h1>{{ issue.title }}</h1>
    <div class="issue-meta">
      <span class="issue-author">由 {{ issue.author }} 创建</span>
      <span class="issue-time">{{ issue.time }}</span>
      <span class="issue-status" :class="issue.status">{{ issue.status }}</span>
    </div>
    <div class="issue-labels">
      <span v-for="label in issue.labels" :key="label" class="issue-label">{{ label }}</span>
    </div>
    <div class="issue-content">
      <div class="issue-description">
        <p>{{ issue.description }}</p>
      </div>
      <div class="issue-comments">
        <h2>评论</h2>
        <div class="comment-list">
          <div v-for="comment in issue.comments" :key="comment.id" class="comment-item">
            <div class="comment-author">{{ comment.author }}</div>
            <div class="comment-content">{{ comment.content }}</div>
            <div class="comment-time">{{ comment.time }}</div>
          </div>
        </div>
        <div class="comment-form">
          <textarea placeholder="添加评论..."></textarea>
          <button class="btn primary">提交</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const issue = ref({
  id: 1,
  title: 'Bug in login',
  author: 'user1',
  time: '2小时前',
  status: 'open',
  labels: ['bug', 'high'],
  description: 'There is a bug in the login functionality. Users cannot log in with valid credentials.',
  comments: [
    {
      id: 1,
      author: 'user2',
      content: 'I can reproduce this issue. Let me investigate.',
      time: '1小时前'
    },
    {
      id: 2,
      author: 'user3',
      content: 'Found the bug! It\'s in the authentication service.',
      time: '30分钟前'
    }
  ]
})
</script>

<style scoped>
.issue-detail {
  padding: 20px;
}

.issue-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin: 16px 0;
  font-size: 14px;
  color: var(--text-secondary);
  transition: color 0.3s;
}

.issue-status {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}

.issue-status.open {
  background: var(--success-color);
  color: #fff;
}

.issue-status.closed {
  background: var(--danger-color);
  color: #fff;
}

.issue-labels {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin: 16px 0;
}

.issue-label {
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  background: var(--border-color);
  color: var(--text-secondary);
  transition: background-color 0.3s, color 0.3s;
}

.issue-content {
  margin-top: 24px;
}

.issue-description {
  background: var(--bg-primary);
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 24px;
  transition: background-color 0.3s;
}

.issue-description p {
  margin: 0;
  line-height: 1.6;
  color: var(--text-primary);
  transition: color 0.3s;
}

.issue-comments h2 {
  margin: 0 0 16px 0;
  font-size: 18px;
  color: var(--text-primary);
  transition: color 0.3s;
}

.comment-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 24px;
}

.comment-item {
  background: var(--bg-secondary);
  padding: 16px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: background-color 0.3s, box-shadow 0.3s;
}

.comment-author {
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
  transition: color 0.3s;
}

.comment-content {
  line-height: 1.6;
  color: var(--text-primary);
  margin-bottom: 8px;
  transition: color 0.3s;
}

.comment-time {
  font-size: 12px;
  color: var(--text-secondary);
  transition: color 0.3s;
}

.comment-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.comment-form textarea {
  width: 100%;
  padding: 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 14px;
  resize: vertical;
  min-height: 100px;
  background: var(--bg-primary);
  color: var(--text-primary);
  transition: border-color 0.3s, background-color 0.3s, color 0.3s;
}

.btn {
  padding: 8px 16px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--bg-primary);
  color: var(--text-primary);
  cursor: pointer;
  transition: all 0.3s;
  align-self: flex-start;
}

.btn.primary {
  background: var(--accent-color);
  color: #fff;
  border-color: var(--accent-color);
}

.btn.primary:hover {
  opacity: 0.9;
}
</style>