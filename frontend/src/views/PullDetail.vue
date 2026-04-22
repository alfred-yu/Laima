<template>
  <div class="pull-detail">
    <h1>{{ pull.title }}</h1>
    <div class="pull-meta">
      <span class="pull-author">由 {{ pull.author }} 创建</span>
      <span class="pull-time">{{ pull.time }}</span>
      <span class="pull-branch">{{ pull.source }} → {{ pull.target }}</span>
      <span class="pull-status" :class="pull.status">{{ pull.status }}</span>
    </div>
    <div class="pull-content">
      <div class="pull-description">
        <p>{{ pull.description }}</p>
      </div>
      <div class="pull-diff">
        <h2>代码变更</h2>
        <div class="diff-file">
          <div class="diff-header">
            <span class="file-name">package.json</span>
            <span class="file-status">修改</span>
          </div>
          <div class="diff-content">
            <pre>{{ pull.diff }}</pre>
          </div>
        </div>
      </div>
      <div class="pull-comments">
        <h2>评论</h2>
        <div class="comment-list">
          <div v-for="comment in pull.comments" :key="comment.id" class="comment-item">
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

const pull = ref({
  id: 1,
  title: 'Add new feature',
  author: 'user1',
  time: '2小时前',
  source: 'feature-1',
  target: 'main',
  status: 'open',
  description: 'This PR adds a new feature to the project.',
  diff: `{\n  "name": "laima-frontend",\n  "private": true,\n  "version": "0.0.0",\n  "type": "module",\n  "scripts": {\n    "dev": "vite",\n    "build": "vue-tsc && vite build",\n    "preview": "vite preview"\n  },\n  "dependencies": {\n    "vue": "^3.5.13",\n    "vue-router": "^4.4.5",\n    "pinia": "^2.3.0",\n    "axios": "^1.7.9",\n+   "monaco-editor": "^0.52.2"\n  },\n  "devDependencies": {\n    "@vitejs/plugin-vue": "^5.2.1",\n    "typescript": "~5.6.2",\n    "vite": "^6.0.5",\n    "vue-tsc": "^2.1.10"\n  }\n}`,
  comments: [
    {
      id: 1,
      author: 'user2',
      content: 'Looks good!',
      time: '1小时前'
    },
    {
      id: 2,
      author: 'user3',
      content: 'Please add tests for this feature.',
      time: '30分钟前'
    }
  ]
})
</script>

<style scoped>
.pull-detail {
  padding: 20px;
}

.pull-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin: 16px 0;
  font-size: 14px;
  color: #666;
}

.pull-status {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}

.pull-status.open {
  background: #d1e7dd;
  color: #0f5132;
}

.pull-status.merged {
  background: #cfe2ff;
  color: #084298;
}

.pull-status.closed {
  background: #f8d7da;
  color: #842029;
}

.pull-content {
  margin-top: 24px;
}

.pull-description {
  background: #f8f9fa;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 24px;
}

.pull-description p {
  margin: 0;
  line-height: 1.6;
  color: #333;
}

.pull-diff {
  margin-bottom: 24px;
}

.pull-diff h2 {
  margin: 0 0 16px 0;
  font-size: 18px;
  color: #333;
}

.diff-file {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.diff-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f8f9fa;
  border-bottom: 1px solid #ddd;
}

.file-name {
  font-weight: 600;
  color: #333;
}

.file-status {
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  background: #d1e7dd;
  color: #0f5132;
}

.diff-content {
  padding: 16px;
  overflow-x: auto;
}

.diff-content pre {
  margin: 0;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.5;
  color: #333;
}

.pull-comments h2 {
  margin: 0 0 16px 0;
  font-size: 18px;
  color: #333;
}

.comment-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 24px;
}

.comment-item {
  background: #fff;
  padding: 16px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.comment-author {
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
}

.comment-content {
  line-height: 1.6;
  color: #333;
  margin-bottom: 8px;
}

.comment-time {
  font-size: 12px;
  color: #999;
}

.comment-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.comment-form textarea {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  resize: vertical;
  min-height: 100px;
}

.btn {
  padding: 8px 16px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: #fff;
  color: #333;
  cursor: pointer;
  transition: all 0.3s;
  align-self: flex-start;
}

.btn.primary {
  background: #0366d6;
  color: #fff;
  border-color: #0366d6;
}

.btn.primary:hover {
  background: #0056b3;
}
</style>