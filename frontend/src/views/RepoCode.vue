<template>
  <div class="repo-code">
    <div class="code-header">
      <div class="branch-selector">
        <select v-model="selectedBranch">
          <option v-for="branch in branches" :key="branch" :value="branch">{{ branch }}</option>
        </select>
      </div>
      <div class="code-actions">
        <button class="btn">上传文件</button>
        <button class="btn">新建文件</button>
        <button class="btn">上传文件夹</button>
      </div>
    </div>
    <div class="code-content">
      <div class="file-tree">
        <div class="file-tree-header">
          <h3>文件</h3>
        </div>
        <ul class="file-list">
          <li v-for="file in files" :key="file.path" class="file-item">
            <span class="file-icon">{{ file.type === 'dir' ? '📁' : '📄' }}</span>
            <span class="file-name">{{ file.name }}</span>
          </li>
        </ul>
      </div>
      <div class="file-content">
        <div class="file-header">
          <h2>{{ currentFile.name }}</h2>
          <div class="file-actions">
            <button class="btn">编辑</button>
            <button class="btn">删除</button>
            <button class="btn">历史</button>
          </div>
        </div>
        <div class="code-editor">
          <pre>{{ currentFile.content }}</pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const selectedBranch = ref('main')
const branches = ref(['main', 'dev', 'feature-1', 'feature-2'])

const files = ref([
  { path: 'src/', name: 'src', type: 'dir' },
  { path: 'src/main.ts', name: 'main.ts', type: 'file' },
  { path: 'src/App.vue', name: 'App.vue', type: 'file' },
  { path: 'src/components/', name: 'components', type: 'dir' },
  { path: 'src/views/', name: 'views', type: 'dir' },
  { path: 'src/stores/', name: 'stores', type: 'dir' },
  { path: 'src/router/', name: 'router', type: 'dir' },
  { path: 'package.json', name: 'package.json', type: 'file' },
  { path: 'tsconfig.json', name: 'tsconfig.json', type: 'file' },
  { path: 'vite.config.ts', name: 'vite.config.ts', type: 'file' }
])

const currentFile = ref({
  name: 'package.json',
  content: `{
  "name": "laima-frontend",
  "private": true,
  "version": "0.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vue-tsc && vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "vue": "^3.5.13",
    "vue-router": "^4.4.5",
    "pinia": "^2.3.0",
    "axios": "^1.7.9"
  },
  "devDependencies": {
    "@vitejs/plugin-vue": "^5.2.1",
    "typescript": "~5.6.2",
    "vite": "^6.0.5",
    "vue-tsc": "^2.1.10"
  }
}`
})
</script>

<style scoped>
.repo-code {
  padding: 20px;
}

.code-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.branch-selector select {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.code-actions {
  display: flex;
  gap: 10px;
}

.btn {
  padding: 6px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: #fff;
  color: #333;
  cursor: pointer;
  transition: all 0.3s;
}

.btn:hover {
  background: #f5f5f5;
}

.code-content {
  display: flex;
  gap: 20px;
  min-height: 600px;
}

.file-tree {
  width: 300px;
  background: #f8f9fa;
  border-radius: 8px;
  padding: 20px;
}

.file-tree-header h3 {
  margin: 0 0 20px 0;
  font-size: 16px;
  color: #333;
}

.file-list {
  list-style: none;
  padding: 0;
}

.file-item {
  padding: 8px 12px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  display: flex;
  align-items: center;
}

.file-item:hover {
  background: #e9ecef;
}

.file-icon {
  margin-right: 8px;
  font-size: 16px;
}

.file-name {
  font-size: 14px;
  color: #333;
}

.file-content {
  flex: 1;
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.file-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid #ddd;
}

.file-header h2 {
  margin: 0;
  font-size: 18px;
  color: #333;
}

.file-actions {
  display: flex;
  gap: 10px;
}

.code-editor {
  background: #f8f9fa;
  border-radius: 4px;
  padding: 20px;
  overflow: auto;
  max-height: 500px;
}

.code-editor pre {
  margin: 0;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.5;
  color: #333;
}
</style>